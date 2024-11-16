package web

type Page struct {
	Header    string
	Name      string
	URL       string
	IsVisible bool
}

type BackupForm struct {
	ID string `form:"ID" binding:"required"`
}

// Форма для сохранения базы данных.
type DatabaseForm struct {
	ID       string `form:"ID"`
	Alias    string `form:"Alias"`
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

type LoginForm struct {
	Username string `form:"Username" binding:"required"`
	Password string `form:"Password" binding:"required"`
}
