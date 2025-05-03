FROM golang:1.24-alpine

ENV PATH="/go/bin:${PATH}"

RUN go install github.com/air-verse/air@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]