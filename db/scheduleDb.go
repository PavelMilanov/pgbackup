package db

import (
	"gorm.io/gorm"
)

// Создание расписания в таблице Schedule.
func (db Schedule) Create(sql *gorm.DB) (uint, error) {
	result := sql.Create(&db)
	if result.Error != nil {
		return 0, result.Error
	}
	return db.ID, nil
}

// Получение всех расписаний c ссылкой на базу данных из таблицы Schedule.
func GetScheduleAll(sql *gorm.DB) ([]Schedule, error) {
	var schedule []Schedule
	result := sql.Find(&schedule)
	if result.Error != nil {
		return schedule, result.Error
	}
	return schedule, nil
}
