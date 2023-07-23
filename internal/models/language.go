package models

import (
	"gorm.io/gorm"
)

type Language struct {
	gorm.Model
	Code        string `gorm:"type:varchar(3);unique;not null"`
	Text        string `gorm:"type:varchar(255);not null"`
	EnglishText string `gorm:"type:varchar(255);not null"`
	Emoji       string `gorm:"type:varchar(10);not null"`
}
