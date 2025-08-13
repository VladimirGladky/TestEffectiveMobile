FROM golang:1.24 AS builder

WORKDIR /newApp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o effective_mobile ./cmd/main.go

EXPOSE 4047

CMD ["./effective_mobile"]