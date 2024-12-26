package tasks

import (
	"time"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Sheduler struct {
	Location time.Location
}

func NewTaskScheduler(location time.Location) *Sheduler {
	return &Sheduler{Location: location}
}

func (s *Sheduler) Start(conn *db.SQLite, timer *cron.Cron) {
	go s.StartSystemTasks(conn)
	go s.CheckDBconnection(conn)
	s.initBackupsTasks(conn, timer)
}

func (s *Sheduler) StartSystemTasks(conn *db.SQLite) {
	ticker := time.NewTicker(1 * time.Hour)
	logrus.Debug("Системный планировщик запущен")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			// запуск в 00 часов
			if now.Hour() == 0 && now.Minute() == 0 {
				settings, _ := db.GetSettings(conn.Sql)
				logrus.Infof("Запущена задача по удалению дампов старше %d дней", settings.BackupCount)
				files := system.ParseOldFiles(float64(settings.BackupCount))
				db.DeleteOldBackups(files, conn)
			}
			// запуск в 01 часов
			if now.Hour() == 1 && now.Minute() == 0 {
				DB := db.GetDbAll(conn.Sql)
				for _, item := range DB {
					size := item.GetDBSize()
					item.Size = size
					conn.Mutex.Lock()
					conn.Sql.Save(&item)
					conn.Mutex.Unlock()
				}
			}
		}
	}
}

func (s *Sheduler) CheckDBconnection(conn *db.SQLite) {
	ticker := time.NewTicker(1 * time.Minute)
	logrus.Debug("Планировщик баз данных запущен")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			DB := db.GetDbAll(conn.Sql)
			for _, item := range DB {
				status := item.CheckConnection()
				if item.Status && !status {
					item.Status = status
					conn.Mutex.Lock()
					conn.Sql.Save(&item)
					conn.Mutex.Unlock()
					logrus.Warnf("Потеряно соединение с базой данных %s", item.Alias)
				} else if !item.Status && status {
					item.Status = status
					conn.Mutex.Lock()
					conn.Sql.Save(&item)
					conn.Mutex.Unlock()
					logrus.Debugf("Восстановлено соединение с базой данных %s", item.Alias)
				}
			}
		}
	}
}
