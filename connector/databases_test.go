package connector

import (
	"testing"
)

func TestCheckConnection(t *testing.T) {
	config := DBConfig{
		Name:     "dev",
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "admin",
	}
	if err := config.checkConnection(); err != nil {
		t.Log(err)
	}
}
