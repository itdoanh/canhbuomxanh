package app

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/platform/database"
	"canhbuomxanh/backend/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
	db  *sql.DB
}

func New() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := database.NewMySQL(cfg.DSN())
	if err != nil {
		return nil, err
	}

	fiberApp := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	fiberApp.Use(requestid.New())
	fiberApp.Use(recover.New())
	fiberApp.Use(logger.New())
	fiberApp.Use(helmet.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	fiberApp.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: parseWindow(cfg.RateLimitWindow),
		KeyGenerator: func(c *fiber.Ctx) string {
			return strings.Join([]string{c.IP(), c.Path()}, "|")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limit exceeded"})
		},
	}))

	router.Register(fiberApp, cfg, db)

	return &Server{app: fiberApp, cfg: cfg, db: db}, nil
}

func (s *Server) Run() error {
	defer s.db.Close()
	return s.app.Listen(fmt.Sprintf(":%s", s.cfg.AppPort))
}

func parseWindow(seconds int) time.Duration {
	if seconds <= 0 {
		return 60 * time.Second
	}
	return time.Duration(seconds) * time.Second
}
