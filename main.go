package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/handlers"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var duration = 1

func init() {
	// создаем дефолтные директории
	if err := os.Mkdir(db.BACKUP_DIR, 0755); err != nil {
	}
	if err := os.Mkdir(db.DEFAULT_BACKUP_DIR, 0755); err != nil {
	}
	if err := os.Mkdir(db.BACKUPDATA_DIR, 0755); err != nil {
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env файл не найден")
	}

	config := db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	location, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.New(cron.WithLocation(location))

	/// фоновые задачи
	go scheduler.Start()
	defer scheduler.Stop()

	tasks := db.GetTaskData()
	for _, task := range tasks {
		if task.Schedule.Run == db.BACKUP_RUN[1] {
			task.CreateCronBackup(scheduler, config)
		}
	}
	///

	handler := handlers.NewHandler(&config, scheduler)

	srv := new(Server)
	go func() {
		if err := srv.Run(handler.InitRouters()); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutdown signal of %d seconds\n", duration)
	if err := srv.Shutdown(time.Duration(duration)); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}
