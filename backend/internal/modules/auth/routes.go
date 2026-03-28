package auth

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(cfg, db)

	group := v1.Group("/auth")
	group.Post("/register", controller.Register)
	group.Post("/login", controller.Login)

	group.Get("/me", middleware.JWTGuard(cfg.JWTSecret), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"userId": c.Locals("userId"),
			"role":   c.Locals("role"),
		})
	})
}
