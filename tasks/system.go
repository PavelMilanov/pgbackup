package tasks

import (
	"time"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Sheduler struct {
	Duration time.Duration
}

func NewTaskScheduler(duration time.Duration) *Sheduler {
	return &Sheduler{Duration: duration}
}

// loc, err := time.LoadLocation("America/New_York")
//
//	if err != nil {
//		panic(err)
//	}
//
// now := time.Now().In(loc)
func (s *Sheduler) StartSystemTasks(sql *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	logrus.Debug("Системный планировщик запущен")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			if now.Hour() == 0 && now.Minute() == 0 {
				settings, _ := db.GetSettings(sql)
				logrus.Infof("Запущена задача по удалению бекапов старше %d дней", settings.BackupCount)
				files := system.ParseOldFiles(float64(settings.BackupCount))
				db.DeleteOldBackups(files, sql)
			}
		}
	}
}
