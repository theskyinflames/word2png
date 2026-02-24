# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**word2png** encrypts sequences of words using AES-256 and encodes the encrypted data into PNG images by mapping each byte value to a specific color. The reverse process decodes PNG images back to the original words. Primary use case: securely storing cryptocurrency seed phrases as encrypted PNG images.

## Commands

```bash
make test          # go test -v -race ./...
make lint          # golangci-lint run + go mod tidy check
make install       # build and install word2png and png2word binaries
make build-wasm    # build WASM module → ./assets/world2png.wasm
make generate      # go generate ./... (regenerates moq mocks)
make tools         # install dev tools (golangci-lint, gofumpt, moq)
```

Run a single test:
```bash
go test -v -race ./lib/ -run TestEncodingDecoding
go test -v -race ./cmd/word2png/ -run TestRemoveWordsByIdx
```

## Architecture

### Packages

- **`lib/`** — Core library. All encoding/decoding logic lives here.
  - `aes256.go` — AES-256-GCM chain encryption. Each word is encrypted using the previous ciphertext as the passphrase (prevents reordering attacks). Key is derived via MD5 to produce 32 bytes.
  - `color.go` — Bijective rune↔color mapping. Uses MD5 of the seed to reorder the WebSafe palette deterministically. Black/white are excluded (used as markers).
  - `byte.go` — Splits bytes into high/low nibbles (4 bits each) and rejoins them. This overcomes the 256-value limit given a 128-color palette.
  - `encoder.go` — Orchestrates words → AES-256 → nibbles → colors → paletted PNG.
  - `decoder.go` — Reverse: PNG → colors → nibbles → AES-256 → words.
  - `zmock_*.go` — Generated mocks (do not edit manually; use `make generate`).

- **`cmd/word2png/`** — Encoder CLI (kingpin flags, masked password input via pterm, supports appending words to existing images, word removal by index).
- **`cmd/png2word/`** — Decoder CLI (kingpin flags, regex filter support, pterm table output).
- **`cmd/wasm/`** — WebAssembly decoder for browser use (exports `decode` to JS; skipped by golangci-lint).

### Encoding Flow

```
words → AES-256 chain encrypt → split bytes into nibbles → map nibbles to colors → PNG (one row per word, black pixel = word terminator)
```

### Decoding Flow

```
PNG → extract colors row-by-row (stop at black pixel) → map colors to nibbles → join nibbles into bytes → AES-256 chain decrypt → words
```

### Key Design Decisions

- **Chain encryption**: each word's ciphertext becomes the passphrase for the next word, so word order is cryptographically enforced.
- **Enumeration**: words are enumerated before encryption to support duplicate entries.
- **Nibble splitting**: each encrypted byte is split into two 4-bit values so 128 palette colors suffice to represent all 256 byte values (two pixels per byte).
- **GCM mode**: provides authenticated encryption; nonce is prepended to ciphertext.
- **WASM exclusion**: the `wasm` directory is excluded from linting (see `.golangci.yml`).

## Testing

Tests use testify assertions and moq-generated mocks. Test fixtures are defined in `lib/fixtures_test.go` (13 words including Unicode, seed `"bartolo"`). The `-race` flag is always used to catch data races.
