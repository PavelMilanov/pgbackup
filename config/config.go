package config

var DURATION = 3

var BACKUP_DIR = "dumps"

var BACKUP_FREQUENCY = map[string]string{
	"ежедневно":   "ежедневно",
	"еженедельно": "еженедельно",
}

var SCHEDULE_STATUS = map[string]string{
	"активно": "активно",
	"вручную": "вручную",
}
