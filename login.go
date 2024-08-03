package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func loginUser(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "pgbackup",
	})
}

func submitLoginForm(c *gin.Context) {
	var formData LoginForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"auth":  true,
		"title": "pgbackup",
	})

}
