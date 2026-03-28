package moderator

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	db *sql.DB
}

type reviewAppealRequest struct {
	Action string `json:"action"`
}

type enforceUserRequest struct {
	Action string `json:"action"`
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (ctl *Controller) AIQueue(c *fiber.Ctx) error {
	if _, err := ensureModeratorRole(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT id, source_type, source_id, risk_level, reason_code, reason_detail, status, created_at FROM flags WHERE status = 'queued' ORDER BY risk_level DESC, id ASC LIMIT 200",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var id, sourceID uint64
		var sourceType, riskLevel, reasonCode, status, createdAt string
		var reasonDetail sql.NullString
		if err := rows.Scan(&id, &sourceType, &sourceID, &riskLevel, &reasonCode, &reasonDetail, &status, &createdAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{
			"id": id,
			"sourceType": sourceType,
			"sourceId": sourceID,
			"riskLevel": riskLevel,
			"reasonCode": reasonCode,
			"reasonDetail": reasonDetail.String,
			"status": status,
			"createdAt": createdAt,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) AppealQueue(c *fiber.Ctx) error {
	if _, err := ensureModeratorRole(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT id, teacher_id, requester_user_id, target_type, target_id, reason, status, created_at FROM appeals WHERE status = 'open' ORDER BY id ASC LIMIT 200",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var id, teacherID, requesterID, targetID uint64
		var targetType, reason, status, createdAt string
		if err := rows.Scan(&id, &teacherID, &requesterID, &targetType, &targetID, &reason, &status, &createdAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{
			"id": id,
			"teacherId": teacherID,
			"requesterUserId": requesterID,
			"targetType": targetType,
			"targetId": targetID,
			"reason": reason,
			"status": status,
			"createdAt": createdAt,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) ReviewAppeal(c *fiber.Ctx) error {
	moderatorID, err := ensureModeratorRole(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	appealID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid appeal id"})
	}

	var req reviewAppealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	status := "rejected"
	if req.Action == "accept" {
		status = "accepted"
	}

	_, err = ctl.db.ExecContext(c.Context(),
		"UPDATE appeals SET status = ?, reviewed_by = ? WHERE id = ?",
		status,
		moderatorID,
		appealID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update failed"})
	}

	return c.JSON(fiber.Map{"message": "appeal reviewed", "status": status})
}

func (ctl *Controller) EnforceUserViolation(c *fiber.Ctx) error {
	_, err := ensureModeratorRole(c)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	userID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req enforceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	status := "suspended"
	if req.Action == "warn" {
		status = "active"
	}
	if req.Action == "ban" {
		status = "deleted"
	}

	_, err = ctl.db.ExecContext(c.Context(),
		"UPDATE users SET status = ? WHERE id = ?",
		status,
		userID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "enforcement failed"})
	}

	return c.JSON(fiber.Map{"message": "violation action applied", "userStatus": status})
}

func ensureModeratorRole(c *fiber.Ctx) (uint64, error) {
	userIDValue := c.Locals("userId")
	roleValue := c.Locals("role")
	role, ok := roleValue.(string)
	if !ok || (role != "moderator" && role != "admin") {
		return 0, fmt.Errorf("forbidden")
	}

	switch typed := userIDValue.(type) {
	case uint64:
		return typed, nil
	case int:
		if typed >= 0 {
			return uint64(typed), nil
		}
	}
	return 0, fmt.Errorf("invalid user id")
}
