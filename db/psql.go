package db

import (
	"fmt"
	"os/exec"
	"time"
)

var BACKUP_DIR = "dumps"
var BACKUPDATA_DIR = "data"

// Проверка подключения к базе данных
func CheckConnection(cfg Config) {
	command := fmt.Sprintf("pg_isready -U %s -d %s -p %d", cfg.User, cfg.DBName, cfg.portToInt(cfg.Port))
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
	// return true
}

// Выполнение задания бекапа базы данных
// cfg - данные для подключения к PostgreSQL.
// db - имя базы данных, которой сделать бекап - указывается пользователем.
func CreateBackup(cfg Config, db string, backupName string) (string, error) {
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s/%s.dump", cfg.Password, cfg.Host, cfg.User, db, BACKUP_DIR, backupName)
	start := time.Now()
	_, err := exec.Command("bash", "-c", command).Output()
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	if err != nil {
		return "0", err
	}
	return elapsed, nil
}

// Выполение задания восстановления базы данных
func Restore(cfg Config, dbBackup string) {
	command := fmt.Sprintf("psql %s < %s/%s", cfg.User, BACKUP_DIR, dbBackup)
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
}
