package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) databasesHandler(c *gin.Context) {
	databases := connector.GetDbAll(h.DB)
	c.HTML(http.StatusOK, "databases.html", gin.H{
		"header":    "Базы данных | PgBackup",
		"databases": databases,
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

func (h *Handler) databaseSaveHandler(c *gin.Context) {
	var data web.DatabaseForm
	if err := c.ShouldBind(&data); err != nil {
		logrus.Error(err)
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	config := connector.DBConfig{
		Name:     data.Name,
		Host:     data.Host,
		Port:     data.Port,
		User:     data.Username,
		Password: data.Password,
	}
	config.Save(h.DB)
	c.Redirect(http.StatusFound, "/databases/")
}

func (h *Handler) createBackupHandler(c *gin.Context) {
	var data web.DatabaseForm
	if err := c.ShouldBind(&data); err != nil {
		logrus.Error(err)
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	schedule := connector.ScheduleConfig{
		DbID:      data.ID,
		Directory: connector.DEFAULT_BACKUP_DIR,
	}
	schedule.SaveManual(h.DB)
	c.Redirect(http.StatusFound, "/databases/")
}
