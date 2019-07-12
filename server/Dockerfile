FROM golang:1.12.6-alpine3.9 AS build

ENV GO111MODULE on
RUN apk add --update --no-cache git

WORKDIR /workspace

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /server


FROM alpine:3.9

WORKDIR /server

RUN apk add --update --no-cache bash

COPY --from=build /server .
COPY ./wait-for-it.sh .

EXPOSE 8080

