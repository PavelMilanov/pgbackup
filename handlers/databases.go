package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) databasesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "databases.html", gin.H{
		"header": "Базы данных | PgBackup",
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: true},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
