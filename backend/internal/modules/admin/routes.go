package admin

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/admin", middleware.JWTGuard(cfg.JWTSecret))

	group.Get("/system/overview", controller.SystemOverview)
	group.Get("/spam-vote/teachers", controller.SpamVoteSignals)
	group.Post("/aiops/config", controller.UpdateAIOpsConfig)
	group.Get("/aiops/config", controller.GetAIOpsConfig)
	group.Post("/access/role", controller.UpdateUserRole)
}
