package db

import (
	"os"
	"testing"
	"time"

	"github.com/PavelMilanov/pgbackup/system"
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

func TestEncrypAndDecrypt(t *testing.T) {
	os.Setenv("AES_KEY", "key3456789012345")
	text := "hello world"
	encrypted := system.Encrypt("hello world")
	decrypted := system.Decrypt(encrypted)
	if text != decrypted {
		t.Fatalf("%s не равен %s", text, decrypted)
	}
}
