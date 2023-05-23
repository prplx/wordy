package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
	"github.com/prplx/wordy/internal/types"
)

type Users interface {
	Create(user *models.User) (uint, error)
	GetByTgId(id uint) (models.User, error)
	Update(user *models.User) error
}

type Expressions interface {
	Create(expression *models.Expression) (uint, error)
	GetByText(text string) (models.Expression, error)
	GetByTextWithTranslationExamplesAudio(text string) (models.Expression, error)
}

type Telegram interface {
	AnswerCallbackQuery(queryId string, text string) error
	SendText(chatId int64, text string, replyMessageId int) (string, error)
	SendReplyKeyboard(chatId int64, buttons []types.KeyboardButton, text string) (string, error)
}

type Translator interface {
	Translate(text, sourceLang, targetLang string) ([]string, error)
	GenerateExamples(text, sourceLang string) ([]string, error)
}

type Languages interface {
	Query() ([]models.Language, error)
	GetByCode(code string) (models.Language, error)
}

type Translations interface {
	QueryByExpressionId(expressionId int) ([]models.Translation, error)
	Create(translations []models.Translation) (int64, error)
}

type Examples interface {
	QueryByExpressionId(expressionId int) ([]models.Example, error)
	Create(examples []models.Example) (int64, error)
}

type TextToSpeech interface {
	Convert(text string, lang string) (string, error)
}

type Audio interface {
	GetByExpressionId(expressionId int) (models.Audio, error)
	Create(audio models.Audio) (int64, error)
}

type Services struct {
	Users        Users
	Telegram     Telegram
	Translator   Translator
	Expressions  Expressions
	Languages    Languages
	Translations Translations
	Examples     Examples
	TextToSpeech TextToSpeech
	Audio        Audio
}

type Deps struct {
	Repositories repositories.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{
		Expressions:  NewExpressionsService(deps.Repositories.Expressions),
		Users:        NewUsersService(deps.Repositories.Users),
		Telegram:     NewTelegramService(),
		Translator:   NewTranslatorService(),
		Languages:    NewLanguagesService(deps.Repositories.Languages),
		Translations: NewTranslationsService(deps.Repositories.Translations),
		Examples:     NewExamplesService(deps.Repositories.Examples),
		Audio:        NewAudioService(deps.Repositories.Audio),
		TextToSpeech: NewTextToSpeechService(),
	}
}
