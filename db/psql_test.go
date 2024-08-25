package db

import (
	"testing"
)

func TestCheckConnection(t *testing.T) {

}

func TestCreateBackup(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "admin",
		DBName:   "dev",
	}
	err := CreateBackup(config, "dev", "dev-2024-08-24")
	if err != nil {
		t.Error(err)
	}
}
