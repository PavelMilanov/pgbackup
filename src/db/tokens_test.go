package db

import (
	"testing"

	"github.com/PavelMilanov/pgbackup/system"
)

func TestGenerate(t *testing.T) {
	token := generateToken(1)
	t.Log(token)
}

func TestRegistration(t *testing.T) {
	var user = User{Username: "user", Password: system.Encrypt("password")}
	if err := user.Save(testDb.Sql); err != nil {
		t.Log(err)
	}
	token := Token{UserID: user.ID}
	token.Save(testDb.Sql)
	t.Log(token)
}

func TestValidation(t *testing.T) {
	var user = User{Username: "user", Password: system.Encrypt("password")}
	token := user.GetToken(testDb.Sql)

	if valid := token.Validate(); !valid {
		t.Log("Токен просрочен")
	}
}

func TestCreateToken(t *testing.T) {
	var token = Token{UserID: 1}
	token.Save(testDb.Sql)
	t.Logf("Создан токен %+v", token)
}

func TestRefreshToken(t *testing.T) {
	var token Token
	testDb.Sql.First(&token)
	t.Logf("Старый токен: %s", token.Hash)
	token.Refresh(testDb.Sql)
	t.Logf("Новый токен: %s", token.Hash)
}
