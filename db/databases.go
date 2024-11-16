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
		return errors.New("Ошибка: " + command)
	}
	cfg.Status = true
	return nil
}

// получение размера базы данных по имени.
func (cfg *Database) getDBSize() error {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s -p %d %s -c \"SELECT pg_size_pretty(pg_database_size('%s'))\"", cfg.Password, cfg.Host, cfg.Username, cfg.Port, cfg.Name, cfg.Name)
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return errors.New("Ошибка: " + command)
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
	logrus.Infof("Добавлена база данных %s", cfg.Alias)
	return nil
}

// Удаляет базу данных и все связанные с ней данные и папки.
func (cfg Database) Delete(sql *gorm.DB) error {
	result := sql.Preload("Schedules").First(&cfg, cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	err := sql.Transaction((func(tx *gorm.DB) error {
		if err := tx.Delete(&cfg).Error; err != nil {
			tx.Rollback()
			return err
		}
		for _, schedule := range cfg.Schedules {
			if err := tx.Preload("Backups").First(&schedule, schedule).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&schedule).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := os.RemoveAll(schedule.Directory); err != nil {
				return err
			}
			for _, backup := range schedule.Backups {
				if err := tx.Delete(&backup).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		logrus.Infof("Удалена база данных %s", cfg.Alias)
		return tx.Commit().Error
	}))
	if err != nil {
		logrus.Infof("Ошибка при выполнении транзакции %s", err)
		return err
	}
	return nil
}

// Получение базы данных по имени из таблицы Databases со связанными Backups и Schedules.
func GetDb(sql *gorm.DB, id int) (Database, error) {
	var db Database
	result := sql.
		Preload("Backups", func(db *gorm.DB) *gorm.DB { // сортировка по последней дате
			return db.Order("date desc")
		}).
		Preload("Schedules").
		First(&db, id)
	if result.Error != nil {
		return db, result.Error
	}
	return db, nil
}

// Возвращает список подключенных баз данных.
func GetDbAll(sql *gorm.DB) []Database {
	var DbList []Database
	result := sql.Find(&DbList)
	if result.Error != nil {
		logrus.Error(result.Error)
	}
	return DbList
}
