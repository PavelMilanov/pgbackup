package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type BackupForm struct {
	SelectedDB      string `form:"backupDBName" binding:"required"`
	SelectedRun     string `form:"backupRun" binding:"required"`
	SelectedComment string `form:"backupComment"`
	SelectedCount   string `form:"backupScheduleCount"`
	SelectedTime    string `form:"backupScheduleTime"`
	SelectedCron    string `form:"backupScheduleCron"`
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
	var data BackupForm
	if err := c.ShouldBind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "backups.html", gin.H{"error": err.Error()})
		return
	}
	fmt.Println(data)
	backupName := data.SelectedDB
	backupRun := data.SelectedRun
	backupComment := data.SelectedComment
	backupCount := data.SelectedCount
	backupTime := data.SelectedTime
	backupCron := data.SelectedCron

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
		_, err := backupModel.CreateManualBackup(*h.CONFIG)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
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
	c.Redirect(http.StatusFound, "/backups")
}

// Скачивание, удаление, восстановление бекапа.
// Действие зависит от нажатой кнопки
func (h *Handler) actionBackupHandler(c *gin.Context) {
	var action = c.PostForm("action")
	var alias = c.PostForm("alias")
	var date = c.PostForm("date")
	switch action {
	case "download":
		backups := db.GetBackupData()
		for _, backup := range backups {
			if backup.Alias == alias && backup.Date == date {
				fileName := backup.Alias + "-" + backup.Date + ".dump"
				filePath := fmt.Sprintf("%s/%s", backup.Directory, fileName)
				fileHeader := fmt.Sprintf("attachment; filename=%s", fileName)
				c.Header("Content-Disposition", fileHeader)
				c.File(filePath)
				log.Printf("%s скачан", filePath)
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"error": "ошибка при скачивании файла",
		})
	case "delete":
		if err := db.DeleteBackupData(alias, date); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		}
	case "restore":
		db.Restore(*h.CONFIG, alias, date)
	}
	c.Redirect(http.StatusFound, "/backups")
}
