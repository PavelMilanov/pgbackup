package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func loginUser(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func submitLoginForm(c *gin.Context) {
	var formData LoginForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	username := os.Getenv("USER_LOGIN")
	password := os.Getenv("USER_PASSWORD")
	if formData.Username != username && formData.Password != password {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"data": "неправильные логин или пароль",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"auth": true,
	})

}
