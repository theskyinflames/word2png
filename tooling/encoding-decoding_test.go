package tooling_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestEncodingDecoding(t *testing.T) {
	seed := "I'm glad to meet you in this dark times."

	// encoding
	encoder := tooling.NewEncoder(seed, tooling.Rune2Color)
	encodedImage, err := encoder.Encode(words)
	require.NoError(t, err)
	require.NotEmpty(t, encodedImage)

	// decoding
	decoder := tooling.NewDecoder(seed, tooling.Rune2Color)
	decodedWords, err := decoder.Decode(encodedImage)
	require.NoError(t, err)
	require.Equal(t, len(decodedWords), len(words))
	require.ElementsMatch(t, words, decodedWords)
}
