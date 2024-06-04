# Use the official Golang image to create a build artifact.
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:3.16
COPY --from=builder /app/main /main
EXPOSE 8080
CMD ["/main"]
