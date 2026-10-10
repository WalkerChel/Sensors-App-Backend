# Build stage
FROM golang:alpine AS builder

# Set necessary environment variables needed for image
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

COPY . .

# Build app
RUN go mod download
RUN go mod tidy
RUN go build -ldflags="-s -w" -trimpath -o bin/SensorsApp cmd/sensors-app/main.go

# Run stage
FROM debian:buster-slim

WORKDIR /build

COPY --from=builder ./build/bin/SensorsApp .
COPY --from=builder ./build/migrations ./migrations
COPY --from=builder ./build/entrypoint.sh .

EXPOSE 8080 

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["/build/entrypoint.sh"]
