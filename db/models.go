package db

import "gorm.io/gorm"

type Database struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Name     string
	Host     string
	Port     int
	Username string
	Password string
	Schedule []Schedule `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Schedule struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey"`
	Directory  string `gorm:"unique;not null"`
	Time       string
	Frequency  string
	DatabaseID uint
	Backups    []Backup `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Backup struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Date       string
	Size       string
	LeadTime   string
	Status     string
	Dump       string `gorm:"unique;not null"`
	ScheduleID uint
}

type Token struct {
	gorm.Model
	ID   uint `gorm:"primaryKey"`
	Hash string
}
