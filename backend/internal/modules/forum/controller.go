package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	db *sql.DB
}

type createPostRequest struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	TeacherID uint64 `json:"teacherId"`
}

type createCommentRequest struct {
	Body     string `json:"body"`
	ParentID uint64 `json:"parentId"`
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (ctl *Controller) ListPosts(c *fiber.Ctx) error {
	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT p.id, p.title, p.body, p.teacher_id, p.created_at, u.full_name FROM posts p JOIN users u ON u.id = p.user_id WHERE p.status = 'active' ORDER BY p.id DESC LIMIT 100",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var id uint64
		var title, body, createdAt, authorName string
		var teacherID sql.NullInt64
		if err := rows.Scan(&id, &title, &body, &teacherID, &createdAt, &authorName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{
			"id":        id,
			"title":     title,
			"body":      body,
			"teacherId": teacherID.Int64,
			"createdAt": createdAt,
			"authorName": authorName,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) CreatePost(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req createPostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Title == "" || req.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title and body are required"})
	}

	result, err := ctl.db.ExecContext(c.Context(),
		"INSERT INTO posts (user_id, teacher_id, title, body) VALUES (?, ?, ?, ?)",
		userID,
		req.TeacherID,
		req.Title,
		req.Body,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "insert failed"})
	}

	id, _ := result.LastInsertId()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (ctl *Controller) ListComments(c *fiber.Ctx) error {
	postID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	rows, err := ctl.db.QueryContext(c.Context(),
		"SELECT c.id, c.parent_id, c.body, c.created_at, u.full_name FROM comments c JOIN users u ON u.id = c.user_id WHERE c.post_id = ? AND c.status = 'active' ORDER BY c.id ASC LIMIT 300",
		postID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
	}
	defer rows.Close()

	items := make([]fiber.Map, 0)
	for rows.Next() {
		var id uint64
		var parentID sql.NullInt64
		var body, createdAt, authorName string
		if err := rows.Scan(&id, &parentID, &body, &createdAt, &authorName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "scan failed"})
		}
		items = append(items, fiber.Map{
			"id": id,
			"parentId": parentID.Int64,
			"body": body,
			"createdAt": createdAt,
			"authorName": authorName,
		})
	}

	return c.JSON(fiber.Map{"items": items})
}

func (ctl *Controller) CreateComment(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	postID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	var req createCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "body is required"})
	}

	result, err := ctl.db.ExecContext(c.Context(),
		"INSERT INTO comments (post_id, user_id, parent_id, body) VALUES (?, ?, ?, ?)",
		postID,
		userID,
		req.ParentID,
		req.Body,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "insert failed"})
	}

	id, _ := result.LastInsertId()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
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

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
