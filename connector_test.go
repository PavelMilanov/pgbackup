package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetDbInfo(t *testing.T) {
	test_db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	emptyDb := [2]string{"postgres", "dev"} // пустая база данных с инициализированной БД  dev
	result := test_db.getDBs()
	if reflect.DeepEqual(result, emptyDb) {
		t.Errorf("Некорретный вывод баз данных %s %s", result, emptyDb)
	}
	fmt.Println(result)
}

func TestCheckConnection(t *testing.T) {
	test_db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	test_db.checkConnection()
}

func TestBackup(t *testing.T) {
	test_db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	test_db.backup()
}

func TestRestore(t *testing.T) {
	test_db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	test_db.restore("backup")
}
