package teacherportal

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/teacher-portal", middleware.JWTGuard(cfg.JWTSecret))

	group.Post("/claim", controller.ClaimProfile)
	group.Post("/opt-out", controller.RequestOptOut)
	group.Post("/appeals", controller.CreateAppeal)
	group.Get("/appeals", controller.ListMyAppeals)
	group.Get("/dashboard", controller.Dashboard)
}
