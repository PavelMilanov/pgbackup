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
	Name     string `form:"Name" binding:"required"`
	Host     string `form:"Host" binding:"required"`
	Port     string `form:"Port" binding:"required"`
	Username string `form:"Username" binding:"required"`
	Password string `form:"Password" binding:"required"`
}

// Форма для сохранения расписания.
type ScheduleForm struct {
	ID        string `form:"ID"`
	Name      string `form:"Name" binding:"required"`
	Frequency string `form:"Frequency" binding:"required"`
	Time      string `form:"Time" binding:"required"`
}
