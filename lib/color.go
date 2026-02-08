package lib

import (
	"crypto/md5"
	"image/color"
	"image/color/palette"
	"os"
)

// Color constants
var (
	ColorsTable = ColorsSource()
	WhiteColor  = palette.WebSafe[len(palette.WebSafe)-1]
	BlackColor  = palette.WebSafe[0]
)

// Rune2Color returns a deterministic rune/color and color/rune mapper factory.
//
// It computes a 128-bit mask from the seed (MD5), uses it to reorder runes
// [0..127] by appending bit=0 positions first and bit=1 positions last,
// and then assigns ColorsTable[i] to the i-th rune in that reordered sequence.
func Rune2Color(seed string) Rune2ColorMapper {
	return func() (map[rune]color.Color, map[color.Color]rune) {
		mask := createMaskFromSeed(seed)

		head := make([]rune, 0, len(mask))
		tail := make([]rune, 0, len(mask))
		for i := range mask {
			r := rune(i)
			if mask[i] == 0 {
				head = append(head, r)
			} else {
				tail = append(tail, r)
			}
		}
		masked := append(head, tail...)
		rune2color := make(map[rune]color.Color, len(masked))
		color2rune := make(map[color.Color]rune, len(masked))
		for i := range masked {
			r := masked[i]
			c := ColorsTable[i]
			rune2color[r] = c
			color2rune[c] = r
		}

		return rune2color, color2rune
	}
}

// createMaskFromSeed returns the 128 bits of the MD5 checksum of seed.
func createMaskFromSeed(seed string) []int8 {
	sum := md5.Sum([]byte(seed))
	return bytes2bits(sum[:])
}

// SaveEncodedImage persists the encoded image bytes in path.
func SaveEncodedImage(encodedImage []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(encodedImage)
	if err != nil {
		return err
	}
	return f.Close()
}

// ColorsSource returns the palette of colors used to encode runes.
func ColorsSource() []color.Color {
	// Black and white are reserved as control colors.
	p := palette.WebSafe[1:]
	return p[:len(p)-1]
}
