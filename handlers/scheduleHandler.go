package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/connector"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) scheduleHandler(c *gin.Context) {
	schedules := connector.GetScheduleAll(h.DB)
	databases := connector.GetDbAll(h.DB)
	c.HTML(http.StatusOK, "schedule.html", gin.H{
		"header":           "Расписание | PgBackup",
		"databases":        databases,
		"backup_frequency": connector.BACKUP_FREQUENCY,
		"schedules":        schedules,
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: true},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

func (h *Handler) scheduleSaveHandler(c *gin.Context) {
	var data web.ScheduleForm
	if err := c.ShouldBind(&data); err != nil {
		logrus.Error(err)
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	config := connector.ScheduleConfig{
		DbName:    data.Name,
		Frequency: data.Frequency,
		Time:      data.Time,
	}
	config.Save(h.DB, h.CRON)
	c.Redirect(http.StatusFound, "/schedule/")
}
