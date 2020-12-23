package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/judascrow/fiber-boilerplate/database"
)

func GetHealthcheck(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success": true,
			"message": "API Is Online",
		})
	}
}
