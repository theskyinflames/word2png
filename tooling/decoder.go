package tooling

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

type Decoder struct {
	c2r map[color.Color]rune
}

func NewDecoder(seed string, c2rMapper Rune2ColorMapper) Decoder {
	_, c2r := c2rMapper(seed)
	return Decoder{
		c2r: c2r,
	}
}

// Decode decodes the words inside a given image to its string form
func (d Decoder) Decode(coded []byte) ([]string, error) {
	buff := &bytes.Buffer{}
	buff.Write(coded)
	img, _, err := image.Decode(buff)
	if err != nil {
		return nil, err
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	readWords := []string{}

	// Reading the coded image.
	// The goal here is to read the
	// coded words inside the image.
	// Each word is delimited with a black colored pixel.
	for y := 0; y < h; y++ {
		var readingWord []color.Color
		readingAWord := false

		for x := 0; x < w; x++ {
			c := img.At(x, y)
			switch {
			case c == BlackColor && !readingAWord:
				// starting word reading
				readingAWord = true
				readingWord = []color.Color{}
				continue

			case c == BlackColor && readingAWord:
				// finishing word reading
				readingWord = append(readingWord, c)
				wordFromColors, err := d.Colors2Word(readingWord)
				if err != nil {
					return nil, err
				}
				readWords = append(readWords, wordFromColors)
				readingAWord = false
				continue

			case readingAWord:
				// read a new letter of the reading word
				readingWord = append(readingWord, c)
			}
		}
	}
	return readWords, nil
}

var ErrMsgNoRuneForColor = "no rune for the color %s"

func (d Decoder) Colors2Word(colors []color.Color) (string, error) {
	w := []rune{}
	for _, color := range colors {
		if color == BlackColor {
			continue
		}
		r, ok := d.c2r[color]
		if !ok {
			return "", fmt.Errorf(ErrMsgNoRuneForColor, fmt.Sprint(color))
		}
		w = append(w, r)
	}
	return string(w), nil
}
