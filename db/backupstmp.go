package db

import (
	"os"

	"gorm.io/gorm"
)

// Get возвращает бекап из БД по его ID.
func (backup *Backup) Get(db *gorm.DB, id string) error {
	result := db.Where("ID = ?", id).First(&backup)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete удаляет бекап из БД и удаляет файл с диска.
func (backup *Backup) Delete(db *gorm.DB) error {
	result := db.Delete(&backup)
	if result.Error != nil {
		return result.Error
	}
	if err := os.Remove(backup.Dump); err != nil {
		return err
	}
	return nil
}
