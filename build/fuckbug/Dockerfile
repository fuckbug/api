FROM golang:1.24-alpine AS build

RUN apk add --no-cache git gcc musl-dev

RUN go install github.com/swaggo/swag/cmd/swag@latest

ENV BIN_FILE /opt/app/bin
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

RUN swag init -g ./cmd/fuckbug/main.go

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} cmd/fuckbug/*

FROM alpine:latest

LABEL SERVICE="fuckbug"
LABEL MAINTAINERS="zaytcev.com"

ENV BIN_FILE /opt/app/bin
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/app/config.json
COPY configs/fuckbug/config.json ${CONFIG_FILE}

EXPOSE 80

CMD ${BIN_FILE} -config ${CONFIG_FILE}
