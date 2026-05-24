package lib_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theskyinflames/word2png/lib"
)

func TestRealCryptoEncodingDecodingRoundTrip(t *testing.T) {
	const (
		passphrase = "I'm glad to meet you in this dark times."
		filePath   = "./result.png"
	)

	aes256 := lib.NewAES256(passphrase)

	encoder := lib.NewEncoder(lib.Rune2Color(passphrase), aes256)
	encodedImage, err := encoder.Encode(words)
	require.NoError(t, err)
	require.NotEmpty(t, encodedImage)

	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.Write(encodedImage)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	defer func() {
		require.NoError(t, os.Remove(filePath))
	}()

	encodedImage, err = os.ReadFile(filePath)
	require.NoError(t, err)
	decoder := lib.NewDecoder(lib.Rune2Color(passphrase), aes256)
	decodedWords, err := decoder.Decode(encodedImage)
	require.NoError(t, err)

	require.Equal(t, words, decodedWords)
}
