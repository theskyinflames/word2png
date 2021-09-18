package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/word2png/lib"
)

func TestRune2Color(t *testing.T) {
	seed := "bartolo"
	r2c, c2r := lib.Rune2Color(seed)()
	require.Equal(t, len(r2c), len(c2r))

	for r, c := range r2c {
		runeFromColor, ok := c2r[c]
		require.True(t, ok)
		require.Equal(t, r, runeFromColor)
	}
}

func TestColorsSource(t *testing.T) {
	p := lib.ColorsSource()
	for _, c := range p {
		require.NotEqual(t, lib.BlackColor, c)
		require.NotEqual(t, lib.WhiteColor, c)
	}
}
