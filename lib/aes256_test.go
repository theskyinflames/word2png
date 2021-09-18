package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/word2png/lib"
)

func TestEncryptDecrypt(t *testing.T) {
	txt := []byte("Make the force ÇÑ be with you 你")

	aes256 := lib.NewAES256(seed)
	encrypted, err := aes256.Encrypt(txt, seed)
	require.NoError(t, err)
	decrypted, err := aes256.Decrypt(encrypted, seed)
	require.NoError(t, err)

	require.Equal(t, txt, decrypted)
}

func TestEncryptDecrypt_Words(t *testing.T) {
	aes256 := lib.NewAES256(seed)

	encrypted, err := aes256.EncryptWords(words)
	require.NoError(t, err)
	decrypted, err := aes256.DecryptWords(encrypted)
	require.NoError(t, err)

	require.ElementsMatch(t, words, decrypted)
}
