package db

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
)

var BACKUP_DIR = "dumps"
var BACKUPDATA_DIR = "data"
var DEFAULT_BACKUP_DIR = BACKUP_DIR + "/" + "default_backup"
var BACKUP_RUN = []string{"вручную", "по расписанию"}

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

// Проверка подключения к базе данных
func CheckConnection(cfg Config) string {
	command := fmt.Sprintf("pg_isready -h %s -U %s -d %s -p %d", cfg.Host, cfg.User, cfg.DBName, cfg.portToInt(cfg.Port))
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Println(err)
	}
	return string(cmd)
}

// Выполнение задания бекапа базы данных
// cfg - данные для подключения к PostgreSQL.
// Входная модель:
// Alias:     backupName,
// Comment:   backupComment,
// Directory: dirName,
//
//	Schedule: db.BackupSchedule{
//		Run:   backupRun,
//		Count: backupCount,
//		Time:  backupTime,
//		Cron:  backupCron,
//	}
func (model *Backup) createBackupSQL(cfg Config) (*Backup, error) {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	backupName := model.Alias + "-" + currTime
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s/%s.dump", cfg.Password, cfg.Host, cfg.User, model.Alias, model.Directory, backupName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Println(err, command)
		return &Backup{}, fmt.Errorf("%s", command)
	}
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	size := model.getBackupSize(backupName)

	model.Date = currTime
	model.Size = size
	model.LeadTime = elapsed
	model.Status = "завершен"
	return model, nil
}

// Выполение задания восстановления базы данных
func Restore(cfg Config, alias, date string) error {
	backup := getBackup(alias, date)
	// 1. очистить базу данных
	// 2. восстановить из бекапа
	commands := []string{"DROP SCHEMA public CASCADE;", "CREATE SCHEMA public;", "GRANT ALL ON SCHEMA public TO  public;"}
	for _, command := range commands {
		run := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c '%s'", cfg.Password, cfg.Host, cfg.User, backup.Alias, command)
		_, err := exec.Command("sh", "-c", run).Output()
		if err != nil {
			log.Println(err, command)
			return fmt.Errorf("%s-%s %s", backup.Alias, backup.Date, command)
		}
		log.Println(backup.Alias, backup.Date, command)
	}
	backupName := backup.Alias + "-" + backup.Date
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s < %s/%s.dump", cfg.Password, cfg.Host, cfg.User, backup.Alias, backup.Directory, backupName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Println(err, command)
		return fmt.Errorf("%s-%s %s", backup.Alias, backup.Date, command)
	}
	log.Println(backup.Alias, backup.Date, command)
	return nil
}

// Выполнение бекапов по расписанию.
func (task *Task) CreateCronBackup(scheduler *cron.Cron, cfg Config) {
	cron := task.toCron()
	scheduler.AddFunc(cron, func() {
		var backupModel = Backup{
			Alias:     task.Alias,
			Comment:   task.Comment,
			Directory: task.Directory,
			Schedule: BackupSchedule{
				Run:   task.Schedule.Run,
				Count: task.Schedule.Count,
				Time:  task.Schedule.Time,
				Cron:  task.Schedule.Cron,
			},
		}
		newBackup, err := backupModel.createBackupSQL(cfg)
		if err != nil {
			log.Println(err)
			return
		}
		newBackup.createBackupData()
		task.deleteOldBackup()
	})
	jobs := scheduler.Entries()
	for _, job := range jobs {
		log.Printf("Job ID: %d, Next Run: %s\n", job.ID, job.Next)
	}
}

// Выполнение бекапа вручную.
func (model *Backup) CreateManualBackup(cfg Config) (*Backup, error) {
	newBackup, err := model.createBackupSQL(cfg)
	if err != nil {
		log.Println(err)
		return &Backup{}, err
	}
	newBackup.createBackupData()
	return newBackup, nil
}
