package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

func (h *Handler) tasksView(c *gin.Context) {
	tasksData := db.GetTaskData()
	c.HTML(http.StatusOK, "tasks.html", gin.H{
		"tasks": tasksData,
	})
}
