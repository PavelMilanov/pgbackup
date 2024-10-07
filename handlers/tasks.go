package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/gin-gonic/gin"
)

func (h *Handler) tasksView(c *gin.Context) {
	tasksData := connector.GetTaskData()
	c.HTML(http.StatusOK, "tasks.html", gin.H{
		"tasks": tasksData,
	})
}

// Удаление задания бекапов.
// Действие зависит от нажатой кнопки
func (h *Handler) actionTaskHandler(c *gin.Context) {
	var action = c.PostForm("action")
	var alias = c.PostForm("alias")
	var dir = c.PostForm("dir")
	switch action {
	case "delete":
		if err := connector.DeleteTaskData(alias, dir); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		}
		connector.DeleteBackupDir(dir)
	}
	c.Redirect(http.StatusFound, "/tasks")
}
