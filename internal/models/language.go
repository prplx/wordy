package models

import (
	"gorm.io/gorm"
)

type Language struct {
	gorm.Model
	Code        string `gorm:"type:varchar(3);unique;not null" json:"code"`
	Text        string `gorm:"type:varchar(255);not null" json:"text"`
	EnglishText string `gorm:"type:varchar(255);not null" json:"english_text"`
	Emoji       string `gorm:"type:varchar(10);not null" json:"emoji"`
}
