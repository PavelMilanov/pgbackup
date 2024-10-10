package db

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Alias     string `gorm:"not null"`
	Directory string `gorm:"unique;not null"`
	Count     string
	Time      string
	Cron      string
	Backups   []Backup `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
