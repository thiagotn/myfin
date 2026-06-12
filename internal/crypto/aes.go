package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen   = 32
	nonceLen  = 12
	argon2mem = 64 * 1024 // 64 MB
	argon2t   = 1
	argon2p   = 4
	keyLen    = 32 // 256 bits for AES-256
)

func Encrypt(plaintext []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key := argon2.IDKey([]byte(passphrase), salt, argon2t, argon2mem, argon2p, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	result := make([]byte, 0, saltLen+len(ciphertext))
	result = append(result, salt...)
	result = append(result, ciphertext...)

	return result, nil
}

func Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	if len(ciphertext) < saltLen+nonceLen+16 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := ciphertext[:saltLen]
	encrypted := ciphertext[saltLen:]

	key := argon2.IDKey([]byte(passphrase), salt, argon2t, argon2mem, argon2p, keyLen)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := encrypted[:nonceLen]
	ciphertextOnly := encrypted[nonceLen:]

	plaintext, err := gcm.Open(nil, nonce, ciphertextOnly, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
