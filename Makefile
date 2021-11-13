default: test lint

install:
	cd cmd/word2png && go install .
	cd cmd/png2word && go install .

build-wasm:
	GOOS=js GOARCH=wasm go build -tags wasm -o ./assets/world2png.wasm ./cmd/wasm/main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

tools: tool-golangci-lint tool-fumpt tool-moq

tool-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.40.1

tool-fumpt:
	go get -u mvdan.cc/gofumpt

tool-moq:
	go get -u github.com/matryer/moq

todo:
	find . -name '*.go' \! -name '*_generated.go' -prune | xargs grep -n TODO

generate:
	@go mod vendor
	go generate ./... | true
	@rm -rf ./vendor

