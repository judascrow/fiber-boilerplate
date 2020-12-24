package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/judascrow/fiber-boilerplate/app/dtos"
	"github.com/judascrow/fiber-boilerplate/app/models"
	"github.com/judascrow/fiber-boilerplate/app/utils"
	"github.com/judascrow/fiber-boilerplate/database"
)

// Return all roles as JSON
func GetAllRoles(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var roles []models.Role

		condition := new(dtos.RoleDto)
		if err := c.QueryParser(condition); err != nil {
			return utils.ResponseError(c, 500, fiber.ErrInternalServerError.Message, err.Error())
		}

		if response := db.Where(condition).Find(&roles); response.Error != nil {
			// panic("Error occurred while retrieving roles from the database: " + response.Error.Error())
			return utils.ResponseError(c, 500, "Error occurred while retrieving roles from the database", err.Error())
		}
		return c.JSON(roles)
	}
}

// Return a single role as JSON
func GetRole(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return utils.ResponseError(c, 400, "ID must be numbers only", err.Error())
		}

		role := new(models.Role)
		if response := db.Find(&role, id); response.Error != nil {
			return utils.ResponseError(c, 500, "An error occurred when retrieving the role", response.Error.Error())
		}

		if role.ID == 0 {
			// Send status not found
			return utils.ResponseError(c, 404, "ID "+c.Params("id")+" not found", "record not found")
		}

		return c.JSON(role)
	}
}

// Add a single role to the database
func AddRole(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		// Map Request Body to dto
		requestBody := new(dtos.RoleDto)
		if err := c.BodyParser(requestBody); err != nil {
			return utils.ResponseError(c, 500, fiber.ErrInternalServerError.Message, err.Error())
		}

		// validate request
		errors := dtos.ValidateStruct(*requestBody)
		if errors != nil {
			return utils.ResponseError(c, 400, "validation error", errors)
		}

		// Map dto to struct
		role := models.Role{
			Name:        requestBody.Name,
			Description: requestBody.Description,
		}

		// Create
		if response := db.Create(&role); response.Error != nil {
			return utils.ResponseError(c, 500, fiber.ErrInternalServerError.Message, err.Error())
		}

		return c.JSON(role)
	}
}

// Edit a single role
func EditRole(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return utils.ResponseError(c, 400, "ID must be numbers only", err.Error())
		}

		editRole := new(dtos.RoleDto)
		role := new(models.Role)

		if err := c.BodyParser(editRole); err != nil {
			return utils.ResponseError(c, 500, "An error occurred when parsing the edited role", err.Error())
		}

		if response := db.Find(&role, id); response.Error != nil {
			return utils.ResponseError(c, 400, "An error occurred when retrieving the existing role", response.Error.Error())
		}

		// Role does not exist
		if role.ID == 0 {
			return utils.ResponseError(c, 404, "ID "+c.Params("id")+" not found", "record not found")
		}

		updateData := models.Role{
			Name:        editRole.Name,
			Description: editRole.Description,
		}

		err = db.Model(&role).Updates(updateData).Error
		if err != nil {
			return utils.ResponseError(c, 500, "role edited error", err.Error())
		}

		return c.JSON(role)
	}
}
