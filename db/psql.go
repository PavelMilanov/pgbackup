package db

// Выполение задания восстановления базы данных
// func Restore(cfg DBConfig, backup db.Backup) error {
// 	// backup := getBackup(alias, date)
// 	// 1. очистить базу данных
// 	// 2. восстановить из бекапа
// 	commands := []string{"DROP SCHEMA public CASCADE;", "CREATE SCHEMA public;", "GRANT ALL ON SCHEMA public TO  public;"}
// 	for _, command := range commands {
// 		run := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c '%s'", cfg.Password, cfg.Host, cfg.User, backup.Alias, command)
// 		_, err := exec.Command("sh", "-c", run).Output()
// 		if err != nil {
// 			logrus.Error(command)
// 			return fmt.Errorf("%s-%s %s", backup.Alias, backup.Date, command)
// 		}
// 		log.Println(backup.Alias, backup.Date, command)
// 	}
// 	// backupName := backup.Alias + "-" + backup.Date
// 	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s < %s", cfg.Password, cfg.Host, cfg.User, backup.Alias, backup.Dump)
// 	_, err := exec.Command("sh", "-c", command).Output()
// 	if err != nil {
// 		logrus.Error(command)
// 		return fmt.Errorf("%s-%s %s", backup.Alias, backup.Date, command)
// 	}
// 	log.Println(backup.Alias, backup.Date, command)
// 	return nil
// }

// Выполнение бекапов по расписанию.
// func CreateCronBackup(scheduler *cron.Cron, cfg DBConfig, sql *gorm.DB, data web.BackupForm) {
// 	dirName := generateRandomBackupDir()
// 	createBackupDir(dirName)
// 	timeToCron := toCron(data.SelectedTime, data.SelectedCron)

// 	task := db.Task{
// 		Alias:     data.SelectedDB,
// 		Directory: dirName,
// 		Count:     data.SelectedCount,
// 		Time:      data.SelectedTime,
// 		Cron:      data.SelectedCron,
// 	}
// 	result := sql.Create(&task)
// 	if result.Error != nil {
// 		logrus.Error(result.Error)
// 		return
// 	}
// 	scheduler.AddFunc(timeToCron, func() {
// 		var backup = Backup{
// 			Alias:     data.SelectedDB,
// 			Directory: DEFAULT_BACKUP_DIR,
// 		}
// 		newBackup, err := backup.createBackupSQL(cfg)
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 		sql.Create(&db.Backup{
// 			Alias:    newBackup.Alias,
// 			Date:     newBackup.Date,
// 			Size:     newBackup.Size,
// 			LeadTime: newBackup.LeadTime,
// 			Run:      data.SelectedRun,
// 			Status:   newBackup.Status,
// 			Comment:  data.SelectedComment,
// 			Dump:     newBackup.Dump,
// 			TaskID:   task.ID,
// 		})
// 		logrus.Infof("Создан бекап %s", newBackup)
// 		// newBackup.createBackupData()
// 		// task.deleteOldBackup()
// 	})
// 	jobs := scheduler.Entries()
// 	for _, job := range jobs {
// 		log.Printf("Job ID: %d, Next Run: %s\n", job.ID, job.Next)
// 	}
// }

// // Выполнение бекапа вручную.
// func CreateManualBackup(cfg DBConfig, sql *gorm.DB, data web.BackupForm) error {
// 	var defaultTask db.Task
// 	sql.Where("Alias = ?", "Default").First(&defaultTask)

// 	var backup = Backup{
// 		Alias:     data.SelectedDB,
// 		Directory: DEFAULT_BACKUP_DIR,
// 	}
// 	newBackup, err := backup.createBackupSQL(cfg)
// 	if err != nil {
// 		logrus.Error(err)
// 		return err
// 	}
// 	model := db.Backup{Alias: newBackup.Alias,
// 		Date:     newBackup.Date,
// 		Size:     newBackup.Size,
// 		LeadTime: newBackup.LeadTime,
// 		Run:      data.SelectedRun,
// 		Status:   newBackup.Status,
// 		Comment:  data.SelectedComment,
// 		Dump:     newBackup.Dump,
// 		TaskID:   defaultTask.ID}
// 	if err := model.Create(sql); err != nil {
// 		logrus.Error(err)
// 		return err
// 	}
// 	logrus.Infof("Создан бекап %v", model)
// 	return nil
// }

// Получение списка всех баз данных в экземпляре PostgreSQL.
// func getDBName(cfg DBConfig) []string {
// 	command := fmt.Sprintf("export PGPASSWORD=%s && psql -h %s -U %s %s -c \"SELECT datname FROM pg_database WHERE datistemplate = false\"", cfg.Password, cfg.Host, cfg.User, cfg.Name)
// 	output, err := exec.Command("sh", "-c", command).Output()
// 	if err != nil {
// 		logrus.Error(command)
// 		return []string{}
// 	}
// 	//datname
// 	//----------
// 	//postgres
// 	//dev
// 	//(2 rows)
// 	startIndex := 23
// 	endIndex := len(string(output)) - 11
// 	data := fmt.Sprint(string(output[startIndex:endIndex]))
// 	dbList := strings.Split(data, "\n")
// 	for i, item := range dbList {
// 		dbList[i] = strings.TrimSpace(item)
// 	}
// 	return dbList
// }

// Вывод информации обо всех базах данных
// func GetDBData(db DBConfig) []PsqlBase {
// 	var dataBases = []PsqlBase{}
// 	dbNames := getDBName(db)
// 	for _, item := range dbNames {
// 		size := getDBSize(db, item)
// 		dataBases = append(dataBases, PsqlBase{Name: item, Size: size})
// 	}
// 	return dataBases
// }

// func GetDBList(sql *gorm.DB) []db.Database {
// 	var result []db.Database
// 	if err := db.GetAll(sql, result); err != nil {
// 		logrus.Error(err)
// 	}
// 	return result
// }
