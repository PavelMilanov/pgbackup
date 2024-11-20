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
	logrus.Info("Соединение с базой данных установлено")
	automigrate(db)
	setDefaultSettings(db)
	initCronTasks(db, timer)
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
	if err := db.FirstOrCreate(&settings, Setting{BackupDays: config.DEFAULT_BACKUP_EXPIRED_DAYS}).Error; err != nil {
		logrus.Fatal("Ошибка при создании или получении настроек")
	}
}

func initCronTasks(sql *gorm.DB, timer *cron.Cron) {
	schedules := GetSchedulesAll(sql)
	for _, schedule := range schedules {
		if schedule.Status == config.SCHEDULE_STATUS["активно"] {
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
	logrus.Infof("Фоновые задачи для бекапов инициализированы %+v", entris)
}
