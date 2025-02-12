package system

import (
	"github.com/AlexanderGrom/componenta/crypt"
	"github.com/PavelMilanov/pgbackup/config"
	"github.com/sirupsen/logrus"
)

// Шифрование строки по алгоритму AES.
func Encrypt(plaintext string) string {
	c, err := crypt.Encrypt(plaintext, string(config.JWT_KEY))
	if err != nil {
		logrus.Fatal("Encrypt:", err)
	}

	return c
}

// Дешифрование строки по алгоритму AES.
func Decrypt(ciphertext string) string {
	s, err := crypt.Decrypt(ciphertext, string(config.JWT_KEY))
	if err != nil {
		logrus.Fatal("Decrypt:", err)
	}
	return s
}
