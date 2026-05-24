package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pterm/pterm"

	"github.com/theskyinflames/word2png/lib"
)

const maxFileWordsSize = 1 * 1024 * 1024 // 1MB

func main() {
	os.Exit(run())
}

func run() int {
	var (
		imagePath   = kingpin.Flag("file", "Save to the especified file if it's filled").Short('f').String()
		words       = kingpin.Flag("words", "list of words to encode").Short('w').Strings()
		fileWords   = kingpin.Flag("file-words", "path to a text file whose content will be encoded as a single word (max 1MB)").String()
		debug       = kingpin.Flag("debug", "writes a debug file").Short('d').Bool()
		b64         = kingpin.Flag("b64", "b64encoded image").String()
		removeWords = kingpin.Flag("remove-word", "remove a word from an image by index number").Short('r').Ints()
		showSeed    = kingpin.Flag("show-seed", "shows the entered seed").Short('s').Bool()
	)
	kingpin.Parse()

	if *fileWords != "" {
		fileWordsList, err := readWordsFromFile(*fileWords)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return 1
		}
		*words = append(*words, fileWordsList...)
	}

	seed, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter your seed")
	if *showSeed {
		pterm.DefaultBasicText.Printf("Entered seed: %s\n", pterm.BgYellow.Sprint(pterm.Black(seed)))
	}

	var debugFile *os.File
	if debug != nil && *debug {
		var err error
		debugFile, err = os.Create("./encrypted-bytes.txt")
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return 1
		}
		defer debugFile.Close()
	}

	aes256 := lib.NewAES256(seed)

	if (*imagePath != "" && imageExists(*imagePath)) || *b64 != "" {
		decoder := lib.NewDecoder(lib.Rune2Color(seed), aes256, lib.DecodeDebugWriterOpt(debugFile))
		beforeWords, err := decoder.DecodeFromSource(*imagePath, *b64)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return 1
		}
		*words = append(beforeWords, *words...)
	}

	*words = RemoveWordsByIdx(*words, *removeWords)

	encoder := lib.NewEncoder(lib.Rune2Color(seed), aes256, lib.EncoderDebugWriterOpt(debugFile))
	b, err := encoder.Encode(*words)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return 1
	}

	switch {
	case *imagePath != "":
		if err := lib.SaveEncodedImage(b, *imagePath); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return 1
		}
	default:
		b64Encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
		if _, err = b64Encoder.Write(b); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return 1
		}
		if err := b64Encoder.Close(); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return 1
		}
	}

	fmt.Println("\ncoding process finished")
	return 0
}

func imageExists(imagePath string) bool {
	if _, err := os.Stat(imagePath); err == nil {
		return true
	}
	return false
}

func readWordsFromFile(path string) ([]string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access file %s: %w", path, err)
	}
	if fi.Size() > maxFileWordsSize {
		return nil, fmt.Errorf("file %s exceeds maximum size of 1MB", path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", path, err)
	}
	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		return []string{}, nil
	}
	return []string{trimmed}, nil
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
