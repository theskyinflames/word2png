package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/theskyinflames/image-coder/tooling"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		file  = kingpin.Flag("file", "Save to the especified file if it's filled").Short('f').String()
		seed  = kingpin.Flag("seed", "coding seed").Short('s').Required().String()
		words = kingpin.Flag("words", "list of words to encode").Short('w').Strings()
	)
	kingpin.Parse()

	debugFile, err := os.Create("./encrypted-bytes.txt")
	exitIfError(err)
	defer func() {
		debugFile.Close()
	}()

	encoder := tooling.NewEncoder(*seed, tooling.Rune2Color, tooling.EncoderDebugWriterOpt(debugFile))
	b, err := encoder.Encode(*words)
	exitIfError(err)

	switch {
	case file != nil && *file != "":
		exitIfError(tooling.SaveEncodedImage(b, *file))
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
