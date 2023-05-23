package v1

import (
	"github.com/gofiber/fiber/v2"
)

func errorResponse(c *fiber.Ctx, statusCode int, message string) error {
	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return status code with error message
	return c.Status(statusCode).SendString(message)
}
