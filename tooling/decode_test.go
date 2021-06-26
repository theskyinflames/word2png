package tooling_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestColors2Word(t *testing.T) {
	var (
		r2c       = make(map[rune]color.Color)
		c2r       = make(map[color.Color]rune)
		firstRune = 0
		lastRune  = 127
	)

	for i := firstRune; i < lastRune; i++ {
		color := tooling.ColorsTable[i]
		r2c[rune(i)] = color
		c2r[color] = rune(i)
	}

	var (
		word   = "bird"
		colors = []color.Color{r2c['b'], r2c['i'], r2c['r'], r2c['d']}
	)

	decoded, err := tooling.Colors2Word(colors, c2r)
	require.NoError(t, err)
	require.Equal(t, word, decoded)
}
