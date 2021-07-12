package tooling_test

import (
	"bytes"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/image-coder/tooling"
)

var words = []string{
	"keyboard",
	"horse",
	"berry",
	"frog",
	"cloud",
	"mouse",
	"monitor",
	"laptop",
	"glass",
	"grass",
	"čÇÑñ",
	"birdÑÇ 你",
}

func TestWords2Colors(t *testing.T) {
	encoder := tooling.NewEncoder("", r2cMapperFixture())
	w2c, err := encoder.Words2colors(words)
	require.NoError(t, err)

	// TODO See how to calculte the lengh in bytes for each crypted word
	for _, colors := range w2c {
		blanks := 0
		for _, c := range colors {
			require.NotNil(t, c)
			if c == tooling.BlackColor {
				blanks++
			}
		}
		require.Equal(t, 1, blanks)

		// the pixel's array forma of a word is started and ended by a black pixel
		// TODO See how to calculte the lengh in bytes for each crypted word
		require.Equal(t, tooling.BlackColor, colors[len(colors)-1])
	}

	// We're not limited to ASCII characters
	s := "你"
	_, err = encoder.Words2colors([]string{s})
	require.NoError(t, err)
}

func TestEncryptWords(t *testing.T) {
	e := tooling.NewEncoder("mySeed", r2cMapperFixture())
	b, err := e.EncryptWords(words)
	require.NoError(t, err)
	require.Len(t, b, len(words))
}

func TestEncode(t *testing.T) {
	encoder := tooling.NewEncoder("I'm glad to meet you in this dark times.", tooling.Rune2Color)
	encodedImage, err := encoder.Encode(words)
	require.NoError(t, err)
	require.NotEmpty(t, encodedImage)

	buff := &bytes.Buffer{}
	buff.Write(encodedImage)
	img, err := png.Decode(buff)
	require.NoError(t, err)

	// One line per word
	require.Equal(t, len(words), img.Bounds().Max.Y-img.Bounds().Min.Y)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		var c color.Color
		blacks := 0
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c = img.At(x, y)
			require.NotNil(t, c) // nil color not allowed
			switch c {
			case tooling.BlackColor:
				// Only one black color is allowed to mark the end of the word
				require.Equal(t, 0, blacks)
				blacks++
			case tooling.WhiteColor:
				// If the current color is white, then the before one only can be
				// black or white
				before := img.At(x-1, y)
				require.True(t, before == tooling.BlackColor || before == tooling.WhiteColor)
			}
		}
		require.Equal(t, 1, blacks) // All words are closed with black color
	}
}

func r2cMapperFixture() tooling.Rune2ColorMapper {
	var (
		r2c         = make(map[rune]color.Color)
		firstString = 0
		lastString  = 127
	)

	for i := firstString; i < lastString; i++ {
		r2c[rune(i)] = tooling.ColorsTable[i]
	}

	return func(seed string) (map[rune]color.Color, map[color.Color]rune) {
		return r2c, nil
	}
}
