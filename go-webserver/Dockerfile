FROM golang:1.18-alpine as builder

RUN mkdir -p /webserver

WORKDIR /webserver

COPY . .

RUN go mod download

RUN go build -o bin/webserver src/main.go

FROM alpine:3.16

ARG BUILD_DATE
ARG VERSION
ARG REVISION

LABEL maintainer="bogd"

RUN addgroup -S app \
    && adduser -S -G app app \
    && apk --no-cache add \
	ca-certificates curl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /webserver/bin/webserver .

RUN mkdir -p static
COPY static_page/ static

RUN chown -R app:app ./

USER app

CMD ["./webserver"]