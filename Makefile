.PHONY: build
build:
	go generate
	go build -race -o bin/moco.exe -ldflags "-s -w"
	upx --lzma --best --compress-icons=2 bin/moco.exe
	copy config.yml bin
