package tooling_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestEncrypDecrypt(t *testing.T) {
	var (
		seed = "mySeed"
		txt  = []byte("Make the force ÇÑ be with you 你")
	)

	crypted := tooling.Encrypt(txt, seed)
	decrypted := tooling.Decrypt(crypted, seed)

	require.Equal(t, txt, decrypted)
}
