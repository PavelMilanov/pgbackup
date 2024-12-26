package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetSettings(sql *gorm.DB) (Setting, error) {
	var cfg Setting
	result := sql.First(&cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return cfg, result.Error
	}
	return cfg, nil
}

func (cfg *Setting) Update(sql *gorm.DB) error {
	result := sql.Raw("UPDATE settings SET backup_count = ?", cfg.BackupCount).Scan(&cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	logrus.Debugf("Настройки приложения изменены: %+v", cfg)
	return nil
}
