package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

type AES256 struct {
	passphrase string
}

func NewAES256(passphrase string) AES256 {
	return AES256{
		passphrase: passphrase,
	}
}

func (a AES256) EncryptWords(words []string) ([][]byte, error) {
	var (
		err            error
		passphrase     = a.passphrase
		encryptedWords = make([][]byte, len(words))
	)
	for i, w := range words {
		encryptedWords[i], err = a.Encrypt([]byte(w), passphrase)
		if err != nil {
			return nil, err
		}
		passphrase = string(encryptedWords[i])
	}
	return encryptedWords, nil
}

// Encrypt encrypts given byte array using AES-256-GCM with an Argon2id-derived key.
func (a AES256) Encrypt(data []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Prepend nonce to ciphertext so it can be extracted during decryption.
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return append(salt, ciphertext...), nil
}

func (a AES256) DecryptWords(encryptedWords [][]byte) ([]string, error) {
	var (
		passphrase     = a.passphrase
		decryptedWords = make([]string, len(encryptedWords))
	)
	for i, ew := range encryptedWords {
		dw, err := a.Decrypt(ew, passphrase)
		if err != nil {
			return nil, err
		}
		passphrase = string(ew)
		decryptedWords[i] = string(dw)
	}
	return decryptedWords, nil
}

// Decrypt decrypts given byte array using AES-256-GCM with an Argon2id-derived key.
func (a AES256) Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) < 16 {
		return nil, errors.New("ciphertext too short")
	}

	salt := data[:16]
	ciphertext := data[16:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	decryptedData, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4
	argon2KeyLen  = 32 // 256 bits for AES-256
)

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
}
