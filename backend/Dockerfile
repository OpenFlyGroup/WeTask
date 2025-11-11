FROM node:20-alpine AS base

WORKDIR /app

# Install dependencies
RUN apk add --no-cache bash coreutils

# Copy package files
COPY package*.json ./
COPY pnpm-lock.yaml* ./

# Install pnpm globally
RUN npm install -g pnpm

# Install dependencies
RUN pnpm install --no-frozen-lockfile

# Copy entrypoint script first
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

# Copy all files (excluding node_modules via .dockerignore)
COPY . .

# Build all services
FROM base AS api-gateway
WORKDIR /app
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["npm", "run", "start:dev:gateway"]

FROM base AS auth-service
WORKDIR /app
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["npm", "run", "start:dev:auth"]

FROM base AS users-service
WORKDIR /app
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["npm", "run", "start:dev:users"]

FROM base AS teams-service
WORKDIR /app
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["npm", "run", "start:dev:teams"]

FROM base AS boards-service
WORKDIR /app
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["npm", "run", "start:dev:boards"]

FROM base AS tasks-service
WORKDIR /app
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["npm", "run", "start:dev:tasks"]

