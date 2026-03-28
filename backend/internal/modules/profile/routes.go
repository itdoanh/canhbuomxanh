package profile

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/profile", middleware.JWTGuard(cfg.JWTSecret))
	group.Get("/me", controller.Me)
	group.Put("/me", controller.UpdateMe)
}
