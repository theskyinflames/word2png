package lib_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/lib"
)

func TestColors2Word(t *testing.T) {
	r2c, c2r := lib.Rune2Color(seed)()

	c2rMapper := func() (map[rune]color.Color, map[color.Color]rune) {
		return r2c, c2r
	}

	// Encrypt the words
	encryptedWords, err := encrypter.EncryptWords(nil)
	require.NoError(t, err)

	// Decode from colors
	decoder := lib.NewDecoder(c2rMapper, decrypter)
	for i, cw := range encryptedWords {
		// build the colors array for the crypted word
		colors := []color.Color{}
		for _, b := range cw {
			high, low := lib.SplitByte(b)
			colors = append(colors, r2c[rune(high)], r2c[rune(low)])
		}
		// decode the color to word, decrypting it using its seed
		cryptedWord, err := decoder.Colors2CryptedWord(colors)
		require.NoError(t, err)
		require.Equal(t, encryptedWords[i], cryptedWord)
	}
}

func TestRemoveEnumerationToken(t *testing.T) {
	word := fmt.Sprintf("%sbar%s%stolo%s%s", lib.EnumerateToken, lib.EnumerateToken, lib.EnumerateToken, lib.EnumerateToken, lib.EnumerateToken)
	enumerated := fmt.Sprintf(lib.EnumerationMask, 0, lib.EnumerateToken, word)
	require.Equal(t, word, lib.RemoveEnumerationToken(enumerated))
}
