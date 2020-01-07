FROM golang:1.13-alpine3.10 AS builder

WORKDIR /tcc.go

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build it
RUN CGO_ENABLED=0 go build -o ./build/blogrest \
    ./cmd/blogrest

# Build the runtime image
FROM alpine:3.10
RUN apk add --no-cache --update ca-certificates
WORKDIR /root

COPY --from=builder /tcc.go/build/blogrest ./service

EXPOSE 8080

ENTRYPOINT ["./service"]