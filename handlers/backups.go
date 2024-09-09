package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type Data struct {
	SelectedDB      string `json:"db"`
	SelectedRun     string `json:"run"`
	SelectedComment string `json:"comment"`
	SelectedCount   string `json:"count"`
	SelectedTime    string `json:"time"`
	SelectedCron    string `json:"cron"`
}

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := db.GetDBData(h.DB)
	backupsInfo := db.GetBackupData()
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
	backupName := json.SelectedDB
	backupRun := json.SelectedRun
	backupComment := json.SelectedComment
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
		var backupModel = db.Backup{
			Alias:     backupName,
			Comment:   backupComment,
			Directory: db.DEFAULT_BACKUP_DIR,
			Schedule: db.BackupSchedule{
				Run:   backupRun,
				Count: backupCount,
				Time:  backupTime,
				Cron:  backupCron,
			},
		}
		newBackup, err := backupModel.CreateBackup(*h.CONFIG)
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
		dirName := db.GenerateRandomBackupDir()
		var backupModel = db.Backup{
			Alias:     backupName,
			Comment:   backupComment,
			Directory: dirName,
			Schedule: db.BackupSchedule{
				Run:   backupRun,
				Count: backupCount,
				Time:  backupTime,
				Cron:  backupCron,
			},
		}
		go backupModel.CreateBackup(*h.CONFIG)
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание создано",
		})
	}
}

// скачивание указанного бекапа
func (h *Handler) downloadBackupHandler(c *gin.Context) {
	var backup = c.Param("backup")
	filePath := fmt.Sprintf("%s/%s", db.BACKUP_DIR, backup)
	c.FileAttachment(filePath, backup)
}
