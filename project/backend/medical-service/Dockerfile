# 1. Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o medicalservice .

# 2. Minimalna finalna slika
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/medicalservice .

EXPOSE 8080

CMD ["./medicalservice"]
