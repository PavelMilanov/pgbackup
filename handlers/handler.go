package handlers

import (
	"text/template"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	cron "github.com/robfig/cron/v3"
)

type Handler struct {
	DB     *gorm.DB
	CONFIG *db.Config
	CRON   *cron.Cron
}

func NewHandler(db *gorm.DB, config *db.Config, scheduler *cron.Cron) *Handler {
	return &Handler{DB: db, CONFIG: config, CRON: scheduler}
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
		web.POST("/bases", h.submitLoginForm)
		web.GET("/bases", h.basesView)
		web.GET("/tasks", h.tasksView)

		backups := web.Group("/backups")
		{
			backups.GET("/", h.backupsView)
			backups.POST("/create", h.backupHandler)
			backups.GET("/download/:alias/:date", h.downloadBackupHandler)
			backups.POST("/delete/:alias/:date", h.deleteBackupHandler)
		}
	}
	api := router.Group("/api")
	{
		api.GET("/check")
	}
	return router
}
