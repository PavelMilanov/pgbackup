package db

import (
	"fmt"
	"os/exec"
	"time"

	"gorm.io/gorm"
)

func (bk *Backup) CreateSQL(cfg Database) error {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	dumpName := bk.Directory + "/" + currTime + ".dump"
	bk.Date = currTime
	bk.Dump = dumpName
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s", cfg.Password, cfg.Host, cfg.Username, cfg.Name, dumpName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return fmt.Errorf("%s", command)
	}
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	bk.LeadTime = elapsed
	size, err := bk.getBackupSize(currTime)
	if err != nil {
		bk.Status = "ошибка"
		return err
	}
	bk.Size = size
	bk.Status = "завершен"
	return nil
}

func (bk *Backup) Save(sql *gorm.DB) error {
	result := sql.Create(&bk)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (bk *Backup) getBackupSize(filename string) (string, error) {
	command := fmt.Sprintf("du -h %s/%s.dump | awk '{print $1}'", bk.Directory, filename)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return "", fmt.Errorf("%s", command)
	}
	return string(cmd), nil
}
