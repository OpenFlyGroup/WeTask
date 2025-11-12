# ? Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# ? Install swag for Swagger documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# ? Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# ? Copy source code
COPY . .

# ? Build argument for service name
ARG SERVICE=api-gateway

# ? Generate Swagger documentation for api-gateway
RUN if [ "${SERVICE}" = "api-gateway" ]; then \
        swag init -g services/api-gateway/main.go -o services/api-gateway/docs; \
    fi

# ? Build the specified service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE} ./services/${SERVICE}

# ? Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# ? Copy the built binary
ARG SERVICE=api-gateway
COPY --from=builder /app/bin/${SERVICE} ./app

# ? Copy Swagger documentation for api-gateway (if needed for file serving)
# Note: Docs are embedded in binary when imported, but keeping for compatibility
COPY --from=builder /app/services/api-gateway/docs ./docs

CMD ["./app"]
