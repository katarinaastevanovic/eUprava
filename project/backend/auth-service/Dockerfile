FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o authservice .

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/authservice .
COPY --from=builder /app/serviceAccountKey.json .

ENV GOOGLE_APPLICATION_CREDENTIALS="/app/serviceAccountKey.json"

EXPOSE 8081

CMD ["./authservice"]
