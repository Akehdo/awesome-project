FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/meeting-service \
    ./cmd/meeting-service

FROM alpine:3.24

RUN apk add --no-cache ca-certificates \
    && adduser -D appuser

USER appuser
WORKDIR /app

COPY --from=builder /out/meeting-service ./meeting-service

EXPOSE 8000

CMD ["./meeting-service"]