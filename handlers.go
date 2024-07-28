package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func loginUser(c *gin.Context) {
	login := os.Getenv("USER_LOGIN")
	password := os.Getenv("USER_PASSWORD")
	var req Login
	c.Bind(&req)
	if req.Username != login || req.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "authorized"})
}
