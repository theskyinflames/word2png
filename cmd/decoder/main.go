package main

import "encoding"

type Decoder interface {
	encoding.TextUnmarshaler
}
