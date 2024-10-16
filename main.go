package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/handlers"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var duration = 1

func init() {
	// создаем дефолтные директории
	os.Mkdir(connector.BACKUP_DIR, 0755)
	// os.Mkdir(connector.DEFAULT_BACKUP_DIR, 0755)
}

func main() {

	// config := connector.Config{
	// 	Host:     os.Getenv("POSTGRES_HOST"),
	// 	Port:     os.Getenv("POSTGRES_PORT"),
	// 	User:     os.Getenv("POSTGRES_USER"),
	// 	Password: os.Getenv("POSTGRES_PASSWORD"),
	// 	DBName:   os.Getenv("POSTGRES_DB"),
	// }
	/// база данных
	/// первичная инициализация задания для ручных бекапов
	sqlite := db.NewDatabase(&db.SQLite{Name: "pgbackup.db"})
	// var defaultTask db.Task
	// result := sqlite.Where("Alias = ?", "Default").First(&defaultTask)
	// if result.Error != nil {
	// 	// Запись не найдена, создаем новую
	// 	if err := sqlite.Create(&db.Task{
	// 		Alias:     "Default",
	// 		Directory: connector.DEFAULT_BACKUP_DIR,
	// 	}).Error; err != nil {
	// 		log.Fatalf("failed to create task: %v", err)
	// 	}
	// }
	///
	/// Фоновые задачи
	location, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.New(cron.WithLocation(location))
	fmt.Println(scheduler.Entries())
	///
	/// логгер
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:00",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("[ %s:%d]", filename, f.Line)
			return "", filename
		},
	})
	///
	/// фоновые задачи
	go scheduler.Start()
	defer scheduler.Stop()

	// tasks := connector.GetTaskData()
	// for _, task := range tasks {
	// 	if task.Schedule.Run == connector.BACKUP_RUN[1] {
	// 		task.CreateCronBackup(scheduler, config)
	// 	}
	// }
	///

	handler := handlers.NewHandler(sqlite, scheduler)

	srv := new(Server)
	go func() {
		if err := srv.Run(handler.InitRouters()); err != nil {
			logrus.Warn(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutdown signal of %d seconds\n", duration)
	if err := srv.Shutdown(time.Duration(duration)); err != nil {
		logrus.WithError(err).Error("ошибка при остановке сервера")
	}
}
