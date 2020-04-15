default: all

all: deps build

deps:
	@echo fetch dependencies
	go mod tidy
	go mod vendor
	@echo done

build:
	@echo build
	go build -o /go/src/sprint-starter/.out/sprint-starter
	@echo done
