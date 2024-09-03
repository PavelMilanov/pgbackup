package handlers

import (
	"log"
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type Data struct {
	SelectedDB    string `json:"db"`
	SelectedRun   string `json:"run"`
	SelectedCount string `json:"count"`
	SelectedTime  string `json:"time"`
	SelectedCron  string `json:"cron"`
}

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := db.GetDBData(h.DB)
	backupsInfo := db.GetBackupData(db.BACKUPDATA_DIR)
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"databases": dbInfo,
		"backups":   backupsInfo,
		"run":       db.BACKUP_RUN,
	})
}

func (h *Handler) backupHandler(c *gin.Context) {
	var json Data
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dbname := json.SelectedDB
	backupRun := json.SelectedRun
	backupCount := json.SelectedCount
	backupTime := json.SelectedTime
	backupCron := json.SelectedCron

	if backupRun == db.BACKUP_RUN[1] && (backupCount == "" || backupCron == "" || backupTime == "") {
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание не может быть пустым",
		})
		return
	}
	switch backupRun {
	case db.BACKUP_RUN[0]: // вручную
		newBackup, err := db.CreateBackup(*h.CONFIG, dbname, backupRun, backupCount, backupTime, backupCron)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": newBackup,
		})
	case db.BACKUP_RUN[1]: // по расписанию
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание создано",
		})
		return
	}
	// newBackup, err := db.CreateBackup(*h.CONFIG, dbname, backupRun, backupCount, backupTime, backupCron)
	// if err != nil {
	// 	log.Fatal(err)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": newBackup,
	// })
}
