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

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// Функция add для шаблонов
func add(x, y int) int {
	return x + y
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")
	web := router.Group("/")
	{
		web.GET("/", h.authView)
		web.GET("/logout", h.authView)
		web.POST("/home", h.submitLoginForm)
		web.GET("/home", h.homeView)
		web.GET("/bases", h.basesView)
		web.GET("/backups", h.backupsView)
	}
	api := router.Group("/api")
	{
		api.GET("/check")
	}
	router.SetFuncMap(template.FuncMap{"add": add})
	return router
}
