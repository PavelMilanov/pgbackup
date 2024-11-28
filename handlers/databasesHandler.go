package handlers

import (
	"net/http"
	"strconv"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

// Handler для главной страницы с базами данных.
func (h *Handler) databasesHandler(c *gin.Context) {
	databases := db.GetDbAll(h.DB.Sql)
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
		return
	}
	port, _ := strconv.Atoi(data.Port)
	cfg := db.Database{
		Alias:    data.Alias,
		Name:     data.Name,
		Host:     data.Host,
		Port:     port,
		Username: data.Username,
		Password: data.Password,
	}
	if err := cfg.Save(h.DB.Sql); err != nil {
		databases := db.GetDbAll(h.DB.Sql)
		c.HTML(http.StatusOK, "databases.html", gin.H{
			"header":    "Базы данных | PgBackup",
			"databases": databases,
			"notification": web.Notify{
				Message: err.Error(),
				Type:    config.NOTIFY_STATUS["ошибка"],
			},
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: false},
				{Name: "Базы данных", URL: "/databases", IsVisible: true},
				{Name: "Настройки", URL: "/settings", IsVisible: false},
			}})
		return
	}
	databases := db.GetDbAll(h.DB.Sql)
	c.HTML(http.StatusOK, "databases.html", gin.H{
		"header":    "Базы данных | PgBackup",
		"databases": databases,
		"notification": web.Notify{
			Message: "База данных добавлена!",
			Type:    config.NOTIFY_STATUS["инфо"],
		},
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

// Handler для удаления базы данных.
func (h *Handler) databaseDeleteHandler(c *gin.Context) {
	var data web.DatabaseForm
	if err := c.ShouldBind(&data); err != nil {
		return
	}
	id, _ := strconv.Atoi(data.ID)
	cfg := db.Database{
		ID: id,
	}
	if err := cfg.Delete(h.DB.Sql); err != nil {
		databases := db.GetDbAll(h.DB.Sql)
		c.HTML(http.StatusOK, "databases.html", gin.H{
			"header":    "Базы данных | PgBackup",
			"databases": databases,
			"notification": web.Notify{
				Message: err.Error(),
				Type:    config.NOTIFY_STATUS["ошибка"],
			},
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: false},
				{Name: "Базы данных", URL: "/databases", IsVisible: true},
				{Name: "Настройки", URL: "/settings", IsVisible: false},
			}})
		return
	}
	databases := db.GetDbAll(h.DB.Sql)
	c.HTML(http.StatusOK, "databases.html", gin.H{
		"header":    "Базы данных | PgBackup",
		"databases": databases,
		"notification": web.Notify{
			Message: "База данных удалена!",
			Type:    config.NOTIFY_STATUS["инфо"],
		},
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

// Handler для создания ручного бекапа для базы данных.
func (h *Handler) createBackupHandler(c *gin.Context) {
	var data web.DatabaseForm
	if err := c.ShouldBind(&data); err != nil {
		return
	}
	id, _ := strconv.Atoi(data.ID)
	cfg := db.Schedule{
		DatabaseID: id,
	}
	if err := cfg.Save(h.DB, h.CRON); err != nil {
		db, _ := db.GetDb(h.DB.Sql, id)
		c.HTML(http.StatusOK, "backups.html", gin.H{
			"header": "Базы данных | PgBackup",
			"db":     db,
			"notification": web.Notify{
				Message: err.Error(),
				Type:    config.NOTIFY_STATUS["ошибка"],
			},
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: false},
				{Name: "Базы данных", URL: "/databases", IsVisible: true},
				{Name: "Настройки", URL: "/settings", IsVisible: false},
			}})
		return
	}
	db, _ := db.GetDb(h.DB.Sql, id)
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"header": "Базы данных | PgBackup",
		"db":     db,
		"notification": web.Notify{
			Message: "Выполнение дампа начато!",
			Type:    config.NOTIFY_STATUS["инфо"],
		},
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

// Handler для вывода всех бекапов для выбранной базы данных.
func (h *Handler) getBackupsHandler(c *gin.Context) {
	var data web.BackupForm
	if err := c.ShouldBind(&data); err != nil {
		return
	}
	id, _ := strconv.Atoi(data.ID)
	db, _ := db.GetDb(h.DB.Sql, id)
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

// Handler для скачивания файла бекапа.
func (h *Handler) downloadBackupHandler(c *gin.Context) {
	var data web.BackupForm
	if err := c.ShouldBind(&data); err != nil {
		return
	}
	id, _ := strconv.Atoi(data.ID)
	config := db.Backup{
		ID: id,
	}
	config.Get(h.DB.Sql)
	filepath := config.Directory + "/" + config.Dump
	c.FileAttachment(filepath, config.Dump)
}

// Handler для удаления файла бекапа.
func (h *Handler) deleteBackupHandler(c *gin.Context) {
	var data web.BackupForm
	if err := c.ShouldBind(&data); err != nil {
		return
	}
	id, _ := strconv.Atoi(data.ID)
	cfg := db.Backup{
		ID: id,
	}
	cfg.Delete(h.DB)
	databases := db.GetDbAll(h.DB.Sql)
	c.HTML(http.StatusOK, "databases.html", gin.H{
		"header":    "Базы данных | PgBackup",
		"databases": databases,
		"notification": web.Notify{
			Message: "Дамп удален!",
			Type:    config.NOTIFY_STATUS["инфо"],
		},
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
