package handlers

import (
	"net/http"
	"text/template"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var expectedHost = "localhost:8080"

type Handler struct {
	DB   *gorm.DB
	CRON *cron.Cron
}

func NewHandler(db *gorm.DB, scheduler *cron.Cron) *Handler {
	return &Handler{DB: db, CRON: scheduler}
}

// Основной middleware для авторизации.
func authMiddleware(sql *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := sessions.Default(c)
		data := sessions.Get("token")
		token := db.GetToken(sql, data.(string))
		if ok := token.Validate(); !ok {
			if err := token.Delete(sql); err != nil {
				logrus.Error(err)
			}
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Базовый middleware безопасности.
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
	store := gormsessions.NewStore(h.DB, true, []byte("mysessions"))

	router.Use(baseSecurityMiddleware(expectedHost))
	router.Use(sessions.Sessions("token", store))

	router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")

	router.GET("/registration", h.registrationHandler)
	router.POST("/registration", h.registrationHandler)
	router.GET("/login", h.loginHandler)
	router.POST("/login", h.loginHandler)

	web := router.Group("/", authMiddleware(h.DB))
	{
		web.GET("/logout", h.logoutHandler)
		web.POST("/logout", h.logoutHandler)
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
