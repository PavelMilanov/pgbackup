package handlers

import (
	"net/http"
	"strconv"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler для главной страницы с базами данных.
func (h *Handler) databasesHandler(c *gin.Context) {
	databases := db.GetDbAll(h.DB)
	// db.GetDbLastBackup(h.DB)
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

// Handler для сохранения базы данных.
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

// Handler для удаления базы данных.
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

// Handler для создания ручного бекапа для базы данных.
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

// Handler для вывода всех бекапов для выбранной базы данных.
func (h *Handler) getBackupsHandler(c *gin.Context) {
	data := c.Param("id")
	id, _ := strconv.Atoi(data)
	db, err := db.GetDb(h.DB, id)
	if err != nil {
		c.JSON(404, gin.H{"message": "not found"})
		return
	}
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"header": "Базы данных | PgBackup",
		"db":     db,
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
