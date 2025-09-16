# 1. Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Kopiramo mod fajlove i povlačimo dependencije
COPY go.mod go.sum ./
RUN go mod download

# Kopiramo ceo kod i kompajliramo binarni fajl
COPY . .
RUN go build -o authservice .

# 2. Minimalna finalna slika
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/authservice .

EXPOSE 8080

CMD ["./authservice"]
