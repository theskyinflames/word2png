package tooling

import (
	"crypto/md5"
	"image/color"
	"image/color/palette"
	"os"
)

var (
	ColorsTable = palette.WebSafe
	whiteColor  = palette.WebSafe[len(palette.WebSafe)-1]
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
//  		* b->[]color.Color[x]
//			* d->[]color.Color[y]
//          * e->[]color.Color[z]
//			* f->[]color.Color[t]
//			...
//
// * Five, building the decoding map f(color.Color)->r
//  		* []color.Color[x]->b
//			* []color.Color[y]->d
//          * []color.Color[z]->e
//			* []color.Color[t]->f
//			...
func Rune2Color(seed string) (map[rune]color.Color, map[color.Color]rune) {
	const numStrings = 128

	mask := createMaskFromSeed(seed)[:numStrings]

	head := make([]rune, 0)
	tail := make([]rune, 0)
	for i := range mask {
		r := rune(i)
		if mask[i] == 0 {
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

func createMaskFromSeed(seed string) []int8 {
	hasher := md5.New()
	hasher.Write([]byte(seed))
	return bytes2bits(hasher.Sum(nil))
}

// Taken from https://stackoverflow.com/questions/52811744/extract-bits-into-a-int-slice-from-byte-slice
func bytes2bits(data []byte) []int8 {
	r := make([]int8, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			r[i*8+j] = int8(b >> uint(7-j) & 0x01)
		}
	}
	return r
}

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
