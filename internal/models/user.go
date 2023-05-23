package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramId       int    `gorm:"type:int;unique;not null" json:"telegram_id"`
	TelegramUsername string `gorm:"type:varchar(255);unique;not null" json:"telegram_username"`
	FirstName        string `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName         string `gorm:"type:varchar(255)" json:"last_name"`
	FirstLanguage    int    `gorm:"type:int" json:"first_language"`
	SecondLanguage   int    `gorm:"type:int" json:"second_language"`
	Expressions      []Expression
}
