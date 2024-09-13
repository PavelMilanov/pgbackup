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

	postgres, err := db.NewPostgreDB(&config)

	jakartaTime, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	/// фоновые задачи
	go scheduler.Start()
	defer scheduler.Stop()

	tasks := db.GetTaskData()
	for _, task := range tasks {
		if task.Schedule.Run == db.BACKUP_RUN[1] {
			cron := task.ToCron()
			scheduler.AddFunc(cron, func() {
				var backupModel = db.Backup{
					Alias:     task.Alias,
					Comment:   task.Comment,
					Directory: task.Directory,
					Schedule: db.BackupSchedule{
						Run:   task.Schedule.Run,
						Count: task.Schedule.Count,
						Time:  task.Schedule.Time,
						Cron:  task.Schedule.Cron,
					},
				}
				newBackup, err := backupModel.CreateBackup(config)
				if err != nil {
					log.Println(err)
				}
				db.CreateBackupData(newBackup)
			})
		}
	}
	jobs := scheduler.Entries()
	for _, job := range jobs {
		log.Printf("Job ID: %d, Next Run: %s\n", job.ID, job.Next)
	}
	///

	handler := handlers.NewHandler(postgres, &config, scheduler)
	if err != nil {
		log.Fatal(err)
	}

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
	db.Close(postgres)
}
