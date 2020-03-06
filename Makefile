VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION) -s -w" -trimpath
SRC=$(shell find . -name '*.go')

.PHONY: binaries clean release install

binaries: alicloud-vault-linux-amd64 alicloud-vault-linux-arm64 alicloud-vault-darwin-amd64 alicloud-vault-windows-386.exe alicloud-vault-freebsd-amd64

clean:
	rm -f alicloud-vault alicloud-vault-linux-amd64 alicloud-vault-linux-arm64 alicloud-vault-darwin-amd64 alicloud-vault-darwin-amd64.dmg alicloud-vault-windows-386.exe alicloud-vault-freebsd-amd64 SHA256SUMS

release: binaries alicloud-vault-darwin-amd64.dmg SHA256SUMS
	@echo "\nTo update homebrew-cask run\n\n    cask-repair -v $(shell echo $(VERSION) | sed 's/v\(.*\)/\1/') alicloud-vault\n"

alicloud-vault-darwin-amd64: $(SRC)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

alicloud-vault-freebsd-amd64: $(SRC)
	GOOS=freebsd GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

alicloud-vault-linux-amd64: $(SRC)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

alicloud-vault-linux-arm64: $(SRC)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $@ .

alicloud-vault-windows-386.exe: $(SRC)
	GOOS=windows GOARCH=386 go build $(BUILD_FLAGS) -o $@ .

alicloud-vault-darwin-amd64.dmg: alicloud-vault-darwin-amd64
	./bin/create-dmg alicloud-vault-darwin-amd64 $@

SHA256SUMS: binaries alicloud-vault-darwin-amd64.dmg
	shasum -a 256 alicloud-vault-freebsd-amd64 alicloud-vault-linux-amd64 alicloud-vault-linux-arm64 alicloud-vault-windows-386.exe alicloud-vault-darwin-amd64.dmg > $@

install:
	rm -f alicloud-vault
	go build $(BUILD_FLAGS) .
	mv alicloud-vault ~/bin/