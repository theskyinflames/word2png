//go:build wasm
// +build wasm

//nolint
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/theskyinflames/word2png/lib"
)

const errMsg = "W2P ERROR: %s"

func decode(b []byte, filter string, seed string) ([]string, error) {
	decrypter := lib.NewAES256(seed)
	decoder := lib.NewDecoder(lib.Rune2Color(seed), decrypter)
	return decoder.Decode(b)
}

func jsDecoder() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return fmt.Sprintf(errMsg, "Invalid no of arguments passed")
		}

		// decode the binary array of the secret image file recieved from JS
		// IMPORTANT: the byte array must be passed as an Uint8Array array from JS
		// see https://github.com/theskyinflames/word2pngUI repo to have an example of usage
		received := make([]byte, args[0].Get("length").Int())
		_ = js.CopyBytesToGo(received, args[0])

		filter := args[1].String()
		seed := args[2].String()

		// decode the secret
		words, err := decode(received, filter, seed)
		if err != nil {
			return fmt.Sprintf(errMsg, err.Error())
		}

		b, _ := json.Marshal(words)
		return string(b)
	})
	return jsonFunc
}

func PassUint8ArrayToGo(this js.Value, args []js.Value) interface{} {
	received := make([]byte, args[0].Get("length").Int())
	_ = js.CopyBytesToGo(received, args[0])
	return nil
}

func main() {
	js.Global().Set("decode", jsDecoder())
	<-make(chan bool)
}
