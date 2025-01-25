FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o geemail .


FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/geemail .

CMD ["./geemail"]

