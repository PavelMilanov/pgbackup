package main

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
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

// Получение списка всех баз данных в экземпляре PostgreSQL.
func (conn *Database) getDBs() []string {
	conStr := fmt.Sprintf("user=%s password=%s host=localhost port=%d dbname=%s sslmode=disable", conn.User, conn.Password, conn.Port, conn.Dbname)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT datname FROM pg_database;")
	if err != nil {
		panic(err)
	}
	dbList := []string{}
	for rows.Next() {
		var datName string
		if err := rows.Scan(&datName); err != nil {
			panic(err.Error())
		}
		if strings.HasPrefix(datName, "template") {
			continue
		}
		dbList = append(dbList, datName)
	}
	defer rows.Close()
	return dbList
}

func (conn *Database) checkConnection() {
	command := fmt.Sprintf("pg_isready -U %s -d %s -p %d", conn.User, conn.Dbname, conn.Port)
	cmd, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd)
	// return true
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
