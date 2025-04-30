FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for dependency fetching
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Fix missing dependencies and run go mod tidy before building
RUN go mod tidy && go build -o url-shortener

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]