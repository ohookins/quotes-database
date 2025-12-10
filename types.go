package main

import (
	"gorm.io/gorm"
)

type Quote struct {
	gorm.Model
	Id   string `gorm:"primaryKey"`
	Data string
}
