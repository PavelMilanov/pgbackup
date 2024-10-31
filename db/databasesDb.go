package db

import (
	"gorm.io/gorm"
)

// Создание базы данных в служебной БД.
func (db Database) Create(sql *gorm.DB) (uint, error) {
	result := sql.Create(&db)
	if result.Error != nil {
		return 0, result.Error
	}
	return db.ID, nil
}

func (db Database) Delete(sql *gorm.DB) error {
	return nil
}

// Получение всех баз данных из таблицы Databases
func GetDbAll(sql *gorm.DB) ([]Database, error) {
	var db []Database
	result := sql.Find(&db)
	if result.Error != nil {
		return db, result.Error
	}
	return db, nil
}

// Получение базы данных по имени из таблицы Databases
func GetDb(sql *gorm.DB, id string) (Database, error) {
	var db Database
	result := sql.First(&db, "ID = ?", id)
	if result.Error != nil {
		return db, result.Error
	}
	return db, nil
}
