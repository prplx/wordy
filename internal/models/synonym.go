package models

import (
	"gorm.io/gorm"
)

type Synonym struct {
	gorm.Model
	Text         string `gorm:"type:varchar(255);not null"`
	ExpressionID uint   `gorm:"not null"`
	Expression   Expression
}

func (t Synonym) GetText() string {
	return t.Text
}
