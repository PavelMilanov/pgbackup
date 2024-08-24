package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type Backup struct {
	Alias  string
	Status string
	Run    string
}

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := db.GetDBInfo(h.DB)
	backupsInfo := getBackupData(db.BACKUPDATA_DIR)
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"databases": dbInfo,
		"backups":   backupsInfo,
	})
}

func (h *Handler) createBackup(c *gin.Context) {
	dbname := c.PostForm("dbname")

	dbInfo := db.GetDBInfo(h.DB)

	currTime := time.Now().Format("2006-01-02") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	backupName := dbname + "-" + currTime

	curBackups := checkBackup(db.BACKUP_DIR)
	backupIsExist := false
	for _, item := range curBackups {
		if strings.Contains(item, backupName) {
			backupIsExist = true
			break
		}

	}
	switch backupIsExist {
	case false:
		newBackupAlias, err := db.CreateBackup(*h.CONFIG, dbname, backupName)
		fmt.Println(newBackupAlias)
		if err != nil {
			log.Fatal(err)
		}
		newBackup := Backup{
			Alias:  newBackupAlias,
			Status: "создан",
			Run:    "без расписания",
		}

		backupsInfo := createBackupData(&newBackup, db.BACKUPDATA_DIR)
		c.HTML(http.StatusOK, "backups.html", gin.H{
			"databases": dbInfo,
			"backups":   backupsInfo,
		})
	case true:
		errStr := fmt.Sprintf("Бекап %s уже существует!", backupName)
		backupsInfo := getBackupData(db.BACKUPDATA_DIR)
		c.HTML(http.StatusOK, "backups.html", gin.H{
			"databases": dbInfo,
			"backups":   backupsInfo,
			"error":     errStr,
		})
	}
}
