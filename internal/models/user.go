package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID       int    `gorm:"unique;not null"`
	TelegramUsername string `gorm:"type:varchar(255);unique;not null"`
	FirstName        string `gorm:"type:varchar(255);not null"`
	LastName         string `gorm:"type:varchar(255)"`
	FirstLanguage    uint
	SecondLanguage   uint
	Expressions      []*Expression `gorm:"many2many:user_expressions;"`
}
