package connector

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Возвращает строку в формате cron для модели Task.
func (task *Task) toCron() string {
	// минуты часы день(*/1 каждый день) * *
	crontime := strings.Split(task.Schedule.Time, ":") // 22:45 => ["22", "45"]
	cron := fmt.Sprintf("%s %s */%s * *", crontime[1], crontime[0], task.Schedule.Cron)
	return cron
}

// Возвращает модель Backup по заданным параметрам
func getBackup(alias, date string) Backup {
	backups := GetBackupData()
	for _, backup := range backups {
		if backup.Alias == alias && backup.Date == date {
			return backup
		}
	}
	return Backup{}
}

// Возвращает модель Task по заданным параметрам
func getTask(alias, dir string) Task {
	tasks := GetTaskData()
	for _, task := range tasks {
		if task.Alias == alias && task.Directory == dir {
			return task
		}
	}
	return Task{}
}

// Получает текущий список структур Backup и добавляет новый.
func (model *Backup) createBackupData() []Backup {
	backups := GetBackupData()
	backups = append(backups, Backup{
		Alias:     model.Alias,
		Date:      model.Date,
		Size:      model.Size,
		LeadTime:  model.LeadTime,
		Status:    model.Status,
		Comment:   model.Comment,
		Directory: model.Directory,
		Schedule:  model.Schedule,
	})
	jsonInfo, err := json.MarshalIndent(backups, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		logrus.Error(err)
		return []Backup{}
	}
	return backups
}

// Получает текущий список структур Backup и удаляет найденный.
func DeleteBackupData(alias, date string) error {
	backup := getBackup(alias, date)
	newBackups := []Backup{}
	if backup.Alias == alias && backup.Date == date {
		fileName := backup.Alias + "-" + backup.Date + ".dump"
		filePath := fmt.Sprintf("%s/%s", backup.Directory, fileName)
		if err := os.Remove(filePath); err != nil {
			logrus.Error(err)
			return err
		}
		log.Printf("%s удален", filePath)
	} else {
		newBackups = append(newBackups, backup)
	}
	jsonInfo, err := json.MarshalIndent(newBackups, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// Удаляет все бекапы в указанной директории.
func deleteBackupsInDir(dir string) {
	backups := GetBackupData()
	newBackups := []Backup{}
	for _, backup := range backups {
		if backup.Directory == dir {
			fileName := backup.Alias + "-" + backup.Date + ".dump"
			filePath := fmt.Sprintf("%s/%s", backup.Directory, fileName)
			if err := os.Remove(filePath); err != nil {
				logrus.Error(err)
			}
			logrus.Infof("%s удален", filePath)
		} else {
			newBackups = append(newBackups, backup)
		}
	}
	jsonInfo, err := json.MarshalIndent(newBackups, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		logrus.Error(err)
	}
}

// Создает указанную директорию.
func СreateBackupDir(dir string) {
	if err := os.Mkdir(dir, 0755); err != nil {
		if !os.IsExist(err) {
			logrus.Infof("%s - директория создана", dir)
		}
	}
}

// Удаляет директорию с бекапами
func DeleteBackupDir(dir string) {
	deleteBackupsInDir(dir)
	os.RemoveAll(dir)
	logrus.Infof("%s - директория удалена", dir)
}

// Генерирует случайную строку из цифр от 0 до 10000.
func GenerateRandomBackupDir() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	dirName := strconv.Itoa(r1.Intn(10000))
	return BACKUP_DIR + "/" + dirName
}

// Парсинт json-файл и возращает список структуры Backup.
func GetBackupData() []Backup {
	var backups []Backup
	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
	jsonInfo, err := os.ReadFile(file)
	if err != nil {
		os.Create(file)
	}
	if err := json.Unmarshal(jsonInfo, &backups); err != nil {
		return []Backup{}
	}
	return backups
}

// Возвращает размер файла бекапа.
// filename - название файла бекапа.
func (model *Backup) getBackupSize(filename string) string {
	command := fmt.Sprintf("du -h %s/%s.dump | awk '{print $1}'", model.Directory, filename)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(command)
	}
	return string(cmd)
}

// Парсинт json-файл и возращает список структуры Task.
func GetTaskData() []Task {
	var tasks []Task
	file := fmt.Sprintf("%s/tasks.json", BACKUPDATA_DIR)
	jsonInfo, err := os.ReadFile(file)
	if err != nil {
		os.Create(file)
	}
	if err := json.Unmarshal(jsonInfo, &tasks); err != nil {
		return []Task{}
	}
	return tasks
}

// Получает текущий список структур Tasks и добавляет новый.
func (model *Task) CreateTaskData() []Task {
	tasks := GetTaskData()
	tasks = append(tasks, Task{
		Alias:     model.Alias,
		Comment:   model.Comment,
		Directory: model.Directory,
		Schedule:  model.Schedule,
	})
	jsonInfo, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	file := fmt.Sprintf("%s/tasks.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		return []Task{}
	}
	return tasks
}

// Получает текущий список структур Tasks и удаляет указанный.
func DeleteTaskData(alias, dir string) error {
	task := getTask(alias, dir)
	tasks := GetTaskData()
	newTasks := []Task{}
	for _, item := range tasks {
		if item.Alias == task.Alias && item.Directory == task.Directory {
			continue
		}
		newTasks = append(newTasks, item)
	}
	jsonInfo, err := json.MarshalIndent(newTasks, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	file := fmt.Sprintf("%s/tasks.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// удаляет старые бекапы согласно расписания.
func (model *Task) deleteOldBackup() error {
	command := fmt.Sprintf("find  %s -name \"*.dump\" -mtime +%s -delete", model.Directory, model.Schedule.Count)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infoln(cmd)
	return nil
}
