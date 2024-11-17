package handlers

import (
	"net/http"
	"text/template"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	cron "github.com/robfig/cron/v3"
)

var expectedHost = "localhost:8080"

type Handler struct {
	DB   *gorm.DB
	CRON *cron.Cron
}

func NewHandler(db *gorm.DB, scheduler *cron.Cron) *Handler {
	return &Handler{DB: db, CRON: scheduler}
}

func authMiddleware(sql *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := db.GetToken(sql)
		if ok := token.Validate(); !ok {
			if err := token.Delete(sql); err != nil {
				logrus.Error(err)
			}
			c.HTML(http.StatusOK, "login.html", gin.H{})
			c.Abort()
		}
		c.Next()
	}
}

func baseSecurityMiddleware(host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host != host {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()
	router.Use(baseSecurityMiddleware(expectedHost))
	router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")
	router.GET("/login", h.loginHandler)
	router.POST("/login", h.loginHandler)
	router.GET("/logout", h.logoutHandler, authMiddleware(h.DB))
	router.POST("/logout", h.logoutHandler)
	web := router.Group("/", authMiddleware(h.DB))
	{
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
	}
	api := router.Group("/api")
	{
		api.GET("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
	return router
}
