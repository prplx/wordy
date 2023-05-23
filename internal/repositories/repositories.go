package repositories

import (
	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type Users interface {
	Create(user *models.User) (uint, error)
	Get(id uint) (models.User, error)
	GetByTgId(id uint) (models.User, error)
	Update(user *models.User) error
}

type Expressions interface {
	Create(expression *models.Expression) (uint, error)
	GetByText(text string) (models.Expression, error)
	GetByTextWithTranslationExamplesAudio(text string) (models.Expression, error)
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

type Audio interface {
	GetByExpressionId(expressionId int) (models.Audio, error)
	Create(audio models.Audio) (int64, error)
}

type Repositories struct {
	Users        Users
	Expressions  Expressions
	Languages    Languages
	Translations Translations
	Examples     Examples
	Audio        Audio
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:        NewUsersRepository(db),
		Expressions:  NewExpressionsRepository(db),
		Languages:    NewLanguagesRepository(db),
		Translations: NewTranslationsRepository(db),
		Examples:     NewExamplesRepository(db),
		Audio:        NewAudioRepository(db),
	}
}
