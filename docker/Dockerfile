# Dockerfile for building wile-extension-example binaries
#
# Build:
#   make docker-build
# Or manually:
#   docker build -f docker/Dockerfile -t wile-extension-example .
#
# Cross-platform:
#   docker build --platform linux/amd64 -f docker/Dockerfile -t wile-extension-example .

FROM golang:1.23-bookworm

ARG TARGETOS=linux
ARG TARGETARCH
ENV TARGETOS=${TARGETOS}
ENV TARGETARCH=${TARGETARCH}

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    make \
    git \
    curl

WORKDIR /app

# Copy go module files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source
COPY . .

# Build all example binaries
RUN make build

# Default command: run tests
CMD ["make", "test"]
