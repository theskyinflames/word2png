package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/theskyinflames/word2png/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

const defaultFilter = ".*"

func main() {
	var (
		file   = kingpin.Flag("file", "Coded image to be used as words code if it's filled").Short('f').String()
		b64    = kingpin.Flag("b64", "b64 string with the coded image").String()
		debug  = kingpin.Flag("debug", "writes a debug file").Short('d').Bool()
		filter = kingpin.Flag("filter", "only shows the words that match the regex expression").String()

		// dangerous zone
		showSeed = kingpin.Flag("show-seed", "shows the entered seed").Short('s').Bool()

		debugFile *os.File
		err       error
	)
	kingpin.Parse()

	seed, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("\nEnter your seed")
	if *showSeed {
		pterm.DefaultBasicText.Printf("Entered seed: %s\n", pterm.BgYellow.Sprintf(pterm.Black(seed)))
	}

	if *filter == "" {
		*filter = defaultFilter
	}
	matchFilter, err := regexp.Compile(*filter)
	exitIfError(err)

	if debug != nil && *debug {
		debugFile, err = os.Create("./decrypted-bytes.txt")
		exitIfError(err)
		defer func() {
			debugFile.Close()
		}()
	}

	decrypter := lib.NewAES256(seed)
	decoder := lib.NewDecoder(lib.Rune2Color(seed), decrypter, lib.DecodeDebugWriterOpt(debugFile))
	words, err := decoder.DecodeFromSource(*file, *b64)
	exitIfError(err)

	fmt.Println("decoding process finished.")
	fmt.Printf("Have been decoded %d words:\n\n", len(words))
	values := pterm.TableData{
		{"Index", "Value"},
	}
	for i := range words {
		if matchFilter.MatchString(words[i]) {
			values = append(values, []string{
				fmt.Sprintf(" %s", strconv.Itoa(i+1)),
				words[i],
			})
		}
	}

	err = pterm.DefaultTable.WithBoxed(true).WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithLeftAlignment().WithData(values).Render()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
		os.Exit(-1)
	}
	os.Exit(0)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(-1)
	}
}
