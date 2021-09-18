package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/theskyinflames/word2png/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		imagePath   = kingpin.Flag("file", "Save to the especified file if it's filled").Short('f').String()
		seed        = kingpin.Flag("seed", "coding seed").Short('s').Required().String()
		words       = kingpin.Flag("words", "list of words to encode").Short('w').Strings()
		debug       = kingpin.Flag("debug", "writes a debug file").Short('d').Bool()
		b64         = kingpin.Flag("b64", "b64encoded image").String()
		removeWords = kingpin.Flag("remove-word", "remove a word from an image by index number").Short('r').Ints()

		debugFile *os.File
		err       error
	)
	kingpin.Parse()

	if debug != nil && *debug {
		debugFile, err = os.Create("./encrypted-bytes.txt")
		exitIfError(err)
		defer func() {
			debugFile.Close()
		}()
	}

	aes256 := lib.NewAES256(*seed)

	// If the image already exists, received words will be appended to the existent ones
	if (*imagePath != "" && imageExists(*imagePath)) || *b64 != "" {
		decoder := lib.NewDecoder(lib.Rune2Color(*seed), aes256, lib.DecodeDebugWriterOpt(debugFile))
		beforeWords, err := decoder.DecodeFromSource(*imagePath, *b64)
		exitIfError(err)
		*words = append(beforeWords, *words...)
	}

	// If remove-words flag has been provided, it's applied
	*words = RemoveWordsByIdx(*words, *removeWords)

	encoder := lib.NewEncoder(lib.Rune2Color(*seed), aes256, lib.EncoderDebugWriterOpt(debugFile))
	b, err := encoder.Encode(*words)
	exitIfError(err)

	switch {
	case *imagePath != "":
		exitIfError(lib.SaveEncodedImage(b, *imagePath))
	default:
		b64Encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
		_, err = b64Encoder.Write(b)
		exitIfError(err)
		exitIfError(b64Encoder.Close())
	}

	fmt.Println("\ncoding process finished")
	os.Exit(0)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(-1)
	}
}

func imageExists(imagePath string) bool {
	// if error is nil, I guess the file exists
	if _, err := os.Stat(imagePath); err == nil {
		return true
	}
	return false
}

func RemoveWordsByIdx(words []string, rmIdxs []int) []string {
	if len(rmIdxs) > 0 {
		ridx := make(map[int]struct{})
		for _, rmIdx := range rmIdxs {
			ridx[rmIdx] = struct{}{}
		}
		remainderWords := make([]string, 0)
		for idx, word := range words {
			if _, ok := ridx[idx+1]; !ok {
				remainderWords = append(remainderWords, word)
			}
		}
		return remainderWords
	}
	return words
}
