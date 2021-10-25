.PHONY: build run compile clean
BUILD_DIR=./out
BINARY_NAME=${BUILD_DIR}/ftranUI
SOURCE_MAIN_NAME=./main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINARY_NAME} ${SOURCE_MAIN_NAME}

compile:
	# 64-Bit
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINARY_NAME}-linux-amd64.bin ${SOURCE_MAIN_NAME}
	# Windows
	env CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows go build -ldflags "-s -w" -o ${BINARY_NAME}-windows-amd64.exe ${SOURCE_MAIN_NAME}
