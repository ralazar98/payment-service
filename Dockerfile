FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .


WORKDIR /app/cmd

RUN go build -o payment-service .

EXPOSE 8081

CMD ["./payment-service"]