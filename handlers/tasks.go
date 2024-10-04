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

// Изменение, удаление задания бекапов.
// Действие зависит от нажатой кнопки
func (h *Handler) actionTaskHandler(c *gin.Context) {
	var action = c.PostForm("action")
	var alias = c.PostForm("alias")
	var dir = c.PostForm("dir")
	switch action {
	case "delete":
		if err := db.DeleteTaskData(alias, dir); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		}
	case "change":
		db.Restore(*h.CONFIG, alias, dir)
	}
	c.Redirect(http.StatusFound, "/tasks")
}
