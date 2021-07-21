package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/lib"
)

func TestEncryptDecrypt(t *testing.T) {
	var (
		seed = "mySeed"
		txt  = []byte("Make the force ÇÑ be with you 你")
	)

	encrypted := lib.Encrypt(txt, seed)
	decrypted := lib.Decrypt(encrypted, seed)

	require.Equal(t, txt, decrypted)
}
