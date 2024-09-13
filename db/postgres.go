package db

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (cfg *Config) portToInt(port string) int {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}
	return intPort
}

func NewPostgreDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.portToInt(cfg.Port))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("база данных подключена")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *gorm.DB) {
	t, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("база данных отключена")
	t.Close()
}
