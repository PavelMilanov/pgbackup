package handlers

import (
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/system"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/gin-gonic/gin"
)

func (h *Handler) mainHandler(c *gin.Context) {
	lastBackups := db.GetLastBackups(h.DB.Sql)
	countBackups := db.GetCountBackups(h.DB.Sql)
	countSchedules := db.GetCountSchedules(h.DB.Sql)
	countData := db.CountBackupsStatus(h.DB.Sql)
	storageData := system.GetStorageInfo()
	cpu := system.GetCPUInfo()
	ram := system.GetMemoryInfo()

	count := struct {
		Completed int64
		Failed    int64
		Total     int64
	}{Completed: countData[0], Failed: countData[1], Total: countData[0] + countData[1]}

	storage := struct {
		Used  string
		Total string
	}{Used: storageData[1], Total: storageData[0]}

	system := struct {
		CPU     int
		RAM     int
		Storage string
	}{CPU: cpu, RAM: ram, Storage: storageData[2]}

	c.HTML(http.StatusOK, "main.html", gin.H{
		"header":          "Главная | PgBackup",
		"storage":         storage,
		"count":           count,
		"system":          system,
		"backups":         lastBackups,
		"backups_count":   countBackups,
		"schedules_count": countSchedules,
		"pages": []web.Page{
			{Name: "Главная", URL: "/", IsVisible: true},
			{Name: "Расписание", URL: "/schedule", IsVisible: false},
			{Name: "Базы данных", URL: "/databases", IsVisible: false},
			{Name: "Настройки", URL: "/settings", IsVisible: false},
		}})
}
