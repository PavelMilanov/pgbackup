package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) mainHandler(c *gin.Context) {
	backups := db.GetDbAll(h.DB)
	c.HTML(http.StatusOK, "main.html", gin.H{
		"header":  "Главная | PgBackup",
		"backups": backups,
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: true},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
