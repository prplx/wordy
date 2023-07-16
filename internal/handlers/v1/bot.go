package v1

import (
	"errors"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/prplx/wordy/internal/helpers"
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/types"
	"github.com/prplx/wordy/pkg/logger"
)

func (h *Handlers) initBotHandlers(api fiber.Router) {
	bot := api.Group("bot")

	{
		bot.Post("/", h.handleBot)
	}
}

func (h *Handlers) handleBot(ctx *fiber.Ctx) error {
	reqIP := net.ParseIP(ctx.IP())
	_, ipnetA, _ := net.ParseCIDR(os.Getenv("TG_SUBNET_A"))
	_, ipnetB, _ := net.ParseCIDR(os.Getenv("TG_SUBNET_B"))
	allowedUsers := []string{"mmystiq", "nastyaknyazhe"}
	var update types.Update
	var lang string

	if !ipnetA.Contains(reqIP) && !ipnetB.Contains(reqIP) {
		logger.Error("Unauthorized request from IP: " + ctx.IP())
		return ctx.SendStatus(http.StatusOK)
	}

	if err := ctx.BodyParser(&update); err != nil {
		logger.Error(err)
		return err
	}

	if update.Message.From.LanguageCode != "" {
		lang = update.Message.From.LanguageCode
	} else if update.CallbackQuery.From.LanguageCode != "" {
		lang = update.CallbackQuery.From.LanguageCode
	} else {
		lang = "en"
	}
	h.services.Localizer.ChangeLanguage(lang)

	if !helpers.StringInSlice(update.Message.From.Username, allowedUsers) && !helpers.StringInSlice(update.CallbackQuery.From.Username, allowedUsers) {
		h.services.Telegram.SendText(update.Message.Chat.Id, h.services.Localizer.L("BotUnderDevelopment"))
		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.MessageId == 0 && update.CallbackQuery.Id == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	languages, err := h.services.Languages.Query()
	if err != nil {
		h.services.Telegram.SendText(update.Message.Chat.Id, h.services.Localizer.L("SomethingWentWrong"))
		return ctx.SendStatus(http.StatusOK)
	}

	user := models.User{
		TelegramId:       update.Message.From.Id,
		TelegramUsername: update.Message.From.Username,
		FirstName:        update.Message.From.FirstName,
		LastName:         update.Message.From.LastName,
	}

	fromId := update.Message.From.Id
	if fromId == 0 {
		fromId = update.CallbackQuery.From.Id
	}
	if fromId == 0 {
		return ctx.SendStatus(http.StatusOK)
	}

	dbUser, err := h.services.Users.GetByTgId(uint(fromId))
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			if update.Message.From.LanguageCode != "" {
				language, err := h.services.Languages.GetByCode(update.Message.From.LanguageCode)
				if err == nil {
					user.FirstLanguage = int(language.ID)
				}
			}

			if _, err := h.services.Users.Create(&user); err != nil {
				logger.Error(err)
				return ctx.SendStatus(http.StatusOK)
			}

			dbUser = user
		} else {
			logger.Error(err)
			return ctx.SendStatus(http.StatusOK)
		}

	}

	if update.Message.Text == "/start" {
		if _, err := h.handleStartCommand(update.Message.Chat.Id); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.Text == "/settings" {
		if _, err := h.handleSettingsCommand(update.Message.Chat.Id); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if update.CallbackQuery.Data == "setLanguagePair" {
		if err := h.handleSetLanguagePair(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageId, h.services.Localizer.L("ChooseFirstLanguage"), "setFirstLanguage", languages); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	setFirstLanguagePattern := `setFirstLanguage: (\w+)`

	re := regexp.MustCompile(setFirstLanguagePattern)
	match := re.FindStringSubmatch(update.CallbackQuery.Data)

	if len(match) > 1 {
		var firstLanguage models.Language
		for _, language := range languages {
			if language.Code == match[1] {
				firstLanguage = language
			}
		}
		dbUser.FirstLanguage = int(firstLanguage.ID)
		if err := h.handleUpdateUserSettings(update.CallbackQuery.Id, &dbUser); err != nil {
			logger.Error(err)
			h.services.Telegram.SendText(update.CallbackQuery.Message.Chat.Id, h.services.Localizer.L("SomethingWentWrong"))
			return ctx.SendStatus(http.StatusOK)
		}

		if err := h.handleSetLanguagePair(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageId, h.services.Localizer.L("ChooseSecondLanguage"), "setSecondLanguage", languages); err != nil {
			logger.Error(err)
		}

		return ctx.SendStatus(http.StatusOK)
	}

	setSecondLanguagePattern := `setSecondLanguage: (\w+)`
	re = regexp.MustCompile(setSecondLanguagePattern)
	match = re.FindStringSubmatch(update.CallbackQuery.Data)

	if len(match) > 1 {
		var secondLanguage models.Language
		for _, language := range languages {
			if language.Code == match[1] {
				secondLanguage = language
			}
		}
		// TODO: check if first language is not equal to second language
		dbUser.SecondLanguage = int(secondLanguage.ID)

		if err := h.handleUpdateUserSettings(update.CallbackQuery.Id, &dbUser); err != nil {
			h.services.Telegram.SendText(update.CallbackQuery.Message.Chat.Id, h.services.Localizer.L("SomethingWentWrong"))
			logger.Error(err)
			return ctx.SendStatus(http.StatusOK)
		}

		if err := h.services.Telegram.DeleteMessage(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageId); err != nil {
			logger.Error(err)
		}

		if _, err := h.services.Telegram.SendText(update.CallbackQuery.Message.Chat.Id, h.services.Localizer.L("SettingsUpdated")); err != nil {
			logger.Error(err)
		}

		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.Text == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	firstUserLanguage, secondUserLanguage := helpers.GetUserFirstAndSecondLanguagesIds(dbUser, languages)
	h.services.Localizer.ChangeLanguage(firstUserLanguage.Code)

	if firstUserLanguage.ID == 0 || secondUserLanguage.ID == 0 {
		if _, err := h.services.Telegram.SendText(update.Message.Chat.Id, h.services.Localizer.L("SetLanguagesWarning")); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	err = h.services.Telegram.SendTypingChatAction(update.Message.Chat.Id)
	if err != nil {
		logger.Error(err)
		return err
	}

	langs := helpers.GetLanguageMap()
	to := langs[firstUserLanguage.Code]
	from := langs[secondUserLanguage.Code]

	if err := h.handleTextTranslation(update.Message.Chat.Id, update.Message.MessageId, int(dbUser.ID), int(secondUserLanguage.ID), strings.TrimSpace(update.Message.Text), from.Text, to.Text, strconv.Itoa(update.Message.From.Id)); err != nil {
		logger.Error(err)
	}

	return ctx.SendStatus(http.StatusOK)
}
