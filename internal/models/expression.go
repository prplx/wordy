package models

import (
	"gorm.io/gorm"
)

type Expression struct {
	gorm.Model
	UserId       int
	User         User
	Text         string `gorm:"type:varchar(255);not null;index" json:"text"`
	LanguageId   int    `gorm:"type:int;not null" json:"language"`
	Audio        []Audio
	Translations []Translation
	Examples     []Example
}

func (e Expression) GetText() string {
	return e.Text
}
