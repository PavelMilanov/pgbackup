package db

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Вывод информации обо всех базах данных
func GetDBData(db *gorm.DB) []PsqlBase {
	var dataBases = []PsqlBase{}
	dbNames := getDBName(db)
	for _, item := range dbNames {
		size := getDBSize(db, item)
		dataBases = append(dataBases, PsqlBase{Name: item, Size: size})
	}
	return dataBases
}

// Получение списка всех баз данных в экземпляре PostgreSQL.
func getDBName(db *gorm.DB) []string {
	var databases []struct {
		Name string `gorm:"column:datname"` // Alias для столбца `datname`
	}
	if err := db.Raw("SELECT datname FROM pg_database WHERE datistemplate = false").Scan(&databases).Error; err != nil {
		log.Fatalf("failed to get databases: %v", err)
	}
	var dbList []string
	for _, db := range databases {
		dbList = append(dbList, db.Name)
	}
	return dbList
}

// получение размера базы данных по имени.
func getDBSize(db *gorm.DB, dbName string) string {
	var size string
	query := fmt.Sprintf("SELECT pg_size_pretty(pg_database_size('%s'))", dbName)
	if err := db.Raw(query).Scan(&size).Error; err != nil {
		log.Fatalf("failed to query database size: %v", err)
	}
	return size
}
