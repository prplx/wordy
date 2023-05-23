package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) initUsersHandlers(api fiber.Router) {
	users := api.Group("users")
	{
		users.Get("/me", h.getCurrentUser)
	}
}

func (h *Handlers) getCurrentUser(ctx *fiber.Ctx) error {
	user := ctx.Get("currentUser")

	if user == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "Cannot get current user")
	}

	return ctx.JSON(user)
}
