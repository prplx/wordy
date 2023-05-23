package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type ExpressionsService struct {
	repository repositories.Expressions
}

func NewExpressionsService(repository repositories.Expressions) *ExpressionsService {
	return &ExpressionsService{
		repository: repository,
	}
}

func (s *ExpressionsService) GetByText(text string) (models.Expression, error) {
	return s.repository.GetByText(text)
}

func (s *ExpressionsService) Create(expression *models.Expression) (uint, error) {
	return s.repository.Create(expression)
}

func (s *ExpressionsService) GetByTextWithTranslationExamplesAudio(text string) (models.Expression, error) {
	return s.repository.GetByTextWithTranslationExamplesAudio(text)
}
