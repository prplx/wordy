package repositories

import (
	"errors"

	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type AudioRepo struct {
	db *gorm.DB
}

func NewAudioRepository(db *gorm.DB) *AudioRepo {
	return &AudioRepo{
		db: db,
	}
}

func (r *AudioRepo) Create(audio models.Audio) (int64, error) {
	result := r.db.Create(&audio)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *AudioRepo) GetByExpressionID(expressionID int) (models.Audio, error) {
	var audio models.Audio
	result := r.db.Where("expression_id = ?", expressionID).First(&audio)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return audio, models.ErrRecordNotFound
	}
	if result.Error != nil {
		return audio, result.Error
	}
	return audio, nil
}
