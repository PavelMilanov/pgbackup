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
	case config.BACKUP_FREQUENCY["ежедневно"]:
		cron = "1"
	case config.BACKUP_FREQUENCY["eженедельно"]:
		cron = "7"
	}
	formatTime := fmt.Sprintf("%s %s */%s * *", crontime[1], crontime[0], cron)
	return formatTime
}

// Генерирует случайную строку из цифр от 0 до 10000 и создает директорию.
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
	// Для статических бекапов
	// Если нет папки для статического бекапа - создать и сделать бекап
	if cfg.Time == "" {
		for _, item := range dbModel.Schedules {
			// если есть расписание без времени - создаем статические бекапы здесь
			if item.Time == cfg.Time {
				backup := Backup{
					Directory:    item.Directory,
					DatabaseName: dbModel.Name,
					ScheduleID:   item.ID,
					DatabaseID:   dbModel.ID,
				}
				backup.Save(dbModel, sql)
				return nil
			}
		}
		cfg.DatabaseName = dbModel.Name
		dir := generateRandomBackupDir()
		cfg.Directory = dir
		cfg.Status = config.SCHEDULE_STATUS["вручную"]
		scheduleId, err := ScheduleCreate(sql, *cfg)
		if err != nil {
			logrus.Error(err)
			return err
		}
		backup := Backup{
			Directory:    cfg.Directory,
			DatabaseName: dbModel.Name,
			ScheduleID:   scheduleId,
			DatabaseID:   dbModel.ID,
		}
		backup.Save(dbModel, sql)
		// для бекапов по расписанию
	} else {
		cfg.DatabaseName = dbModel.Name
		dir := generateRandomBackupDir()
		cfg.Directory = dir
		cfg.Status = config.SCHEDULE_STATUS["активно"]
		scheduleId, err := ScheduleCreate(sql, *cfg)
		if err != nil {
			logrus.Error(err)
			return err
		}
		cronTime := toCron(cfg.Time, cfg.Frequency)
		entryID, _ := timer.AddFunc(cronTime, func() {
			backup := Backup{
				Directory:    cfg.Directory,
				DatabaseName: dbModel.Name,
				ScheduleID:   scheduleId,
				DatabaseID:   dbModel.ID,
			}
			backup.Save(dbModel, sql)
		})
		entry := timer.Entry(entryID)
		logrus.Infof("Добавлено расписание %v для бекапа", entry)
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
	err := sql.Transaction((func(tx *gorm.DB) error {
		tx.Delete(cfg)
		for _, backup := range cfg.Backups {
			tx.Delete(&backup)
		}
		return nil
	}))
	if err != nil {
		logrus.Infof("Ошибка при удалении расписания %s", cfg.DatabaseName)
		return err
	} else {
		os.RemoveAll(cfg.Directory)
		logrus.Infof("Удалено расписание %s", cfg.DatabaseName)
		return nil

	}
}

// Возвращает список конфигураций расписаний, которые запускаются по расписанию
func GetSchedulesAll(sql *gorm.DB) []Schedule {
	var scheduleList []Schedule
	result := sql.Preload("Backups", func(db *gorm.DB) *gorm.DB { // выбираем последений по дате бекап для отображения
		return db.Order("date desc").Limit(1)
	}).Where("Status = ?", config.SCHEDULE_STATUS["активно"]).Find(&scheduleList)
	if result.Error != nil {
		logrus.Error(result.Error)
		return scheduleList
	}
	return scheduleList
}
