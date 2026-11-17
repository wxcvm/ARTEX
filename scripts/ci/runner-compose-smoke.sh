#!/usr/bin/env bash
# Test real Compose Broker+Executor service startup without live credentials.
set -euo pipefail
cd "$(dirname "$0")/../.."

tmp="$(mktemp -d)"
project="artexrunner${GITHUB_RUN_ID:-0}${RANDOM}"
export POSTGRES_PASSWORD=fixture-only ARTEX_IMAGE=artex-runner-ci ARTEX_TAG=local
cp docker-compose.yml docker-compose.runner.yml "$tmp/"
mkdir -p "$tmp/data/tasks" "$tmp/data/sessions"
# Disposable CI fixture only; not a production permission recommendation.
chmod 0777 "$tmp/data/tasks" "$tmp/data/sessions"

compose=(docker compose -p "$project" -f "$tmp/docker-compose.yml" -f "$tmp/docker-compose.runner.yml")
cleanup() {
  POSTGRES_PASSWORD=fixture-only ARTEX_IMAGE=artex-runner-ci ARTEX_TAG=local \
    "${compose[@]}" down -v --remove-orphans >/dev/null 2>&1 || true
  rm -rf "$tmp"
}
trap cleanup EXIT

# A previous CI step built this minimal image from the current Go source.
docker image inspect artex-runner-ci:local >/dev/null
# No control service, database or real recording proxy is started.
POSTGRES_PASSWORD=fixture-only ARTEX_IMAGE=artex-runner-ci ARTEX_TAG=local \
  "${compose[@]}" up -d --no-deps --wait tool-executor tool-runner

for svc in tool-runner tool-executor; do
  cid="$("${compose[@]}" ps -q "$svc")"
  test -n "$cid"
  test "$(docker inspect -f '{{.State.Health.Status}}' "$cid")" = healthy
done
test -z "$("${compose[@]}" ps -q artex)"

"${compose[@]}" exec -T tool-executor sh -ec '
  test -S /run/artex-executor/executor.sock
  test ! -e /run/artex-runner/runner.sock
  test ! -e /app/keys/jwt.key
  test -z "$(printenv ARTEX_PG_DSN || true)"
'
"${compose[@]}" exec -T tool-runner sh -ec '
  test -S /run/artex-runner/runner.sock
  test -S /run/artex-executor/executor.sock
  test ! -e /app/data/tasks/compose.txt
'

# Client has a separate PID namespace and no network access. Broker relays
# the request through the private execution Socket to the Executor container.
docker run --rm --network none \
  -v "${project}_runner_socket:/run/artex-runner:ro" \
  --entrypoint sh artex-runner-ci:local -ec '
    curl --fail --silent --show-error --max-time 20 \
      --unix-socket /run/artex-runner/runner.sock \
      -H "Content-Type: application/json" \
      --data "{\"session\":\"compose\",\"name\":\"Write\",\"working_dir\":\"/app/data/tasks\",\"input\":{\"file_path\":\"compose.txt\",\"content\":\"compose-works\"}}" \
      http://localhost/call
  ' > "$tmp/write.json"
docker run --rm --network none \
  -v "${project}_runner_socket:/run/artex-runner:ro" \
  --entrypoint sh artex-runner-ci:local -ec '
    curl --fail --silent --show-error --max-time 20 \
      --unix-socket /run/artex-runner/runner.sock \
      -H "Content-Type: application/json" \
      --data "{\"session\":\"compose\",\"name\":\"Read\",\"working_dir\":\"/app/data/tasks\",\"input\":{\"file_path\":\"compose.txt\"}}" \
      http://localhost/call
  ' > "$tmp/read.json"

python3 - "$tmp" <<'PY'
import json, pathlib, sys
p = pathlib.Path(sys.argv[1])
for name in ("write", "read"):
    x = json.loads((p / (name + ".json")).read_text())
    assert not x.get("error") and not x["result"].get("IsError"), (name, x)
assert "compose-works" in (p / "read.json").read_text()
assert (p / "data/tasks/compose.txt").read_text() == "compose-works"
print("Compose Broker/Executor: healthy startup, isolated sockets and file RPC OK")
PY
