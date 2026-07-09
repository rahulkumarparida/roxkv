# syntax=docker/dockerfile:1.7

FROM node:22-alpine AS base
WORKDIR /app
COPY WebDashBoard/roxkv-dashboard/package.json WebDashBoard/roxkv-dashboard/package-lock.json ./
RUN npm i

FROM base AS development
COPY WebDashBoard/roxkv-dashboard/ ./
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]

FROM base AS build
COPY WebDashBoard/roxkv-dashboard/ ./
RUN npm run build

FROM node:22-alpine AS production
WORKDIR /app
RUN npm install --global serve
COPY --from=build /app/dist ./dist
EXPOSE 5173
HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=5 CMD wget -q --spider http://127.0.0.1:5173/ || exit 1
CMD ["serve", "-s", "dist", "-l", "5173"]
