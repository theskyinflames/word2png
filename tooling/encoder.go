package tooling

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"io"
)

type EncoderOption func(*Encoder)

func EncoderDebugWriterOpt(dw io.Writer) EncoderOption {
	return func(e *Encoder) {
		e.debugWriter = dw
	}
}

type Encoder struct {
	r2c         map[rune]color.Color
	passphrase  string
	debugWriter io.Writer
}

type Rune2ColorMapper func(seed string) (map[rune]color.Color, map[color.Color]rune)

func NewEncoder(seed string, r2cMapper Rune2ColorMapper, opts ...EncoderOption) Encoder {
	r2c, _ := r2cMapper(seed)
	e := Encoder{
		r2c:        r2c,
		passphrase: seed,
	}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}

var errMsgNoColorsForWord = "no colors for the word %s"

// Encode encodes a list of words in an image based on the rune-2-color slice
func (e Encoder) Encode(words []string) ([]byte, error) {
	w2c, err := e.Words2colors(words)
	// Add words to the image
	if err != nil {
		return nil, err
	}

	longestWord := longestWord(w2c) + 2 // BlackColor as a mark of start/end of the word

	// Image to encode the words into it
	img := image.NewPaletted(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{longestWord, len(words)},
	}, palette.WebSafe)

	// Paint canvas background
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, WhiteColor)
		}
	}

	for y, word := range words {
		wordColors, ok := w2c[word]
		if !ok {
			return nil, fmt.Errorf(errMsgNoColorsForWord)
		}
		for x, wc := range wordColors {
			img.Set(x, y, wc)
		}
	}

	buff := &bytes.Buffer{}
	png.Encode(buff, img)
	return buff.Bytes(), nil
}

func longestWord(words map[string][]color.Color) int {
	l := 0
	for _, w := range words {
		if len(w) > l {
			l = len(w)
		}
	}
	return l
}

func (e Encoder) EncryptWords(words []string) ([][]byte, error) {
	crypted := make([][]byte, 0)
	passphrase := []byte(e.passphrase)
	for _, w := range words {
		b := Encrypt([]byte(w), string(passphrase))
		crypted = append(crypted, b)
		passphrase = b
	}
	return crypted, nil
}

var ErrMsgNoColorForRune = "no color for the rune %d"

// Words2colors return for each word, its representation as an array of colors
func (e Encoder) Words2colors(words []string) (map[string][]color.Color, error) {
	encryptedWords, err := e.EncryptWords(words)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]color.Color)
	for i, word := range words {
		m[word] = []color.Color{}

		if e.debugWriter != nil {
			e.debugWriter.Write([]byte(fmt.Sprintf("\nword: %d \n", i)))
		}

		for _, wordByte := range encryptedWords[i] {
			// MD5 checksum signature limits us to having only 128 available colors.
			// But with each byte we have 256 (2^8) possibilities for the color.
			// So I take one color for the first 4 bits (low) of each byte
			// and another one for the next 4 bits (high)
			// By doing that, for example, a byte 10111001 is splited in two
			// bytes: 00001011 (high part) and 00001001 (low part). Each of these
			// parts will have its own color. So, at decode time, these two bytes
			// will be reconbined againg to get back the original byte 10111001.
			high, low := SplitByte(wordByte)

			if e.debugWriter != nil {
				e.debugWriter.Write([]byte(fmt.Sprintf("%08b, %08b - %08b\n", wordByte, high, low)))
			}

			for _, b := range []byte{high, low} {
				r := rune(b)
				color, ok := e.r2c[r]
				if !ok {
					return nil, fmt.Errorf(ErrMsgNoColorForRune, r)
				}
				m[word] = append(m[word], color)
			}
		}
		m[word] = append(m[word], BlackColor)
	}
	return m, nil
}
