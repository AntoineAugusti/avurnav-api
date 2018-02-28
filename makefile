.PHONY: all build

all: build

BUILD_TAG = $(shell git log --pretty=format:'%h' -n 1)
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

build:
	@go build -ldflags "-X main.buildTag=$(BUILD_TAG) -X main.buildDate=$(BUILD_DATE)"
	@echo "Build complete - $(BUILD_TAG) on $(BUILD_DATE)"

clean:
	rm -f moduluschecking-api
