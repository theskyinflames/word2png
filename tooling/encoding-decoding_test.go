package tooling_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestEncodingDecoding(t *testing.T) {
	seed := "I'm glad to meet you in this dark times."

	// encoding
	r2c, _ := tooling.Rune2Color(seed)
	encodedImage, err := tooling.Encode(words, r2c)
	require.NoError(t, err)
	require.NotEmpty(t, encodedImage)

	// decoding
	_, c2r := tooling.Rune2Color(seed)
	decodedWords, err := tooling.Decode(encodedImage, c2r)
	require.NoError(t, err)
	require.Equal(t, len(decodedWords), len(words))
	require.ElementsMatch(t, words, decodedWords)
}
