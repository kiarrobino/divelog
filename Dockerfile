# Stage 1 — build
FROM golang:1.25 AS builder

WORKDIR /app

# copy dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o divelog ./cmd/divelog

# Stage 2 — run
FROM debian:bookworm-slim

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/divelog .

# copy the web UI
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./divelog", "serve"]