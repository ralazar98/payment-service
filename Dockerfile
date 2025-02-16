FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY configs/config.yaml /app/cmd/configs/


RUN go mod download

COPY .. .


WORKDIR /app/cmd

RUN go build -o payment-service .

EXPOSE 8081

CMD ["./payment-service"]