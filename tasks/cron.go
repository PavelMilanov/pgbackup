package tasks

import (
	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitBackupsTasks(sql *gorm.DB, timer *cron.Cron) {
	schedules := db.GetSchedulesAll(sql)
	for _, schedule := range schedules {
		if schedule.Status == config.SCHEDULE_STATUS["активно"] {
			dbModel, _ := db.GetDb(sql, schedule.DatabaseID)
			cronTime := system.ToCron(schedule.Time, schedule.Frequency)
			timer.AddFunc(cronTime, func() {
				backup := db.Backup{
					Directory:  schedule.Directory,
					ScheduleID: schedule.ID,
					DatabaseID: schedule.DatabaseID,
				}
				backup.Save(dbModel, sql)
			})
		}
	}
	//entris := timer.Entries()
	logrus.Debug("Планировщик дампов запущен")
}
