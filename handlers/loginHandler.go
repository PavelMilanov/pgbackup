package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) loginHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	} else if c.Request.Method == "POST" {
		sessions := sessions.Default(c)
		var data web.LoginForm
		if err := c.ShouldBind(&data); err != nil {
			return
		}
		user := db.User{
			Username: data.Username,
			Password: system.Encrypt((data.Password)),
		}
		if check := user.IsRegister(h.DB); !check {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "неправильные логин или пароль",
			})
			return
		}
		sessions.Set("token", user.ID)
		sessions.Save()
		c.Redirect(http.StatusFound, "/")
	}
}

func (h *Handler) registrationHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "registration.html", gin.H{})
	} else if c.Request.Method == "POST" {
		var data web.LoginForm
		if err := c.ShouldBind(&data); err != nil {
			return
		}
		user := db.User{
			Username: data.Username,
			Password: system.Encrypt((data.Password)),
		}
		if err := user.Save(h.DB); err != nil {
			c.HTML(http.StatusOK, "registration.html", gin.H{"error": err.Error()})
		}
		token := db.Token{UserID: user.ID}
		token.Save(h.DB)
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}
