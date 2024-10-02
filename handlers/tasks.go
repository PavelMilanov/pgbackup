package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

func (h *Handler) tasksView(c *gin.Context) {
	dbData := db.GetDBData(*h.CONFIG)
	c.HTML(http.StatusOK, "tasks.html", gin.H{
		"databases": dbData,
	})
}
