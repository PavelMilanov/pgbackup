package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/handlers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env файл не найден")
	}

	config := db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	postgres, err := db.NewPostgreDB(config)
	handler := handlers.NewHandler(postgres)
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
	if err := srv.Shutdown(1); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
	db.Close(postgres)
}
