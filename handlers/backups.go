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

type Data struct {
	SelectedDB       string `json:"db"`
	SelectedSchedule string `json:"schedule"`
}

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := db.GetDBInfo(h.DB)
	backupsInfo := db.GetBackupData(db.BACKUPDATA_DIR)
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"databases": dbInfo,
		"backups":   backupsInfo,
		"schedules": db.BACKUP_SCHEDULE,
	})
}

func (h *Handler) createBackup(c *gin.Context) {
	var json Data
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dbname := json.SelectedDB
	schedule := json.SelectedSchedule
	fmt.Println(dbname, schedule)

	currTime := time.Now().Format("2006-01-02") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	backupName := dbname + "-" + currTime

	curBackups := db.CheckBackup(db.BACKUP_DIR)
	for _, item := range curBackups {
		if strings.Contains(item, backupName) {
			errStr := fmt.Sprintf("Бекап %s уже существует!", backupName)
			log.Println(errStr)
			c.JSON(http.StatusOK, gin.H{
				"error": errStr,
			})
			return
		}

	}
	timer, err := db.CreateBackup(*h.CONFIG, dbname, backupName)
	if err != nil {
		log.Fatal(err)
	}
	size := db.GetBackupSize(db.BACKUP_DIR, backupName)
	newBackup := db.Backup{
		Alias:    dbname,
		Date:     currTime,
		Size:     size,
		LeadTime: timer,
		Status:   "создан",
		Run:      schedule,
	}

	backupsInfo := db.CreateBackupData(&newBackup, db.BACKUPDATA_DIR)
	log.Println(backupsInfo)

	c.JSON(http.StatusOK, gin.H{
		"message": newBackup,
	})
}
