# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o app .

# Runtime stage
FROM alpine:latest
WORKDIR /root/

# Copy binary
COPY --from=builder /app/app .

# Copy .env jika ada
COPY --from=builder /app/.env .

# ✅ Copy folder assets ke image runtime
COPY --from=builder /app/assets ./assets

EXPOSE 8080

CMD ["./app"]
