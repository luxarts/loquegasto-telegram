FROM golang:alpine3.16 AS builder

RUN mkdir /build
WORKDIR /build
COPY . .
RUN go build -o output ./cmd/main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/output /app/
WORKDIR /app
CMD ["./output"]