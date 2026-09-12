.PHONY: build

BUILD_PATH=build
BINARY_NAME=symlinker

build:
	@CGO_ENABLED=0 go build -o ${BUILD_PATH}/${BINARY_NAME} .

test:
	@./${BUILD_PATH}/${BINARY_NAME} config_test.json

clean:
	@rm -rf ${BUILD_PATH}/${BINARY_NAME}
	@echo "done cleaning..."
