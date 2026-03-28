package router

import (
	"database/sql"
	"time"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/modules/admin"
	"canhbuomxanh/backend/internal/modules/aiops"
	"canhbuomxanh/backend/internal/modules/auth"
	"canhbuomxanh/backend/internal/modules/forum"
	"canhbuomxanh/backend/internal/modules/moderator"
	"canhbuomxanh/backend/internal/modules/profile"
	"canhbuomxanh/backend/internal/modules/teacher"
	"canhbuomxanh/backend/internal/modules/teacherportal"
	"canhbuomxanh/backend/internal/modules/voting"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, cfg *config.Config, db *sql.DB) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "up"
		if err := db.PingContext(c.Context()); err != nil {
			dbStatus = "down"
		}
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": cfg.AppName,
			"db": dbStatus,
			"time": time.Now().Format(time.RFC3339),
		})
	})

	v1.Get("/metrics", func(c *fiber.Ctx) error {
		var users int64
		var posts int64
		var votes int64
		var queuedFlags int64

		_ = db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM users").Scan(&users)
		_ = db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM posts").Scan(&posts)
		_ = db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM votes").Scan(&votes)
		_ = db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM flags WHERE status = 'queued'").Scan(&queuedFlags)

		return c.JSON(fiber.Map{
			"users": users,
			"posts": posts,
			"votes": votes,
			"queuedFlags": queuedFlags,
			"service": cfg.AppName,
			"time": time.Now().Format(time.RFC3339),
		})
	})

	auth.RegisterRoutes(v1, cfg, db)
	forum.RegisterRoutes(v1, cfg, db)
	teacher.RegisterRoutes(v1, db)
	profile.RegisterRoutes(v1, cfg, db)
	voting.RegisterRoutes(v1, cfg, db)
	teacherportal.RegisterRoutes(v1, cfg, db)
	moderator.RegisterRoutes(v1, cfg, db)
	admin.RegisterRoutes(v1, cfg, db)
	aiops.RegisterRoutes(v1, cfg, db)
}
