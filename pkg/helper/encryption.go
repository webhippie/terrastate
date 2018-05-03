package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
)

// Encrypt encrypts a given payload by given key.
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher")
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create gcm")
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "failed to create nonce")
	}

	ciphertext, err := gcm.Seal(nonce, nonce, plaintext, nil), nil

	if err != nil {
		return nil, errors.Wrap(err, "failed to seal content")
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
		return nil, errors.Wrap(err, "failed to create cipher")
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create gcm")
	}

	nonceSize := gcm.NonceSize()

	if len(decoded) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
