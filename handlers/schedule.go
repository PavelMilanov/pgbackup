package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) scheduleHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "schedule.html", gin.H{
		"header": "Расписание | PgBackup",
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: true},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
