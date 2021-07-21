package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/lib"
)

func TestEncrypDecrypt(t *testing.T) {
	var (
		seed = "mySeed"
		txt  = []byte("Make the force ÇÑ be with you 你")
	)

	crypted := lib.Encrypt(txt, seed)
	decrypted := lib.Decrypt(crypted, seed)

	require.Equal(t, txt, decrypted)
}
