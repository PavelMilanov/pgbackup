package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/handlers"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	// создаем директорию для бекапов
	os.Mkdir(config.BACKUP_DIR, 0755)
	// создаем директорию для базы
	os.Mkdir(config.DATA_DIR, 0755)

	jwt := os.Getenv("JWT_KEY")
	aes := os.Getenv("AES_KEY")
	if jwt == "" || aes == "" {
		logrus.Fatalf("Не указана переменная окружения JWT_KEY или AES_KEY")
	}
}

func main() {
	/// логгер
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:00",
	})

	/// Фоновые задачи
	location, _ := time.LoadLocation(os.Getenv("TZ"))
	c := cron.New(cron.WithLocation(location))

	/// первичная инициализация задания для ручных бекапов
	sqliteFIle := filepath.Join(config.DATA_DIR, "pgbackup.db")
	sqlite := db.NewDatabase(sqliteFIle, c)
	defer db.CloseDatabase(sqlite.Sql)
	c.AddFunc("*/1 * * * *", func() { // 0 0 * * *
		db.StartSystemTasks(&sqlite)
	})
	logrus.Debug("Системный планировщик запущен")

	c.AddFunc("@every 1m", func() {
		db.CheckDBconnection(&sqlite)
	})
	logrus.Debug("Планировщик баз данных запущен")

	c.Start()
	defer c.Stop()

	handler := handlers.NewHandler(&sqlite, c)
	srv := new(Server)
	go func() {
		if err := srv.Run(handler.InitRouters()); err != nil {
			logrus.Warn(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Infof("Сигнал остановки сервера через %d секунды\n", config.DURATION)
	if err := srv.Shutdown(time.Duration(config.DURATION)); err != nil {
		logrus.WithError(err).Error("ошибка при остановке сервера")
	}
}
