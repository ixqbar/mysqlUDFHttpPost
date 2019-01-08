SHELL:=/bin/bash
TARGET=mysqlHttp

all: linux

win:
	GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o ./bin/${TARGET}.dll ./src

linux:
	GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o ./bin/${TARGET}.so ./src

mac:
	GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o ./bin/${TARGET}.so ./src

clean:
	rm -rf ./bin/${TARGET}*

