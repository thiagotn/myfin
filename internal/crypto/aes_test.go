package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte(`{"portfolio": {"total_value": 10000.00}}`)
	passphrase := "test-passphrase-123"

	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(ciphertext) < len(plaintext) {
		t.Errorf("ciphertext too short: got %d, want at least %d", len(ciphertext), len(plaintext))
	}

	decrypted, err := Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted != plaintext\nwant: %s\ngot:  %s", plaintext, decrypted)
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	plaintext := []byte(`{"secret": "data"}`)
	passphrase := "correct-passphrase"

	ciphertext, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	wrongPassphrase := "wrong-passphrase"
	_, err = Decrypt(ciphertext, wrongPassphrase)
	if err == nil {
		t.Error("Decrypt with wrong passphrase should fail, but didn't")
	}
}

func TestEncryptIsDeterministic(t *testing.T) {
	plaintext := []byte(`{"data": "same"}`)
	passphrase := "passphrase"

	cipher1, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("first Encrypt failed: %v", err)
	}

	cipher2, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("second Encrypt failed: %v", err)
	}

	if string(cipher1) == string(cipher2) {
		t.Error("two encryptions should not be identical (salt+nonce are random)")
	}
}
