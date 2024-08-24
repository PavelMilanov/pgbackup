package handlers

import (
	"fmt"
	"os"
	"testing"

	"github.com/PavelMilanov/pgbackup/db"
)

func TestCreateBackupData(t *testing.T) {
	backup := Backup{
		Alias:  "test",
		Status: "running",
		Run:    "manual",
	}
	data := createBackupData(&backup, ".")
	fmt.Println(data)
}

func TestGetBackupData(t *testing.T) {
	data := getBackupData(".")
	fmt.Println(data)
}

func TestCheckBackup(t *testing.T) {
	path := fmt.Sprintf("../%s", db.BACKUP_DIR)
	_ = os.Chdir(path)
	checkBackup(".")
}
