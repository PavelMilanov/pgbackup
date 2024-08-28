package db

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Получает текущий список структур Backup и добавляет новый.
func CreateBackupData(backup *Backup, dir string) []Backup {
	backups := GetBackupData(dir)
	backups = append(backups, Backup{
		Alias:    backup.Alias,
		Date:     backup.Date,
		Size:     backup.Size,
		LeadTime: backup.LeadTime,
		Status:   backup.Status,
		Run:      backup.Run,
	})
	jsonInfo, err := json.MarshalIndent(backups, "", "\t")
	if err != nil {
		fmt.Println("Ошибка записи данных:", err)
	}
	file := fmt.Sprintf("%s/backups.json", dir)
	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
		return backups
	}
	return backups
}

// Парсинт json-файл и возращает список структуры Backup.
func GetBackupData(dir string) []Backup {
	var backups []Backup
	file := fmt.Sprintf("%s/backups.json", dir)
	jsonInfo, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(jsonInfo, &backups); err != nil {
		fmt.Println(err)
		return backups
	}
	return backups
}

// Возвращает список названий бекапов.
func CheckBackup(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var backups []string
	for _, entry := range files {
		backup := strings.Split(entry.Name(), ".")[0] // dev-2024-08024.dump > dev-2024-08024
		backups = append(backups, backup)
	}
	return backups
}

func GetBackupSize(dir string, filename string) string {
	command := fmt.Sprintf("du -h %s/%s.dump | awk '{print $1}'", dir, filename)
	cmd, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(command, cmd)
	return string(cmd)
}
