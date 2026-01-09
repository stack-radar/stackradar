# Multi-stage build for StackRadar (Go)
FROM golang:1.24-alpine AS builder

LABEL maintainer="StackRadar Contributors"
LABEL description="Enterprise tech stack detector for CI/CD automation"

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o techstack .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    git \
    ruby \
    ruby-dev \
    build-base \
    cmake \
    icu-dev \
    && gem install github-linguist --no-document \
    && apk del build-base ruby-dev cmake icu-dev \
    && rm -rf /var/cache/apk/*

WORKDIR /workspace

# Copy binary from builder
COPY --from=builder /build/techstack /usr/local/bin/techstack

ENTRYPOINT ["techstack"]
CMD ["get", "--help"]