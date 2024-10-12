package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) settingsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", gin.H{})
}
