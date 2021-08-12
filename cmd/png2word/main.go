package main

import (
	"fmt"
	"os"

	"github.com/theskyinflames/image-coder/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		file  = kingpin.Flag("file", "Coded image to be used as words code if it's filled").Short('f').String()
		b64   = kingpin.Flag("b64", "b64 string with the coded image").String()
		seed  = kingpin.Flag("seed", "coding seed").Short('s').Required().String()
		debug = kingpin.Flag("debug", "writes a debug file").Short('d').Bool()

		debugFile *os.File
		err       error
	)
	kingpin.Parse()

	if debug != nil && *debug {
		debugFile, err = os.Create("./decrypted-bytes.txt")
		exitIfError(err)
		defer func() {
			debugFile.Close()
		}()
	}

	decrypter := lib.NewAES256(*seed)
	decoder := lib.NewDecoder(lib.Rune2Color(*seed), decrypter, lib.DecodeDebugWriterOpt(debugFile))
	words, err := decoder.DecodeFromSource(*file, *b64)
	exitIfError(err)

	fmt.Println("decoding process finished.")
	fmt.Printf("Have been decoded %d words:", len(words))
	for i := range words {
		fmt.Printf("\n%d - %s", i+1, words[i])
	}
	os.Exit(0)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(-1)
	}
}
