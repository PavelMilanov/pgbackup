package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLite struct {
	Name string
}

func NewDatabase(sql *SQLite) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sql.Name), &gorm.Config{})
	if err != nil {
		logrus.Error("failed to connect database")
	}
	automigrate(db)
	return db
}

func automigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Database{}, &Schedule{}, &Backup{}, &Token{}); err != nil {
		logrus.Error("failed to migrate")
	}
}
