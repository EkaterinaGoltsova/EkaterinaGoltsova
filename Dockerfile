FROM golang:1.14-alpine as build

RUN apk --update add --no-cache git make

WORKDIR /go/src/sprint-starter
ADD . /go/src/sprint-starter

RUN make

FROM alpine

WORKDIR /go/src/sprint-starter

COPY --from=build /go/src/sprint-starter .

RUN chmod +x /go/src/sprint-starter/.out/sprint-starter

EXPOSE 8080

ENTRYPOINT ["/go/src/sprint-starter/.out/sprint-starter"]
