# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/do-ssh

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o do-ssh ./...

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /usr/src/do-ssh/do-ssh /usr/local/bin/do-ssh

CMD ["/usr/local/bin/do-ssh"]
