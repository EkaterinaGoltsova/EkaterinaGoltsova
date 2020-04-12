FROM golang:1.10-alpine

RUN apk --update add --no-cache git make

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/sprint-starter
ADD . /go/src/sprint-starter

RUN dep ensure -v
RUN go build -o /go/src/sprint-starter/.out/sprint-starter

RUN chmod +x /go/src/sprint-starter/.out/sprint-starter

EXPOSE 8080
ENTRYPOINT ["/go/src/sprint-starter/.out/sprint-starter"]
