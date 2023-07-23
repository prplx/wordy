package models

import (
	"gorm.io/gorm"
)

type Audio struct {
	gorm.Model
	Url          string `gorm:"type:varchar(255);not null"`
	ExpressionID uint   `gorm:"not null"`
	Expression   Expression
}
