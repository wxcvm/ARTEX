# syntax=docker/dockerfile:1
# Build this fork's embedded ARTEX binary entirely inside Docker.
# Runtime image construction then reuses the existing Dockerfile/toolchain.
FROM node:22-bookworm AS frontend
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build:static

FROM golang:1.26.3-bookworm AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/web/out /src/server/webui/dist
RUN mkdir -p /out && CGO_ENABLED=0 GOOS=linux go build -tags embedui -trimpath -o /out/artex ./cmd/artex

FROM scratch AS artifact
COPY --from=backend /out/artex /artex
