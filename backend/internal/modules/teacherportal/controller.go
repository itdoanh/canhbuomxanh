package teacherportal

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	db *sql.DB
}

type claimProfileRequest struct {
	TeacherID uint64 `json:"teacherId"`
}

type optOutRequest struct {
	TeacherID uint64 `json:"teacherId"`
}

type createAppealRequest struct {
	TeacherID uint64 `json:"teacherId"`
	TargetType string `json:"targetType"`
	TargetID uint64 `json:"targetId"`
	Reason string `json:"reason"`
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (ctl *Controller) ClaimProfile(c *fiber.Ctx) error {
	userID, role, err := authContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if role != "teacher" && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req claimProfileRequest
	if err := c.BodyParser(&req); err != nil || req.TeacherID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	_, err = ctl.db.ExecContext(c.Context(),
		"UPDATE teachers SET user_id = ?, claim_status = 'claimed', is_public = 1 WHERE id = ?",
		userID,
		req.TeacherID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "claim profile failed"})
	}

	return c.JSON(fiber.Map{"message": "profile claimed and now public"})
}

func (ctl *Controller) RequestOptOut(c *fiber.Ctx) error {
	_, role, err := authContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if role != "teacher" && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req optOutRequest
	if err := c.BodyParser(&req); err != nil || req.TeacherID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	_, err = ctl.db.ExecContext(c.Context(),
		"UPDATE teachers SET is_public = 0, opt_out_requested_at = NOW() WHERE id = ?",
		req.TeacherID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "opt-out request failed"})
	}

	return c.JSON(fiber.Map{"message": "opt-out requested; deletion workflow target is 24h"})
}

func (ctl *Controller) CreateAppeal(c *fiber.Ctx) error {
	userID, role, err := authContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if role != "teacher" && role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req createAppealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.TeacherID == 0 || req.TargetType == "" || req.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields"})
	}

	_, err = ctl.db.ExecContext(c.Context(),
		"INSERT INTO appeals (teacher_id, requester_user_id, target_type, target_id, reason) VALUES (?, ?, ?, ?, ?)",
		req.TeacherID,
		userID,
		req.TargetType,
		req.TargetID,
		req.Reason,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "create appeal failed"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "appeal submitted"})
}

func (ctl *Controller) ListMyAppeals(c *fiber.Ctx) error {
	userID, _, err := authContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT id, teacher_id, target_type, target_id, reason, status, created_at FROM appeals WHERE requester_user_id = ? ORDER BY id DESC LIMIT 100",
		userID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var id, teacherID, targetID uint64
		var targetType, reason, status, createdAt string
		if err := rows.Scan(&id, &teacherID, &targetType, &targetID, &reason, &status, &createdAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{
			"id": id,
			"teacherId": teacherID,
			"targetType": targetType,
			"targetId": targetID,
			"reason": reason,
			"status": status,
			"createdAt": createdAt,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) Dashboard(c *fiber.Ctx) error {
	userID, _, err := authContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var totalAppeals int64
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM appeals WHERE requester_user_id = ?", userID).Scan(&totalAppeals)

	var pendingAppeals int64
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM appeals WHERE requester_user_id = ? AND status = 'open'", userID).Scan(&pendingAppeals)

	return c.JSON(fiber.Map{
		"totalAppeals": totalAppeals,
		"pendingAppeals": pendingAppeals,
	})
}

func authContext(c *fiber.Ctx) (uint64, string, error) {
	userIDValue := c.Locals("userId")
	roleValue := c.Locals("role")

	var userID uint64
	switch typed := userIDValue.(type) {
	case uint64:
		userID = typed
	case int:
		if typed >= 0 {
			userID = uint64(typed)
		}
	default:
		return 0, "", fmt.Errorf("invalid user id")
	}

	role, ok := roleValue.(string)
	if !ok || role == "" {
		return 0, "", fmt.Errorf("invalid role")
	}

	return userID, role, nil
}
