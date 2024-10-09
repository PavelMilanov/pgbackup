package connector

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/PavelMilanov/pgbackup/db"
	"github.com/PavelMilanov/pgbackup/web"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
// Directory: dirName,
func (model *Backup) createBackupSQL(cfg Config) (*Backup, error) {
	start := time.Now()
	currTime := start.Format("2006-01-02-15:04") // шаблон GO для формата ГГГГ-мм-дд "2006-01-02 15:04:05" со временем
	backupName := model.Alias + "-" + currTime
	dumpName := model.Directory + "/" + backupName + ".dump"
	command := fmt.Sprintf("export PGPASSWORD=%s && pg_dump -h %s -U %s %s > %s", cfg.Password, cfg.Host, cfg.User, model.Alias, dumpName)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
		return &Backup{}, fmt.Errorf("%s", command)
	}
	timer := time.Since(start).Seconds()
	elapsed := fmt.Sprintf("%.3f сек", timer)
	size := model.getBackupSize(backupName)

	model.Date = currTime
	model.Size = size
	model.LeadTime = elapsed
	model.Status = "завершен"
	model.Dump = dumpName
	return model, nil
}

// Выполение задания восстановления базы данных
func Restore(cfg Config, backup db.Backup) error {
	// backup := getBackup(alias, date)
	// 1. очистить базу данных
	// 2. восстановить из бекапа
	commands := []string{"DROP SCHEMA public CASCADE;", "CREATE SCHEMA public;", "GRANT ALL ON SCHEMA public TO  public;"}
	for _, command := range commands {
		run := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c '%s'", cfg.Password, cfg.Host, cfg.User, backup.Alias, command)
		_, err := exec.Command("sh", "-c", run).Output()
		if err != nil {
			logrus.Error(command)
			return fmt.Errorf("%s-%s %s", backup.Alias, backup.Date, command)
		}
		log.Println(backup.Alias, backup.Date, command)
	}
	// backupName := backup.Alias + "-" + backup.Date
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s < %s", cfg.Password, cfg.Host, cfg.User, backup.Alias, backup.Dump)
	_, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
		return fmt.Errorf("%s-%s %s", backup.Alias, backup.Date, command)
	}
	log.Println(backup.Alias, backup.Date, command)
	return nil
}

// Выполнение бекапов по расписанию.
func CreateCronBackup(scheduler *cron.Cron, cfg Config, sql *gorm.DB, data web.BackupForm) {
	dirName := generateRandomBackupDir()
	createBackupDir(dirName)
	timeToCron := toCron(data.SelectedTime, data.SelectedCron)

	task := db.Task{
		Alias:     data.SelectedDB,
		Directory: dirName,
		Count:     data.SelectedCount,
		Time:      data.SelectedTime,
		Cron:      data.SelectedCron,
	}
	result := sql.Create(&task)
	fmt.Println(result.Error)
	scheduler.AddFunc(timeToCron, func() {
		var backup = Backup{
			Alias:     data.SelectedDB,
			Directory: DEFAULT_BACKUP_DIR,
		}
		newBackup, err := backup.createBackupSQL(cfg)
		if err != nil {
			logrus.Error(err)
		}
		sql.Create(&db.Backup{
			Alias:    newBackup.Alias,
			Date:     newBackup.Date,
			Size:     newBackup.Size,
			LeadTime: newBackup.LeadTime,
			Run:      data.SelectedRun,
			Status:   newBackup.Status,
			Comment:  data.SelectedComment,
			Dump:     newBackup.Dump,
			TaskID:   task.ID,
		})
		logrus.Infof("Создан бекап %s", newBackup)
		// newBackup.createBackupData()
		// task.deleteOldBackup()
	})
	jobs := scheduler.Entries()
	for _, job := range jobs {
		log.Printf("Job ID: %d, Next Run: %s\n", job.ID, job.Next)
	}
}

// Выполнение бекапа вручную.
func CreateManualBackup(cfg Config, sql *gorm.DB, data web.BackupForm) error {
	var defaultTask db.Task
	sql.Where("Alias = ?", "Default").First(&defaultTask)

	var backup = Backup{
		Alias:     data.SelectedDB,
		Directory: DEFAULT_BACKUP_DIR,
	}
	newBackup, err := backup.createBackupSQL(cfg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	sql.Create(&db.Backup{
		Alias:    newBackup.Alias,
		Date:     newBackup.Date,
		Size:     newBackup.Size,
		LeadTime: newBackup.LeadTime,
		Run:      data.SelectedRun,
		Status:   newBackup.Status,
		Comment:  data.SelectedComment,
		Dump:     newBackup.Dump,
		TaskID:   defaultTask.ID,
	})
	logrus.Infof("Создан бекап %s", newBackup)
	return nil
}

// получение размера базы данных по имени.
func getDBSize(cfg Config, dbName string) string {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c \"SELECT pg_size_pretty(pg_database_size('%s'))\"", cfg.Password, cfg.Host, cfg.User, cfg.DBName, dbName)
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
		return ""
	}
	//pg_size_pretty
	//----------------
	//7453 kB
	//(1 row)
	startIndex := 35
	endIndex := len(string(output)) - 10
	size := fmt.Sprint(string(output)[startIndex:endIndex]) // -> 7453 kB
	return size
}

// Получение списка всех баз данных в экземпляре PostgreSQL.
func getDBName(cfg Config) []string {
	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c \"SELECT datname FROM pg_database WHERE datistemplate = false\"", cfg.Password, cfg.Host, cfg.User, cfg.DBName)
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
		return []string{}
	}
	//datname
	//----------
	//postgres
	//dev
	//(2 rows)
	startIndex := 23
	endIndex := len(string(output)) - 11
	data := fmt.Sprint(string(output[startIndex:endIndex]))
	dbList := strings.Split(data, "\n")
	for i, item := range dbList {
		dbList[i] = strings.TrimSpace(item)
	}
	return dbList
}

// Вывод информации обо всех базах данных
func GetDBData(db Config) []PsqlBase {
	var dataBases = []PsqlBase{}
	dbNames := getDBName(db)
	for _, item := range dbNames {
		size := getDBSize(db, item)
		dataBases = append(dataBases, PsqlBase{Name: item, Size: size})
	}
	return dataBases
}
