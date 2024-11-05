package db

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Проверка подключения к базе данных
func (cfg *Database) checkConnection() error {
	command := fmt.Sprintf("pg_isready -h %s -U %s -d %s -p %d", cfg.Host, cfg.Username, cfg.Name, cfg.Port)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		cfg.Status = false
		return errors.New(command)
	}
	cfg.Status = true
	return nil
}

// получение размера базы данных по имени.
func (cfg *Database) getDBSize() error {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s -p %d %s -c \"SELECT pg_size_pretty(pg_database_size('%s'))\"", cfg.Password, cfg.Host, cfg.Username, cfg.Port, cfg.Name, cfg.Name)
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return errors.New(string(output))
	}
	//pg_size_pretty
	//----------------
	//7453 kB
	//(1 row)
	startIndex := 35
	endIndex := len(string(output)) - 10
	size := fmt.Sprint(string(output)[startIndex:endIndex]) // -> 7453 kB
	cfg.Size = size
	return nil
}

// Добавляет данные о базе данных в служебную БД.
// Перед добавлением в таблицу проверяется подключение.
func (cfg *Database) Save(sql *gorm.DB) error {
	if err := cfg.checkConnection(); err != nil {
		logrus.Error(err)
		return err
	}
	if err := cfg.getDBSize(); err != nil {
		logrus.Error(err)
		return err
	}
	result := sql.Create(&cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	logrus.Infof("Добавлена база данных %v", cfg)
	return nil
}

func (cfg Database) Delete(sql *gorm.DB) error {
	result := sql.Preload("Schedules").First(&cfg, cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	tx := sql.Begin()
	tx.Delete(&cfg)
	for _, schedule := range cfg.Schedules {
		sql.Preload("Backups").First(&schedule, schedule)
		tx.Delete(&schedule)
		if err := os.Remove(schedule.Directory); err != nil {
			tx.Rollback()
		}
		for _, backup := range schedule.Backups {
			sql.Delete(&backup)
		}
	}
	tx.Commit()
	logrus.Infof("Удалена база данных и все связанные данные %v", cfg)
	return nil
}

// Получение базы данных по имени из таблицы Databases
func GetDb(sql *gorm.DB, id int) (Database, error) {
	var db Database
	result := sql.First(&db, "ID = ?", id)
	if result.Error != nil {
		return db, result.Error
	}
	return db, nil
}

// Возвращает список подключенных баз данных.
func GetDbAll(sql *gorm.DB) []Database {
	var DbList []Database
	result := sql.Model(&Database{}).Preload("Schedules").Find(&DbList)
	if result.Error != nil {
		logrus.Error(result.Error)
	}
	return DbList
}
