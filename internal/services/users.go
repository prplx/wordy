package services

import (
	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/repositories"
)

type UsersService struct {
	repository repositories.Users
}

func NewUsersService(repository repositories.Users) *UsersService {
	return &UsersService{
		repository: repository,
	}
}

func (s *UsersService) Create(user *models.User) (uint, error) {
	return s.repository.Create(user)
}

func (s *UsersService) GetByTgID(id uint) (models.User, error) {
	return s.repository.GetByTgID(id)
}

func (s *UsersService) Update(user *models.User) error {
	return s.repository.Update(user)
}
