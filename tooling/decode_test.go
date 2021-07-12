package tooling_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestColors2Word(t *testing.T) {
	var (
		words     = [][]byte{[]byte("birdÑÇ 你"), []byte("barcelona1"), []byte("sevilla")}
		firstSeed = "mySeed"
	)

	r2c, c2r := tooling.Rune2Color(firstSeed)

	c2rMapper := func(seed string) (map[rune]color.Color, map[color.Color]rune) {
		return r2c, c2r
	}

	// Crypt - decrypt cycle
	cryptedWords := make([][]byte, len(words))
	seeds := []string{firstSeed}
	for i := range words {
		cryptedWords[i] = tooling.Encrypt(words[i], seeds[i])
		seeds = append(seeds, string(cryptedWords[i]))
	}
	for i := range words {
		decrypted := tooling.Decrypt(cryptedWords[i], seeds[i])
		require.Equal(t, words[i], decrypted)
	}

	// Decode from colors
	decoder := tooling.NewDecoder(firstSeed, c2rMapper)
	for i, cw := range cryptedWords {
		// build the colors array for the crypted word
		colors := []color.Color{}
		for _, b := range cw {
			high, low := tooling.SplitByte(b)
			colors = append(colors, r2c[rune(high)], r2c[rune(low)])
		}
		// decode the color to word, decrypting it using its seed
		cryptedWord, err := decoder.Colors2CryptedWord(colors)
		decoded := tooling.Decrypt(cryptedWord, seeds[i])

		require.NoError(t, err)
		require.Equal(t, string(words[i]), string(decoded))
	}
}
