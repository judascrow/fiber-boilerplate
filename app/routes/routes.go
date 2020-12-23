package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/judascrow/fiber-boilerplate/app/controllers"
	"github.com/judascrow/fiber-boilerplate/database"
)

func RegisterAPI(api fiber.Router, db *database.Database) {
	registerHealthcheck(api, db)
	registerRoles(api, db)
	registerUsers(api, db)
}

func registerHealthcheck(api fiber.Router, db *database.Database) {
	healthcheck := api.Group("/healthcheck")

	healthcheck.Get("/", controllers.GetHealthcheck(db))
}

func registerRoles(api fiber.Router, db *database.Database) {
	roles := api.Group("/roles")

	roles.Get("/", controllers.GetAllRoles(db))
	roles.Get("/:id", controllers.GetRole(db))
	// roles.Post("/", Controller.AddRole(db))
	// roles.Put("/:id", Controller.EditRole(db))
	// roles.Delete("/:id", Controller.DeleteRole(db))
}

func registerUsers(api fiber.Router, db *database.Database) {
	users := api.Group("/users")

	users.Get("/", controllers.GetAllUsers(db))
}
