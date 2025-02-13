package system

import "testing"

func TestEncrypt(t *testing.T) {
	key := []byte("32-char-key-for-AES-256!")

	testText := []string{"hello world", "world", "hello world !", "hello world", "admin", "admin"}
	for _, text := range testText {
		cryptoStr, err := Encrypt(text, key)
		if err != nil {
			t.Error(err)
		}
		decryptStr, err := Decrypt(cryptoStr, key)
		if err != nil {
			t.Error(err)
		}
		if decryptStr != text {
			t.Errorf("Encrypt: %s != Decrypt: %s", decryptStr, text)
		}
		t.Logf("Text: %s, Encrypted: %s, Decrypted: %s", text, cryptoStr, decryptStr)
	}
}
