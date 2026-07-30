.PHONY: build

MAIN_GO=./cmd/main.go
BUILD_PATH=build
BINARY_NAME=symlinker

build:
	@CGO_ENABLED=0 go build -o ${BUILD_PATH}/${BINARY_NAME} ${MAIN_GO}

clean:
	@rm -rf ${BUILD_PATH}/${BINARY_NAME}
	@echo "done cleaning..."
