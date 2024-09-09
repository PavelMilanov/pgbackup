package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Получает текущий список структур Backup и добавляет новый.
func (model *Backup) сreateBackupData() []Backup {
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

func (model *Backup) createBackupDir(dir string) {
	if err := os.Mkdir(dir, 0755); err != nil {
		if !os.IsExist(err) {
			log.Printf("%s - директория проинициализирована", dir)
		} else {
			log.Printf("%s - директория создана", dir)
		}
	}
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

// // Возвращает список названий бекапов.
// func checkBackup(dir string) []string {
// 	files, err := os.ReadDir(dir)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var backups []string
// 	for _, entry := range files {
// 		backup := strings.Split(entry.Name(), ".")[0] // dev-2024-08024.dump > dev-2024-08024
// 		backups = append(backups, backup)
// 	}
// 	return backups
// }

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
