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
	schedules := db.GetSchedulesAll(h.DB.Sql)
	databases := db.GetDbAll(h.DB.Sql)
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
		return
	}
	id, _ := strconv.Atoi(data.Name)
	cfg := db.Schedule{
		DatabaseID: id,
		Frequency:  data.Frequency,
		Time:       data.Time,
	}
	if err := cfg.Save(h.DB, h.CRON); err != nil {
		schedules := db.GetSchedulesAll(h.DB.Sql)
		databases := db.GetDbAll(h.DB.Sql)
		c.HTML(http.StatusOK, "schedule.html", gin.H{
			"header":           "Расписание | PgBackup",
			"databases":        databases,
			"backup_frequency": config.BACKUP_FREQUENCY,
			"schedules":        schedules,
			"notification": web.Notify{
				Message: err.Error(),
				Type:    config.NOTIFY_STATUS["ошибка"],
			},
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: true},
				{Name: "Базы данных", URL: "/databases", IsVisible: false},
				{Name: "Настройки", URL: "/settings", IsVisible: false},
			}})
	}
	schedules := db.GetSchedulesAll(h.DB.Sql)
	databases := db.GetDbAll(h.DB.Sql)
	c.HTML(http.StatusOK, "schedule.html", gin.H{
		"header":           "Расписание | PgBackup",
		"databases":        databases,
		"backup_frequency": config.BACKUP_FREQUENCY,
		"schedules":        schedules,
		"notification": web.Notify{
			Message: "Расписание добавлено!",
			Type:    config.NOTIFY_STATUS["инфо"],
		},
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: true},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}

// Handler для удаления расписания.
func (h *Handler) scheduleDeleteHandler(c *gin.Context) {
	var data web.ScheduleForm
	if err := c.ShouldBind(&data); err != nil {
		return
	}
	id, _ := strconv.Atoi(data.ID)
	cfg := db.Schedule{
		ID: id,
	}
	if err := cfg.Delete(h.DB.Sql, h.CRON); err != nil {
		schedules := db.GetSchedulesAll(h.DB.Sql)
		databases := db.GetDbAll(h.DB.Sql)
		c.HTML(http.StatusOK, "schedule.html", gin.H{
			"header":           "Расписание | PgBackup",
			"databases":        databases,
			"backup_frequency": config.BACKUP_FREQUENCY,
			"schedules":        schedules,
			"notification": web.Notify{
				Message: err.Error(),
				Type:    config.NOTIFY_STATUS["ошибка"],
			},
			"pages": []web.Page{
				{Name: "Главная", URL: "/", IsVisible: false},
				{Name: "Расписание", URL: "/schedule", IsVisible: true},
				{Name: "Базы данных", URL: "/databases", IsVisible: false},
				{Name: "Настройки", URL: "/settings", IsVisible: false},
			}})
	}
	schedules := db.GetSchedulesAll(h.DB.Sql)
	databases := db.GetDbAll(h.DB.Sql)
	c.HTML(http.StatusOK, "schedule.html", gin.H{
		"header":           "Расписание | PgBackup",
		"databases":        databases,
		"backup_frequency": config.BACKUP_FREQUENCY,
		"schedules":        schedules,
		"notification": web.Notify{
			Message: "Расписание удалено!",
			Type:    config.NOTIFY_STATUS["инфо"],
		},
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: false},
			{Name: "Расписание", URL: "/schedule", IsVisible: true},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
