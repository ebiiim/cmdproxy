.PHONY: all build install

all: build

build:
	go build ./cmd/cmdprx
	go build ./cmd/cmdprx-cli

install:
	go install "github.com/ebiiim/cmdproxy/cmd/cmdprx"
	go install "github.com/ebiiim/cmdproxy/cmd/cmdprx-cli"
