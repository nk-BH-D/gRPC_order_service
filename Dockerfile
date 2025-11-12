# Этап 1: билд
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app-bin ./cmd/server/main.go

# Этап 2: запуск
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/app-bin /app-bin
USER nonroot
ENTRYPOINT [ "/app-bin" ]
