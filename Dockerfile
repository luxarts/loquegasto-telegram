# Build
FROM golang:alpine AS build

# Destination of copy
WORKDIR /build

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build
RUN go build -v -o /build/bin ./cmd/main.go

# Deploy
FROM debian

COPY --from=build /build/bin /app/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

ENTRYPOINT ["./bin"]