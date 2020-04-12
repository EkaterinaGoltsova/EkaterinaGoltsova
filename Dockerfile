FROM golang:1.10-alpine AS build

RUN apk --update add --no-cache git make

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/sprint-starter
ADD . /go/src/sprint-starter

RUN dep ensure -v
RUN go build -o .out/sprint-starter

COPY --from=build /go/src/sprint-starter/.out/sprint-starter /usr/bin/sprint-starter
RUN chmod +x /usr/bin/sprint-starter


EXPOSE 8080
ENTRYPOINT ["/usr/bin/sprint-starter"]
