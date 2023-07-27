package v1

import (
	"errors"
	"fmt"
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
		h.services.Telegram.SendText(update.Message.Chat.ID, h.services.Localizer.L("BotUnderDevelopment"))
		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.MessageID == 0 && update.CallbackQuery.ID == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	languages, err := h.services.Languages.Query()
	if err != nil {
		h.services.Telegram.SendText(update.Message.Chat.ID, h.services.Localizer.L("SomethingWentWrong"))
		return ctx.SendStatus(http.StatusOK)
	}

	user := models.User{
		TelegramID:       update.Message.From.ID,
		TelegramUsername: update.Message.From.Username,
		FirstName:        update.Message.From.FirstName,
		LastName:         update.Message.From.LastName,
	}

	fromId := update.Message.From.ID
	if fromId == 0 {
		fromId = update.CallbackQuery.From.ID
	}
	if fromId == 0 {
		return ctx.SendStatus(http.StatusOK)
	}

	dbUser, err := h.services.Users.GetByTgId(uint(fromId))
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			if update.Message.From.LanguageCode != "" {
				var lng models.Language
				for _, l := range languages {
					if l.Code == update.Message.From.LanguageCode {
						lng = l
						break
					}
				}
				if lng.Code == "" {
					for _, l := range languages {
						if l.Code == "en" {
							lng = l
							break
						}
					}
				}
				user.FirstLanguage = lng.ID
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
		if _, err := h.handleStartCommand(update.Message.Chat.ID); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.Text == "/settings" {
		if _, err := h.handleSettingsCommand(update.Message.Chat.ID); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if update.CallbackQuery.Data == "settings" {
		if _, err := h.handleSettingsCommand(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	command := update.CallbackQuery.Data
	setLanguagePattern := `setLanguage \((\d+)\)$`
	re := regexp.MustCompile(setLanguagePattern)
	match := re.FindStringSubmatch(command)
	if len(match) == 2 {
		var lang string
		if match[1] == "1" {
			lang = "First"
		} else if match[1] == "2" {
			lang = "Second"
		}
		if err := h.handleSetLanguagePair(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, h.services.Localizer.L(fmt.Sprintf("Choose%sLanguage", lang)), command, "settings", languages); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	setFirstLanguagePattern := `setLanguage \((\d+)\) \((\w+)\)`
	re = regexp.MustCompile(setFirstLanguagePattern)
	match = re.FindStringSubmatch(command)
	firstUserLanguage, secondUserLanguage := helpers.GetUserFirstAndSecondLanguagesIds(dbUser, languages)
	var toCompareWith uint

	if len(match) == 3 {
		var lang uint
		for _, l := range languages {
			if l.Code == match[2] {
				lang = l.ID
				break
			}
		}
		var isSettingFirstLanguage = match[1] == "1"

		if isSettingFirstLanguage {
			toCompareWith = secondUserLanguage.ID
		} else {
			toCompareWith = firstUserLanguage.ID
		}

		if lang == toCompareWith {
			if err := h.handleSetLanguagePair(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "⚠️ "+h.services.Localizer.L("LanguagesMustBeDifferent")+"\n"+h.services.Localizer.L("ChooseFirstLanguage"), "setFirstLanguage", "settings", languages); err != nil {
				logger.Error(err)
			}
			return ctx.SendStatus(http.StatusOK)
		}

		if isSettingFirstLanguage {
			dbUser.FirstLanguage = lang
		} else {
			dbUser.SecondLanguage = lang
		}

		if err := h.handleUpdateUserSettings(update.CallbackQuery.ID, &dbUser); err != nil {
			h.services.Telegram.SendText(update.CallbackQuery.Message.Chat.ID, h.services.Localizer.L("SomethingWentWrong"))
			logger.Error(err)
			return ctx.SendStatus(http.StatusOK)
		}

		if err := h.services.Telegram.DeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID); err != nil {
			logger.Error(err)
		}

		if _, err := h.services.Telegram.SendText(update.CallbackQuery.Message.Chat.ID, h.services.Localizer.L("SettingsUpdated")); err != nil {
			logger.Error(err)
		}

		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.Text == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	h.services.Localizer.ChangeLanguage(firstUserLanguage.Code)

	if firstUserLanguage.ID == 0 {
		if _, err := h.services.Telegram.SendText(update.Message.Chat.ID, h.services.Localizer.L("SetFirstLanguageWarning")); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if secondUserLanguage.ID == 0 {
		if _, err := h.services.Telegram.SendText(update.Message.Chat.ID, h.services.Localizer.L("SetSecondLanguageWarning")); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	err = h.services.Telegram.SendTypingChatAction(update.Message.Chat.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if detectedLanguage, exists := h.services.LanguageDetector.Detect(strings.TrimSpace(update.Message.Text)); exists {
		var from, to models.Language

		if detectedLanguage == firstUserLanguage.EnglishText {
			from = firstUserLanguage
			to = secondUserLanguage
		} else if detectedLanguage == secondUserLanguage.EnglishText {
			from = secondUserLanguage
			to = firstUserLanguage
		} else {
			// TODO: send a message to the user that the detected language is not in the language pair
			h.services.Telegram.SendText(update.Message.Chat.ID, h.services.Localizer.L("SomethingWentWrong"))
			return ctx.SendStatus(http.StatusOK)
		}

		if err := h.handleTextTranslation(update.Message.Chat.ID, update.Message.MessageID, dbUser, strings.TrimSpace(update.Message.Text), from, to, strconv.Itoa(update.Message.From.ID)); err != nil {
			logger.Error(err)
		}

	} else {
		logger.Info("Language not found")
		h.services.Telegram.SendText(update.Message.Chat.ID, h.services.Localizer.L("SomethingWentWrong"))
	}

	return ctx.SendStatus(http.StatusOK)
}
