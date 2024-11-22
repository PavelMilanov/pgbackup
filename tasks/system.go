package tasks

import (
	"time"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Sheduler struct {
	Location time.Location
}

func NewTaskScheduler(location time.Location) *Sheduler {
	return &Sheduler{Location: location}
}

func (s *Sheduler) Start(sql *gorm.DB) {
	go s.StartSystemTasks(sql)
	go s.CheckDBconnection(sql)
}

func (s *Sheduler) StartSystemTasks(sql *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	logrus.Debug("Системный планировщик запущен")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			// запуск в 00 часов
			if now.Hour() == 0 && now.Minute() == 0 {
				settings, _ := db.GetSettings(sql)
				logrus.Infof("Запущена задача по удалению дампов старше %d дней", settings.BackupCount)
				files := system.ParseOldFiles(float64(settings.BackupCount))
				db.DeleteOldBackups(files, sql)
			}
			// запуск в 01 часов
			if now.Hour() == 1 && now.Minute() == 0 {
				DB := db.GetDbAll(sql)
				for _, item := range DB {
					size := item.GetDBSize()
					item.Size = size
					sql.Save(&item)
				}
			}
		}
	}
}

func (s *Sheduler) CheckDBconnection(sql *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	logrus.Debug("Планировщик баз данных запущен")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			DB := db.GetDbAll(sql)
			for _, item := range DB {
				status := item.CheckConnection()
				if item.Status && !status {
					item.Status = status
					sql.Save(&item)
					logrus.Warnf("Потеряно соединение с базой данных %s", item.Alias)
				} else if !item.Status && status {
					item.Status = status
					sql.Save(&item)
					logrus.Infof("Восстановлено соединение с базой данных %s", item.Alias)
				}
			}
		}
	}
}
