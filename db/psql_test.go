package db

import (
	"fmt"
	"testing"
)

func TestCheckConnection(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "admin",
		DBName:   "dev",
	}
	conn := CheckConnection(config)
	fmt.Println(conn)
}

func TestCreateBackup(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "admin",
		DBName:   "dev",
	}
	timer, err := CreateBackup(config, "dev", "dev-2024-08-24")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(timer)
}
