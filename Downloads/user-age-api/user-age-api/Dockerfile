FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server .
COPY db/migrations ./db/migrations

EXPOSE 3000

CMD ["./server"]
