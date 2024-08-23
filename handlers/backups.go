package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type BackupData struct {
	DBnames []string
	Backups []Backup
}

type Backup struct {
	Alias  string
	Status string
}

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := db.GetDBInfo(h.DB)
	data := []BackupData{}
	data = append(data, BackupData{Backups: []Backup{}, DBnames: dbInfo})
	c.HTML(http.StatusOK, "backups.html", data)
}

func (h *Handler) createBackup(c *gin.Context) {
	dbname := c.PostForm("dbname")
	dbInfo := db.GetDBInfo(h.DB)
	// newBackup, err := db.CreateBackup(*h.CONFIG, dbname)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// newBackupData := cmd.CreateBackupData()
	data := []BackupData{}
	data = append(data, BackupData{Backups: []Backup{{Alias: dbname, Status: "created"}}, DBnames: dbInfo})

	c.HTML(http.StatusOK, "backups.html", data)
}
