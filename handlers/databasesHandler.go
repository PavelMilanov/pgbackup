package handlers

import (
	"net/http"
	"strconv"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) databasesHandler(c *gin.Context) {
	databases := db.GetDbAll(h.DB)
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
	port, _ := strconv.Atoi(data.Port)
	config := db.Database{
		Name:     data.Name,
		Host:     data.Host,
		Port:     port,
		Username: data.Username,
		Password: data.Password,
	}
	config.Save(h.DB)
	c.Redirect(http.StatusFound, "/databases/")
}

func (h *Handler) databaseDeleteHandler(c *gin.Context) {
	var data web.DatabaseForm
	if err := c.ShouldBind(&data); err != nil {
		logrus.Error(err)
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(data.ID)
	config := db.Database{
		ID: id,
	}
	config.Delete(h.DB)
	c.Redirect(http.StatusFound, "/databases/")
}

func (h *Handler) createBackupHandler(c *gin.Context) {
	var data web.DatabaseForm
	if err := c.ShouldBind(&data); err != nil {
		logrus.Error(err)
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(data.ID)
	config := db.Schedule{
		DatabaseID: id,
	}
	config.Save(h.DB, h.CRON)
	c.Redirect(http.StatusFound, "/databases/")
}
