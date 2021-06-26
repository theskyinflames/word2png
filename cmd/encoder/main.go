package main

import "encoding"

type Encoder interface {
	encoding.TextMarshaler
}
