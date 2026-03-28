package profile

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	db *sql.DB
}

type updateProfileRequest struct {
	FullName  string `json:"fullName"`
	AvatarURL string `json:"avatarUrl"`
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (ctl *Controller) Me(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	row := ctl.db.QueryRowContext(c.Context(),
		"SELECT id, email, full_name, role, avatar_url, status, created_at FROM users WHERE id = ? LIMIT 1",
		userID,
	)

	var id uint64
	var email, fullName, role, status string
	var avatarURL sql.NullString
	var createdAt string
	if err := row.Scan(&id, &email, &fullName, &role, &avatarURL, &status, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}

	return c.JSON(fiber.Map{
		"id":        id,
		"email":     email,
		"fullName":  fullName,
		"role":      role,
		"avatarUrl": avatarURL.String,
		"status":    status,
		"createdAt": createdAt,
	})
}

func (ctl *Controller) UpdateMe(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fullName is required"})
	}

	_, err = ctl.db.ExecContext(c.Context(),
		"UPDATE users SET full_name = ?, avatar_url = ? WHERE id = ?",
		req.FullName,
		req.AvatarURL,
		userID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update failed"})
	}

	return c.JSON(fiber.Map{"message": "profile updated"})
}

func getUserID(c *fiber.Ctx) (uint64, error) {
	value := c.Locals("userId")
	switch typed := value.(type) {
	case uint64:
		return typed, nil
	case int:
		if typed >= 0 {
			return uint64(typed), nil
		}
	}
	return 0, fmt.Errorf("invalid user id")
}
