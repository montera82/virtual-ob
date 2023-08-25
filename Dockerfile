FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum .env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o virtual-orb ./cmd/virtual-orb/main.go

FROM alpine:3.14

WORKDIR /root/

COPY --from=builder /app/virtual-orb .

COPY --from=builder /app/.env .

EXPOSE 8002

CMD ["./virtual-orb"]
