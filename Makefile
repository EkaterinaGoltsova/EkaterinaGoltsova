default: all

all: build

build:
	@echo build
	go build -o /go/src/sprint-starter/.out/sprint-starter ./cmd
	@echo done
