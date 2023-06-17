package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type SynonymsService struct {
	repository repositories.Synonyms
}

func NewSynonymsService(repository repositories.Synonyms) *SynonymsService {
	return &SynonymsService{
		repository: repository,
	}
}

func (s *SynonymsService) Create(synonyms []models.Synonym) (int64, error) {
	return s.repository.Create(synonyms)
}
