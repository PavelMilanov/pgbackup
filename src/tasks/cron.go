package tasks

import (
	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func (s Sheduler) initBackupsTasks(conn *db.SQLite, timer *cron.Cron) {
	schedules := db.GetSchedulesAll(conn.Sql)
	for _, schedule := range schedules {
		if schedule.Status == config.SCHEDULE_STATUS["активно"] {
			dbModel, _ := db.GetDb(conn.Sql, schedule.DatabaseID)
			cronTime := system.ToCron(schedule.Time, schedule.Frequency)
			timer.AddFunc(cronTime, func() {
				backup := db.Backup{
					Directory:  schedule.Directory,
					ScheduleID: schedule.ID,
					DatabaseID: schedule.DatabaseID,
				}
				backup.Save(dbModel, conn)
			})
		}
	}
	//entris := timer.Entries()
	logrus.Debug("Планировщик дампов запущен")
}
