package lib

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"os"
	"strings"
)

//go:generate moq -out zmock_decoder_test.go -pkg lib_test . Decrypter

type Decrypter interface {
	DecryptWords(encryptedWords [][]byte) ([]string, error)
}

// DecoderOption a constructor option
type DecoderOption func(*Decoder)

// Decoder decodes encoded words in a PNG image
type Decoder struct {
	c2r         map[color.Color]rune
	debugWriter io.Writer

	decrypter Decrypter
}

// DecodeDebugWriterOpt provides an output for debug messages
func DecodeDebugWriterOpt(dw io.Writer) DecoderOption {
	return func(d *Decoder) {
		d.debugWriter = dw
	}
}

// NewDecoder is a constructor
func NewDecoder(c2rMapper Rune2ColorMapper, decrypter Decrypter, opts ...DecoderOption) Decoder {
	_, c2r := c2rMapper()
	d := Decoder{
		c2r:       c2r,
		decrypter: decrypter,
	}

	for _, opt := range opts {
		opt(&d)
	}

	return d
}

// Decode decodes the words inside a given image
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
	// Each word is ended with a black colored pixel.
	//
	// There is one word per line of pixels
	encryptedWords := make([][]byte, 0)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		var readingWord []color.Color

		if d.debugWriter != nil {
			_, _ = d.debugWriter.Write([]byte(fmt.Sprintf("\nword: %d \n", y)))
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
				encryptedWords = append(encryptedWords, encryptedWord)

				break nextword
			default:
				readingWord = append(readingWord, c)
			}
		}
	}

	enumeratedWords, err := d.decrypter.DecryptWords(encryptedWords)
	if err != nil {
		return nil, err
	}
	for _, ew := range enumeratedWords {
		w := RemoveEnumerationToken(ew) // This removes the enumeration added at encoding time
		readWords = append(readWords, w)
	}

	return readWords, nil
}

// ErrMsgNoRuneForColor is self described
var ErrMsgNoRuneForColor = "no rune for the color %s"

// Colors2CryptedWord translates the array of colors to the original word
func (d Decoder) Colors2CryptedWord(colors []color.Color) ([]byte, error) {
	// Each 2 Colors is equivalent to the high and low parts of the
	// byte that belongs to the encrypted form of the original word.
	// So, first we need to rebuild each byte by join its high and low parts.
	// After that we'll have the bytes array corresponding to the crpted form
	// of the original word. So, all we'll have to do to get the word is decrypt
	// this sequence of bytes.

	// Rebuild the encrypted form of the word
	cryptedWord := []byte{}
	for i := 0; i < len(colors); i += 2 {
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
			_, _ = d.debugWriter.Write([]byte(fmt.Sprintf("%08b, %08b - %08b\n", cw, cHigh, cLow)))
		}
	}

	return cryptedWord, nil
}

// RemoveEnumerationToken is self described
func RemoveEnumerationToken(word string) string {
	tokenAt := strings.Index(word, EnumerateToken)
	return word[tokenAt+1:]
}

// DecodeFromSource decodes an image from the provided source
// TODO add test
func (d Decoder) DecodeFromSource(imagePath string, b64 string) ([]string, error) {
	var (
		buff []byte
		err  error
	)

	switch {
	case imagePath == "" && b64 == "":
		return nil, errors.New("must be specified either a file or a b64 encoded string")
	case imagePath != "":
		buff, err = os.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}
	default:
		buff, err = base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, err
		}
	}

	return d.Decode(buff)
}
