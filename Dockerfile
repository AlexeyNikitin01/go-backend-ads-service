FROM golang:alpine as build

ENV BIN_FILE /opt/ads/app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

RUN CGO_ENABLED=0 go build \
        -ldflags "-s -w" \
        -o ${BIN_FILE} cmd/main/main.go

FROM alpine:latest

ENV BIN_FILE "/opt/ads/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
