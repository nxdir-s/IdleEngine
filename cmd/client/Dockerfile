FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bootstrap -ldflags="-s -w" cmd/client/main.go

FROM arm64v8/alpine:latest

WORKDIR /app

COPY --from=builder /app/bootstrap /app/bootstrap

EXPOSE 3000
EXPOSE 6060

CMD ["/app/bootstrap", "ws://127.0.0.1:3000/ws"]
