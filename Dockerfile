# Build stage
FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -trimpath -o bin/SensorsApp cmd/sensors-app/main.go

# Run stage
FROM debian:bookworm-slim

# Create a non-root user and group
RUN groupadd -g 1000 appgroup && \
    useradd -u 1000 -g appgroup -s /bin/sh appuser

WORKDIR /build

COPY --from=builder /build/bin/SensorsApp .
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/entrypoint.sh .

RUN chmod +x ./entrypoint.sh && \
    chown -R appuser:appgroup /build

# Switch to non-root user
USER appuser

EXPOSE 8080 

ENTRYPOINT ["/build/entrypoint.sh"]
