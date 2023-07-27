package handlers

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/prplx/wordy/internal/handlers/v1"
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

func (h *Handlers) Init(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	h.initAPI(app)
}

func (h *Handlers) initAPI(app *fiber.App) {
	handlerV1 := v1.NewHandlers(h.services)

	api := app.Group("/api")

	{
		handlerV1.Init(api)
	}
}
