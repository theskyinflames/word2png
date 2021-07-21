package lib_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/lib"
)

func TestColors2Word(t *testing.T) {
	var (
		words     = [][]byte{[]byte("birdÑÇ 你"), []byte("barcelona1"), []byte("sevilla")}
		firstSeed = "mySeed"
	)

	r2c, c2r := lib.Rune2Color(firstSeed)

	c2rMapper := func(seed string) (map[rune]color.Color, map[color.Color]rune) {
		return r2c, c2r
	}

	// Encrypt the words
	encryptedWords := make([][]byte, len(words))
	seeds := []string{firstSeed}
	for i := range words {
		encryptedWords[i] = lib.Encrypt(words[i], seeds[i])
		seeds = append(seeds, string(encryptedWords[i]))
	}

	// Decode from colors
	decoder := lib.NewDecoder(firstSeed, c2rMapper)
	for i, cw := range encryptedWords {
		// build the colors array for the crypted word
		colors := []color.Color{}
		for _, b := range cw {
			high, low := lib.SplitByte(b)
			colors = append(colors, r2c[rune(high)], r2c[rune(low)])
		}
		// decode the color to word, decrypting it using its seed
		cryptedWord, err := decoder.Colors2CryptedWord(colors)
		decoded := lib.Decrypt(cryptedWord, seeds[i])

		require.NoError(t, err)
		require.Equal(t, string(words[i]), string(decoded))
	}
}
