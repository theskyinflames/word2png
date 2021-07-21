default: test lint

install:
	cd cmd/encoder && go install .
	cd cmd/decoder && go install .

test:
	go test -v -race ./...

lint:
	golangci-lint run
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

tools: tool-golangci-lint tool-fumpt

tool-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.40.1

tool-fumpt:
	go get -u mvdan.cc/gofumpt
	go get -u mvdan.cc/gofumpt/gofumports

todo:
	find . -name '*.go' \! -name '*_generated.go' -prune | xargs grep -n TODO

