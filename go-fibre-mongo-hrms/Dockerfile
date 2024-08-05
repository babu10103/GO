# builder
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go_fiber main.go

# runner
FROM ubuntu:22.04

WORKDIR /root/

COPY --from=builder /app/go_fiber .

EXPOSE 9999

CMD ["./go_fiber"]
