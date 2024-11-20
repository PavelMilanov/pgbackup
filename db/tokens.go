package db

import (
	"strconv"
	"time"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Генерирует токен доступа при успешной авторизации.
func (t *Token) Generate() error {

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.TOKEN_EXPIRED_TIME) * time.Hour)),
		Subject:   strconv.Itoa(t.UserID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWT_KEY)
	if err != nil {
		logrus.Error(err)
		return err
	}
	t.Hash = tokenString
	return nil
}

// Валидирует токен аутентификации.
func (t *Token) Validate() bool {
	token, err := jwt.ParseWithClaims(t.Hash, &t.RegisteredClaims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})
	if err != nil {
		logrus.Error(err)
		return false
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		now := time.Now()
		difference := claims.ExpiresAt.Time.Sub(now)
		if difference < 0 {
			return false
		}
	} else {
		logrus.Debug("Токен не валиден")
		return false
	}
	return true
}

func (t *Token) Save(sql *gorm.DB) error {
	result := sql.Create(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (t *Token) Delete(sql *gorm.DB) error {
	result := sql.Where("Hash = ?", t.Hash).Delete(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func GetToken(sql *gorm.DB, data string) Token {
	var t Token
	result := sql.Select("Hash", "UserID").Where("Hash = ?", data).First(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return t
	}
	return t
}
