.PHONY: all test

all:
	go build -o bin/sora_exporter main.go

test:
	go test -race -v .