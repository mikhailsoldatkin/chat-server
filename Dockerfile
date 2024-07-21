FROM golang:1.22 AS builder

WORKDIR /chat

RUN apt-get update &&  \
    apt-get upgrade -y &&  \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o chat_server ./cmd/grpc_server/main.go

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

WORKDIR /chat

COPY --from=builder /chat/chat_server .
COPY .env .

CMD ["./chat_server"]
