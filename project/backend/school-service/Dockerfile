# 1. Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o schoolservice .

# 2. Minimalna finalna slika
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/schoolservice .

EXPOSE 8080

CMD ["./schoolservice"]
