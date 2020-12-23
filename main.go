package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/judascrow/fiber-boilerplate/app/middlewares"
	"github.com/judascrow/fiber-boilerplate/app/models"
	"github.com/judascrow/fiber-boilerplate/app/routes"
	configuration "github.com/judascrow/fiber-boilerplate/config"
	"github.com/judascrow/fiber-boilerplate/database"
)

type App struct {
	*fiber.App

	DB *database.Database
}

func main() {
	config := configuration.New()

	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}

	app.registerMiddlewares(config)

	// Initialize database
	db, err := database.New(&database.DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER"),
		Host:     config.GetString("DB_HOST"),
		Username: config.GetString("DB_USERNAME"),
		Password: config.GetString("DB_PASSWORD"),
		Port:     config.GetInt("DB_PORT"),
		Database: config.GetString("DB_DATABASE"),
	})

	// Auto-migrate database models
	if err != nil {
		fmt.Println("failed to connect to database:", err.Error())
	} else {
		if db == nil {
			fmt.Println("failed to connect to database: db variable is nil")
		} else {
			app.DB = db
			err := app.DB.AutoMigrate(&models.Role{})
			if err != nil {
				fmt.Println("failed to automigrate role model:", err.Error())
				return
			}
			err = app.DB.AutoMigrate(&models.User{})
			if err != nil {
				fmt.Println("failed to automigrate user model:", err.Error())
				return
			}
		}
	}

	api := app.Group("/api")
	apiv1 := api.Group("/v1")
	routes.RegisterAPI(apiv1, app.DB)

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

	// Middleware - Favicon
	if config.GetBool("MW_FIBER_FAVICON_ENABLED") {
		app.Use(favicon.New(favicon.Config{
			File:         config.GetString("MW_FIBER_FAVICON_FILE"),
			CacheControl: config.GetString("MW_FIBER_FAVICON_CACHECONTROL"),
		}))
	}

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
	} else {
		app.Use(logger.New())
	}

	// Middleware - Force HTTPS
	if config.GetBool("MW_FORCE_HTTPS_ENABLED") {
		app.Use(middlewares.ForceHTTPS())
	}

	// Middleware - Force trailing slash
	if config.GetBool("MW_FORCE_TRAILING_SLASH_ENABLED") {
		app.Use(middlewares.ForceTrailingSlash())
	}

	// Middleware - Recover
	if config.GetBool("MW_FIBER_RECOVER_ENABLED") {
		app.Use(recover.New())
	}

}

// Stop the Fiber application
func (app *App) exit() {
	_ = app.Shutdown()
}
