package handler

import (
	"context"
	"log/slog"

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

type Handle struct {
	DB *pgxpool.Pool
}

func (h *Handle) CreateUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	user.ID = uuid.New().String()

	_, err := h.DB.Exec(context.Background(),
		"INSERT INTO users (id, name, age) VALUES ($1, $2, $3)",
		user.ID, user.Name, user.Age,
	)
	if err != nil {
		slog.Error(("Failed to insert into database: "), slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "insert DB fail"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handle) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	slog.Info("Get request by ID:", "id", id) //slog.info?

	var user User

	err := h.DB.QueryRow(context.Background(),
		"SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}
		slog.Error(("Failed to execute SELECT query"), slog.Any("error", err), "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(user)
}

func (h *Handle) GetAllUsers(c *fiber.Ctx) error {
	var users []User

	rows, err := h.DB.Query(context.Background(), "SELECT id, name, age FROM users")
	if err != nil {
		slog.Error(("Error while selecting user: "), slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get the users"})
	}
	defer rows.Close()

	for rows.Next() {
		var u User

		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			slog.Error("Error while reading string", err)
			continue
		}
		users = append(users, u)

	}
	return c.JSON(users)
}

func (h *Handle) DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	slog.Info("Delete request by id", "id", id)

	del, err := h.DB.Exec(context.Background(),
		"DELETE FROM users WHERE id = $1", id,
	)

	if err != nil {
		slog.Error("Failed to execute delete query", slog.Any("error", err), "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error"})
	}

	if del.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})

	}

	return c.SendStatus(fiber.StatusNoContent)
}
