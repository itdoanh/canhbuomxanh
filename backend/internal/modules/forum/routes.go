package forum

import (
	"database/sql"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, cfg *config.Config, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/forum")
	group.Get("/posts", controller.ListPosts)
	group.Get("/posts/:id/comments", controller.ListComments)
	group.Post("/posts", middleware.JWTGuard(cfg.JWTSecret), controller.CreatePost)
	group.Post("/posts/:id/comments", middleware.JWTGuard(cfg.JWTSecret), controller.CreateComment)
}
