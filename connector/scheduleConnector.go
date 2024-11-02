package connector

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Основная модель расписание в приложении.
// Дополнение к модели Schedule из служебной БД.
type ScheduleConfig struct {
	ID         string
	DbID       string
	DbName     string
	Directory  string
	Frequency  string
	Time       string
	LastBackup string
	Status     string
}

// Возвращает строку в формате cron для модели Task.
func toCron(time, frequency string) string {
	// минуты часы день(*/1 каждый день) * *
	crontime := strings.Split(time, ":") // 22:45 => ["22", "45"]
	var cron string
	switch frequency {
	case BACKUP_FREQUENCY[0]:
		cron = "1"
	case BACKUP_FREQUENCY[1]:
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
	if err := os.Mkdir(BACKUP_DIR+"/"+dir, 0755); err != nil {
		if !os.IsExist(err) {
			logrus.Infof("%s - директория создана", dir)
		}
	}
	return BACKUP_DIR + "/" + dir
}

// // Удаляет директорию с бекапами
// func DeleteBackupDir(dir string) {
// 	deleteBackupsInDir(dir)
// 	os.RemoveAll(dir)
// 	logrus.Infof("%s - директория удалена", dir)
// }

// Сохраняет расписание
func (cfg *ScheduleConfig) Save(sql *gorm.DB, timer *cron.Cron) error {
	dbModel, err := db.GetDb(sql, cfg.DbName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	dir := generateRandomBackupDir()
	schedule := db.Schedule{
		Directory:  dir,
		Time:       cfg.Time,
		Frequency:  cfg.Frequency,
		DatabaseID: dbModel.ID,
	}
	scheduleId, err := db.ScheduleCreate(sql, schedule)
	if err != nil {
		logrus.Error(err)
		return err
	}
	cronTime := toCron(cfg.Time, cfg.Frequency)
	entryID, _ := timer.AddFunc(cronTime, func() {
		dbCfg := DBConfig{
			Name:     dbModel.Name,
			Host:     dbModel.Host,
			Port:     fmt.Sprintf("%d", dbModel.Port),
			User:     dbModel.Username,
			Password: dbModel.Password,
		}
		backup := BackupConfig{
			Directory:  schedule.Directory,
			ScheduleID: scheduleId,
		}
		backup.CreateSQL(dbCfg)
		backup.Save(sql)
	})
	entry := timer.Entry(entryID)
	logrus.Infof("Добавлен планировщик %v", entry)
	cfg.ID = fmt.Sprintf("%d", scheduleId)
	logrus.Infof("Добавлено расписание %v", cfg)
	return nil
}

func (cfg *ScheduleConfig) Change(sql *gorm.DB, timer *cron.Cron) error {
	id, _ := strconv.ParseUint(cfg.ID, 10, 32)
	schedule := db.Schedule{
		ID:        uint(id),
		Time:      cfg.Time,
		Frequency: cfg.Frequency,
	}
	result, err := db.ScheduleUpdate(sql, schedule)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("Изменено расписание %v", result)
	return nil
}

// Возвращает список конфигураций расписаний
func GetScheduleAll(sql *gorm.DB) []ScheduleConfig {
	var scheduleList []ScheduleConfig
	schedules, err := db.GetScheduleAll(sql)
	if err != nil {
		logrus.Error(err)
		return scheduleList
	}
	for _, item := range schedules {
		id := fmt.Sprintf("%d", item.DatabaseID)
		db, err := db.GetDb(sql, id)
		if err != nil {
			logrus.Error(err)
			return scheduleList
		}
		config := &ScheduleConfig{
			ID:         fmt.Sprintf("%d", item.ID),
			DbID:       id,
			DbName:     db.Name,
			Directory:  item.Directory,
			Frequency:  item.Frequency,
			Time:       item.Time,
			LastBackup: "",
			Status:     "",
		}
		scheduleList = append(scheduleList, *config)
	}
	return scheduleList
}
