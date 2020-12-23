package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/judascrow/fiber-boilerplate/app/models"
	"github.com/judascrow/fiber-boilerplate/database"
)

// Return all roles as JSON
func GetAllRoles(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var roles []models.Role
		if response := db.Find(&roles); response.Error != nil {
			// panic("Error occurred while retrieving roles from the database: " + response.Error.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   response.Error.Error(),
				"message": "Error occurred while retrieving roles from the database",
			})
		}
		err := c.JSON(roles)
		if err != nil {
			// panic("Error occurred when returning JSON of roles: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error occurred when returning JSON of roles",
			})
		}
		return err
	}
}

// Return a single role as JSON
func GetRole(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "ID must be numbers only",
			})
		}
		role := new(models.Role)
		if response := db.Find(&role, id); response.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   response.Error.Error(),
				"message": "An error occurred when retrieving the role",
			})
		}
		if role.ID == 0 {
			// Send status not found
			err := c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "record not found",
				"message": "ID " + c.Params("id") + " not found",
			})
			return err
		}

		c.JSON(role)
		return err
	}
}
