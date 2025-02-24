.PHONY: build clean init deps run

BINARY_NAME=PwnHash
BINARY_DIR=PwnagotchiHashTool
MAIN_LOC=cmd/PwnHash/main.go

build: clean init deps
	go build -o $(BINARY_DIR)/${BINARY_NAME} $(MAIN_LOC)

clean:
	go clean
	rm -rf ${BINARY_DIR}

init:
	mkdir -p ${BINARY_DIR}

run:
	cd ${BINARY_DIR} && ./${BINARY_NAME}

deps:
	go mod tidy
	go mod download