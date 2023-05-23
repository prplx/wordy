package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prplx/wordy/internal/services"
)

type Handlers struct {
	services *services.Services
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		services: services,
	}
}

func (h *Handlers) Init(router fiber.Router) {
	v1 := router.Group("/v1")
	{
		h.initUsersHandlers(v1)
		h.initBotHandlers(v1)
	}
}
