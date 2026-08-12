# syntax=docker/dockerfile:1

# ─── Stage 1: Build the dashboard bundle ─────────────────────────────────────
FROM node:22-alpine AS dashboard-builder

WORKDIR /app/WebDashBoard/roxkv-dashboard

COPY WebDashBoard/roxkv-dashboard/package.json WebDashBoard/roxkv-dashboard/package-lock.json ./
RUN npm ci

COPY WebDashBoard/roxkv-dashboard/ ./
RUN npm run build

# ─── Stage 2: Build the Go binary ────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for module fetching (some modules need it)
RUN apk add --no-cache git

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically-linked binary
COPY . .
COPY --from=dashboard-builder /app/server/dashboard/dist ./server/dashboard/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/roxkv ./cmd/roxkv

# ─── Stage 3: Minimal runtime image ─────────────────────────────────────────
FROM alpine:3.21

WORKDIR /app

# wget is needed for HEALTHCHECK; ca-certificates for HTTPS (Ollama API)
RUN apk add --no-cache ca-certificates wget

# Copy the compiled binary
COPY --from=builder /out/roxkv /usr/local/bin/roxkv

# Copy the entrypoint that creates persistence directories
COPY scripts/backend-entrypoint.sh /usr/local/bin/backend-entrypoint.sh
RUN chmod +x /usr/local/bin/backend-entrypoint.sh

# RoxKV exposes 5 servers:
#   6969 — Native CLI TCP
#   6970 — AI Chat TCP
#   6971 — HTTP / SSE API
#   6972 — AI Chat HTTP
#   6973 — RESP-compatible TCP (redis-cli)
#   8080 — Embedded dashboard
EXPOSE 6969 6970 6971 6972 6973 8080

# Health check: both HTTP servers must respond
HEALTHCHECK --interval=15s --timeout=5s --start-period=20s --retries=5 \
  CMD wget -qO- http://127.0.0.1:6971/healthz >/dev/null && \
      wget -qO- http://127.0.0.1:6972/healthz >/dev/null || exit 1

# The entrypoint creates ~/.roxkv/* directories then exec's the CMD
ENTRYPOINT ["backend-entrypoint.sh"]
CMD ["roxkv", "tcp"]
