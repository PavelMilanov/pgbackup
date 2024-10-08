package db

import (
	"gorm.io/gorm"
)

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

type Backup struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Alias    string
	Date     string
	Size     string
	LeadTime string
	Run      string `gorm:"not null"`
	Status   string
	Comment  string
	Dump     string `gorm:"unique;not null"`
	TaskID   uint
}

type Token struct {
	gorm.Model
	ID   uint `gorm:"primaryKey"`
	Hash string
}
