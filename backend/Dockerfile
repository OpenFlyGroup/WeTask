# ? syntax=docker/dockerfile:1.4
FROM golang:1.21-alpine AS builder

WORKDIR /app

# ? Install swag
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# ? Copy go mod
COPY go.mod go.sum ./
RUN go mod download

# ? Copy source
COPY . .

# ? Build arg
ARG SERVICE=api-gateway

# ? Generate docs only for api-gateway
RUN if [ "$SERVICE" = "api-gateway" ]; then \
    swag init -g services/api-gateway/main.go -o services/api-gateway/docs; \
    fi

# ? Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE} ./services/${SERVICE}


# ? ============ DOCS EXTRACTOR (Optional Stage) ============
# ? This stage exists ONLY to extract docs if they exist
FROM builder AS docs-extractor
RUN --mount=type=cache,target=/tmp \
    if [ -d "/app/services/api-gateway/docs" ]; then \
    mkdir -p /extracted-docs && \
    cp -r /app/services/api-gateway/docs/* /extracted-docs/ 2>/dev/null || true; \
    else \
    mkdir -p /extracted-docs; \
    fi


# ? ============ RUNTIME ============
FROM alpine:latest AS runtime

RUN apk --no-cache add ca-certificates
WORKDIR /app

ARG SERVICE=api-gateway

# ? Copy binary
COPY --from=builder /app/bin/${SERVICE} ./app
RUN chmod +x ./app

# ? Copy docs ONLY if they exist (from docs-extractor stage)
# ? This path ALWAYS exists (even if empty)
COPY --from=docs-extractor /extracted-docs ./docs

CMD ["./app"]