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
		if check := user.IsRegister(h.DB.Sql); !check {
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
		var data web.RegistrationForm
		if err := c.ShouldBind(&data); err != nil {
			return
		}
		if data.Password != data.ConfirmPassword {
			c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "пароли не совпадают"})
			return
		}
		user := db.User{
			Username: data.Username,
			Password: system.Encrypt((data.Password)),
		}
		check := user.IsRegister(h.DB.Sql)
		if !check {
			c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "пользователь уже зарегистрирован"})
			return
		}
		if err := user.Save(h.DB.Sql); err != nil {
			c.HTML(http.StatusOK, "registration.html", gin.H{"error": err.Error()})
		}
		token := db.Token{UserID: user.ID}
		token.Save(h.DB.Sql)
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}
