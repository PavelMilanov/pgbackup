package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
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
}

func main() {
	/// Фоновые задачи
	location, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.New(cron.WithLocation(location))

	/// база данных
	/// первичная инициализация задания для ручных бекапов
	sqlite := db.NewDatabase(&db.SQLite{Name: "pgbackup.db"}, scheduler)

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
	logrus.Infof("Shutdown signal of %d seconds\n", config.DURATION)
	if err := srv.Shutdown(time.Duration(config.DURATION)); err != nil {
		logrus.WithError(err).Error("ошибка при остановке сервера")
	}
}
