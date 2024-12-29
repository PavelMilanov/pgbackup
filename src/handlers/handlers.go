package handlers

import (
	"net/http"
	"text/template"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	DB   *db.SQLite
	CRON *cron.Cron
}

func NewHandler(conn *db.SQLite, scheduler *cron.Cron) *Handler {
	return &Handler{DB: conn, CRON: scheduler}
}

// Основной middleware для авторизации.
func authMiddleware(sql *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := sessions.Default(c)
		data := sessions.Get("token")
		// если пришел пустое значение сессии.
		if data == nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{})
			c.Abort()
			return
		}
		token := db.GetToken(sql, data.(int))
		// валидириуем сам токен.
		if ok := token.Validate(); !ok {
			if err := token.Refresh(sql); err != nil {
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
		if host == "*" {
			return
		} else if c.Request.Host != host {
			logrus.Debug("Host invalid: ", c.Request.Host)
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
	store := gormsessions.NewStore(h.DB.Sql, true, []byte("mysessions"))

	router.Use(baseSecurityMiddleware(config.HOST))
	router.Use(sessions.Sessions("token", store))

	router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static/", "./static")

	router.GET("/registration", h.registrationHandler)
	router.POST("/registration", h.registrationHandler)
	router.GET("/login", h.loginHandler)
	router.POST("/login", h.loginHandler)

	web := router.Group("/", authMiddleware(h.DB.Sql))
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
			databases.POST("/backup/restore", h.restoreBackupHandler)
			databases.GET("/backup/download", h.downloadBackupHandler)
			databases.GET("/backups", h.getBackupsHandler)
		}
		web.GET("/settings", h.settingsHandler)
		web.POST("/settings", h.settingsHandler)
	}
	api := router.Group("/api")
	{
		api.GET("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
	return router
}
