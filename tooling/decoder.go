package tooling

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

// Decode decodes the words inside a given image to its string form
func Decode(coded []byte, c2r map[color.Color]rune) ([]string, error) {
	buff := &bytes.Buffer{}
	buff.Write(coded)
	img, _, err := image.Decode(buff)
	if err != nil {
		return nil, err
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	words := []string{}

	fmt.Println("reading...")
	for y := 0; y < h; y++ {
		var currentWord []color.Color
		readingAWord := false

		for x := 0; x < w; x++ {
			c := img.At(x, y)

			if c == BlackColor && !readingAWord {
				// starting word reading
				readingAWord = true
				currentWord = []color.Color{}
				continue
			}

			if c == BlackColor && readingAWord {
				// finishing word reading
				currentWord = append(currentWord, c)
				w, err := Colors2Word(currentWord, c2r)
				if err != nil {
					return nil, err
				}
				words = append(words, w)
				readingAWord = false
				continue
			}

			if readingAWord {
				// read a new letter of the word
				currentWord = append(currentWord, c)
			}
		}
	}

	return words, nil
}

var ErrMsgNoRuneForColor = "no rune for the color %s"

func Colors2Word(colors []color.Color, c2r map[color.Color]rune) (string, error) {
	w := []rune{}
	for _, color := range colors {
		if color == BlackColor {
			continue
		}
		r, ok := c2r[color]
		if !ok {
			return "", fmt.Errorf(ErrMsgNoRuneForColor, fmt.Sprint(color))
		}
		w = append(w, r)
	}
	return string(w), nil
}
