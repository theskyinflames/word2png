package tooling

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func Encrypt(data []byte, passphrase string) []byte {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(createHash(passphrase))

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {

	key, _ := hex.DecodeString(createHash(passphrase))

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	//Decrypt the data
	decryptedData, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return decryptedData
}

//func Encrypt(data []byte, passphrase string) []byte {
//	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
//	gcm, err := cipher.NewGCM(block)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	nonce := make([]byte, gcm.NonceSize())
//	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
//		panic(err.Error())
//	}
//	ciphertext := gcm.Seal(nonce, nonce, data, nil)
//	return ciphertext
//}

//func Decrypt(data []byte, passphrase string) []byte {
//	key := []byte(createHash(passphrase))
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		panic(err.Error())
//	}
//	gcm, err := cipher.NewGCM(block)
//	if err != nil {
//		panic(err.Error())
//	}
//	nonceSize := gcm.NonceSize()
//	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
//	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
//	if err != nil {
//		panic(err.Error())
//	}
//	return plaintext
//}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
