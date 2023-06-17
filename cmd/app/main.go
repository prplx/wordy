package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/prplx/wordy/internal/handlers"
	"github.com/prplx/wordy/internal/helpers"
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
	"github.com/prplx/wordy/internal/services"
	"github.com/prplx/wordy/pkg/logger"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := run(context.Background()); err != nil {
		logger.Fatal(err)
	}
}

func run(ctx context.Context) error {
	var tun ngrok.Tunnel

	if !helpers.IsProduction() {
		err := godotenv.Load()
		if err != nil {
			return err
		}

		tun, err = ngrok.Listen(ctx,
			config.HTTPEndpoint(),
			ngrok.WithAuthtokenFromEnv(),
		)
		if err != nil {
			return err
		}
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return err
	}

	autoMigrate(db)
	if err := seed(db); err != nil {
		logger.Error(err)
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("i18n/active.en.toml")
	bundle.MustLoadMessageFile("i18n/active.ru.toml")

	repositories := repositories.NewRepositories(db)
	services := services.NewServices(services.Deps{
		Repositories:    *repositories,
		LocalizerBundle: bundle,
	})
	handlers := handlers.NewHandlers(services)

	app := fiber.New()
	app.Use(recover.New())
	handlers.Init(app)

	if helpers.IsProduction() {
		port := fmt.Sprintf(":%s", helpers.Getenv("PORT", "3000"))
		return app.Listen(port)
	} else {
		helpers.SetWebhookUrl(tun.URL())
		return http.Serve(tun, adaptor.FiberApp(app))
	}
}

func autoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Expression{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Example{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Language{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Audio{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Translation{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Synonym{}); err != nil {
		log.Fatal(err)
	}
}

func seed(db *gorm.DB) error {
	languages := []models.Language{
		{
			Code:  "en",
			Text:  "English",
			Emoji: "🇬🇧",
		},
		{
			Code:  "ru",
			Text:  "Русский",
			Emoji: "🇷🇺",
		},
		{
			Code:  "nl",
			Text:  "Nederlands",
			Emoji: "🇳🇱",
		},
	}
	result := db.Create(&languages)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
