FROM golang:1.12.6-alpine3.9

ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN apk add --update --no-cache git
WORKDIR /test

COPY . .
RUN go mod download
