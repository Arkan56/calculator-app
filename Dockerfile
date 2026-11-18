# syntax=docker/dockerfile:1
#
# One Dockerfile, two runnable targets:
#   docker build --target backend  -t calc-backend  .
#   docker build --target frontend -t calc-frontend .
# docker-compose.yml builds and runs both together.
#
# Assumed repo layout (adjust the paths below if yours differs):
#   backend/   -> go.mod, go.sum, cmd/api/main.go, internal/...
#   frontend/  -> package.json, package-lock.json, src/, ...

ARG GO_VERSION=1.26
ARG NODE_VERSION=22

# ---------- Backend: build ----------
FROM golang:${GO_VERSION}-alpine AS backend-build
WORKDIR /src

# Download dependencies first so this layer is cached until go.mod/go.sum change
COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY backend/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

# ---------- Backend: runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot AS backend
COPY --from=backend-build /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]

# ---------- Frontend: build ----------
FROM node:${NODE_VERSION}-alpine AS frontend-build
WORKDIR /app

COPY frontend/package.json frontend/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY frontend/ ./
# Vite inlines env vars at BUILD time. A relative URL means the browser calls
# the same origin and nginx proxies /api to the backend (no CORS needed).
ARG VITE_API_URL=/api/v1
ENV VITE_API_URL=${VITE_API_URL}
RUN npm run build

# ---------- Frontend: runtime ----------
FROM nginx:alpine AS frontend
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=frontend-build /app/dist /usr/share/nginx/html
EXPOSE 80