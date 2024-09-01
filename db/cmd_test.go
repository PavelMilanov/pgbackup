package db

import (
	"fmt"
	"testing"
)

func TestCreateBackupData(t *testing.T) {
	backup := Backup{
		Alias:    "test",
		Status:   "running",
		Schedule: BackupSchedule{"manual", "test", "test", "test"},
	}
	data := CreateBackupData(&backup, "../data")
	fmt.Println(data)
}

func TestGetBackupData(t *testing.T) {
	data := GetBackupData("../data")
	fmt.Println(data)
}

func TestCheckBackup(t *testing.T) {
	CheckBackup("../data")
}

func TestGetBackupSize(t *testing.T) {
	data := GetBackupSize("../dumps", "postgres-2024-08-28.dump")
	fmt.Print(data)
}
