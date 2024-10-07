package connector

import (
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
	t.Log(conn)
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

func TestGetDBSize(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "admin",
		DBName:   "postgres",
	}
	size := getDBSize(config, "dev")
	t.Log(size)
}

func TestGetDBName(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "admin",
		DBName:   "postgres",
	}
	data := getDBName(config)
	t.Log(data)
}
