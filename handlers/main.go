package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) mainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", gin.H{})
}
