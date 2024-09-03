package db

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var BACKUP_DIR = "dumps"
var BACKUPDATA_DIR = "data"
var BACKUP_RUN = []string{"вручную", "по расписанию"}

// Модель бекапа.
type Backup struct {
	Alias    string
	Date     string
	Size     string
	LeadTime string
	Status   string
	Schedule BackupSchedule
}

// Модель для расписания бекапа в формате cron.
type BackupSchedule struct {
	Run   string
	Count string
	Time  string
	Cron  string
}

// Модель Базы данных psql.
type PsqlBase struct {
	Name string
	Size string
}

// Проверка подключения к базе данных
func CheckConnection(cfg Config) string {
	command := fmt.Sprintf("pg_isready -h %s -U %s -d %s -p %d", cfg.Host, cfg.User, cfg.DBName, cfg.portToInt(cfg.Port))
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	return string(cmd)
}

// Выполнение задания бекапа базы данных
// cfg - данные для подключения к PostgreSQL.
// db - имя базы данных, которой сделать бекап - указывается пользователем.
// []string{время выполения, размер бекапа}
func CreateBackup(cfg Config, dbname, backupRun, backupCount, backupTime, backupCron string) (Backup, error) {
	start := time.Now()
	currTime := start.Format("2006-01-02") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	backupName := dbname + "-" + currTime
	curBackups := checkBackup(BACKUP_DIR)
	for _, item := range curBackups {
		if strings.Contains(item, backupName) {
			errStr := fmt.Sprintf("Бекап %s уже существует!", backupName)
			return Backup{}, fmt.Errorf(errStr)
		}
	}
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s/%s.dump", cfg.Password, cfg.Host, cfg.User, dbname, BACKUP_DIR, backupName)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Println(err)
		return Backup{}, err
	}
	log.Println(cmd)
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	size := getBackupSize(BACKUP_DIR, backupName)

	model := Backup{
		Alias:    dbname,
		Date:     currTime,
		Size:     size,
		LeadTime: elapsed,
		Status:   "завершен",
		Schedule: BackupSchedule{
			Run:   backupRun,
			Count: backupCount,
			Time:  backupTime,
			Cron:  backupCron,
		},
	}
	сreateBackupData(&model, BACKUPDATA_DIR)
	return model, nil
}

// Выполение задания восстановления базы данных
func Restore(cfg Config, dbBackup string) {
	command := fmt.Sprintf("psql %s < %s/%s", cfg.User, BACKUP_DIR, dbBackup)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	log.Println(cmd)
}
