package connector

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Основная модель базы данных в приложении.
// Дополнение к модели Database из служебной БД.
type DBConfig struct {
	ID       uint
	Name     string
	Host     string
	Port     string
	User     string
	Password string
	Size     string
	Status   string
}

// Добавляет данные о базе данных в служебную БД.
// Перед добавлением в таблицу проверяется подключение.
func (cfg *DBConfig) Save(sql *gorm.DB) error {
	if err := cfg.checkConnection(); err != nil {
		logrus.Error(err)
		return err
	}
	port := portToInt(cfg.Port)
	db := db.Database{
		Name:     cfg.Name,
		Host:     cfg.Host,
		Port:     port,
		Username: cfg.User,
		Password: cfg.Password,
	}
	id, err := db.Create(sql)
	if err != nil {
		logrus.Error(err)
		return err
	}
	cfg.ID = id
	logrus.Infof("Добавлена база данных %v", cfg)
	return nil
}

func GetDbAll(sql *gorm.DB) []DBConfig {
	var DbList []DBConfig
	databases, err := db.GetDbAll(sql)
	if err != nil {
		logrus.Error(err)
		return DbList
	}
	for _, item := range databases {
		config := &DBConfig{
			ID:       item.ID,
			Name:     item.Name,
			Host:     item.Host,
			Port:     strconv.Itoa(item.Port),
			User:     item.Username,
			Password: item.Password,
		}
		if err := config.checkConnection(); err != nil {
			logrus.Error(err)
		}
		if err := config.getDBSize(); err != nil {
			logrus.Error(err)
		}
		DbList = append(DbList, *config)
	}
	return DbList
}

func portToInt(port string) int {
	num, err := strconv.Atoi(port)
	if err != nil {
		logrus.Error("Error converting string to int:", err)
		return 0
	}
	return num
}

// Проверка подключения к базе данных
func (cfg *DBConfig) checkConnection() error {
	command := fmt.Sprintf("pg_isready -h %s -U %s -d %s -p %d", cfg.Host, cfg.User, cfg.Name, portToInt(cfg.Port))
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		cfg.Status = "ошибка"
		return errors.New(command)
	}
	fmt.Print(string(cmd))
	cfg.Status = "активно"
	return nil
}

// получение размера базы данных по имени.
func (cfg *DBConfig) getDBSize() error {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c \"SELECT pg_size_pretty(pg_database_size('%s'))\"", cfg.Password, cfg.Host, cfg.User, cfg.Name, cfg.Name)
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(output)
		return errors.New(string(output))
	}
	//pg_size_pretty
	//----------------
	//7453 kB
	//(1 row)
	startIndex := 35
	endIndex := len(string(output)) - 10
	size := fmt.Sprint(string(output)[startIndex:endIndex]) // -> 7453 kB
	cfg.Size = size
	return nil
}
