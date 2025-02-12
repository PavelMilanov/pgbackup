package system

import "testing"

func TestEncrypt(t *testing.T) {
	str := "hello world"
	cryptoStr := Encrypt(str)
	decryptStr := Decrypt(cryptoStr)
	t.Logf("Encrypt: %s != Decrypt: %s", str, decryptStr)
	if decryptStr != str {
		t.Errorf("Encrypt: %s != Decrypt: %s", decryptStr, str)
		return
	}
}
