package connector

var BACKUP_DIR = "dumps"

var BACKUPDATA_DIR = "data"
var DEFAULT_BACKUP_DIR = BACKUP_DIR + "/" + "default_backup"
var BACKUP_RUN = []string{"вручную", "по расписанию"}

// // Модель расписания.
// type Task struct {
// 	Alias     string
// 	Comment   string
// 	Directory string
// }

// Модель бекапа.
type Backup struct {
	Alias     string
	Date      string
	Size      string
	LeadTime  string
	Status    string
	Directory string
	Dump      string
}

// Модель Базы данных psql.
type PsqlBase struct {
	Name string
	Size string
}
