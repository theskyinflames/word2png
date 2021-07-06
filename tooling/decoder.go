package tooling

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"io"
)

type DecoderOption func(*Decoder)

type Decoder struct {
	c2r         map[color.Color]rune
	passphrase  string
	debugWriter io.Writer
}

func DecodeDebugWriterOpt(dw io.Writer) DecoderOption {
	return func(d *Decoder) {
		d.debugWriter = dw
	}
}

func NewDecoder(passphrase string, c2rMapper Rune2ColorMapper, opts ...DecoderOption) Decoder {
	_, c2r := c2rMapper(passphrase)
	d := Decoder{
		c2r:        c2r,
		passphrase: passphrase,
	}

	for _, opt := range opts {
		opt(&d)
	}

	return d
}

// Decode decodes the words inside a given image to its string form.
// It expects that the image is in PNG format.
func (d Decoder) Decode(coded []byte) ([]string, error) {
	buff := &bytes.Buffer{}
	buff.Write(coded)
	img, err := png.Decode(buff)
	if err != nil {
		return nil, err
	}

	readWords := []string{}

	// Reading the coded image.
	// The goal here is to read the
	// coded words inside the image.
	// Each word is delimited with a black colored pixel.
	//
	// There is one word per line of pixels
	passphrase := []byte(d.passphrase)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		var readingWord []color.Color

		if d.debugWriter != nil {
			d.debugWriter.Write([]byte(fmt.Sprintf("\nword: %d \n", y)))
		}

	nextword:
		// Each 2 colors are the high and low parts of the crypted byte
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y)
			switch c {
			case nil:
				return nil, errors.New("nil color not allowed")
			case BlackColor:
				// finishing word reading
				encryptedWord, err := d.Colors2CryptedWord(readingWord)
				if err != nil {
					return nil, err
				}

				w := Decrypt(encryptedWord, string(passphrase))
				readWords = append(readWords, string(w))

				// setting the seed to decrypt next word
				passphrase = encryptedWord

				break nextword
			default:
				readingWord = append(readingWord, c)
			}
		}
	}
	return readWords, nil
}

var ErrMsgNoRuneForColor = "no rune for the color %s"

// Colors2CryptedWord translates the array of colors to the original word
func (d Decoder) Colors2CryptedWord(colors []color.Color) ([]byte, error) {
	// Each 2 Colors is equivalent to the high and low parts of the
	// byte that belongs to the cripted form of the original word.
	// So, first we need to rebuild each byte by join its high and low parts.
	// After that we'll have the bytes array corresponding to the crpted form
	// of the original word. So, all we'll have to do to get the word is decrypt
	// this sequence of bytes.

	// Rebuild the crypted form of the word
	cryptedWord := []byte{}
	for i := 0; i < len(colors); i = i + 2 {
		cHigh, ok := d.c2r[colors[i]]
		if !ok {
			return nil, fmt.Errorf(ErrMsgNoRuneForColor, fmt.Sprintf("%#v", cHigh))
		}
		cLow, ok := d.c2r[colors[i+1]]
		if !ok {
			return nil, fmt.Errorf(ErrMsgNoRuneForColor, fmt.Sprintf("%#v", cHigh))
		}

		originalByte := JoinByte(byte(cHigh), byte(cLow))
		cryptedWord = append(cryptedWord, originalByte)
	}

	if d.debugWriter != nil {
		for _, cw := range cryptedWord {
			cHigh, cLow := SplitByte(cw)
			d.debugWriter.Write([]byte(fmt.Sprintf("%08b, %08b - %08b\n", cw, cHigh, cLow)))
		}
	}

	// Decrypt
	return cryptedWord, nil
}
