package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt encrypts a given payload by given key.
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to create nonce: %w", err)
	}

	ciphertext, err := gcm.Seal(nonce, nonce, plaintext, nil), nil

	if err != nil {
		return nil, fmt.Errorf("failed to seal content: %w", err)
	}

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt decrypts a given payload by given key.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(ciphertext))

	if err != nil {
		return decoded, nil
	}

	c, err := aes.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()

	if len(decoded) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
