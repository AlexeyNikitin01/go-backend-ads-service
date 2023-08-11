FROM golang:alpine

WORKDIR /build

COPY . .

RUN go mod tidy
