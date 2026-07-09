# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS base
WORKDIR /app
RUN apk add --no-cache ca-certificates git wget
COPY go.mod go.sum ./
RUN go mod download

FROM base AS development
COPY . .
COPY scripts/backend-entrypoint.sh /usr/local/bin/backend-entrypoint.sh
RUN chmod +x /usr/local/bin/backend-entrypoint.sh
EXPOSE 6969 6970 6971 6972
ENTRYPOINT ["backend-entrypoint.sh"]
CMD ["go", "run", "./cmd/roxkv", "roxkv-tcp"]

FROM base AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/roxkv ./cmd/roxkv

FROM alpine:3.22 AS production
WORKDIR /app
RUN apk add --no-cache ca-certificates wget
COPY --from=build /out/roxkv /usr/local/bin/roxkv
COPY scripts/backend-entrypoint.sh /usr/local/bin/backend-entrypoint.sh
RUN chmod +x /usr/local/bin/backend-entrypoint.sh
EXPOSE 6969 6970 6971 6972
HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=5 CMD wget -qO- http://127.0.0.1:6971/healthz >/dev/null && wget -qO- http://127.0.0.1:6972/healthz >/dev/null || exit 1
ENTRYPOINT ["backend-entrypoint.sh"]
CMD ["roxkv", "roxkv-tcp"]
