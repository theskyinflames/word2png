package lib

import (
	"crypto/md5"
	"image/color"
	"image/color/palette"
	"os"
)

var (
	ColorsTable = ColorsSource()
	WhiteColor  = palette.WebSafe[len(palette.WebSafe)-1]
	BlackColor  = palette.WebSafe[0]
)

// Rune2Color returns an bijective function f(rune): color
//
// The goal here is to have a function that,
// given a seed provides an unique application to color
// f(seed,[]string,color.Color[]) :-> unique(r,color.Color)
// To do that,
// * first, get the binary mask from the mask
// * second, create two string slices: head and tail
// * next, iterate seed mask, for each position:
//     if there is a '0' in the mask,
//          the string for this position is added to the tail slice
//     else, the string for this position is added to the head slice.
//     At the end, we concatenate head + tail slices
//
//     So, for a mask: "01011101...", we'll end up with an slice like "b,d,e,f,h,....,a,c,g, ...."
//
// * Four, building the encoding map f(r)->color.Color . Taking the above slide,
//     for each string we apply the color which is in the same position with a map.
//     Taking our before example:
//  		* b->[]color.Color[0]
//			* d->[]color.Color[1]
//          * e->[]color.Color[2]
//			* f->[]color.Color[3]
//			...
//
// * Five, building the decoding map f(color.Color)->r
//  		* []color.Color[0]->b
//			* []color.Color[1]->d
//          * []color.Color[2]->e
//			* []color.Color[3]->f
//			...
func Rune2Color(seed string) (map[rune]color.Color, map[color.Color]rune) {
	// MD5 checksum provides an 256 bits lengh signature
	// So if we want to pair each rune to a color using
	// the MD5 checksum mask of the seed as mapper.
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

// ColorSource return the palette of colors used to encode words
func ColorsSource() []color.Color {
	// We do not use the black nor white colors to encode runes
	p := palette.WebSafe[1:]
	return p[:len(p)-1]
}