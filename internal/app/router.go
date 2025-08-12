package app

import (
	"github.com/alkhanov95/api-gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func setupUserRoutes(h *handler.Handle) *fiber.App {
	app := fiber.New() //heap stack read!

	app.Post("/users", h.CreateUser)
	app.Get("/users/:id", h.GetUserByID)
	app.Get("/users", h.GetAllUsers)
	app.Delete("/users/:id", h.DeleteUserByID)
	app.Put("/users/:id", h.UpdateUser)

	return app
}
