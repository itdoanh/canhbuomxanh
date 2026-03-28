package admin

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	db *sql.DB
}

type updateAIOpsConfigRequest struct {
	BannedKeywords []string `json:"bannedKeywords"`
	SuggestWeight  float64  `json:"suggestWeight"`
}

type updateRoleRequest struct {
	UserID uint64 `json:"userId"`
	Role   string `json:"role"`
}

var aiopsConfig = fiber.Map{
	"bannedKeywords": []string{"boi_nho", "ha_nhoc", "so_sanh_doc_hai"},
	"suggestWeight":  1.0,
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (ctl *Controller) SystemOverview(c *fiber.Ctx) error {
	if err := ensureAdmin(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var users, teachers, posts, votes, flags int64
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM users").Scan(&users)
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM teachers").Scan(&teachers)
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM posts").Scan(&posts)
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM votes").Scan(&votes)
	_ = ctl.db.QueryRowContext(c.Context(), "SELECT COUNT(*) FROM flags WHERE status = 'queued'").Scan(&flags)

	return c.JSON(fiber.Map{
		"users": users,
		"teachers": teachers,
		"posts": posts,
		"votes": votes,
		"queuedFlags": flags,
	})
}

func (ctl *Controller) SpamVoteSignals(c *fiber.Ctx) error {
	if err := ensureAdmin(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT teacher_id, COUNT(*) AS vote_count FROM votes GROUP BY teacher_id ORDER BY vote_count DESC LIMIT 20",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var teacherID uint64
		var voteCount int64
		if err := rows.Scan(&teacherID, &voteCount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{"teacherId": teacherID, "voteCount": voteCount})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) UpdateAIOpsConfig(c *fiber.Ctx) error {
	if err := ensureAdmin(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req updateAIOpsConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.BannedKeywords != nil {
		aiopsConfig["bannedKeywords"] = req.BannedKeywords
	}
	if req.SuggestWeight > 0 {
		aiopsConfig["suggestWeight"] = req.SuggestWeight
	}

	return c.JSON(fiber.Map{"message": "aiops config updated", "config": aiopsConfig})
}

func (ctl *Controller) GetAIOpsConfig(c *fiber.Ctx) error {
	if err := ensureAdmin(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	return c.JSON(aiopsConfig)
}

func (ctl *Controller) UpdateUserRole(c *fiber.Ctx) error {
	if err := ensureAdmin(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req updateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.UserID == 0 || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "userId and role are required"})
	}

	_, err := ctl.db.ExecContext(c.Context(), "UPDATE users SET role = ? WHERE id = ?", req.Role, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update role failed"})
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}

func ensureAdmin(c *fiber.Ctx) error {
	roleValue := c.Locals("role")
	role, ok := roleValue.(string)
	if !ok || role != "admin" {
		return fmt.Errorf("forbidden")
	}
	return nil
}
