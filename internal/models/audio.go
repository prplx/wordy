package models

import (
	"gorm.io/gorm"
)

type Audio struct {
	gorm.Model
	Url          string `gorm:"type:varchar(255);not null" json:"url"`
	ExpressionId int
	Expression   Expression
}
