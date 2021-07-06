package tooling

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
)

type Encoder struct {
	r2c map[rune]color.Color
}

type Rune2ColorMapper func(seed string) (map[rune]color.Color, map[color.Color]rune)

func NewEncoder(seed string, r2cMapper Rune2ColorMapper) Encoder {
	r2c, _ := r2cMapper(seed)
	return Encoder{
		r2c: r2c,
	}
}

var errMsgNoColorsForWord = "no colors for the word %s"

// Encode encodes a list of words in an image based on the rune-2-color slice
func (e Encoder) Encode(words []string) ([]byte, error) {
	longestWord := LongestWord(words) + 2 // BlackColor as a mark of start/end of the word

	// Image to encode the words into it
	img := image.NewPaletted(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{longestWord, len(words)},
	}, palette.WebSafe)

	// Paint canvas background
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, whiteColor)
		}
	}

	// Add words to the image
	w2c, err := e.Words2colors(words)
	if err != nil {
		return nil, err
	}
	y := 0
	for _, word := range words {
		wordColors, ok := w2c[word]
		if !ok {
			return nil, fmt.Errorf(errMsgNoColorsForWord)
		}
		for x, wc := range wordColors {
			img.Set(x, y, wc)
		}
		y++
	}

	buff := &bytes.Buffer{}
	png.Encode(buff, img)
	return buff.Bytes(), nil
}

func LongestWord(words []string) int {
	l := 0
	for _, w := range words {
		if len(w) > l {
			l = len(w)
		}
	}
	return l
}

var ErrMsgNoColorForRune = "no color for the string %d"

// Words2colors return for each word, its representation as an array of colors
func (e Encoder) Words2colors(words []string) (map[string][]color.Color, error) {
	m := make(map[string][]color.Color)
	for _, word := range words {
		m[word] = []color.Color{}
		m[word] = append(m[word], BlackColor)
		for _, r := range word {
			color, ok := e.r2c[r]
			if !ok {
				return nil, fmt.Errorf(ErrMsgNoColorForRune, r)
			}
			m[word] = append(m[word], color)
		}
		m[word] = append(m[word], BlackColor)
	}
	return m, nil
}
