# Build
FROM golang:alpine AS build

# Destination of copy
WORKDIR /build

# Download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build
RUN go build -o bin ./cmd/main.go

# Deploy
FROM alpine

RUN adduser -S -D -H -h /app appuser
USER appuser

COPY --from=build /build/bin /app/

EXPOSE 8080

WORKDIR /app

CMD ["./bin"]