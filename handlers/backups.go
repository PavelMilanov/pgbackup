package handlers

import (
	"fmt"
	"net/http"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	dbInfo := connector.GetDBData(*h.CONFIG)
	backupsInfo := connector.GetBackupData()
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"databases": dbInfo,
		"backups":   backupsInfo,
		"run":       connector.BACKUP_RUN,
		"pages": []web.Page{
			{Name: "Бэкапы", URL: "/backups", IsVisible: true},
			{Name: "Задания", URL: "/tasks", IsVisible: false},
		}})
}

func (h *Handler) backupHandler(c *gin.Context) {
	var data BackupForm
	if err := c.ShouldBind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "backups.html", gin.H{"error": err.Error()})
		return
	}
	backupName := data.SelectedDB
	backupRun := data.SelectedRun
	backupComment := data.SelectedComment
	backupCount := data.SelectedCount
	backupTime := data.SelectedTime
	backupCron := data.SelectedCron

	if backupRun == connector.BACKUP_RUN[1] && (backupCount == "" || backupCron == "" || backupTime == "") {
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание не может быть пустым",
		})
		return
	}
	switch backupRun {
	case connector.BACKUP_RUN[0]: // вручную
		connector.СreateBackupDir(connector.DEFAULT_BACKUP_DIR)
		var backupModel = connector.Backup{
			Alias:     backupName,
			Comment:   backupComment,
			Directory: connector.DEFAULT_BACKUP_DIR,
			Schedule: connector.BackupSchedule{
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
	case connector.BACKUP_RUN[1]: // по расписанию
		dirName := connector.GenerateRandomBackupDir()
		connector.СreateBackupDir(dirName)
		var task = connector.Task{
			Alias:     backupName,
			Comment:   backupComment,
			Directory: dirName,
			Schedule: connector.BackupSchedule{
				Run:   backupRun,
				Count: backupCount,
				Time:  backupTime,
				Cron:  backupCron,
			},
		}
		task.CreateTaskData()
		task.CreateCronBackup(h.CRON, *h.CONFIG)
		// c.JSON(http.StatusOK, gin.H{
		// 	"error": "расписание создано",
		// })
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
		backups := connector.GetBackupData()
		for _, backup := range backups {
			if backup.Alias == alias && backup.Date == date {
				fileName := backup.Alias + "-" + backup.Date + ".dump"
				filePath := fmt.Sprintf("%s/%s", backup.Directory, fileName)
				fileHeader := fmt.Sprintf("attachment; filename=%s", fileName)
				c.Header("Content-Disposition", fileHeader)
				c.File(filePath)
				logrus.Infof("%s скачан", filePath)
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"error": "ошибка при скачивании файла",
		})
	case "delete":
		if err := connector.DeleteBackupData(alias, date); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		}
	case "restore":
		connector.Restore(*h.CONFIG, alias, date)
	}
	c.Redirect(http.StatusFound, "/backups")
}
