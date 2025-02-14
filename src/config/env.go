package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Env struct {
	URL    string `mapstructure:"URL"`
	JwtKey string `mapstructure:"JWT_KEY"`
	AesKey string `mapstructure:"AES_KEY"`
	TZ     string `mapstructure:"TZ"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("не найден файл .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		logrus.Fatal("не загружен файл: ", err)
	}
	return &env
}
