package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func loginUser(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Hello world",
	})
}
