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

// Rune2Color returns an bijective function f(rune)->color
//
// The goal here is to have a function that, given a seed,
// it provides an unique application to color: f(seed,[]string,color.Color[])->[]{r,color.Color}
//
// To achieve that, it follows these steps:
//
// 1. We get the binary mask from the seed
// 2. Create two string slices: head and tail
// 3. Iterate binary seed mask, and for each position:
//     - we take the rune that corresponds to the position number: rune(position)
//     - if there is a '0' in the mask,
//          the rune is added to the tail slice
//       else, the rune is added to the head slice.
//     - at the end, we concatenate head + tail slices
//     So, for a seed binary mask "01011101...", we'll end up with an slice like "b,d,e,f,h,....,a,c,g, ...."
// 4. Building the encoding map f(rune)->color.Color . Taking the above slide,
//     for each string we apply the color which is in the same position with a map.
//     Taking our before example:
//  		* b->[]color.Color[0]
//			* d->[]color.Color[1]
//          * e->[]color.Color[2]
//			* f->[]color.Color[3]
//			...
// 5. Building the decoding map f(color.Color)->r
//  		* []color.Color[0]->b
//			* []color.Color[1]->d
//          * []color.Color[2]->e
//			* []color.Color[3]->f
//			...
func Rune2Color(seed string) (map[rune]color.Color, map[color.Color]rune) {
	md5BinaryMask := createMaskFromSeed(seed)

	head := make([]rune, 0)
	tail := make([]rune, 0)
	for i := range md5BinaryMask {
		r := rune(i)
		if md5BinaryMask[i] == 0 {
			head = append(head, r)
		} else {
			tail = append(tail, r)
		}
	}
	masked := append(head, tail...)
	rune2color := make(map[rune]color.Color)
	color2rune := make(map[color.Color]rune)
	for i := range masked {
		r := masked[i]
		c := ColorsTable[i]
		rune2color[r] = c
		color2rune[c] = r
	}

	return rune2color, color2rune
}

// creteMaskFromSeed returns a 32 byte array
// which is the MD5 checksum of the seed
func createMaskFromSeed(seed string) []int8 {
	hasher := md5.New()
	hasher.Write([]byte(seed))
	return bytes2bits(hasher.Sum(nil))
}

// SaveEncodedImage is self described
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

// ColorsSource return the palette of colors used to encode words
func ColorsSource() []color.Color {
	// We do not use the black nor white colors to encode runes
	p := palette.WebSafe[1:]
	return p[:len(p)-1]
}
