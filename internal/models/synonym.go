package models

import (
	"gorm.io/gorm"
)

type Synonym struct {
	gorm.Model
	Text         string `gorm:"type:varchar(255);not null" json:"text"`
	ExpressionId int
	Expression   Expression
}

func (t Synonym) GetText() string {
	return t.Text
}
