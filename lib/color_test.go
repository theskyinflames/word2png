package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/lib"
)

func TestString2Color(t *testing.T) {
	seed := "bartolo"
	r2c, c2r := lib.Rune2Color(seed)
	require.Equal(t, len(r2c), len(c2r))

	for myString, color := range r2c {
		StringFromColor, ok := c2r[color]
		require.True(t, ok)
		require.Equal(t, myString, StringFromColor)
	}
}

func TestColorsPalette(t *testing.T) {
	p := lib.ColorsSource()
	for _, c := range p {
		require.NotEqual(t, lib.BlackColor, c)
		require.NotEqual(t, lib.WhiteColor, c)
	}
}
