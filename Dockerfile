# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git (needed for some Go modules) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/carrpigeo ./cmd/api/main.go

# ---- Runtime Stage ----
FROM alpine:3.22

WORKDIR /app

# Install ca-certificates for HTTPS/SMTP TLS connections
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/carrpigeo .

# Create a directory for logs
RUN mkdir -p /app/storage

EXPOSE 2500

ENTRYPOINT ["./carrpigeo"]
