package web

type Page struct {
	Header    string
	Name      string
	URL       string
	IsVisible bool
}

type BackupForm struct {
	SelectedDB      string `form:"backupDBName" binding:"required"`
	SelectedRun     string `form:"backupRun" binding:"required"`
	SelectedComment string `form:"backupComment"`
	SelectedCount   string `form:"backupScheduleCount"`
	SelectedTime    string `form:"backupScheduleTime"`
	SelectedCron    string `form:"backupScheduleCron"`
}

// Форма для сохранения базы данных.
type DatabaseForm struct {
	ID       string `form:"ID"`
	Name     string `form:"Name"`
	Host     string `form:"Host"`
	Port     string `form:"Port"`
	Username string `form:"Username"`
	Password string `form:"Password"`
}

// Форма для сохранения расписания.
type ScheduleForm struct {
	ID        string `form:"ID"`
	Name      string `form:"Name"`
	Frequency string `form:"Frequency"`
	Time      string `form:"Time"`
}
