package db

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	ID   uint `gorm:"primaryKey"`
	Hash string
}
