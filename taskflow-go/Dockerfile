# Builder stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api cmd/main.go

# Final stage
FROM alpine:3.19
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/api .
COPY .env.example .env
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./api"]
