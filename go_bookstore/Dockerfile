# builder
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dist/main cmd/main/main.go

# runtime
FROM ubuntu:22.04

WORKDIR /root/

COPY --from=builder /app/dist/main /app/wait-for-it.sh .

RUN ls -l

EXPOSE 9010

CMD ["./main"]

