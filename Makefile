.PHONY:install
.ONESHELL:

SRC := $(wildcard cmd/*.go)

my:$(SRC)
	go build ./cmd/my  

install:
	go install ./cmd/my
	rm -f ~/.local/bin/my
	ln -s $$(go env GOPATH)/bin/my ~/.local/bin/my
	