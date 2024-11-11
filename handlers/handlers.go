package handlers

import (
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	cron "github.com/robfig/cron/v3"
)

type Handler struct {
	DB   *gorm.DB
	CRON *cron.Cron
}

func NewHandler(db *gorm.DB, scheduler *cron.Cron) *Handler {
	return &Handler{DB: db, CRON: scheduler}
}

func authMiddleware(c *gin.Context) {
	// fmt.Print("ping!")
	c.Next()
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()

	// gin.SetMode(gin.ReleaseMode)
	router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")
	web := router.Group("/")
	{
		web.GET("/login", h.loginHandler)
		web.GET("/", h.mainHandler)
		schedule := web.Group("/schedule")
		{
			schedule.GET("/", h.scheduleHandler)
			schedule.POST("/save", h.scheduleSaveHandler)
			schedule.POST("/delete", h.scheduleDeleteHandler)
		}
		databases := web.Group("/databases")
		{
			databases.GET("/", h.databasesHandler)
			databases.POST("/save", h.databaseSaveHandler)
			databases.POST("/delete", h.databaseDeleteHandler)
			databases.POST("/backup/create", h.createBackupHandler)
			databases.POST("/backup/delete", h.deleteBackupHandler)
			databases.GET("/backup/download", h.downloadBackupHandler)
			databases.GET("/backups", h.getBackupsHandler)
		}
		web.GET("/settings", h.settingsHandler)
		web.GET("/logout", h.logoutHandler)
	}
	api := router.Group("/api")
	{
		api.GET("/check")
	}
	return router
}
