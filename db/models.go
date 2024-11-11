package db

import (
	"gorm.io/gorm"
)

// Модель базы данных для обслуживания.
type Database struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	Name      string
	Host      string
	Port      int
	Username  string
	Password  string
	Status    bool
	Size      string
	Schedules []Schedule `gorm:"constraint:OnDelete:CASCADE;"`
	Backups   []Backup   `gorm:"constraint:OnDelete:CASCADE;"`
}

// Модель расписания выполнения бекапов для выбранной базы даннах.
type Schedule struct {
	gorm.Model
	ID           int    `gorm:"primaryKey"`
	Directory    string `gorm:"unique"`
	Time         string
	Frequency    string
	Status       string
	DatabaseName string
	DatabaseID   int
	Backups      []Backup `gorm:"constraint:OnDelete:CASCADE;"`
}

// Модель с метаданными бекапа.
type Backup struct {
	gorm.Model
	ID           int `gorm:"primaryKey"`
	Date         string
	Size         string
	LeadTime     string
	Status       string
	Directory    string
	Dump         string `gorm:"unique"`
	DatabaseName string
	ScheduleID   int
	DatabaseID   int
}

// Модель токена аутентицикации
type Token struct {
	gorm.Model
	ID   int `gorm:"primaryKey"`
	Hash string
}
