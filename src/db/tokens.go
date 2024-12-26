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
func generateToken(userID int) string {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.TOKEN_EXPIRED_TIME) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   strconv.Itoa(userID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWT_KEY)
	if err != nil {
		logrus.Error(err)
		return err.Error()
	}
	return tokenString
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
	t.Hash = generateToken(t.UserID)
	result := sql.Create(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (t *Token) Refresh(sql *gorm.DB) error {
	newHash := generateToken(t.UserID)
	result := sql.Model(&t).Where("Hash = ?", t.Hash).Update("Hash", newHash)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	logrus.Debug("Обновлен токен авторизации")
	return nil
}

func GetToken(sql *gorm.DB, user int) Token {
	var t Token
	result := sql.Raw("Select hash,user_id FROM tokens WHERE user_id = ?", user).Scan(&t)
	if result.Error != nil {
		logrus.Error(result.Error)
		return t
	}
	return t
}
