package tooling_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

func TestString2Color(t *testing.T) {
	seed := "bartolo"
	r2c, c2r := tooling.Rune2Color(seed)
	require.Equal(t, len(r2c), len(c2r))

	for myString, color := range r2c {
		StringFromColor, ok := c2r[color]
		require.True(t, ok)
		require.Equal(t, myString, StringFromColor)
	}
}

func TestColorsPalette(t *testing.T) {
	p := tooling.ColorsSource()
	for _, c := range p {
		require.NotEqual(t, tooling.BlackColor, c)
		require.NotEqual(t, tooling.WhiteColor, c)
	}
}
