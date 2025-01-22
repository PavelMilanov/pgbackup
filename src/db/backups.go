package db

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Создает файл бекапа базы данных и сохраняет методанных в служебную БД.
func (bk *Backup) Save(cfg Database, conn *SQLite) {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	bk.Dump = currTime + ".dump.gz"
	bk.Date = start.Format(time.RFC1123)
	command := fmt.Sprintf("export PGPASSWORD=\"%s\" && pg_dump -h %s -U %s -p %d %s | gzip > %s", cfg.Password, cfg.Host, cfg.Username, cfg.Port, cfg.Name, filepath.Join(bk.Directory, bk.Dump))
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		errCommand := fmt.Sprintf("touch %s", filepath.Join(bk.Directory, bk.Dump))
		exec.Command("sh", "-c", errCommand).Output()
		timer := time.Since(start).Seconds()
		elapsed := fmt.Sprintf("%.3f сек", timer)
		bk.LeadTime = elapsed
		bk.getBackupSize(bk.Dump)
		bk.Status = false
		conn.Mutex.Lock()
		result := conn.Sql.Create(&bk)
		conn.Mutex.Unlock()
		if result.Error != nil {
			logrus.Error(result.Error)
		}
		logrus.Warnf("Ошибка при создании дампа %s", filepath.Join(bk.Directory, bk.Dump))
		return
	}
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	bk.LeadTime = elapsed
	bk.getBackupSize(bk.Dump)
	bk.Status = true
	conn.Mutex.Lock()
	result := conn.Sql.Create(&bk)
	conn.Mutex.Unlock()
	if result.Error != nil {
		logrus.Error(result.Error)
		os.Remove(filepath.Join(bk.Directory, bk.Dump))
		return
	}
	logrus.Infof("Создан дамп %s", filepath.Join(bk.Directory, bk.Dump))
}

func (bk *Backup) Restore(cfg Database) error {
	command := fmt.Sprintf("export PGPASSWORD=\"%s\" && gunzip -c %s pg_dump -h %s -U %s -p %d %s", cfg.Password, filepath.Join(bk.Directory, bk.Dump), cfg.Host, cfg.Username, cfg.Port, cfg.Name)
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return errors.New(string(out))
	}
	fmt.Println(string(out))
	return nil
}

// Удаление файла бекапа и его метаданных.
func (bk *Backup) Delete(conn *SQLite) {
	conn.Mutex.Lock()
	err := conn.Sql.Transaction((func(tx *gorm.DB) error {
		tx.First(&bk)
		tx.Delete(&bk)
		return nil
	}))
	conn.Mutex.Unlock()
	if err != nil {
		logrus.Errorf("Ошибка при удалении бекапа %s %v", bk.Dump, err)
	} else {
		os.Remove(bk.Dump)
		logrus.Infof("Удален бекап %s", filepath.Join(bk.Directory, bk.Dump))
	}
}

// Получение модели бекапа.
func (bk *Backup) Get(sql *gorm.DB) {
	sql.Find(&bk)
}

// Получение до 5 моделей бекапов с зависимостями.
func GetLastBackups(sql *gorm.DB) []Backup {
	var bkList []Backup
	sql.Raw("SELECT * FROM backups ORDER BY date DESC LIMIT 5").Scan(&bkList)
	return bkList
}

// Получение общего числа дампов.
func GetCountBackups(sql *gorm.DB) int64 {
	var count int64
	sql.Raw("SELECT count(*) FROM backups").Scan(&count)
	return count
}

// Получение количества успешных и неуспешных бекапов.
func CountBackupsStatus(sql *gorm.DB) []int64 {
	var success int64
	var failed int64
	var bkList []Backup
	sql.Where("Status == ?", 1).Find(&bkList).Count(&success)
	sql.Where("Status == ?", 0).Find(&bkList).Count(&failed)
	return []int64{success, failed}
}

// Получение размера файла бекапа на диске.
func (bk *Backup) getBackupSize(filename string) {
	command := fmt.Sprintf("du -h %s/%s | awk '{print $1}'", bk.Directory, filename)
	cmd, _ := exec.Command("sh", "-c", command).Output()
	bk.Size = string(cmd)
}

// Удаление дампов из базы данных по условию.
// Используется для планировщик для удалениея старых дампов.
func DeleteOldBackups(files []string, conn *SQLite) {
	var model Backup
	for _, file := range files {
		conn.Mutex.Lock()
		result := conn.Sql.Raw("DELETE FROM backups WHERE dump = ?", file).Scan(&model)
		conn.Mutex.Unlock()
		if result.Error != nil {
			logrus.Error(result.Error)
			continue
		}
		os.Remove(filepath.Join(model.Directory, model.Dump))
	}
	logrus.Infof("Удалено %d бекапа", len(files))
}
