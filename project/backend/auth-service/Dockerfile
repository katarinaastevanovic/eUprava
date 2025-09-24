# 1. Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Kopiramo mod fajlove i povlačimo dependencije
COPY go.mod go.sum ./
RUN go mod download

# Kopiramo ceo kod i serviceAccountKey.json
COPY . .

# Build binarni fajl
RUN go build -o authservice .

# 2. Minimalna finalna slika
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Kopiramo binarni fajl i serviceAccountKey.json
COPY --from=builder /app/authservice .
COPY --from=builder /app/serviceAccountKey.json .

# Optional: environment varijabla ako tvoj kod koristi path
ENV GOOGLE_APPLICATION_CREDENTIALS="/app/serviceAccountKey.json"

EXPOSE 8080

CMD ["./authservice"]
