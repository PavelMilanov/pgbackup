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
