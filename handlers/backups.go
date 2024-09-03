package handlers

import (
	"log"
	"net/http"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/gin-gonic/gin"
)

type Data struct {
	SelectedDB    string `json:"db"`
	SelectedRun   string `json:"run"`
	SelectedCount string `json:"count"`
	SelectedTime  string `json:"time"`
	SelectedCron  string `json:"cron"`
}

func (h *Handler) backupsView(c *gin.Context) {
	dbInfo := db.GetDBData(h.DB)
	backupsInfo := db.GetBackupData(db.BACKUPDATA_DIR)
	c.HTML(http.StatusOK, "backups.html", gin.H{
		"databases": dbInfo,
		"backups":   backupsInfo,
		"run":       db.BACKUP_RUN,
	})
}

func (h *Handler) backupHandler(c *gin.Context) {
	var json Data
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dbname := json.SelectedDB
	backupRun := json.SelectedRun
	backupCount := json.SelectedCount
	backupTime := json.SelectedTime
	backupCron := json.SelectedCron

	log.Println(dbname, backupRun, backupCount, backupTime, backupCron)
	if backupRun == db.BACKUP_RUN[1] && (backupCount == "" || backupCron == "" || backupTime == "") {
		c.JSON(http.StatusOK, gin.H{
			"error": "расписание не может быть пустым",
		})
		return
	}

	// curBackups := db.CheckBackup(db.BACKUP_DIR)
	// for _, item := range curBackups {
	// 	if strings.Contains(item, backupName) {
	// 		errStr := fmt.Sprintf("Бекап %s уже существует!", backupName)
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"error": errStr,
	// 		})
	// 		return
	// 	}
	// }
	newBackup, err := db.CreateBackup(*h.CONFIG, dbname, backupRun, backupCount, backupTime, backupCron)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	// newBackup := db.Backup{
	// 	Alias:    dbname,
	// 	Date:     currTime,
	// 	Size:     backupData["size"],
	// 	LeadTime: backupData["time"],
	// 	Status:   backupData["status"],
	// 	Schedule: db.BackupSchedule{
	// 		Run:   backupRun,
	// 		Count: backupCount,
	// 		Time:  backupTime,
	// 		Cron:  backupCron,
	// 	},
	// }

	// backupsInfo := db.CreateBackupData(&newBackup, db.BACKUPDATA_DIR)
	// log.Println(backupsInfo)
	c.JSON(http.StatusOK, gin.H{
		"message": newBackup,
	})
}
