package main

import (
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/judascrow/fiber-boilerplate/app/middlewares"
	configuration "github.com/judascrow/fiber-boilerplate/config"
)

type App struct {
	*fiber.App
}

func main() {
	config := configuration.New()

	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}

	app.registerMiddlewares(config)

	// Initialize database
	var err error

	api := app.Group("/api")
	apiv1 := api.Group("/v1")
	apiv1.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API Is Online",
		})

	})

	// Custom 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		if err := c.SendStatus(fiber.StatusNotFound); err != nil {
			panic(err)
		}
		if err := c.Render("errors/404", fiber.Map{}); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return err
	})

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()

	// Start listening on the specified address
	err = app.Listen(config.GetString("APP_ADDR"))
	if err != nil {
		app.exit()
	}
}

func (app *App) registerMiddlewares(config *configuration.Config) {
	// Middleware - Custom Access Logger based on zap
	if config.GetBool("MW_ACCESS_LOGGER_ENABLED") {
		app.Use(middlewares.AccessLogger(&middlewares.AccessLoggerConfig{
			Type:        config.GetString("MW_ACCESS_LOGGER_TYPE"),
			Environment: config.GetString("APP_ENV"),
			Filename:    config.GetString("MW_ACCESS_LOGGER_FILENAME"),
			MaxSize:     config.GetInt("MW_ACCESS_LOGGER_MAXSIZE"),
			MaxAge:      config.GetInt("MW_ACCESS_LOGGER_MAXAGE"),
			MaxBackups:  config.GetInt("MW_ACCESS_LOGGER_MAXBACKUPS"),
			LocalTime:   config.GetBool("MW_ACCESS_LOGGER_LOCALTIME"),
			Compress:    config.GetBool("MW_ACCESS_LOGGER_COMPRESS"),
		}))
	}
}

// Stop the Fiber application
func (app *App) exit() {
	_ = app.Shutdown()
}
