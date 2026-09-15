FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/hws ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/hws .

EXPOSE 7070

ENTRYPOINT ["./hws"]
