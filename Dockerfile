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
RUN go build -o bin/SensorsApp cmd/sensors-app/main.go

# Run stage
FROM debian:buster-slim

WORKDIR /build

COPY --from=builder ./build/bin/SensorsApp .
COPY --from=builder ./build/.env .
COPY --from=builder ./build/migrations ./migrations
COPY --from=builder ./build/entrypoint.sh .

ENV GOLANG_MIGRATE_VERSION=v4.16.2

# Install curl and migrate
RUN apt-get update --yes && \
    apt-get install --yes --no-install-recommends \
    curl \
    ca-certificates \
    make && \
    rm -rf /var/lib/apt/lists/*

# Install the migrator binary
RUN curl -LO https://github.com/golang-migrate/migrate/releases/download/${GOLANG_MIGRATE_VERSION}/migrate.linux-amd64.tar.gz &&\
    tar -C /usr/local/bin -xzvf migrate.linux-amd64.tar.gz &&\
    rm -f migrate.linux-amd64.tar.gz &&\
    # mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate &&\
    chmod 0755 /usr/local/bin/migrate &&\
    ln -s /usr/local/bin/migrate /usr/bin/migrate &&\
    migrate -version 

EXPOSE 8080 

ENTRYPOINT ["/build/entrypoint.sh"]
