package handlers

import (
	"text/template"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/gin-gonic/gin"

	cron "github.com/robfig/cron/v3"
)

type Handler struct {
	CONFIG *connector.Config
	CRON   *cron.Cron
}

func NewHandler(config *connector.Config, scheduler *cron.Cron) *Handler {
	return &Handler{CONFIG: config, CRON: scheduler}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")
	web := router.Group("/")
	{
		web.GET("/", h.authView)
		web.GET("/logout", h.authView)

		tasks := web.Group("/tasks")
		{
			tasks.GET("/", h.tasksView)
			tasks.POST("/action", h.actionTaskHandler)
		}
		backups := web.Group("/backups")
		{
			backups.GET("/", h.backupsView)
			backups.POST("/create", h.backupHandler)
			backups.POST("/action", h.actionBackupHandler)
		}
	}
	api := router.Group("/api")
	{
		api.GET("/check")
	}
	return router
}
