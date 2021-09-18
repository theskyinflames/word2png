package lib_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/word2png/lib"
)

func TestEncodingDecoding(t *testing.T) {
	const (
		seed     = "I'm glad to meet you in this dark times."
		filePath = "./result.png"
	)

	// encoding
	encoder := lib.NewEncoder(lib.Rune2Color(seed), encrypter)
	encodedImage, err := encoder.Encode(words)
	require.NoError(t, err)
	require.NotEmpty(t, encodedImage)

	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.Write(encodedImage)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// comment this if you want to keep the png image
	defer func() {
		require.NoError(t, os.Remove(filePath))
	}()

	// decoding
	encodedImage, err = ioutil.ReadFile(filePath)
	require.NoError(t, err)
	decoder := lib.NewDecoder(lib.Rune2Color(seed), decrypter)
	decodedWords, err := decoder.Decode(encodedImage)
	require.NoError(t, err)
	require.Equal(t, len(decodedWords), len(words))
	require.ElementsMatch(t, words, decodedWords)
}
