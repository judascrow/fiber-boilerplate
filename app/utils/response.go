package utils

import "github.com/gofiber/fiber/v2"

func ResponseError(c *fiber.Ctx, statusCode int, message string, errors interface{}) (err error) {
	return c.Status(statusCode).JSON(fiber.Map{
		"errors":  errors,
		"message": message,
	})
}
