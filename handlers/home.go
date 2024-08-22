package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

func (h *Handler) homeView(c *gin.Context) {
	dbInfo := db.GetDBInfo(h.DB)
	c.HTML(http.StatusOK, "home.html", gin.H{"databases": dbInfo})
}
