package handler

import (
	"errors"
	"log/slog"

	"github.com/alkhanov95/api-gateway/internal/apperr"
	"github.com/alkhanov95/api-gateway/internal/models"
	"github.com/alkhanov95/api-gateway/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handle struct {
	uc usecase.UserProvider
	//поменять usecase на usecase.UserProvider
}

func New(uc usecase.UserProvider) *Handle {
	return &Handle{uc: uc}
}
func (h *Handle) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// генерим UUID для нового юзера
	user.ID = uuid.New().String()

	// дергаем usecase → repo → db
	id, err := h.uc.CreateUser(c.UserContext(), &user)
	if err != nil {
		slog.Error("Failed to create user", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	// возвращаем 201 + id
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (h *Handle) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	slog.Info("Get request by ID:", "id", id)

	user, err := h.uc.GetUserByID(c.UserContext(), id)
	if err != nil {
		//проверить чтобы в postman возвращалось 404
		if errors.Is(err, apperr.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		slog.Error("Failed to get user by ID", slog.Any("error", err), "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

func (h *Handle) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.uc.List(c.UserContext())
	if err != nil {
		slog.Error(("Error while selecting user: "), slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get the users"})

	}
	return c.JSON(users)
}

func (h *Handle) DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.uc.Delete(c.UserContext(), id)
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

	err := h.uc.Update(c.UserContext(), user)
	if err != nil {
		slog.Error("Failed to execute update query", slog.Any("error", err), "id", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
 

