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
		logrus.Error("Ошибка при подключении к базе данных")
	}
	automigrate(db)
	initCronTasks(db, timer)
	return db
}

func automigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Database{}, &Schedule{}, &Backup{}, &Token{}); err != nil {
		logrus.Errorf("%s", err)
		return
	}
}

func initCronTasks(sql *gorm.DB, timer *cron.Cron) {
	schedules := GetSchedulesAll(sql)
	for _, schedule := range schedules {
		if schedule.Status == config.SCHEDULE_STATUS[0] {
			dbModel, _ := GetDb(sql, schedule.DatabaseID)
			cronTime := toCron(schedule.Time, schedule.Frequency)
			timer.AddFunc(cronTime, func() {
				backup := Backup{
					Directory:  schedule.Directory,
					ScheduleID: schedule.ID,
					DatabaseID: schedule.DatabaseID,
				}
				backup.Save(dbModel, sql)
			})
		}
	}
	entris := timer.Entries()
	logrus.Infof("Фоновые задачи для бекапов инициализированы %v", entris)
}
