package services

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
	"github.com/prplx/wordy/internal/types"
	"github.com/prplx/wordy/pkg/jsonlog"
)

type Users interface {
	Create(user *models.User) (uint, error)
	GetByTgID(id uint) (models.User, error)
	Update(user *models.User) error
}

type Expressions interface {
	Create(expression *models.Expression) (uint, error)
	GetByText(text string) (models.Expression, error)
	GetByTextWithAllData(text string) (models.Expression, error)
	GetUserByID(expression *models.Expression, user *models.User) error
	AddUser(expression *models.Expression, user *models.User) error
}

type Telegram interface {
	SendText(chatID int64, text string, replyMessageID ...int) (string, error)
	SendReplyKeyboard(chatID int64, buttons []types.KeyboardButton, text string) (string, error)
	SendTypingChatAction(chatID int64) error
	EditMessage(chatID int64, messageID int, text string, buttons ...[]types.KeyboardButton) error
	EditReplyMarkup(chatID int64, messageID int) error
	DeleteMessage(chatID int64, messageID int) error
}

type Translator interface {
	Translate(text, sourceLang, targetLang string) ([]string, error)
	GenerateExamples(text, sourceLang string) ([]string, error)
	GenerateSynonyms(text, sourceLang string) ([]string, error)
}

type Languages interface {
	Query() ([]models.Language, error)
	GetByCode(code string) (models.Language, error)
}

type Translations interface {
	QueryByExpressionID(expressionID int) ([]models.Translation, error)
	Create(translations []models.Translation) (int64, error)
}

type Examples interface {
	QueryByExpressionID(expressionID int) ([]models.Example, error)
	Create(examples []models.Example) (int64, error)
}

type TextToSpeech interface {
	Convert(text, lang, userID string) (string, error)
}

type Audio interface {
	GetByExpressionID(expressionID int) (models.Audio, error)
	Create(audio models.Audio) (int64, error)
}

type Synonyms interface {
	Create(synonyms []models.Synonym) (int64, error)
}

type Localizer interface {
	L(id string, count ...interface{}) string
	ChangeLanguage(lang string)
}

type LanguageDetector interface {
	Detect(text string) (string, bool)
}

type Services struct {
	Users            Users
	Telegram         Telegram
	Translator       Translator
	Expressions      Expressions
	Languages        Languages
	Translations     Translations
	Examples         Examples
	TextToSpeech     TextToSpeech
	Synonyms         Synonyms
	Audio            Audio
	Localizer        Localizer
	LanguageDetector LanguageDetector
	Logger           *jsonlog.Logger
}

type Deps struct {
	Repositories    repositories.Repositories
	LocalizerBundle *i18n.Bundle
	Logger          *jsonlog.Logger
}

func NewServices(deps Deps) *Services {
	return &Services{
		Expressions:      NewExpressionsService(deps.Repositories.Expressions),
		Users:            NewUsersService(deps.Repositories.Users),
		Telegram:         NewTelegramService(),
		Translator:       NewTranslatorService(),
		Languages:        NewLanguagesService(deps.Repositories.Languages),
		Translations:     NewTranslationsService(deps.Repositories.Translations),
		Examples:         NewExamplesService(deps.Repositories.Examples),
		Audio:            NewAudioService(deps.Repositories.Audio),
		TextToSpeech:     NewTextToSpeechService(),
		Synonyms:         NewSynonymsService(deps.Repositories.Synonyms),
		Localizer:        NewLocalizerService(deps.LocalizerBundle),
		LanguageDetector: NewLanguageDetectorService(),
		Logger:           deps.Logger,
	}
}
