package handlers

import (
	"fmt"
	"net/http"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := connector.GetDBData(*h.CONFIG)
	var backupsInfo []db.Backup
	h.DB.Find(&backupsInfo)
	// backupsInfo := connector.GetBackupData()
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
	var data web.BackupForm
	if err := c.ShouldBind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "backups.html", gin.H{"error": err.Error()})
		return
	}

	if data.SelectedRun == connector.BACKUP_RUN[1] && (data.SelectedCount == "" || data.SelectedCron == "" || data.SelectedTime == "") {
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание не может быть пустым",
		})
		return
	}
	switch data.SelectedRun {
	case connector.BACKUP_RUN[0]: // вручную
		err := connector.CreateManualBackup(*h.CONFIG, h.DB, data)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Redirect(http.StatusFound, "/backups/")
	case connector.BACKUP_RUN[1]: // по расписанию
		connector.CreateCronBackup(h.CRON, *h.CONFIG, h.DB, data)
		// c.JSON(http.StatusOK, gin.H{
		// 	"error": "расписание создано",
		// })
		c.Redirect(http.StatusFound, "/tasks/")
	}
}

// Скачивание, удаление, восстановление бекапа.
// Действие зависит от нажатой кнопки
func (h *Handler) actionBackupHandler(c *gin.Context) {
	var action = c.PostForm("action")
	var id = c.PostForm("id")
	var backup db.Backup
	err := backup.Get(h.DB, id)
	if err != nil {
		logrus.Error(err)
	}
	switch action {
	case "download":
		fileHeader := fmt.Sprintf("attachment; filename=%s-%s.dump", backup.Alias, backup.Date)
		c.Header("Content-Disposition", fileHeader)
		c.File(backup.Dump)
		logrus.Infof("%s скачан", backup.Dump)
		c.JSON(http.StatusOK, gin.H{
			"error": "ошибка при скачивании файла",
		})
	case "delete":
		if err := backup.Delete(h.DB); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		}
	case "restore":
		err := connector.Restore(*h.CONFIG, backup)
		if err != nil {
			logrus.Error(err)
		}
	}
	c.Redirect(http.StatusFound, "/backups")
}
