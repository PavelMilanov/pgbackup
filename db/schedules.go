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

// Обновление расписания в таблице Schedule.
func ScheduleUpdate(sql *gorm.DB, db Schedule) (Schedule, error) {
	var model Schedule
	result := sql.First(&model)
	if result.Error != nil {
		return model, result.Error
	}
	model.Frequency = db.Frequency
	model.Time = db.Time
	sql.Save(&model)
	return model, nil
}

// Сохраняет расписание
func (cfg *Schedule) Save(sql *gorm.DB, timer *cron.Cron) error {
	dbModel, err := GetDb(sql, cfg.DatabaseID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	cfg.DatabaseName = dbModel.Name
	cfg.Status = "активно"
	dir := generateRandomBackupDir()
	cfg.Directory = dir
	scheduleId, err := ScheduleCreate(sql, *cfg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	cronTime := toCron(cfg.Time, cfg.Frequency)
	entryID, _ := timer.AddFunc(cronTime, func() {
		backup := Backup{
			Directory:  cfg.Directory,
			ScheduleID: scheduleId,
		}
		if err := backup.CreateSQL(dbModel); err != nil {
			logrus.Error(err)
		}
		if err := backup.Save(sql); err != nil {
			logrus.Error(err)
		}
	})
	entry := timer.Entry(entryID)
	logrus.Infof("Добавлен планировщик %v", entry)
	logrus.Infof("Добавлено расписание %v", cfg)
	return nil
}

func (cfg *Schedule) SaveManual(sql *gorm.DB) error {
	dbModel, err := GetDb(sql, cfg.DatabaseID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	schedule := Schedule{
		Directory:  config.DEFAULT_BACKUP_DIR,
		DatabaseID: dbModel.ID,
	}
	scheduleId, err := ScheduleCreate(sql, schedule)
	if err != nil {
		logrus.Error(err)
		return err
	}
	fmt.Println(scheduleId)
	// dbCfg := DBConfig{
	// 	Name:     dbModel.Name,
	// 	Host:     dbModel.Host,
	// 	Port:     fmt.Sprintf("%d", dbModel.Port),
	// 	User:     dbModel.Username,
	// 	Password: dbModel.Password,
	// }
	// backup := BackupConfig{
	// 	Directory:  schedule.Directory,
	// 	ScheduleID: scheduleId,
	// }
	// backup.CreateSQL(dbCfg)
	// backup.Save(sql)
	return nil
}

func (cfg *Schedule) Delete(sql *gorm.DB, timer *cron.Cron) error {
	result := sql.Delete(&cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	logrus.Infof("Удалено расписание %v", cfg)
	return nil
}

// Возвращает список конфигураций расписаний
func GetScheduleAll(sql *gorm.DB) []Schedule {
	var scheduleList []Schedule
	result := sql.Find(&scheduleList)
	if result.Error != nil {
		logrus.Error(result.Error)
		return scheduleList
	}
	// for _, item := range scheduleList {
	// 	id := fmt.Sprintf("%d", item.DatabaseID)
	// 	db, err := GetDb(sql, id)
	// 	if err != nil {
	// 		logrus.Error(err)
	// 		return scheduleList
	// 	}
	// 	config := &Schedule{
	// 		ID:         fmt.Sprintf("%d", item.ID),
	// 		DbID:       id,
	// 		DbName:     db.Name,
	// 		Directory:  item.Directory,
	// 		Frequency:  item.Frequency,
	// 		Time:       item.Time,
	// 		LastBackup: "",
	// 		Status:     "",
	// 	}
	// 	scheduleList = append(scheduleList, *config)
	// }
	return scheduleList
}
