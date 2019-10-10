SHELL := /bin/bash

BINARY_NAME=ctest

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o ${BINARY_NAME}

clean:
	go clean
	rm -f ${BINARY_NAME}

deploy: build
	scp ${BINARY_NAME} root@192.168.88.188:/usr/bin