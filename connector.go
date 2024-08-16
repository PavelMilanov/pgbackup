package main

import (
	"fmt"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
)

var BACKUP_DIR = "/backups"

// Модель для подключения к базе данных
type Database struct {
	User     string
	Password string
	Dbname   string
	Port     int
}

// Вывод все информации о базах данных через psql
func (conn *Database) getDbInfo() {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s -d %s -c '\\l+'", conn.Password, "localhost", conn.User, conn.Dbname)
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
}

// Выполнение задания бекапа базы данных
func (conn *Database) backup() {
	currTime := time.Now().Format("2006-01-02") // шаблон GO для формата ГГГГ-мм-дд
	command := fmt.Sprintf("pg_dump -U %s %s > %s/%s.dump", conn.User, conn.Dbname, BACKUP_DIR, conn.Dbname+"-"+currTime)
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
}

// Выполение задания восстановления базы данных
func (conn *Database) restore(dbBackup string) {
	command := fmt.Sprintf("psql %s < %s/%s", conn.User, BACKUP_DIR, dbBackup)
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
}
