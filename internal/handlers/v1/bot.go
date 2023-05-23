package v1

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

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
	var update types.Update

	if err := ctx.BodyParser(&update); err != nil {
		logger.Error(err)
		return err
	}

	if update.Message.From.Username != "mmystiq" && update.CallbackQuery.From.Username != "mmystiq" {
		h.services.Telegram.SendText(update.Message.Chat.Id, "The bot is currently under development. Please, come back later.", 0)
		return ctx.SendStatus(http.StatusOK)
	}

	if update.Message.MessageId == 0 && update.CallbackQuery.Id == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	languages, err := h.services.Languages.Query()
	if err != nil {
		h.services.Telegram.SendText(update.Message.Chat.Id, "Something went wrong, please try again later", 0)
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
			log.Println("User not found, creating a new one")

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

	if update.Message.Text == "/settings" {
		if _, err := h.handleSettingsCommand(update.Message.From, update.Message.Chat.Id); err != nil {

			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if update.CallbackQuery.Data == "setSourceLanguage" || update.CallbackQuery.Data == "setTargetLanguage" {
		if _, err := h.handleSetLanguage(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Id, update.CallbackQuery.Data, languages); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	setFirstLanguagePattern := `setSourceLanguage: (\w+)`

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
		}

		return ctx.SendStatus(http.StatusOK)
	}

	setSecondLanguagePattern := `setTargetLanguage: (\w+)`

	re = regexp.MustCompile(setSecondLanguagePattern)
	match = re.FindStringSubmatch(update.CallbackQuery.Data)

	if len(match) > 1 {
		var secondLanguage models.Language
		for _, language := range languages {
			if language.Code == match[1] {
				secondLanguage = language
			}
		}
		dbUser.SecondLanguage = int(secondLanguage.ID)

		if err := h.handleUpdateUserSettings(update.CallbackQuery.Id, &dbUser); err != nil {
			logger.Error(err)
		}
	}

	if update.Message.Text == "" {
		return ctx.SendStatus(http.StatusOK)
	}

	firstUserLanguage, secondUserLanguage := helpers.GetUserFirstAndSecondLanguagesIds(dbUser, languages)

	if firstUserLanguage.ID == 0 || secondUserLanguage.ID == 0 {
		if _, err := h.services.Telegram.SendText(update.Message.Chat.Id, "Please, set source and target languages using /settings first", 0); err != nil {
			logger.Error(err)
		}
		return ctx.SendStatus(http.StatusOK)
	}

	if err := h.handleTextTranslation(update.Message.Chat.Id, update.Message.MessageId, int(dbUser.ID), int(secondUserLanguage.ID), strings.TrimSpace(update.Message.Text), secondUserLanguage.Code, firstUserLanguage.Code); err != nil {
		logger.Error(err)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (h *Handlers) handleSettingsCommand(user types.User, chatId int64) (string, error) {
	return h.services.Telegram.SendReplyKeyboard(chatId, []types.KeyboardButton{{Text: "Set source language", CallbackData: "setSourceLanguage"}, {Text: "Set target language", CallbackData: "setTargetLanguage"}}, "Bot settings")
}

func (h *Handlers) handleUpdateUserSettings(queryId string, user *models.User) error {
	if err := h.services.Users.Update(user); err != nil {
		return err
	}

	return h.services.Telegram.AnswerCallbackQuery(queryId, "Settings updated")
}

func (h *Handlers) handleSetLanguage(chatId int64, queryId string, command string, languages []models.Language) (string, error) {
	if err := h.services.Telegram.AnswerCallbackQuery(queryId, ""); err != nil {
		return "", err
	}
	var buttons []types.KeyboardButton
	for _, language := range languages {
		buttons = append(buttons, types.KeyboardButton{Text: language.Text + " " + language.Emoji, CallbackData: command + ": " + language.Code})
	}

	return h.services.Telegram.SendReplyKeyboard(chatId, buttons, "Choose language")
}

func (h *Handlers) handleTextTranslation(chatId int64, replyMessageId int, userId int, languageId int, text string, from string, to string) error {
	wg := sync.WaitGroup{}
	wg.Add(3)
	translations, examples, audio := "", "", ""
	topErr := error(nil)

	dbExpression, err := h.services.Expressions.GetByTextWithTranslationExamplesAudio(text)
	// Happy case, everything is in the database
	if helpers.IsExpressionFulfilledWithTranslationsExamplesAudio(dbExpression) {
		translations = helpers.BuildMessageFromSliceOfTexted(dbExpression.Translations)
		examples = helpers.BuildMessageFromSliceOfTexted(dbExpression.Examples)
		audio = dbExpression.Audio[0].Url
		if _, err := h.services.Telegram.SendText(chatId, helpers.BuildMessage(translations, examples, audio), replyMessageId); err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		expression := models.Expression{
			Text: text, UserId: userId, LanguageId: languageId,
		}
		if errors.Is(err, models.ErrRecordNotFound) {
			if _, err := h.services.Expressions.Create(&expression); err != nil {
				return err
			}
			dbExpression = expression
		} else {
			return err
		}
	}

	go func(tr []models.Translation, wg *sync.WaitGroup) {
		if len(tr) > 0 {
			translations = helpers.BuildMessageFromSliceOfTexted(tr)
			wg.Done()
			return
		}

		dbTranslations, err := h.createTranslations(int(dbExpression.ID), from, to, text)
		if err != nil {
			topErr = err
		}
		if len(dbTranslations) > 0 {
			translations = helpers.BuildMessageFromSliceOfTexted(dbTranslations)
		}

		wg.Done()
	}(dbExpression.Translations, &wg)

	go func(ex []models.Example, wg *sync.WaitGroup) {
		if len(ex) > 0 {
			examples = helpers.BuildMessageFromSliceOfTexted(ex)
			wg.Done()
			return
		}

		dbExamples, err := h.createExamples(int(dbExpression.ID), from, text)
		if err != nil {
			topErr = err
		}
		if len(dbExamples) > 0 {
			examples = helpers.BuildMessageFromSliceOfTexted(dbExamples)
		}

		wg.Done()
	}(dbExpression.Examples, &wg)

	go func(a []models.Audio, wg *sync.WaitGroup) {
		if len(a) > 0 {
			audio = a[0].Url
			wg.Done()
			return
		}

		dbAudio, err := h.createAudio(int(dbExpression.ID), from, text)
		if err != nil {
			topErr = err
		}
		if dbAudio.Url != "" {
			audio = dbAudio.Url
		}

		wg.Done()
	}(dbExpression.Audio, &wg)

	wg.Wait()

	var messageToSend string

	if topErr != nil {
		logger.Error(topErr)
		messageToSend = helpers.BuildMessage("Something went wrong")
	} else {
		messageToSend = helpers.BuildMessage(translations, examples, audio)
	}

	if messageToSend != "" {
		if _, err := h.services.Telegram.SendText(chatId, messageToSend, replyMessageId); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handlers) createTranslations(expressionId int, from, to, text string) ([]models.Translation, error) {
	translations, err := h.services.Translator.Translate(text, from, to)
	if err != nil {
		return nil, err
	}

	translationsToCreate := make([]models.Translation, 0, len(translations))
	for _, line := range translations {
		translationsToCreate = append(translationsToCreate, models.Translation{Text: line, ExpressionId: expressionId})
	}
	if _, err := h.services.Translations.Create(translationsToCreate); err != nil {
		return nil, err
	}

	return translationsToCreate, nil
}

func (h *Handlers) createExamples(expressionId int, language, text string) ([]models.Example, error) {
	examples, err := h.services.Translator.GenerateExamples(text, language)
	if err != nil {
		return nil, err
	}

	examplesToCreate := make([]models.Example, 0, len(examples))
	for _, example := range examples {
		examplesToCreate = append(examplesToCreate, models.Example{Text: example, ExpressionId: expressionId})
	}

	if _, err := h.services.Examples.Create(examplesToCreate); err != nil {
		return nil, err
	}

	return examplesToCreate, nil
}

func (h *Handlers) createAudio(expressionId int, language, text string) (models.Audio, error) {
	audioUrl, err := h.services.TextToSpeech.Convert(text, language)
	if err != nil {
		return models.Audio{}, err
	}

	audioToCreate := models.Audio{Url: audioUrl, ExpressionId: expressionId}

	if _, err := h.services.Audio.Create(audioToCreate); err != nil {
		return audioToCreate, err
	}

	return audioToCreate, nil
}
