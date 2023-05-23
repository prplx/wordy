package models

import (
	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	Text         string `gorm:"type:varchar(255);not null" json:"text"`
	ExpressionId int
	Expression   Expression
}

func (t Translation) GetText() string {
	return t.Text
}
