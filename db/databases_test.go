package db

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

var testConfig = Database{
	Name:     "dev",
	Host:     "localhost",
	Port:     5433,
	Username: "test",
	Password: "test",
}
var location, _ = time.LoadLocation("Europe/Moscow")
var testScheduler = cron.New(cron.WithLocation(location))

var testDb = NewDatabase(&SQLite{Name: "test.db"}, testScheduler)

func TestCheckConnection(t *testing.T) {
	if err := testConfig.checkConnection(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetDBSize(t *testing.T) {
	if err := testConfig.getDBSize(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetDb(t *testing.T) {
	db, err := GetDb(testDb, 1)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(db)
}
