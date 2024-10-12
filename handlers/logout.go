package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) logoutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "logout.html", gin.H{})
}
