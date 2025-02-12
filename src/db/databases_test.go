package db

import (
	"os"
	"testing"
	"time"

	"github.com/PavelMilanov/pgbackup/system"
	"github.com/robfig/cron/v3"
)

var testConfig = Database{
	Alias:    "test",
	Name:     "test",
	Host:     "localhost",
	Port:     5433,
	Username: "test",
	Password: "test",
}
var location, _ = time.LoadLocation("Europe/Moscow")
var testScheduler = cron.New(cron.WithLocation(location))

var testDb = NewDatabase("test.db", testScheduler)

func TestCheckConnection(t *testing.T) {
	testConfig.checkConnection()
}

func TestGetDBSize(t *testing.T) {
	testConfig.getDBSize()
}

func TestGetDb(t *testing.T) {
	db, err := GetDb(testDb.Sql, 1)
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
