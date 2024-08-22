package db

import (
	"log"

	"gorm.io/gorm"
)

// Получение списка всех баз данных в экземпляре PostgreSQL.
func GetDBInfo(db *gorm.DB) []string {
	var databases []struct {
		Name string `gorm:"column:datname"` // Alias для столбца `datname`
	}
	if err := db.Raw("SELECT datname FROM pg_database WHERE datistemplate = false").Scan(&databases).Error; err != nil {
		log.Fatalf("failed to get databases: %v", err)
	}
	var dbList []string
	// Печать имен баз данных
	for _, db := range databases {
		dbList = append(dbList, db.Name)
	}
	return dbList
}
