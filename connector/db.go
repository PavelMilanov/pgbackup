package connector

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

var BACKUP_DIR = "dumps"
var BACKUPDATA_DIR = "data"
var DEFAULT_BACKUP_DIR = BACKUP_DIR + "/" + "default_backup"
var BACKUP_RUN = []string{"вручную", "по расписанию"}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (cfg *Config) portToInt(port string) int {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		logrus.Fatal(err)
	}
	return intPort
}

// Модель расписания.
type Task struct {
	Alias     string
	Comment   string
	Directory string
	Schedule  BackupSchedule
}

// Модель бекапа.
type Backup struct {
	Alias     string
	Date      string
	Size      string
	LeadTime  string
	Status    string
	Comment   string
	Directory string
	Schedule  BackupSchedule
}

// Модель для расписания бекапа в формате cron.
type BackupSchedule struct {
	Run   string
	Count string
	Time  string
	Cron  string
}

// Модель Базы данных psql.
type PsqlBase struct {
	Name string
	Size string
}
