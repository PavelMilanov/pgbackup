package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) backupsView(c *gin.Context) {
	// dbInfo := db.GetDBInfo(h.DB)

	c.HTML(http.StatusOK, "backups.html", gin.H{
		"backups": []string{""},
	})
}
