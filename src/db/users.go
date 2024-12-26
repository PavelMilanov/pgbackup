package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Регистрация нового пользователя
func (cfg *User) Save(sql *gorm.DB) error {
	result := sql.Create(&cfg)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	logrus.Infof("Зарегистрирован новый пользователь %s", cfg.Username)
	return nil
}

// Проверка, зарегистрирован ли пользователь
func (cfg *User) IsRegister(sql *gorm.DB) bool {
	result := sql.Select("ID").Where("Username = ? AND Password = ?", cfg.Username, cfg.Password).First(&cfg)
	if result.Error != nil || result.RowsAffected == 0 {
		logrus.Error(result.Error)
		return false
	}
	return true
}

// Возвращает модель токена пользователя.
func (cfg *User) GetToken(sql *gorm.DB) Token {
	result := sql.Preload("Token").Where("Username = ? AND Password = ?", cfg.Username, cfg.Password).First(&cfg)
	if result.Error != nil {
		logrus.Debug(gorm.ErrRecordNotFound)
		return cfg.Token
	}
	return cfg.Token
}
