package voting

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(db *sql.DB) *Controller {
	return &Controller{service: NewService(db)}
}

func (ctl *Controller) CastVote(c *fiber.Ctx) error {
	userID, role, err := readAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req castVoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := ctl.service.CastVote(c.Context(), userID, role, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (ctl *Controller) ReleaseSemesterVotes(c *fiber.Ctx) error {
	_, role, err := readAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if !isPrivilegedRole(role) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req releaseVotesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	affected, err := ctl.service.ReleaseSemesterVotes(c.Context(), req.SemesterKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"updatedRows": affected})
}

func (ctl *Controller) ReportForcedVote(c *fiber.Ctx) error {
	userID, _, err := readAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req forcedVoteAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := ctl.service.ReportForcedVote(c.Context(), userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "alert created"})
}

func (ctl *Controller) RecomputeBadges(c *fiber.Ctx) error {
	_, role, err := readAuthContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if !isPrivilegedRole(role) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req recomputeBadgesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := ctl.service.RecomputeBadges(c.Context(), req.TeacherID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func (ctl *Controller) GetTeacherBadgeSummary(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid teacher id"})
	}

	summary, err := ctl.service.GetTeacherBadgeSummary(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(summary)
}

func readAuthContext(c *fiber.Ctx) (uint64, string, error) {
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

func isPrivilegedRole(role string) bool {
	return role == "admin" || role == "moderator"
}
