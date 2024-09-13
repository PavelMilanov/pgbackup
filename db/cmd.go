package db

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
)

// Возвращает строку в формате cron для модели Task.
func (task *Task) ToCron() string {
	// минуты часы день(*/1 каждый день) * *
	crontime := strings.Split(task.Schedule.Time, ":") // 22:45 => ["22", "45"]
	cron := fmt.Sprintf("%s %s */%s * *", crontime[1], crontime[0], task.Schedule.Cron)
	return cron
}

// Получает текущий список структур Backup и добавляет новый.
func CreateBackupData(model *Backup) []Backup {
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
		log.Println("Ошибка записи данных:", err)
	}
	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		return []Backup{}
	}
	return backups
}

// Создает указанную директорию.
func СreateBackupDir(dir string) {
	if err := os.Mkdir(dir, 0755); err != nil {
		if !os.IsExist(err) {
			log.Printf("%s - директория проинициализирована", dir)
		} else {
			log.Printf("%s - директория создана", dir)
		}
	}
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
		panic(err)
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
func CreateTaskData(model *Task) []Task {
	tasks := GetTaskData()
	tasks = append(tasks, Task{
		Alias:     model.Alias,
		Comment:   model.Comment,
		Directory: model.Directory,
		Schedule:  model.Schedule,
	})
	jsonInfo, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		log.Println("Ошибка записи данных:", err)
	}
	file := fmt.Sprintf("%s/tasks.json", BACKUPDATA_DIR)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		return []Task{}
	}
	return tasks
}
