package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type ExamplesService struct {
	repository repositories.Examples
}

func NewExamplesService(repository repositories.Examples) *ExamplesService {
	return &ExamplesService{
		repository: repository,
	}
}

func (s *ExamplesService) QueryByExpressionID(expressionID int) ([]models.Example, error) {
	return s.repository.QueryByExpressionID(expressionID)
}

func (s *ExamplesService) Create(examples []models.Example) (int64, error) {
	return s.repository.Create(examples)
}
