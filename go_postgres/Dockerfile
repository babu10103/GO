FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go_pgapi .

# run stage
FROM ubuntu:22.04

WORKDIR /root/

COPY --from=builder /app/go_pgapi .

EXPOSE 9030

CMD ["./go_pgapi"]
