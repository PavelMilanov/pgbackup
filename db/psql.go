package db

import (
	"fmt"
	"os/exec"
	"time"
)

var BACKUP_DIR = "/backups"

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
func Backup(cfg Config) {
	currTime := time.Now().Format("2006-01-02") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	command := fmt.Sprintf("pg_dump -U %s %s > %s/%s.dump", cfg.User, cfg.DBName, BACKUP_DIR, cfg.DBName+"-"+currTime)
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
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
