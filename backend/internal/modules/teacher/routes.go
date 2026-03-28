package teacher

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(v1 fiber.Router, db *sql.DB) {
	controller := NewController(db)
	group := v1.Group("/teachers")
	group.Get("/", controller.ListPublic)
	group.Get("/:id", controller.GetPublic)
}
