package repositories

import (
	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type Users interface {
	Create(user *models.User) (uint, error)
	Get(id uint) (models.User, error)
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

type Synonyms interface {
	Create(synonyms []models.Synonym) (int64, error)
}

type Audio interface {
	GetByExpressionID(expressionID int) (models.Audio, error)
	Create(audio models.Audio) (int64, error)
}

type Repositories struct {
	Users        Users
	Expressions  Expressions
	Languages    Languages
	Translations Translations
	Examples     Examples
	Synonyms     Synonyms
	Audio        Audio
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:        NewUsersRepository(db),
		Expressions:  NewExpressionsRepository(db),
		Languages:    NewLanguagesRepository(db),
		Translations: NewTranslationsRepository(db),
		Examples:     NewExamplesRepository(db),
		Synonyms:     NewSynonymsRepository(db),
		Audio:        NewAudioRepository(db),
	}
}
