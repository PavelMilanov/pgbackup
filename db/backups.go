package db

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (bk *Backup) Save(cfg Database, sql *gorm.DB) error {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	dumpName := bk.Directory + "/" + currTime + ".dump"
	bk.Date = currTime
	bk.Dump = dumpName
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s -p %d %s > %s", cfg.Password, cfg.Host, cfg.Username, cfg.Port, cfg.Name, dumpName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		bk.Status = "ошибка"
		return errors.New(command)
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
	defer sql.Create(&bk)
	logrus.Infof("Создан бекап %s", bk.Dump)
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
