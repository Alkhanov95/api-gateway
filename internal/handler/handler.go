package handler

import (
	"log/slog"

	"github.com/alkhanov95/api-gateway/internal/repository"
	"github.com/alkhanov95/api-gateway/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handle struct {
	repo *repository.UserRepo
}

func New(repo *repository.UserRepo) *Handle {
	return &Handle{repo: repo}
}

func (h *Handle) CreateUser(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	user.ID = uuid.New().String()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": user.ID})
}

func (h *Handle) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	slog.Info("Get request by ID:", "id", id)

	user, err := h.repo.GetByID(c.UserContext(), id)
	if err != nil {
		slog.Error("Failed to get user by ID", slog.Any("error", err), "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

func (h *Handle) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.repo.List(c.UserContext())
	if err != nil {
		slog.Error(("Error while selecting user: "), slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get the users"})

	}
	return c.JSON(users)
}

func (h *Handle) DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.repo.Delete(c.UserContext(), id)
	if err != nil {
		slog.Error("Failed to execute delete query", slog.Any("error", err), "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handle) UpdateUser(c *fiber.Ctx) error {
	var user *models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	user.ID = c.Params("id")

	err := h.repo.Update(c.UserContext(), user)
	if err != nil {
		slog.Error("Failed to execute update query", slog.Any("error", err), "id", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
