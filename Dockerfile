# Builder Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=auto

# Install curl to download the standalone Tailwind CLI binary
RUN apk add --no-cache curl ca-certificates

# Download standalone Tailwind CLI binary for Linux x64
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# Cache Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile minified production CSS into ./static/css/styles.css
RUN mkdir -p static/css \
    && tailwindcss -i ./input.css -o ./static/css/styles.css --minify

# Build Go application
RUN go build -o afyatrack .

# Final Runtime Stage
FROM alpine:latest
WORKDIR /app

# Copy compiled application binary and static assets
COPY --from=builder /app/afyatrack .
COPY --from=builder /app/static ./static
COPY --from=builder /app/index.html .

EXPOSE 8080
CMD ["./afyatrack"]