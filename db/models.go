package db

import "gorm.io/gorm"

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
}

type Schedule struct {
	gorm.Model
	ID           int    `gorm:"primaryKey"`
	Directory    string `gorm:"unique;not null"`
	Time         string
	Frequency    string
	DatabaseID   int
	DatabaseName string
	LastBackup   string
	Status       string
	Backups      []Backup `gorm:"constraint:OnDelete:CASCADE;"`
}

type Backup struct {
	gorm.Model
	ID         int `gorm:"primaryKey"`
	Date       string
	Size       string
	LeadTime   string
	Status     string
	Directory  string
	Dump       string `gorm:"unique;not null"`
	ScheduleID int
}

type Token struct {
	gorm.Model
	ID   int `gorm:"primaryKey"`
	Hash string
}
