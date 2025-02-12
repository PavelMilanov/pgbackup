package db

import (
	"sync"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLite struct {
	Sql   *gorm.DB
	Mutex *sync.Mutex
}

func NewDatabase(sql string, timer *cron.Cron) SQLite {
	conn, err := gorm.Open(sqlite.Open(sql), &gorm.Config{PrepareStmt: true})
	if err != nil {
		logrus.Fatal("Ошибка при подключении к базе данных")
	}
	var mutex sync.Mutex
	db := SQLite{Sql: conn, Mutex: &mutex}

	logrus.Info("Соединение с базой данных установлено")
	automigrate(db.Sql)
	setDefaultSettings(db.Sql)
	initBackupsTasks(&db, timer)
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

func initBackupsTasks(conn *SQLite, timer *cron.Cron) {
	schedules := GetSchedulesAll(conn.Sql)
	for _, schedule := range schedules {
		if schedule.Status == config.SCHEDULE_STATUS["активно"] {
			dbModel, _ := GetDb(conn.Sql, schedule.DatabaseID)
			cronTime := system.ToCron(schedule.Time, schedule.Frequency)
			timer.AddFunc(cronTime, func() {
				backup := Backup{
					Directory:  schedule.Directory,
					ScheduleID: schedule.ID,
					DatabaseID: schedule.DatabaseID,
				}
				backup.Save(dbModel, conn)
			})
		}
	}
	logrus.Debug("Планировщик дампов запущен")
}

// Задача удаления старых дампов и актуализации размера баз.
func StartSystemTasks(conn *SQLite) {
	settings, _ := GetSettings(conn.Sql)
	files := system.ParseOldFiles(float64(settings.BackupCount))
	deleteOldBackups(files, conn)
	DB := GetDbAll(conn.Sql)
	for _, item := range DB {
		item.Username = system.Decrypt(item.Username)
		item.Password = system.Decrypt(item.Password)
		size := item.getDBSize()
		item.Size = size
		conn.Mutex.Lock()
		conn.Sql.Save(&item)
		conn.Mutex.Unlock()
	}
}

// Задача проверки соединения с базами.
func CheckDBconnection(conn *SQLite) {
	DB := GetDbAll(conn.Sql)
	for _, item := range DB {
		status := item.checkConnection()
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
			logrus.Infof("Восстановлено соединение с базой данных %s", item.Alias)
		}
	}

}
