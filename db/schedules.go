package db

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Возвращает строку в формате cron для модели Task.
func toCron(time, frequency string) string {
	// минуты часы день(*/1 каждый день) * *
	crontime := strings.Split(time, ":") // 22:45 => ["22", "45"]
	var cron string
	switch frequency {
	case config.BACKUP_FREQUENCY[0]:
		cron = "1"
	case config.BACKUP_FREQUENCY[1]:
		cron = "7"
	}
	formatTime := fmt.Sprintf("%s %s */%s * *", crontime[1], crontime[0], cron)
	return formatTime
}

// Генерирует случайную строку из цифр от 0 до 10000.
func generateRandomBackupDir() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	dir := strconv.Itoa(r1.Intn(10000))
	if err := os.Mkdir(config.BACKUP_DIR+"/"+dir, 0755); err != nil {
		if !os.IsExist(err) {
			logrus.Infof("%s - директория создана", dir)
		}
	}
	return config.BACKUP_DIR + "/" + dir
}

// Создание расписания в таблице Schedule.
func ScheduleCreate(sql *gorm.DB, db Schedule) (int, error) {
	result := sql.Create(&db)
	if result.Error != nil {
		return 0, result.Error
	}
	return db.ID, nil
}

// Сохраняет расписание
func (cfg *Schedule) Save(sql *gorm.DB, timer *cron.Cron) error {
	dbModel, err := GetDb(sql, cfg.DatabaseID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	cfg.DatabaseName = dbModel.Name
	dir := generateRandomBackupDir()
	cfg.Directory = dir
	scheduleId, err := ScheduleCreate(sql, *cfg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if cfg.Time != "" {
		cfg.Status = config.SCHEDULE_STATUS[0]
		cronTime := toCron(cfg.Time, cfg.Frequency)
		entryID, _ := timer.AddFunc(cronTime, func() {
			backup := Backup{
				Directory:  cfg.Directory,
				ScheduleID: scheduleId,
			}
			if err := backup.Save(dbModel, sql); err != nil {
				logrus.Error(err)
			}
		})
		entry := timer.Entry(entryID)
		logrus.Infof("Добавлен планировщик %v", entry)
		logrus.Infof("Добавлено расписание %s", cfg.DatabaseName)
	} else {
		for _, item := range dbModel.Schedules {
			// если есть расписание без времени - создаем статические бекапы здесь
			if item.Time == "" {
				backup := Backup{
					Directory:  item.Directory,
					ScheduleID: item.ID,
				}
				if err := backup.Save(dbModel, sql); err != nil {
					logrus.Error(err)
				}
				break
			}
		}
		backup := Backup{
			Directory:  cfg.Directory,
			ScheduleID: scheduleId,
		}
		if err := backup.Save(dbModel, sql); err != nil {
			logrus.Error(err)
		}
	}
	return nil
}

// Удаляет расписание и все связанные с ним данные и папки.
func (cfg *Schedule) Delete(sql *gorm.DB, timer *cron.Cron) error {
	result := sql.Preload("Backups").First(&cfg, cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	tx := sql.Begin()
	tx.Delete(cfg)
	if err := os.Remove(cfg.Directory); err != nil {
		tx.Rollback()
	}
	for _, backup := range cfg.Backups {
		sql.Delete(&backup)
	}
	tx.Commit()
	logrus.Infof("Удалено расписание %s", cfg.DatabaseName)
	return nil
}

// Возвращает список конфигураций расписаний
func GetSchedules(sql *gorm.DB) []Schedule {
	var scheduleList []Schedule
	result := sql.Where("Status = ?", config.SCHEDULE_STATUS[0]).Find(&scheduleList)
	if result.Error != nil {
		logrus.Error(result.Error)
		return scheduleList
	}
	return scheduleList
}
