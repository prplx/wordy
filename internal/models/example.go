package models

import (
	"gorm.io/gorm"
)

type Example struct {
	gorm.Model
	Text         string `gorm:"type:varchar(255);not null" json:"text"`
	ExpressionID uint   `gorm:"not null"`
	Expression   Expression
}

func (e Example) GetText() string {
	return e.Text
}
