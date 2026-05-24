# AGENTS.md

## Gotchas & Sharp Edges

- ~~**`make generate` vendor dance:** It runs `go mod vendor` before `go generate ./...`, then removes `./vendor`. This may not work if vendor directory is gitignored or has stale contents.~~ ✅ Fixed — moq is now a `tool` dependency in `go.mod`, `//go:generate` uses `go run`, and `make generate` is just `go generate ./...`.
- **Formatters are enforced:** `gofumpt` + `goimports` run as linters. Run `golangci-lint run` locally (or `make lint`) before pushing — it also verifies `go mod tidy` didn't change anything (`git diff --quiet go.mod go.sum`).
- **No real-crypto e2e test:** `lib/fixtures_test.go` uses generated mocks (`EncrypterMock`/`DecrypterMock`), so encryption/decryption round-trips are never tested with real AES-256. Don't assume integration coverage exists.
- **WASM build typo:** `make build-wasm` outputs `assets/world2png.wasm` (missing 'd'). Preserve filename for backward compatibility with `word2pngUI`.
- **`os.Exit()` + non-standard code:** Both CLIs (`cmd/word2png/`, `cmd/png2word/`) call `os.Exit(-1)` on failure, skipping deferred cleanup. Handle with care.
- **`cmd/wasm/` excluded from golangci-lint** (see `.golangci.yml` `build-tags: [infra]` and WASM build constraint).

## Commands Quick Reference

```bash
make test          # go test -v -race ./...
make lint          # golangci-lint run + go mod tidy -v && git diff --quiet go.mod go.sum
make install       # builds + installs word2png and png2word binaries
make build-wasm    # GOOS=js GOARCH=wasm → assets/world2png.wasm
make generate      # go generate ./...  (moq is a tool dep in go.mod)
make tools         # install golangci-lint v1.40.1, gofumpt, moq
```
