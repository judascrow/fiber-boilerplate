package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/judascrow/fiber-boilerplate/app/models"
	"github.com/judascrow/fiber-boilerplate/database"
)

// Return all users as JSON
func GetAllUsers(db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var Users []models.User
		if response := db.Find(&Users); response.Error != nil {
			panic("Error occurred while retrieving users from the database: " + response.Error.Error())
		}
		// Match roles to users
		for index, User := range Users {
			if User.RoleID != 0 {
				Role := new(models.Role)
				if response := db.Find(&Role, User.RoleID); response.Error != nil {
					panic("An error occurred when retrieving the role: " + response.Error.Error())
				}
				if Role.ID != 0 {
					Users[index].Role = *Role
				}
			}
		}
		err := ctx.JSON(Users)
		if err != nil {
			panic("Error occurred when returning JSON of users: " + err.Error())
		}
		return err
	}
}
