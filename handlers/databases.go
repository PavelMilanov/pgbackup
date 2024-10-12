package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) databasesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "databases.html", gin.H{})
}
