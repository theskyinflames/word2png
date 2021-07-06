package tooling_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestEncodingDecoding(t *testing.T) {
	const (
		seed     = "I'm glad to meet you in this dark times."
		filePath = "./result.png"
	)

	// encoding
	encoder := tooling.NewEncoder(seed, tooling.Rune2Color)
	encodedImage, err := encoder.Encode(words)
	require.NoError(t, err)
	require.NotEmpty(t, encodedImage)

	f, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = f.Write(encodedImage)
	require.NoError(t, err)

	// decoding
	encodedImage, err = ioutil.ReadFile(filePath)
	require.NoError(t, err)
	decoder := tooling.NewDecoder(seed, tooling.Rune2Color)
	decodedWords, err := decoder.Decode(encodedImage)
	require.NoError(t, err)
	require.Equal(t, len(decodedWords), len(words))
	require.ElementsMatch(t, words, decodedWords)
}
