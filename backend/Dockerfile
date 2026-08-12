# syntax=docker/dockerfile:1.7

############################
# Dashboard build
############################
FROM node:22-alpine AS dashboard-build

WORKDIR /app/WebDashBoard/roxkv-dashboard

COPY WebDashBoard/roxkv-dashboard/package.json WebDashBoard/roxkv-dashboard/package-lock.json ./
RUN npm ci

COPY WebDashBoard/roxkv-dashboard/ ./
RUN npm run build

############################
# Base
############################
FROM golang:1.25-alpine AS base

WORKDIR /app

RUN apk update && \
    apk add --no-cache git ca-certificates wget && \
    update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

############################
# Development
############################
FROM base AS development

WORKDIR /app

COPY . .
COPY --from=dashboard-build /app/server/dashboard/dist ./server/dashboard/dist

COPY scripts/backend-entrypoint.sh /usr/local/bin/backend-entrypoint.sh
RUN chmod +x /usr/local/bin/backend-entrypoint.sh

EXPOSE 6969 6970 6971 6972 6973 8080

ENTRYPOINT ["/usr/local/bin/backend-entrypoint.sh"]
CMD ["go", "run", "./cmd/roxkv", "tcp"]

############################
# Build
############################
FROM base AS build

WORKDIR /app

COPY . .
COPY --from=dashboard-build /app/server/dashboard/dist ./server/dashboard/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/roxkv ./cmd/roxkv

############################
# Production
############################
FROM alpine:3.22

WORKDIR /app

RUN apk update && \
    apk add --no-cache ca-certificates wget && \
    update-ca-certificates

COPY --from=build /out/roxkv /usr/local/bin/roxkv
COPY scripts/backend-entrypoint.sh /usr/local/bin/backend-entrypoint.sh

RUN chmod +x /usr/local/bin/backend-entrypoint.sh

EXPOSE 6969 6970 6971 6972 6973 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=5 \
CMD wget -qO- http://127.0.0.1:6971/healthz >/dev/null && \
    wget -qO- http://127.0.0.1:6972/healthz >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/backend-entrypoint.sh"]
CMD ["roxkv", "tcp"]
