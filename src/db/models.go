package db

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Модель базы данных для обслуживания.
type Database struct {
	gorm.Model
	ID        int `gorm:"primaryKey"`
	Alias     string
	Name      string
	Host      string
	Port      int
	Username  string
	Password  string
	Status    bool
	Size      string
	Schedules []Schedule `gorm:"constraint:OnDelete:CASCADE;"`
	Backups   []Backup   `gorm:"constraint:OnDelete:CASCADE;"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime:false"`
}

// Модель расписания выполнения бекапов для выбранной базы даннах.
type Schedule struct {
	gorm.Model
	ID            int    `gorm:"primaryKey"`
	Directory     string `gorm:"unique"`
	Time          string
	Frequency     string
	Status        string
	DatabaseAlias string
	DatabaseID    int
	Backups       []Backup  `gorm:"constraint:OnDelete:CASCADE;"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:false"`
}

// Модель с метаданными бекапа.
type Backup struct {
	gorm.Model
	ID            int `gorm:"primaryKey"`
	Date          string
	Size          string
	LeadTime      string
	Status        bool
	Directory     string
	Dump          string `gorm:"unique"`
	DatabaseAlias string
	ScheduleID    int
	DatabaseID    int
	UpdatedAt     time.Time `gorm:"autoUpdateTime:false"`
}

// Модель пользователя приложения.
type User struct {
	gorm.Model
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	Token     Token
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

// Модель токена аутентицикации.
type Token struct {
	gorm.Model
	ID                   int `gorm:"primaryKey"`
	Hash                 string
	UserID               int
	jwt.RegisteredClaims `gorm:"-"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime:false"`
}

// Модель настроек приложения.
type Setting struct {
	gorm.Model
	ID          int `gorm:"primaryKey"`
	BackupCount int
	Version     string `gorm:"-"`
}
