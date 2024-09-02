package db

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

var BACKUP_DIR = "dumps"
var BACKUPDATA_DIR = "data"
var BACKUP_RUN = []string{"в ручную", "расписание"}
var BACKUP_SCHEDULE = []string{"ежедневно", "еженедельно", "ежемесячно"}

type Backup struct {
	Alias    string
	Date     string
	Size     string
	LeadTime string
	Status   string
	Schedule BackupSchedule
}

type BackupSchedule struct {
	Run   string
	Count string
	Time  string
	Cron  string
}

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
func CreateBackup(cfg Config, db string, backupName string) (string, error) {
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s/%s.dump", cfg.Password, cfg.Host, cfg.User, db, BACKUP_DIR, backupName)
	start := time.Now()
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return "0", err
	}
	log.Println(cmd)
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	return elapsed, nil
}

// Выполение задания восстановления базы данных
func Restore(cfg Config, dbBackup string) {
	command := fmt.Sprintf("psql %s < %s/%s", cfg.User, BACKUP_DIR, dbBackup)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
}
