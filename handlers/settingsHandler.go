package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) settingsHandler(c *gin.Context) {
	config, _ := db.GetSettings(h.DB)
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "settings.html", gin.H{
			"header": "Настройки | PgBackup",
			"config": config,
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
		config := db.Setting{BackupCount: count}
		if err := config.Update(h.DB); err != nil {
			fmt.Println(err)
		}
		c.HTML(http.StatusOK, "settings.html", gin.H{
			"header": "Настройки | PgBackup",
			"config": config,
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: false},
				{Name: "Базы данных", URL: "/databases", IsVisible: false},
				{Name: "Настройки", URL: "/settings", IsVisible: true},
			}})
	}
}
