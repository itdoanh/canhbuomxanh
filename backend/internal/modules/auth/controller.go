package auth

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(cfg *config.Config, db *sql.DB) *Controller {
	return &Controller{
		service: NewService(cfg, db),
	}
}

func (ctl *Controller) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := ctl.service.Register(c.Context(), req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "register success"})
}

func (ctl *Controller) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	res, err := ctl.service.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
