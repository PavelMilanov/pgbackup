package db

import (
	"gorm.io/gorm"
)

// Создание расписания в таблице Schedule.
func ScheduleCreate(sql *gorm.DB, db Schedule) (uint, error) {
	result := sql.Create(&db)
	if result.Error != nil {
		return 0, result.Error
	}
	return db.ID, nil
}

// Обновление расписания в таблице Schedule.
func ScheduleUpdate(sql *gorm.DB, db Schedule) (Schedule, error) {
	var model Schedule
	result := sql.First(&model)
	if result.Error != nil {
		return model, result.Error
	}
	model.Frequency = db.Frequency
	model.Time = db.Time
	sql.Save(&model)
	return model, nil
}

// Удаление расписания в таблице Schedule.
func ScheduleDelete(sql *gorm.DB, db Schedule) error {
	result := sql.Delete(&db)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
