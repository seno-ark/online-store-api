# Build stage
FROM golang:1.22.3-alpine3.20 AS builder

WORKDIR /app
COPY go.mod go.sum ./
# RUN go mod download
COPY . .
RUN go build cmd/api/main.go

# Run stage
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 9000
CMD ["./main"]