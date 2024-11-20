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
