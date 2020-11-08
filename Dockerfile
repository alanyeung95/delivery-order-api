## Dockerfile is being used to build an image

# Builder
FROM golang:1.13-alpine3.10 AS builder

# Dependency
RUN apk add --no-cache git gcc g++ make

# Directory inside container
WORKDIR /app

# Copy host file to /app inside container
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

# Runner
FROM  golang:1.13-alpine3.10
WORKDIR /app
COPY --from=builder /app .

CMD ["./main"]