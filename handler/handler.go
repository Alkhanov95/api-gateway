package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func SetupUserRoutes(app *fiber.App, db *pgxpool.Pool) {
	app.Post("/users", func(c *fiber.Ctx) error {
		go test var user User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "невалидный JSON"})
		}

		user.ID = uuid.New().String()

		_, err := db.Exec(context.Background(),
			"INSERT INTO users (id, name, age) VALUES ($1, $2, $3)",
			user.ID, user.Name, user.Age,
		)
		if err != nil {
			log.Printf("❌ Ошибка при вставке в БД: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "вставка не удалась"})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		log.Println("📥 GET запрос на ID:", id)

		var user User

		err := db.QueryRow(
			context.Background(),
			"SELECT id, name, age FROM users WHERE id = $1",
			id,
		).Scan(&user.ID, &user.Name, &user.Age)

		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "не найден"})
			}
			log.Printf("❌ Ошибка при SELECT: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ошибка сервера"})
		}

		return c.JSON(user)
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		var users []User

		rows, err := db.Query(context.Background(), "SELECT id, name, age FROM users")
		if err != nil {
			log.Printf("ошибка при выборе пользователя : ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get users"})
		}
		defer rows.Close()

		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
				log.Printf("Ошибка при чтении строки", err)
				continue
			}
			users = append(users, u)
		}
		return c.JSON(users)
	})
}
