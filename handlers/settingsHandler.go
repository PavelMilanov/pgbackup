package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) settingsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"header": "Настройки | PgBackup",
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: true},
		}})
}
