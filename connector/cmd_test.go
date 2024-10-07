package connector

import (
	"fmt"
	"testing"
)

// func TestCreateBackupData(t *testing.T) {
// 	backup := Backup{
// 		Alias:    "test",
// 		Status:   "running",
// 		Schedule: BackupSchedule{"manual", "test", "test", "test"},
// 	}
// 	data := сreateBackupData(&backup, "../data")
// 	fmt.Println(data)
// }

func TestGetBackupData(t *testing.T) {
	data := GetBackupData()
	fmt.Println(data)
}

// func TestCheckBackup(t *testing.T) {
// 	checkBackup("../data")
// }

// func TestGetBackupSize(t *testing.T) {
// 	data := getBackupSize("../dumps", "postgres-2024-08-28.dump")
// 	fmt.Print(data)
// }
