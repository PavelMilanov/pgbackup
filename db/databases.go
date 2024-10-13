package db

import (
	"fmt"

	"gorm.io/gorm"
)

// Создание базы данных в служебной БД.
func (db Database) Create(sql *gorm.DB) (uint, error) {
	result := sql.Create(&db)
	if result.Error != nil {
		return 0, result.Error
	}
	fmt.Println(db.ID)
	return db.ID, nil
}

func (db Database) Delete(sql *gorm.DB) error {
	return nil
}

func GetAll(sql *gorm.DB) ([]Database, error) {
	var db []Database
	result := sql.Find(&db)
	if result.Error != nil {
		return db, result.Error
	}
	return db, nil
}
