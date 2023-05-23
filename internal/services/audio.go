package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type AudioService struct {
	repository repositories.Audio
}

func NewAudioService(repository repositories.Audio) *AudioService {
	return &AudioService{
		repository: repository,
	}
}

func (s *AudioService) GetByExpressionId(expressionId int) (models.Audio, error) {
	return s.repository.GetByExpressionId(expressionId)
}

func (s *AudioService) Create(audio models.Audio) (int64, error) {
	return s.repository.Create(audio)
}
