package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/example/user-age-api/internal/handler"
)

// RegisterUserRoutes registers all user-related routes.
func RegisterUserRoutes(app *fiber.App, h *handler.UserHandler) {
	users := app.Group("/users")

	users.Post("/", h.CreateUser)
	users.Get("/", h.ListUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}
