package db

import (
	"os"
	"strconv"
	"time"

	"math/rand"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Генерирует случайную строку из цифр от 0 до 10000 и создает директорию.
func generateRandomBackupDir() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	dir := strconv.Itoa(r1.Intn(10000))
	if err := os.Mkdir(config.BACKUP_DIR+"/"+dir, 0755); err != nil {
		if !os.IsExist(err) {
			logrus.Debugf("%s - директория создана", dir)
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
					Directory:     item.Directory,
					DatabaseAlias: dbModel.Name,
					ScheduleID:    item.ID,
					DatabaseID:    dbModel.ID,
				}
				go backup.Save(dbModel, sql)
				return nil
			}
		}
		cfg.DatabaseAlias = dbModel.Alias
		dir := generateRandomBackupDir()
		cfg.Directory = dir
		cfg.Status = config.SCHEDULE_STATUS["вручную"]
		scheduleId, err := ScheduleCreate(sql, *cfg)
		if err != nil {
			logrus.Error(err)
			return err
		}
		backup := Backup{
			Directory:     cfg.Directory,
			DatabaseAlias: dbModel.Alias,
			ScheduleID:    scheduleId,
			DatabaseID:    dbModel.ID,
		}
		go backup.Save(dbModel, sql)
		// для бекапов по расписанию
	} else {
		cfg.DatabaseAlias = dbModel.Alias
		dir := generateRandomBackupDir()
		cfg.Directory = dir
		cfg.Status = config.SCHEDULE_STATUS["активно"]
		scheduleId, err := ScheduleCreate(sql, *cfg)
		if err != nil {
			logrus.Error(err)
			return err
		}
		cronTime := system.ToCron(cfg.Time, cfg.Frequency)
		entryID, _ := timer.AddFunc(cronTime, func() {
			backup := Backup{
				Directory:     cfg.Directory,
				DatabaseAlias: dbModel.Alias,
				ScheduleID:    scheduleId,
				DatabaseID:    dbModel.ID,
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
		logrus.Infof("Ошибка при удалении расписания %s", cfg.DatabaseAlias)
		return err
	} else {
		os.RemoveAll(cfg.Directory)
		logrus.Infof("Удалено расписание %s", cfg.DatabaseAlias)
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

// Получение общего числа расписаний.
func GetCountSchedules(sql *gorm.DB) int64 {
	var count int64
	sql.Raw("SELECT count(*) FROM schedules WHERE status = ?", config.SCHEDULE_STATUS["активно"]).Scan(&count)
	return count
}
