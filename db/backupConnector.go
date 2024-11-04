package db

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Основная модель бэкапа в приложении.
// Дополнение к модели Backup из служебной БД.
type BackupConfig struct {
	ID         uint
	ScheduleID uint
	Alias      string
	Date       string
	Size       string
	LeadTime   string
	Status     string
	Directory  string
	Dump       string
}

func (bk *BackupConfig) CreateSQL(cfg Database) error {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	dumpName := bk.Directory + "/" + currTime + ".dump"
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s", cfg.Password, cfg.Host, cfg.Username, cfg.Name, dumpName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
		return fmt.Errorf("%s", command)
	}
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	size := bk.getBackupSize(currTime)

	bk.Date = currTime
	bk.Size = size
	bk.LeadTime = elapsed
	bk.Status = "завершен"
	bk.Dump = dumpName
	return nil
}

func (bk *BackupConfig) Save(sql *gorm.DB) error {
	backup := Backup{
		Date:       bk.Date,
		Size:       bk.Size,
		LeadTime:   bk.LeadTime,
		Status:     bk.Status,
		Dump:       bk.Dump,
		ScheduleID: bk.ScheduleID,
	}
	if err := backup.Create(sql); err != nil {
		return err
	}
	return nil
}

func (bk *BackupConfig) getBackupSize(filename string) string {
	command := fmt.Sprintf("du -h %s/%s.dump | awk '{print $1}'", bk.Directory, filename)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
	}
	return string(cmd)
}
