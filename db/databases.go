package db

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/PavelMilanov/pgbackup/system"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Проверка подключения к базе данных
func (cfg *Database) CheckConnection() bool {
	command := fmt.Sprintf("pg_isready -h %s -U %s -d %s -p %d", cfg.Host, cfg.Username, cfg.Name, cfg.Port)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error("Ошибка: " + command)
		return false
	}
	return true
}

// получение размера базы данных по имени.
func (cfg *Database) GetDBSize() string {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s -p %d %s -c \"SELECT pg_size_pretty(pg_database_size('%s'))\"", cfg.Password, cfg.Host, cfg.Username, cfg.Port, cfg.Name, cfg.Name)
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error("Ошибка: " + command)
		return ""
	}
	//pg_size_pretty
	//----------------
	//7453 kB
	//(1 row)
	startIndex := 35
	endIndex := len(string(output)) - 10
	size := fmt.Sprint(string(output)[startIndex:endIndex]) // -> 7453 kB
	return size
}

// Добавляет данные о базе данных в служебную БД.
// Перед добавлением в таблицу проверяется подключение.
func (cfg *Database) Save(sql *gorm.DB) error {
	status := cfg.CheckConnection()
	size := cfg.GetDBSize()
	if !status || size == "" {
		return errors.New("Не удалось подключиться к базе данных " + cfg.Alias)
	}
	cfg.Status = status
	cfg.Size = size
	encryptedUsername := system.Encrypt(cfg.Username)
	cfg.Username = encryptedUsername
	encryptedPassword := system.Encrypt(cfg.Password)
	cfg.Password = encryptedPassword
	result := sql.Create(&cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return errors.New("Ошибка при сохранении базы данных: " + cfg.Alias)
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
			return db.Order("created_at desc")
		}).
		Preload("Schedules").
		First(&db, id)
	if result.Error != nil {
		return db, result.Error
	}
	descriptedUsername := system.Decrypt(db.Username)
	db.Username = descriptedUsername
	descriptedPassword := system.Decrypt(db.Password)
	db.Password = descriptedPassword
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
