package v1

import (
	"errors"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
		// TODO: implement
		return ctx.SendStatus(http.StatusOK)
	}

	if update.CallbackQuery.Data == "setFirstLanguage" || update.CallbackQuery.Data == "setSecondLanguage" {
		if _, err := h.handleSetLanguage(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Id, update.CallbackQuery.Data, languages); err != nil {
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
		dbUser.SecondLanguage = int(secondLanguage.ID)

		if err := h.handleUpdateUserSettings(update.CallbackQuery.Id, &dbUser); err != nil {
			logger.Error(err)
		}
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

func (h *Handlers) handleStartCommand(chatId int64) (string, error) {
	return h.services.Telegram.SendText(chatId, h.services.Localizer.L("Greeting"))
}

func (h *Handlers) handleSettingsCommand(chatId int64) (string, error) {
	return h.services.Telegram.SendReplyKeyboard(chatId, []types.KeyboardButton{{Text: h.services.Localizer.L("SetLanguages"), CallbackData: "setLanguagePair"}, {Text: h.services.Localizer.L("SetFirstLanguage"), CallbackData: "setFirstLanguage"}, {Text: h.services.Localizer.L("SetSecondLanguage"), CallbackData: "setSecondLanguage"}}, h.services.Localizer.L("BotSettings"))
}

func (h *Handlers) handleUpdateUserSettings(queryId string, user *models.User) error {
	if err := h.services.Users.Update(user); err != nil {
		return err
	}

	return h.services.Telegram.AnswerCallbackQuery(queryId, h.services.Localizer.L("SettingsUpdated"))
}

func (h *Handlers) handleSetLanguage(chatId int64, queryId string, command string, languages []models.Language) (string, error) {
	if err := h.services.Telegram.AnswerCallbackQuery(queryId, ""); err != nil {
		return "", err
	}

	var buttons []types.KeyboardButton
	for _, language := range languages {
		buttons = append(buttons, types.KeyboardButton{Text: language.Text + " " + language.Emoji, CallbackData: command + ": " + language.Code})
	}

	return h.services.Telegram.SendReplyKeyboard(chatId, buttons, h.services.Localizer.L("ChooseLanguage"))
}

func (h *Handlers) handleTextTranslation(chatId int64, replyMessageId int, userId int, languageId int, text, from, to, tgUserId string) error {
	wg := sync.WaitGroup{}
	wg.Add(4)
	translations, examples, audio, synonyms := "", "", "", ""
	topErr := error(nil)

	dbExpression, err := h.services.Expressions.GetByTextWithAllData(text)

	// Happy case, everything is in the database
	if helpers.IsExpressionWithAllData(dbExpression) {
		translations = h.buildTranslationsBlock(dbExpression.Translations)
		synonyms = h.buildSynonymsBlock(dbExpression.Synonyms)
		examples = h.buildExamplesBlock(dbExpression.Examples)
		audio = dbExpression.Audio[0].Url
		message := helpers.BuildMessage(translations, synonyms, examples, audio)
		if _, err := h.services.Telegram.SendText(chatId, message, replyMessageId); err != nil {
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
			translations = h.buildTranslationsBlock(dbTranslations)
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
			examples = h.buildExamplesBlock(dbExamples)
		}

		wg.Done()
	}(dbExpression.Examples, &wg)

	if len(strings.Fields(text)) == 1 {
		go func(s []models.Synonym, wg *sync.WaitGroup) {
			if len(s) > 0 {
				synonyms = helpers.BuildMessageFromSliceOfTexted(s)
				wg.Done()
				return
			}

			dbSynonyms, err := h.createSynonyms(int(dbExpression.ID), from, text)
			if err != nil {
				topErr = err
			}
			if len(dbSynonyms) > 0 {
				synonyms = h.buildSynonymsBlock(dbSynonyms)
			}

			wg.Done()
		}(dbExpression.Synonyms, &wg)
	} else {
		wg.Done()
	}

	go func(a []models.Audio, userId string, wg *sync.WaitGroup) {
		if len(a) > 0 {
			audio = a[0].Url
			wg.Done()
			return
		}

		dbAudio, err := h.createAudio(int(dbExpression.ID), from, text, userId)
		if err != nil {
			topErr = err
		}
		if dbAudio.Url != "" {
			audio = dbAudio.Url
		}

		wg.Done()
	}(dbExpression.Audio, tgUserId, &wg)

	wg.Wait()

	var messageToSend string

	if topErr != nil {
		logger.Error(topErr)
		messageToSend = helpers.BuildMessage(h.services.Localizer.L("SomethingWentWrong"))
	} else {
		messageToSend = helpers.BuildMessage(translations, synonyms, examples, audio)
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

func (h *Handlers) createSynonyms(expressionId int, language, text string) ([]models.Synonym, error) {
	synonyms, err := h.services.Translator.GenerateSynonyms(text, language)
	if err != nil {
		return nil, err
	}

	synonymsToCreate := make([]models.Synonym, 0, len(synonyms))
	for _, synonym := range synonyms {
		synonymsToCreate = append(synonymsToCreate, models.Synonym{Text: synonym, ExpressionId: expressionId})
	}

	if _, err := h.services.Synonyms.Create(synonymsToCreate); err != nil {
		return nil, err
	}

	return synonymsToCreate, nil
}

func (h *Handlers) createAudio(expressionId int, language, text string, userId string) (models.Audio, error) {
	audioUrl, err := h.services.TextToSpeech.Convert(text, language, userId)
	if err != nil {
		return models.Audio{}, err
	}

	audioToCreate := models.Audio{Url: audioUrl, ExpressionId: expressionId}

	if _, err := h.services.Audio.Create(audioToCreate); err != nil {
		return audioToCreate, err
	}

	return audioToCreate, nil
}

func (h *Handlers) buildTranslationsBlock(translations []models.Translation) string {
	return helpers.AddBlockTitleToText(h.services.Localizer.L("Translation", "2.5"), helpers.BuildMessageFromSliceOfTexted(translations))
}

func (h *Handlers) buildExamplesBlock(examples []models.Example) string {
	return helpers.AddBlockTitleToText(h.services.Localizer.L("Example", "2.5"), helpers.BuildMessageFromSliceOfTexted(examples))
}

func (h *Handlers) buildSynonymsBlock(synonyms []models.Synonym) string {
	return helpers.AddBlockTitleToText(h.services.Localizer.L("Synonym", "2.5"), helpers.BuildMessageFromSliceOfTexted(synonyms))
}
