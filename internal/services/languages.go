package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type LanguagesService struct {
	repository repositories.Languages
}

func NewLanguagesService(repository repositories.Languages) *LanguagesService {
	return &LanguagesService{
		repository: repository,
	}
}

func (s *LanguagesService) Query() ([]models.Language, error) {
	return s.repository.Query()
}

func (s *LanguagesService) GetByCode(code string) (models.Language, error) {
	return s.repository.GetByCode(code)
}
