package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper

	errorHandler fiber.ErrorHandler
	fiber        *fiber.Config
}

var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Set error message
	message := err.Error()

	// Check if it's a fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// TODO: Check return type for the client, JSON, HTML, YAML or any other (API vs web)

	// Return HTTP response
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	c.Status(code)

	// Render default error view
	err = c.Render("errors/"+strconv.Itoa(code), fiber.Map{"message": message})
	if err != nil {
		return c.SendString(message)
	}
	return err
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath(".")

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	config.SetErrorHandler(defaultErrorHandler)

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	return config
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_ADDR", ":8080")
	config.SetDefault("APP_ENV", "local")

	// Set default Custom Access Logger middleware configuration
	config.SetDefault("MW_ACCESS_LOGGER_ENABLED", true)
	config.SetDefault("MW_ACCESS_LOGGER_TYPE", "console")
	config.SetDefault("MW_ACCESS_LOGGER_FILENAME", "access.log")
	config.SetDefault("MW_ACCESS_LOGGER_MAXSIZE", 500)
	config.SetDefault("MW_ACCESS_LOGGER_MAXAGE", 28)
	config.SetDefault("MW_ACCESS_LOGGER_MAXBACKUPS", 3)
	config.SetDefault("MW_ACCESS_LOGGER_LOCALTIME", false)
	config.SetDefault("MW_ACCESS_LOGGER_COMPRESS", false)

	// Set default Fiber Favicon middleware configuration
	config.SetDefault("MW_FIBER_FAVICON_ENABLED", false)
	config.SetDefault("MW_FIBER_FAVICON_FILE", "")
	config.SetDefault("MW_FIBER_FAVICON_CACHECONTROL", "public, max-age=31536000")
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:                   config.GetBool("FIBER_PREFORK"),
		ServerHeader:              config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:             config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:             config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:                 config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:              config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                      config.GetBool("FIBER_ETAG"),
		BodyLimit:                 config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:               config.GetInt("FIBER_CONCURRENCY"),
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
	}
}
