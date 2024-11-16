package db

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Генерирует токен доступа при успешной авторизации.
func (t *Token) Generate() error {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // вернуть на 72 часа
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logrus.Error(err)
		return err
	}
	t.Hash = tokenString
	return nil
}

// Валидирует токен аутентификации.
func (t *Token) Validate() bool {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(t.Hash, &t.RegisteredClaims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
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
		logrus.Error("Invalid token")
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
	result := sql.Delete(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func GetToken(sql *gorm.DB) Token {
	var t Token
	result := sql.First(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return t
	}
	return t
}
