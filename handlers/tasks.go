package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) tasksView(c *gin.Context) {
	var tasksInfo []db.Task
	h.DB.Where("Alias <> ?", "Default").Find(&tasksInfo) // получаем все задачи, кроме дефолтной
	c.HTML(http.StatusOK, "tasks.html", gin.H{
		"tasks": tasksInfo,
		"pages": []web.Page{
			{Name: "Бэкапы", URL: "/backups", IsVisible: false},
			{Name: "Задания", URL: "/tasks", IsVisible: true},
		}})
}

// Удаление задания бекапов.
// Действие зависит от нажатой кнопки
func (h *Handler) actionTaskHandler(c *gin.Context) {
	var action = c.PostForm("action")
	var id = c.PostForm("id")
	switch action {
	case "delete":
		// if err := connector.DeleteTaskData(alias, dir); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"error": err.Error(),
		// 	})
		// }
		// connector.DeleteBackupDir(dir)
		h.DB.Delete(&db.Task{}, id)
	}
	c.Redirect(http.StatusFound, "/tasks")
}
