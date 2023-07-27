package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/prplx/wordy"
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

func Run(ctx context.Context) {
	var tun ngrok.Tunnel

	if !helpers.IsProduction() {
		err := godotenv.Load()
		if err != nil {
			logger.Error(err)
			return
		}

		tun, err = ngrok.Listen(ctx,
			config.HTTPEndpoint(),
			ngrok.WithAuthtokenFromEnv(),
		)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Error(err)
		return
	}

	autoMigrate(db)
	if err := seed(db); err != nil {
		logger.Error(err)
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("i18n/active.en.toml")
	bundle.MustLoadMessageFile("i18n/active.ru.toml")
	bundle.MustLoadMessageFile("i18n/active.nl.toml")

	repositories := repositories.NewRepositories(db)
	services := services.NewServices(services.Deps{
		Repositories:    *repositories,
		LocalizerBundle: bundle,
	})
	handlers := handlers.NewHandlers(services)

	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})
	app.Use(recover.New())
	handlers.Init(app)

	server := new(wordy.Server)

	go func() {
		var err error
		port := helpers.Getenv("PORT", "3000")
		if helpers.IsProduction() {
			err = server.Run(port, adaptor.FiberApp(app))
		} else {
			helpers.SetWebhookUrl(tun.URL())
			err = server.Run(port, adaptor.FiberApp(app), tun)
		}
		if err != nil {
			logger.Fatalf("An error occured while starting server: %s", err.Error())
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func autoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.Expression{}, &models.Example{}, &models.Language{}, &models.Audio{}, &models.Translation{}, &models.Synonym{}); err != nil {
		log.Fatal(err)
	}
}

func seed(db *gorm.DB) error {
	languages := []models.Language{}

	for _, language := range helpers.GetLanguageMap() {
		languages = append(languages, models.Language{
			Code:        language.Code,
			Text:        language.Title,
			EnglishText: language.EnglishTitle,
			Emoji:       language.Emoji,
		})
	}

	result := db.Create(&languages)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
