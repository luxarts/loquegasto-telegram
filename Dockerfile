FROM golang:alpine AS builder

RUN mkdir /build
WORKDIR /build
COPY . .
RUN go build -o output ./cmd/main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser # System user, no password, no home dir,
USER appuser
COPY --from=builder /build/output /app/
WORKDIR /app
CMD ["./output"]