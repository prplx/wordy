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

func (s *ExpressionsService) GetByTextWithAllData(text string) (models.Expression, error) {
	return s.repository.GetByTextWithAllData(text)
}

func (s *ExpressionsService) GetUserByID(expression *models.Expression, user *models.User) error {
	return s.repository.GetUserByID(expression, user)
}

func (s *ExpressionsService) AddUser(expression *models.Expression, user *models.User) error {
	return s.repository.AddUser(expression, user)
}
