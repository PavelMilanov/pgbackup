package main

import (
	"testing"
)

func TestGetDbInfo(t *testing.T) {
	test_db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	test_db.getDbInfo()
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
