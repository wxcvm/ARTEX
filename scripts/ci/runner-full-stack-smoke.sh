#!/usr/bin/env bash
# Full opt-in Compose stack with fake database credentials and a local HTTP target.
set -euo pipefail
cd "$(dirname "$0")/../.."
tmp="$(mktemp -d)"
project="artexstack$$"
target="artex-target-$$"
export POSTGRES_PASSWORD=ci-fixture-only POSTGRES_USER=artex POSTGRES_DB=artex
export ARTEX_IMAGE=artex-runner-ci ARTEX_TAG=local
cp docker-compose.yml docker-compose.runner.yml "$tmp/"
mkdir -p "$tmp/data/tasks" "$tmp/data/sessions" "$tmp/keys" "$tmp/skills" "$tmp/www"
printf 'local-proxy-fixture-ok\n' > "$tmp/www/index.html"
chmod 0777 "$tmp/data" "$tmp/data/tasks" "$tmp/data/sessions" "$tmp/keys"
cat > "$tmp/override.yml" <<'YAML'
services:
  artex:
    entrypoint: ["/app/artex"]
    command: ["-addr", ":8787", "-proxy", ":8788"]
    healthcheck:
      test: ["CMD-SHELL", "curl -fsS http://127.0.0.1:8787/api/health >/dev/null"]
      interval: 2s
      timeout: 3s
      retries: 30
YAML
dc() {
  docker compose -p "$project" -f "$tmp/docker-compose.yml" \
    -f "$tmp/docker-compose.runner.yml" -f "$tmp/override.yml" "$@"
}
cleanup() {
  docker rm -f "$target" >/dev/null 2>&1 || true
  dc down -v --remove-orphans >/dev/null 2>&1 || true
  # The disposable capture CA tree is created by root inside the test
  # container (0700); only the temporary CI fixture is removed with sudo.
  sudo rm -rf -- "$tmp"
}
trap cleanup EXIT
docker image inspect artex-runner-ci:local >/dev/null
# Full stack: PostgreSQL, ARTEX API/MITM proxy, Broker, Executor and relay.
dc up -d --wait --wait-timeout 120
for svc in postgres artex tool-runner tool-executor artex-proxy-relay; do
  cid="$(dc ps -q "$svc")"
  test -n "$cid"
  test "$(docker inspect -f '{{.State.Running}}' "$cid")" = true
done
curl -fsS --max-time 10 http://127.0.0.1:8787/api/health > "$tmp/health.json"
python3 - "$tmp/health.json" <<'PY'
import json,sys
x=json.load(open(sys.argv[1]))
assert x.get("ok") is True and x.get("service") == "artex", x
print("ARTEX API and database startup OK")
PY
dc exec -T artex sh -ec 'test -S /run/artex-runner/runner.sock && test -n "$ARTEX_PG_DSN"'
dc exec -T tool-executor sh -ec '
  test ! -e /run/artex-runner/runner.sock
  test -S /run/artex-executor/executor.sock
  test -z "$(printenv ARTEX_PG_DSN || true)"
  test ! -e /app/keys/jwt.key
'
# Target exists only on controller network. Access from Executor must cross
# its fixed-port proxy relay and the real ARTEX HTTP recording proxy.
docker run -d --name "$target" --network "$project"_runner-control \
  --network-alias artex-ci-target --cap-drop ALL \
  -v "$tmp/www:/srv:ro" python:3.12-alpine \
  python3 -m http.server 8080 --bind 0.0.0.0 --directory /srv >/dev/null
for attempt in $(seq 1 30); do
  if dc exec -T tool-executor sh -ec \
    'curl -fsS --noproxy "" --max-time 3 --proxy http://artex-proxy-relay:8788 http://artex-ci-target:8080/' \
    > "$tmp/proxied.txt" 2>/dev/null; then
    if grep -qx "local-proxy-fixture-ok" "$tmp/proxied.txt"; then break; fi
  fi
  sleep .3
done
grep -qx "local-proxy-fixture-ok" "$tmp/proxied.txt"
dc exec -T tool-executor sh -ec '
  if curl -fsS --noproxy "*" --max-time 2 http://artex-ci-target:8080/ >/dev/null 2>&1; then
    echo "execution container bypassed the proxy/network boundary" >&2
    exit 1
  fi
'
echo "Full Compose startup and real local proxy path OK"
