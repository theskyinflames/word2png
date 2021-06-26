package tooling_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

var words = []string{
	"berry",
	"horse",
	"frog",
	"cloud",
	"keyboard",
	"mouse",
	"monitor",
	"laptop",
	"glass",
	"grass",
}

func TestLongestWord(t *testing.T) {
	require.Equal(t, len("keyboard"), tooling.LongestWord(words))
}

func TestWords2Colors(t *testing.T) {
	var (
		r2p         = make(map[rune]color.Color)
		firstString = 0
		lastString  = 127
	)

	for i := firstString; i < lastString; i++ {
		r2p[rune(i)] = tooling.ColorsTable[i]
	}

	w2p, err := tooling.Words2colors(words, r2p)
	require.NoError(t, err)

	require.Len(t, w2p, len(words))
	for word, colors := range w2p {
		// the pixel's array forma of a word is started and ended by a black pixel
		require.Equal(t, len(word)+2, len(colors))
		require.Equal(t, tooling.BlackColor, colors[0])
		require.Equal(t, tooling.BlackColor, colors[len(colors)-1])
	}

	delete(r2p, rune('b'))
	_, err = tooling.Words2colors(words, r2p)
	require.EqualError(t, err, fmt.Errorf(tooling.ErrMsgNoColorForRune, rune('b')).Error())
}
