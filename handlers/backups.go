package handlers

import (
	"fmt"
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
		db.СreateBackupDir(db.DEFAULT_BACKUP_DIR)
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
		newBackup, err := backupModel.CreateManualBackup(*h.CONFIG)
		if err != nil {
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
		db.СreateBackupDir(dirName)
		var task = db.Task{
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
		task.CreateTaskData()
		task.CreateCronBackup(h.CRON, *h.CONFIG)
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание создано",
		})
	}
}

// скачивание указанного бекапа
func (h *Handler) downloadBackupHandler(c *gin.Context) {
	var alias = c.Param("alias")
	var date = c.Param("date")
	backups := db.GetBackupData()
	for _, backup := range backups {
		if backup.Alias == alias && backup.Date == date {
			filename := backup.Alias + "-" + backup.Date + ".dump"
			filePath := fmt.Sprintf("%s/%s", backup.Directory, filename)
			fileHeader := fmt.Sprintf("attachment; filename=%s", filename)
			c.Header("Content-Disposition", fileHeader)
			c.File(filePath)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"error": "ошибка при скачивании файла",
	})
}
