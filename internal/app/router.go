package app

import (
	"github.com/alkhanov95/api-gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, h *handler.Handle) {
	app.Post("/users", h.CreateUser)
	app.Get("/users/:id", h.GetUserByID)
	app.Get("/users", h.GetAllUsers)
	app.Delete("/users/:id", h.DeleteUserByID)
}
