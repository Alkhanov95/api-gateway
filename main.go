
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = make(map[string]User)

func main() {
	app := fiber.New()

	//post
	app.Post("/user", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err !=  {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		user.ID = uuid.New().String()
		users[user.ID] = user
		return c.Status(201).JSON(user)
	})

	//GET
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, exists := users[id]

		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})

	//PUT
	app.Put("/user", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		if user.ID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "ID is required"})
		}
		_, exists := users[user.ID]
		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		users[user.ID] = user
		return c.JSON(user)
	})

	app.Delete("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, exists := users[id]
		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		delete(users, id)
		return c.JSON(fiber.Map{"message": "user deleted"})
	})

	app.Listen(":3000")
}

