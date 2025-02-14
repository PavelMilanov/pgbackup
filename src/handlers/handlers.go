package handlers

import (
	"net/http"
	"time"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	DB   *db.SQLite
	CRON *cron.Cron
	ENV  *config.Env
}

func NewHandler(conn *db.SQLite, scheduler *cron.Cron, env *config.Env) *Handler {
	return &Handler{DB: conn, CRON: scheduler, ENV: env}
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

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{h.ENV.URL},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
	// store := gormsessions.NewStore(h.DB.Sql, true, []byte("mysessions"))

	// router.Use(sessions.Sessions("token", store))

	// router.SetFuncMap(template.FuncMap{"add": func(x, y int) int { return x + y }})
	// router.LoadHTMLGlob("templates/**/*")
	// router.Static("/static/", "./static")

	router.GET("/registration", h.registrationHandler)
	router.POST("/registration", h.registrationHandler)
	router.GET("/login", h.loginHandler)
	router.POST("/login", h.loginHandler)

	v2 := router.Group("/v2/")
	{
		v2.GET("/", h.mainHandler)
	}
	web := router.Group("/")
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
