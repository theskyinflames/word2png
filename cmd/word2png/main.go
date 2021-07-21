package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/theskyinflames/image-coder/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		file  = kingpin.Flag("file", "Save to the especified file if it's filled").Short('f').String()
		seed  = kingpin.Flag("seed", "coding seed").Short('s').Required().String()
		words = kingpin.Flag("words", "list of words to encode").Short('w').Strings()
		debug = kingpin.Flag("debug", "writes a debug file").Short('d').Bool()

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

	encoder := lib.NewEncoder(*seed, lib.Rune2Color, lib.EncoderDebugWriterOpt(debugFile))
	b, err := encoder.Encode(*words)
	exitIfError(err)

	switch {
	case *file != "":
		exitIfError(lib.SaveEncodedImage(b, *file))
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
