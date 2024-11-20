package system

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/sirupsen/logrus"
)

// Шифрование строки по алгоритму AES.
func Encrypt(plaintext string) string {
	bc, err := aes.NewCipher(config.AES_KEY)
	if err != nil {
		logrus.Error(err)
	}
	paddedText := pad([]byte(plaintext), aes.BlockSize)
	dst := make([]byte, len(paddedText))
	cipher.NewCBCEncrypter(bc, config.AES_KEY[:aes.BlockSize]).CryptBlocks(dst, paddedText)
	return base64.StdEncoding.EncodeToString(dst)
}

// Дешифрование строки по алгоритму AES.
func Decrypt(ciphertext string) string {
	bc, err := aes.NewCipher(config.AES_KEY)
	if err != nil {
		logrus.Fatal(err)
	}
	ciphertextBytes, _ := base64.StdEncoding.DecodeString(ciphertext)
	res := make([]byte, len(ciphertextBytes))
	cipher.NewCBCDecrypter(bc, config.AES_KEY[:aes.BlockSize]).CryptBlocks(res, ciphertextBytes)
	unpaddedText, err := unpad(res)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(unpaddedText)
}

// Функция добавляет padding по стандарту PKCS#7
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Функция удаляет padding по стандарту PKCS#7
func unpad(src []byte) ([]byte, error) {
	padding := src[len(src)-1]
	length := len(src) - int(padding)
	if length < 0 {
		return nil, errors.New("invalid padding")
	}
	return src[:length], nil
}
