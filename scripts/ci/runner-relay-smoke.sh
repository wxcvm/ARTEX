#!/usr/bin/env bash
# Verify the actual ARTEX relay binary forwards traffic across the same two
# Docker network roles used by docker-compose.runner.yml. No live targets or
# credentials: a Python HTTP fixture stands in for the recording proxy.
set -euo pipefail
cd "$(dirname "$0")/../.."

tmp="$(mktemp -d)"
control_net="artex-relay-control-$$"
execution_net="artex-relay-execution-$$"
controller="artex-relay-fixture-$$"
relay="artex-relay-service-$$"
cleanup() {
  docker rm -f "$relay" "$controller" >/dev/null 2>&1 || true
  docker network rm "$execution_net" "$control_net" >/dev/null 2>&1 || true
  rm -rf "$tmp"
}
trap cleanup EXIT

# The previous CI step builds this minimal image from the current Go source.
docker image inspect artex-runner-ci:local >/dev/null
mkdir -p "$tmp/www"
printf 'recorded-proxy-relay-ok\n' > "$tmp/www/index.html"

docker network create "$control_net" >/dev/null
docker network create "$execution_net" >/dev/null

# Dummy control-plane listener: reachable only on the control network.
docker run -d --name "$controller" --network "$control_net" \
  --network-alias artex --cap-drop ALL \
  -v "$tmp/www:/srv:ro" python:3.12-alpine \
  python3 -m http.server 8788 --bind 0.0.0.0 --directory /srv >/dev/null

# Real relay binary, on the control network first, then attach only the relay
# (not the controller) to the runner/execution network.
docker run -d --name "$relay" --network "$control_net" \
  --cap-drop ALL --security-opt no-new-privileges --read-only \
  artex-runner-ci:local -runner-relay-listen :8788 \
  -runner-relay-target artex:8788 >/dev/null
docker network connect --alias artex-proxy-relay "$execution_net" "$relay"

# The isolated client must reach the control endpoint through the relay.
docker run --rm --network "$execution_net" --cap-drop ALL \
  --entrypoint sh artex-runner-ci:local -ec '
    for i in $(seq 1 30); do
      if curl --noproxy "*" --silent --show-error --fail --max-time 2 \
        http://artex-proxy-relay:8788/ 2>/dev/null | grep -qx "recorded-proxy-relay-ok"; then
        echo "Relay reached the expected fixed control endpoint"
        exit 0
      fi
      sleep .3
    done
    exit 1
  '

# The exact same execution network must NOT reach the controller directly.
docker run --rm --network "$execution_net" --cap-drop ALL \
  --entrypoint sh artex-runner-ci:local -ec '
    if curl --noproxy "*" --silent --fail --max-time 2 \
      http://artex:8788/ >/dev/null 2>&1; then
      echo "BUG: isolated runner can reach the control host directly" >&2
      exit 1
    fi
    echo "Runner execution network cannot directly resolve/reach controller"
  '
