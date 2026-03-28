package voting

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/voting", middleware.JWTGuard(cfg.JWTSecret))

	group.Post("/votes", controller.CastVote)
	group.Post("/votes/release", controller.ReleaseSemesterVotes)
	group.Post("/alerts/forced-vote", controller.ReportForcedVote)
	group.Post("/badges/recompute", controller.RecomputeBadges)
	group.Get("/badges/teachers/:id", controller.GetTeacherBadgeSummary)
}
