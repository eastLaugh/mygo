.PHONY:install my
.ONESHELL:

my:
	go build ./cmd/my  

install:
	go install ./cmd/my
	rm -f ~/.local/bin/my
	ln -s $$(go env GOPATH)/bin/my ~/.local/bin/my
	