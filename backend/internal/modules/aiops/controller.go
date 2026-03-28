package aiops

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(db *sql.DB) *Controller {
	return &Controller{service: NewService(db)}
}

func (ctl *Controller) ReloadLexicon(c *fiber.Ctx) error {
	if err := ensureAIOpsAccess(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	if err := ctl.service.ReloadLexicon(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "lexicon reloaded", "summary": ctl.service.Summary()})
}

func (ctl *Controller) LexiconSummary(c *fiber.Ctx) error {
	if err := ensureAIOpsAccess(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	return c.JSON(ctl.service.Summary())
}

func (ctl *Controller) AnalyzeText(c *fiber.Ctx) error {
	if err := ensureAIOpsAccess(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req analyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "text is required"})
	}

	risk, score, hits := ctl.service.Analyze(req)
	return c.JSON(fiber.Map{"riskLevel": risk, "score": score, "hits": hits})
}

func (ctl *Controller) AnalyzeAndFlag(c *fiber.Ctx) error {
	if err := ensureAIOpsAccess(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	var req analyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Text == "" || req.SourceType == "" || req.SourceID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "sourceType, sourceId, text are required"})
	}
	if req.SourceType != "post" && req.SourceType != "comment" && req.SourceType != "teacher_profile" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid sourceType"})
	}

	result, err := ctl.service.AnalyzeAndFlag(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (ctl *Controller) EnrichedQueue(c *fiber.Ctx) error {
	if err := ensureAIOpsAccess(c); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	items, err := ctl.service.EnrichedQueue(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"items": items})
}

func ensureAIOpsAccess(c *fiber.Ctx) error {
	roleValue := c.Locals("role")
	role, ok := roleValue.(string)
	if !ok {
		return fmt.Errorf("invalid role")
	}
	if role != "admin" && role != "moderator" {
		return fmt.Errorf("forbidden")
	}
	return nil
}
