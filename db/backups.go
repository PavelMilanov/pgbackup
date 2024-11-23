package db

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Создает файл бекапа базы данных и сохраняет методанных в служебную БД.
func (bk *Backup) Save(cfg Database, sql *gorm.DB) {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	dumpName := currTime + ".dump"
	bk.Date = start.Format(time.RFC1123)
	bk.Dump = dumpName
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s -p %d %s > %s", cfg.Password, cfg.Host, cfg.Username, cfg.Port, cfg.Name, bk.Directory+"/"+dumpName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		errCommand := fmt.Sprintf("touch %s", bk.Directory+"/"+dumpName)
		exec.Command("sh", "-c", errCommand).Output()
		timer := time.Since(start).Seconds()
		elapsed := fmt.Sprintf("%.3f сек", timer)
		bk.LeadTime = elapsed
		bk.getBackupSize(dumpName)
		bk.Status = false
		result := sql.Create(&bk)
		if result.Error != nil {
			logrus.Error(result.Error)
		}
		logrus.Warnf("Ошибка при создании дампа %s", bk.Directory+"/"+bk.Dump)
		return
	}
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	bk.LeadTime = elapsed
	bk.getBackupSize(dumpName)
	bk.Status = true
	result := sql.Create(&bk)
	if result.Error != nil {
		logrus.Error(result.Error)
		os.Remove(bk.Directory + "/" + bk.Dump)
		return
	}
	logrus.Infof("Создан дамп %s", bk.Directory+"/"+bk.Dump)
}

// Удаление файла бекапа и его метаданных.
func (bk *Backup) Delete(sql *gorm.DB) {
	err := sql.Transaction((func(tx *gorm.DB) error {
		tx.First(&bk)
		tx.Delete(&bk)
		return nil
	}))
	if err != nil {
		logrus.Errorf("Ошибка при удалении бекапа %s %v", bk.Dump, err)
	} else {
		os.Remove(bk.Dump)
		logrus.Infof("Удален бекап %s", bk.Directory+"/"+bk.Dump)
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

func DeleteOldBackups(files []string, sql *gorm.DB) {
	var model Backup
	for _, file := range files {
		result := sql.Raw("DELETE FROM backups WHERE dump = ?", file).Scan(&model)
		if result.Error != nil {
			logrus.Error(result.Error)
			continue
		}
		os.Remove(model.Directory + "/" + model.Dump)
	}
	logrus.Infof("Удалено %d бекапа", len(files))
}
