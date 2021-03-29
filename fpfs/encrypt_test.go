package fpfs

import (
	"testing"
)

func TestEncrypt_RoundTrip(t *testing.T) {
	key := []byte("pass")
	data := "Hello World"
	encrypted := Encrypt(key, data)
	decrypted := Decrypt(key, encrypted)
	if data != decrypted {
		t.Errorf("Unexpected decrypted value. Got: %q, Want: %q", decrypted, data)
	}
}

func TestDecrypt_Empty(t *testing.T) {
	if out := Decrypt([]byte("key"), ""); out != "" {
		t.Errorf("Unexpected output. Got: %q, Want: %q", out, "")
	}
}

func TestEncrypt_IncorrectPass(t *testing.T) {
	key := []byte("pass")
	data := "Hello World"
	encrypted := Encrypt(key, data)

	defer func() { recover() }()
	Decrypt([]byte("incorrect"), encrypted)
	t.FailNow() // Should not be hit.
}
