package models

import (
	"gorm.io/gorm"
)

type Expression struct {
	gorm.Model
	UserID         uint `gorm:"not null"`
	User           User
	Text           string `gorm:"type:varchar(255);not null;index"`
	FromLanguageID uint   `gorm:"not null"`
	ToLanguageID   uint   `gorm:"not null"`
	Audio          []Audio
	Translations   []Translation
	Examples       []Example
	Synonyms       []Synonym
}

func (e Expression) GetText() string {
	return e.Text
}
