package teacher

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	db *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (ctl *Controller) ListPublic(c *fiber.Ctx) error {
	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT id, display_name, school_name, subject_name, claim_status FROM teachers WHERE is_public = 1 ORDER BY id DESC LIMIT 100",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var id uint64
		var displayName string
		var schoolName, subjectName sql.NullString
		var claimStatus string
		if err := rows.Scan(&id, &displayName, &schoolName, &subjectName, &claimStatus); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{
			"id":          id,
			"displayName": displayName,
			"schoolName":  schoolName.String,
			"subjectName": subjectName.String,
			"claimStatus": claimStatus,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) GetPublic(c *fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid teacher id"})
	}

	row := ctl.db.QueryRowContext(c.Context(),
		"SELECT id, display_name, school_name, subject_name, claim_status, is_public FROM teachers WHERE id = ? LIMIT 1",
		id,
	)

	var teacherID uint64
	var displayName string
	var schoolName, subjectName sql.NullString
	var claimStatus string
	var isPublic bool
	if err := row.Scan(&teacherID, &displayName, &schoolName, &subjectName, &claimStatus, &isPublic); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "teacher not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}

	if !isPublic {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "teacher profile is not public"})
	}

	return c.JSON(fiber.Map{
		"id":          teacherID,
		"displayName": displayName,
		"schoolName":  schoolName.String,
		"subjectName": subjectName.String,
		"claimStatus": claimStatus,
	})
}
