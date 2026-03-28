package moderator

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/moderator", middleware.JWTGuard(cfg.JWTSecret))

	group.Get("/queue", controller.AIQueue)
	group.Get("/appeals", controller.AppealQueue)
	group.Post("/appeals/:id/review", controller.ReviewAppeal)
	group.Post("/violations/user/:id", controller.EnforceUserViolation)
}
