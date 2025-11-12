# ? Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# ? Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# ? Copy source code
COPY . .

# ? Build argument for service name
ARG SERVICE=api-gateway

# ? Build the specified service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE} ./services/${SERVICE}

# ? Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# ? Copy the built binary
ARG SERVICE=api-gateway
COPY --from=builder /app/bin/${SERVICE} ./app

CMD ["./app"]
