package connector

// // Возвращает модель Backup по заданным параметрам
// func getBackup(alias, date string) Backup {
// 	backups := GetBackupData()
// 	for _, backup := range backups {
// 		if backup.Alias == alias && backup.Date == date {
// 			return backup
// 		}
// 	}
// 	return Backup{}
// }

// // Получает текущий список структур Backup и удаляет найденный.
// func DeleteBackupData(alias, date string) error {
// 	backup := getBackup(alias, date)
// 	newBackups := []Backup{}
// 	if backup.Alias == alias && backup.Date == date {
// 		fileName := backup.Alias + "-" + backup.Date + ".dump"
// 		filePath := fmt.Sprintf("%s/%s", backup.Directory, fileName)
// 		if err := os.Remove(filePath); err != nil {
// 			logrus.Error(err)
// 			return err
// 		}
// 		log.Printf("%s удален", filePath)
// 	} else {
// 		newBackups = append(newBackups, backup)
// 	}
// 	jsonInfo, err := json.MarshalIndent(newBackups, "", "\t")
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
// 	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
// 		logrus.Error(err)
// 		return err
// 	}
// 	return nil
// }

// // Удаляет все бекапы в указанной директории.
// func deleteBackupsInDir(dir string) {
// 	backups := GetBackupData()
// 	newBackups := []Backup{}
// 	for _, backup := range backups {
// 		if backup.Directory == dir {
// 			fileName := backup.Alias + "-" + backup.Date + ".dump"
// 			filePath := fmt.Sprintf("%s/%s", backup.Directory, fileName)
// 			if err := os.Remove(filePath); err != nil {
// 				logrus.Error(err)
// 			}
// 			logrus.Infof("%s удален", filePath)
// 		} else {
// 			newBackups = append(newBackups, backup)
// 		}
// 	}
// 	jsonInfo, err := json.MarshalIndent(newBackups, "", "\t")
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
// 	if err := os.WriteFile(file, jsonInfo, 0640); err != nil {
// 		logrus.Error(err)
// 	}
// }

// // Парсинт json-файл и возращает список структуры Backup.
// func GetBackupData() []Backup {
// 	var backups []Backup
// 	file := fmt.Sprintf("%s/backups.json", BACKUPDATA_DIR)
// 	jsonInfo, err := os.ReadFile(file)
// 	if err != nil {
// 		os.Create(file)
// 	}
// 	if err := json.Unmarshal(jsonInfo, &backups); err != nil {
// 		return []Backup{}
// 	}
// 	return backups
// }

// Возвращает размер файла бекапа.
// filename - название файла бекапа.
// func (model *Backup) getBackupSize(filename string) string {
// 	command := fmt.Sprintf("du -h %s/%s.dump | awk '{print $1}'", model.Directory, filename)
// 	cmd, err := exec.Command("sh", "-c", command).Output()
// 	if err != nil {
// 		logrus.Error(command)
// 	}
// 	return string(cmd)
// }

// // удаляет старые бекапы согласно расписания.
// func (model *Task) deleteOldBackup() error {
// 	command := fmt.Sprintf("find  %s -name \"*.dump\" -mtime +%s -delete", model.Directory, model.Schedule.Count)
// 	cmd, err := exec.Command("sh", "-c", command).Output()
// 	if err != nil {
// 		logrus.Error(err)
// 		return err
// 	}
// 	logrus.Infoln(cmd)
// 	return nil
// }
