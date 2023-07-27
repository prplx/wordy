package v1

import (
	"errors"
	"strings"
	"sync"

	"github.com/prplx/wordy/internal/helpers"
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/types"
)

func (h *Handlers) handleStartCommand(chatId int64) (string, error) {
	return h.services.Telegram.SendText(chatId, h.services.Localizer.L("Greeting"))
}

func (h *Handlers) handleSettingsCommand(chatId int64, messageId ...int) (string, error) {
	buttons := []types.KeyboardButton{{Text: h.services.Localizer.L("SetFirstLanguage"), CallbackData: "setLanguage (1)"}, {Text: h.services.Localizer.L("SetSecondLanguage"), CallbackData: "setLanguage (2)"}}
	if len(messageId) == 0 {
		return h.services.Telegram.SendReplyKeyboard(chatId, buttons, h.services.Localizer.L("BotSettings"))
	} else {
		return "", h.services.Telegram.EditMessage(chatId, messageId[0], h.services.Localizer.L("BotSettings"), buttons)
	}
}

func (h *Handlers) handleSetLanguagePair(chatId int64, messageId int, text, command, menuBackCommand string, languages []models.Language) error {
	var buttons []types.KeyboardButton
	for _, language := range languages {
		buttons = append(buttons, types.KeyboardButton{Text: language.Text + " " + language.Emoji, CallbackData: command + " (" + language.Code + ")"})
	}
	buttons = append(buttons, types.KeyboardButton{Text: "← " + h.services.Localizer.L("Back"), CallbackData: menuBackCommand})

	return h.services.Telegram.EditMessage(chatId, messageId, text, buttons)
}

func (h *Handlers) handleUpdateUserSettings(queryId string, user *models.User) error {
	return h.services.Users.Update(user)
}

func (h *Handlers) handleTextTranslation(chatId int64, replyMessageId int, user models.User, text string, from, to models.Language, tgUserId string) error {
	wg := sync.WaitGroup{}
	isUsersFirstLanguage := uint(user.FirstLanguage) == from.ID
	if isUsersFirstLanguage {
		wg.Add(1)
	} else {
		wg.Add(4)
	}
	translations, examples, audio, synonyms := "", "", "", ""
	topErr := error(nil)

	dbExpression, err := h.services.Expressions.GetByTextWithAllData(text)

	// Happy case, everything is in the database
	if dbExpression.FromLanguageID == from.ID && dbExpression.ToLanguageID == to.ID && helpers.IsExpressionWithAllData(dbExpression) {
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
			Text: text, UserID: user.ID, FromLanguageID: from.ID, ToLanguageID: to.ID,
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

		dbTranslations, err := h.createTranslations(dbExpression.ID, from.Text, to.Text, text)
		if err != nil {
			topErr = err
		}

		if len(dbTranslations) > 0 {
			translations = h.buildTranslationsBlock(dbTranslations)
		}

		wg.Done()
	}(dbExpression.Translations, &wg)

	if !isUsersFirstLanguage {
		go func(ex []models.Example, wg *sync.WaitGroup) {
			if len(ex) > 0 {
				examples = helpers.BuildMessageFromSliceOfTexted(ex)
				wg.Done()
				return
			}

			dbExamples, err := h.createExamples(dbExpression.ID, from.Text, text)
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

				dbSynonyms, err := h.createSynonyms(dbExpression.ID, from.Text, text)
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

			dbAudio, err := h.createAudio(dbExpression.ID, from.Text, text, userId)
			if err != nil {
				topErr = err
			}
			if dbAudio.Url != "" {
				audio = dbAudio.Url
			}

			wg.Done()
		}(dbExpression.Audio, tgUserId, &wg)

	}

	wg.Wait()

	var messageToSend string

	if topErr != nil {
		h.services.Logger.Error(topErr)
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

func (h *Handlers) createTranslations(expressionID uint, from, to, text string) ([]models.Translation, error) {
	translations, err := h.services.Translator.Translate(text, from, to)
	if err != nil {
		return nil, err
	}

	translationsToCreate := make([]models.Translation, 0, len(translations))
	for _, line := range translations {
		translationsToCreate = append(translationsToCreate, models.Translation{Text: line, ExpressionID: expressionID})
	}
	if _, err := h.services.Translations.Create(translationsToCreate); err != nil {
		return nil, err
	}

	return translationsToCreate, nil
}

func (h *Handlers) createExamples(expressionID uint, language, text string) ([]models.Example, error) {
	examples, err := h.services.Translator.GenerateExamples(text, language)
	if err != nil {
		return nil, err
	}

	examplesToCreate := make([]models.Example, 0, len(examples))
	for _, example := range examples {
		examplesToCreate = append(examplesToCreate, models.Example{Text: example, ExpressionID: expressionID})
	}

	if _, err := h.services.Examples.Create(examplesToCreate); err != nil {
		return nil, err
	}

	return examplesToCreate, nil
}

func (h *Handlers) createSynonyms(expressionID uint, language, text string) ([]models.Synonym, error) {
	synonyms, err := h.services.Translator.GenerateSynonyms(text, language)
	if err != nil {
		return nil, err
	}

	synonymsToCreate := make([]models.Synonym, 0, len(synonyms))
	for _, synonym := range synonyms {
		synonymsToCreate = append(synonymsToCreate, models.Synonym{Text: synonym, ExpressionID: expressionID})
	}

	if _, err := h.services.Synonyms.Create(synonymsToCreate); err != nil {
		return nil, err
	}

	return synonymsToCreate, nil
}

func (h *Handlers) createAudio(expressionID uint, language, text string, userID string) (models.Audio, error) {
	audioUrl, err := h.services.TextToSpeech.Convert(text, language, userID)
	if err != nil {
		return models.Audio{}, err
	}

	audioToCreate := models.Audio{Url: audioUrl, ExpressionID: expressionID}

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
