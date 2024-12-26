package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) logoutHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "logout.html", gin.H{})
	} else if c.Request.Method == "POST" {
		sessions := sessions.Default(c)
		sessions.Set("token", nil)
		sessions.Save()
		c.Redirect(http.StatusFound, "/login")
	}
}
