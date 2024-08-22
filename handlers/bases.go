package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

func (h *Handler) basesView(c *gin.Context) {
	dbInfo := db.GetDBInfo(h.DB)
	c.HTML(http.StatusOK, "bases.html", gin.H{
		"databases": dbInfo,
	})
}
