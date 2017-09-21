# メタ情報
NAME := radigo
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
 -X 'main.revision=$(REVISION)'

# 必要なツール類をセットアップする
## Setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/alecthomas/gometalinter
	gometalinter --install
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help

# テストを実行する
## Run tests
test: 
	go test $$(glide novendor)

# glideを使って依存パッケージをインストールする
## Install dependencies
deps: setup
	glide install

## Update dependencies
update: setup
	glide update

## Lint
lint: 
	gometalinter --deadline 30s $$(glide novendor)


## Format source codes 
fmt: 
	goimports -w $$(glide nv -x)

## build binaries ex. make bin/radigo
bin/%: cmd/%/main.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ $<

## Show help
help:
	@make2help $(MAKEFILE_ LIST)

.PHONY: setup deps update test lint help
