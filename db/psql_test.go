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

// func TestCreateBackup(t *testing.T) {
// 	config := Config{
// 		Host:     "localhost",
// 		Port:     "5432",
// 		User:     "admin",
// 		Password: "admin",
// 		DBName:   "dev",
// 	}
// 	backup, err := CreateBackup(config)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	fmt.Println(backup)
// }
