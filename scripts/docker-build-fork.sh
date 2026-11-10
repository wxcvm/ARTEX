#!/usr/bin/env bash
# Build a local ARTEX runtime image from this checkout, not a third-party tag.
set -euo pipefail
cd "$(dirname "$0")/.."

if ! docker buildx version >/dev/null 2>&1; then
  echo "Docker Buildx is required" >&2
  exit 1
fi

arch="${ARTEX_DOCKER_ARCH:-$(docker info --format '{{.Architecture}}')}"
case "$arch" in
  amd64|x86_64) arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) echo "Unsupported architecture: $arch (set ARTEX_DOCKER_ARCH=amd64 or arm64)" >&2; exit 1 ;;
esac

# Use Docker-managed Node/Go builders. Neither system Node nor Go is needed.
mkdir -p "dist/$arch"
docker buildx build --platform "linux/$arch" --file Dockerfile.builder \
  --output "type=local,dest=dist/$arch" .
# The existing runtime Dockerfile ships this just-built binary and its tools.
docker buildx build --platform "linux/$arch" --file Dockerfile \
  --tag artex:local --load .

echo "Built artex:local from the current repository source ($arch)."
echo "Start it without editing .env: ARTEX_IMAGE=artex ARTEX_TAG=local docker compose up -d"
