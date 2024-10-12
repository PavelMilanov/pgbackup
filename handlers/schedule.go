package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) scheduleHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "schedule.html", gin.H{})
}
