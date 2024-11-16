package handlers

import (
	"net/http"
	"os"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) loginHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	} else if c.Request.Method == "POST" {
		var data web.LoginForm
		if err := c.ShouldBind(&data); err != nil {
			return
		}
		username := os.Getenv("LOGIN")
		password := os.Getenv("PASSWORD")
		if data.Username != username || data.Password != password {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "неправильные логин или пароль",
			})
			return
		}
		auth := db.Token{}
		if err := auth.Generate(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": err.Error()})
			return
		}
		auth.Save(h.DB)
		c.Redirect(http.StatusFound, "/")
	}
}
