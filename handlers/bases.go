package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

// type Base struct {
// 	Name string
// 	Size string
// }

func (h *Handler) basesView(c *gin.Context) {
	dbData := db.GetDBData(*h.CONFIG)
	c.HTML(http.StatusOK, "bases.html", gin.H{
		"databases": dbData,
	})
}
