package aiops

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/aiops", middleware.JWTGuard(cfg.JWTSecret))

	group.Post("/lexicon/reload", controller.ReloadLexicon)
	group.Get("/lexicon/summary", controller.LexiconSummary)
	group.Post("/analyze", controller.AnalyzeText)
	group.Post("/analyze-and-flag", controller.AnalyzeAndFlag)
	group.Get("/queue/enriched", controller.EnrichedQueue)
}
