FROM golang:1.14-alpine as build

RUN apk --update add --no-cache git make

WORKDIR /go/src/sprint-starter
ADD . /go/src/sprint-starter

RUN make

FROM alpine:3.11.6

WORKDIR /go/src/sprint-starter

COPY --from=build /go/src/sprint-starter/.out/sprint-starter /usr/bin/sprint-starter
COPY --from=build /go/src/sprint-starter/configs/config.yaml /etc/sprint-starter/config.yaml
COPY --from=build /go/src/sprint-starter/web/templates /etc/sprint-starter/templates

RUN chmod +x /usr/bin/sprint-starter

EXPOSE 8080

ENTRYPOINT ["/usr/bin/sprint-starter", "--config=/etc/sprint-starter/config.yaml", "--templates=/etc/sprint-starter/templates/", "--port=8080"]
