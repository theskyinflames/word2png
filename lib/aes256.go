package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
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

// Encrypt encrypts given byte array using AES-256 algorithm using the passphrase
func (a AES256) Encrypt(data []byte, passphrase string) ([]byte, error) {
	// AES-256 needs a 32 bytes key. So it's taken form
	// the passphrase MD5 checksum
	key, _ := hex.DecodeString(createHash(passphrase))

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data.
	// The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
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

// Decrypt decrypts given byte array using AES-256 algorithm using the passphrase
func (a AES256) Decrypt(data []byte, passphrase string) ([]byte, error) {
	// AES-256 needs a 32 bytes key. So it's taken form
	// the passphrase MD5 checksum
	key, _ := hex.DecodeString(createHash(passphrase))

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the data
	decryptedData, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
