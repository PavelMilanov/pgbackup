package handlers

import (
	"fmt"
	"testing"
)

func TestCreateBackupData(t *testing.T) {
	backup := Backup{
		Alias:  "test",
		Status: "running",
		Run:    "manual",
	}
	data := createBackupData(&backup, "../data")
	fmt.Println(data)
}

func TestGetBackupData(t *testing.T) {
	data := getBackupData("../data")
	fmt.Println(data)
}

func TestCheckBackup(t *testing.T) {
	checkBackup("../data")
}
