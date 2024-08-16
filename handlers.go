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

func authView(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func submitLoginForm(c *gin.Context) {
	var formData LoginForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err.Error()})
		return
	}
	username := os.Getenv("USER_LOGIN")
	password := os.Getenv("USER_PASSWORD")
	if formData.Username != username || formData.Password != password {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"data": "неправильные логин или пароль",
		})
		return
	}
	db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	databases := db.getDBs()
	c.HTML(http.StatusOK, "home.html", gin.H{"databases": databases})
}

func homeView(c *gin.Context) {
	db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	databases := db.getDBs()
	c.HTML(http.StatusOK, "home.html", gin.H{"databases": databases})
}

func basesView(c *gin.Context) {
	db := Database{
		User:     "admin",
		Password: "admin",
		Dbname:   "dev",
		Port:     5432,
	}
	databases := db.getDBs()
	c.HTML(http.StatusOK, "bases.html", gin.H{
		"databases": databases,
	})
}

func backupsView(c *gin.Context) {
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"data": nil,
	})
}
