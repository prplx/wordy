package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type TranslationsService struct {
	repository repositories.Translations
}

func NewTranslationsService(repository repositories.Translations) *TranslationsService {
	return &TranslationsService{
		repository: repository,
	}
}

func (s *TranslationsService) QueryByExpressionId(expressionId int) ([]models.Translation, error) {
	return s.repository.QueryByExpressionId(expressionId)
}

func (s *TranslationsService) Create(translations []models.Translation) (int64, error) {
	return s.repository.Create(translations)
}
