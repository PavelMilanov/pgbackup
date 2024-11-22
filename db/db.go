package db

import (
	"github.com/PavelMilanov/pgbackup/config"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLite struct {
	Name string
}

func NewDatabase(sql *SQLite, timer *cron.Cron) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sql.Name), &gorm.Config{PrepareStmt: true})
	if err != nil {
		logrus.Fatal("Ошибка при подключении к базе данных")
	}
	logrus.Info("Соединение с базой данных установлено")
	automigrate(db)
	setDefaultSettings(db)
	return db
}

func CloseDatabase(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		logrus.Fatal("Ошибка при закрытии соединения с базой данных:", err)
	}
}

func automigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Database{}, &Schedule{}, &Backup{}, &User{}, &Token{}, &Setting{}); err != nil {
		logrus.Fatalf("%s", err)
		return
	}
}

func setDefaultSettings(db *gorm.DB) {
	var settings Setting
	if err := db.FirstOrCreate(&settings, Setting{BackupCount: config.DEFAULT_BACKUP_EXPIRED_DAYS}).Error; err != nil {
		logrus.Fatal("Ошибка при создании или получении настроек")
	}
}
