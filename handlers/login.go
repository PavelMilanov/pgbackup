package handlers

import (
	"net/http"
	"os"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (h *Handler) authView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *Handler) submitLoginForm(c *gin.Context) {
	var formData LoginForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	username := os.Getenv("USER_LOGIN")
	password := os.Getenv("USER_PASSWORD")
	if formData.Username != username || formData.Password != password {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"data": "неправильные логин или пароль",
		})
		return
	}
	dbInfo := db.GetDBData(h.DB)
	backupsInfo := db.GetBackupData(db.BACKUPDATA_DIR)
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"databases": dbInfo,
		"backups":   backupsInfo,
		"run":       db.BACKUP_RUN,
	})
}
