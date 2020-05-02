default: all

all: build

build:
	@echo build
	go build -o /go/src/sprint-starter/.out/sprint-starter ./cmd/sprint-starter
	@echo done

start:
	@echo start docker build
	docker build -t="sprint-starter" .
	@echo start docker container
	docker run -p 8080:8080 sprint-starter