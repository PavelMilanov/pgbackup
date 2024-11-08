package handlers

import (
	"net/http"
	"strconv"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

// Handler для главной страницы с расписаниями.
func (h *Handler) scheduleHandler(c *gin.Context) {
	schedules := db.GetSchedules(h.DB)
	databases := db.GetDbAll(h.DB)
	c.HTML(http.StatusOK, "schedule.html", gin.H{
		"header":           "Расписание | PgBackup",
		"databases":        databases,
		"backup_frequency": config.BACKUP_FREQUENCY,
		"schedules":        schedules,
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: true},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

// Handler для сохранения расписания.
func (h *Handler) scheduleSaveHandler(c *gin.Context) {
	var data web.ScheduleForm
	if err := c.ShouldBind(&data); err != nil {
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(data.Name)
	config := db.Schedule{
		DatabaseID: id,
		Frequency:  data.Frequency,
		Time:       data.Time,
	}
	if err := config.Save(h.DB, h.CRON); err != nil {
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/schedule/")
}

// Handler для удаления расписания.
func (h *Handler) scheduleDeleteHandler(c *gin.Context) {
	var data web.ScheduleForm
	if err := c.ShouldBind(&data); err != nil {
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(data.ID)
	config := db.Schedule{
		ID: id,
	}
	if err := config.Delete(h.DB, h.CRON); err != nil {
		//c.HTML(http.StatusBadRequest, "databases.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/schedule/")
}
