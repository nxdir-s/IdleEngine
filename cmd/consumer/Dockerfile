FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/consumer/main.go

FROM arm64v8/alpine:latest

WORKDIR /app

COPY --from=builder /app/bootstrap /app/bootstrap

CMD ["/app/bootstrap"]
