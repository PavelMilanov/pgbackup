package handlers

import (
	"net/http"
	"strconv"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) settingsHandler(c *gin.Context) {
	settings, _ := db.GetSettings(h.DB.Sql)
	settings.Version = config.VERSION
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "settings.html", gin.H{
			"header": "Настройки | PgBackup",
			"config": settings,
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: false},
				{Name: "Базы данных", URL: "/databases", IsVisible: false},
				{Name: "Настройки", URL: "/settings", IsVisible: true},
			}})
	} else if c.Request.Method == "POST" {
		var data web.GeneralSettingsForm
		if err := c.ShouldBind(&data); err != nil {
			return
		}
		count, _ := strconv.Atoi(data.BackupCount)
		settings.BackupCount = count
		if err := settings.Update(h.DB.Sql); err != nil {
			c.HTML(http.StatusOK, "settings.html", gin.H{
				"header": "Настройки | PgBackup",
				"config": settings,
				"notification": web.Notify{
					Message: err.Error(),
					Type:    config.NOTIFY_STATUS["ошибка"],
				},
				"pages": []web.Page{
					{Name: "Главная", URL: "/", IsVisible: false},
					{Name: "Расписание", URL: "/schedule", IsVisible: false},
					{Name: "Базы данных", URL: "/databases", IsVisible: false},
					{Name: "Настройки", URL: "/settings", IsVisible: true},
				}})
		}
		c.HTML(http.StatusOK, "settings.html", gin.H{
			"header": "Настройки | PgBackup",
			"config": settings,
			"notification": web.Notify{
				Message: "Настройки обновлены!",
				Type:    config.NOTIFY_STATUS["инфо"],
			},
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: false},
				{Name: "Базы данных", URL: "/databases", IsVisible: false},
				{Name: "Настройки", URL: "/settings", IsVisible: true},
			}})
	}
}
