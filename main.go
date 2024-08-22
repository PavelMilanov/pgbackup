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
)

var duration = 1

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

	postgres, err := db.NewPostgreDB(&config)
	handler := handlers.NewHandler(postgres, &config)
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
