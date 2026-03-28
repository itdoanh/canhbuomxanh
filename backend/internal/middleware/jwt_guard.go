package middleware

import (
	"errors"
	"strings"

	"canhbuomxanh/backend/internal/platform/security"

	"github.com/gofiber/fiber/v2"
)

func JWTGuard(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		tokenString, err := parseBearerToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		claims, err := security.ParseJWT(tokenString, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("userId", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func parseBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("invalid authorization format")
	}
	return parts[1], nil
}
