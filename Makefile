.PHONY: build install wordlist

install:
	go install github.com/tendermint/mintkey/cmd/...

build: wordlist
	go build github.com/tendermint/mintkey/cmd/...

wordlist:
	go-bindata -ignore ".*\.go" -o wordlist/wordlist.go -pkg "wordlist" wordlist/...
