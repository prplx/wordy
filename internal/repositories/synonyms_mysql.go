package repositories

import (
	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type SynonymsRepo struct {
	db *gorm.DB
}

func NewSynonymsRepository(db *gorm.DB) *SynonymsRepo {
	return &SynonymsRepo{
		db: db,
	}
}

func (r *SynonymsRepo) Create(synonyms []models.Synonym) (int64, error) {
	result := r.db.Create(synonyms)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
