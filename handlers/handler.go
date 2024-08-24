package handlers

import (
	"text/template"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	CONFIG *db.Config
}

func NewHandler(db *gorm.DB, config *db.Config) *Handler {
	return &Handler{DB: db, CONFIG: config}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")
	web := router.Group("/")
	{
		web.GET("/", h.authView)
		web.GET("/logout", h.authView)
		web.POST("/home", h.submitLoginForm)
		web.GET("/home", h.homeView)
		web.GET("/bases", h.basesView)

		backups := web.Group("/backups")
		{
			backups.GET("/", h.backupsView)
			backups.POST("/", h.createBackup)
		}
	}
	api := router.Group("/api")
	{
		api.GET("/check")
	}
	return router
}
