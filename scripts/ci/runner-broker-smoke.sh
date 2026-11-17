#!/usr/bin/env bash
# Verify a controller-owned Unix Socket broker dispatches SDK calls to a
# separate Docker execution container WITHOUT mounting the socket there.
set -euo pipefail
cd "$(dirname "$0")/../.."
tmp="$(mktemp -d)"
exec_net="artex-tool-exec-$$"
ctrl_net="artex-tool-control-$$"
broker="artex-tool-broker-$$"
executor="artex-tool-executor-$$"
cleanup() {
  docker rm -f "$broker" "$executor" >/dev/null 2>&1 || true
  docker network rm "$exec_net" "$ctrl_net" >/dev/null 2>&1 || true
  rm -rf "$tmp"
}
trap cleanup EXIT
mkdir -p "$tmp/socket" "$tmp/work" "$tmp/private"
# Only these disposable CI fixture mounts are world-writable.
chmod 0777 "$tmp/socket" "$tmp/work"
echo FAKE_CONTROLLER_KEY > "$tmp/private/jwt.key"
docker image inspect artex-runner-ci:local >/dev/null
docker network create "$ctrl_net" >/dev/null
docker network create "$exec_net" >/dev/null

docker run -d --name "$executor" --network "$exec_net" --network-alias tool-executor \
  --read-only --cap-drop ALL --security-opt no-new-privileges \
  --tmpfs /tmp:rw,nosuid,size=64m \
  -v "$tmp/work:/app/data/tasks" artex-runner-ci:local \
  -runner-executor-listen :8789 >/dev/null
for i in $(seq 1 100); do
  if docker exec "$executor" curl -fsS http://127.0.0.1:8789/health >/dev/null 2>&1; then break; fi
  sleep .1
done
docker exec "$executor" curl -fsS http://127.0.0.1:8789/health >/dev/null

docker run -d --name "$broker" --network "$ctrl_net" \
  --read-only --cap-drop ALL --security-opt no-new-privileges \
  -e ARTEX_RUNNER_REMOTE_PEER_ONLY=1 \
  -v "$tmp/socket:/run/artex-runner" artex-runner-ci:local \
  -runner-broker-socket /run/artex-runner/runner.sock \
  -runner-broker-target http://tool-executor:8789 >/dev/null
docker network connect "$exec_net" "$broker"
for i in $(seq 1 100); do
  if test -S "$tmp/socket/runner.sock"; then break; fi
  if ! docker ps --format '{{.Names}}' | grep -qx "$broker"; then
    docker logs "$broker"
    exit 1
  fi
  sleep .1
done
test -S "$tmp/socket/runner.sock"
call() {
  sudo curl --silent --show-error --fail --max-time 30 \
    --unix-socket "$tmp/socket/runner.sock" \
    -H 'Content-Type: application/json' --data-binary @- http://localhost/call
}
printf '%s' '{"session":"split","name":"Write","input":{"file_path":"split.txt","content":"broker-executor-ok"},"working_dir":"/app/data/tasks"}' | call > "$tmp/write.json"
printf '%s' '{"session":"split","name":"Read","input":{"file_path":"split.txt"},"working_dir":"/app/data/tasks"}' | call > "$tmp/read.json"
printf '%s' '{"session":"split","name":"Bash","input":{"command":"test ! -e /run/artex-runner/runner.sock && printf socket-not-mounted"},"working_dir":"/app/data/tasks"}' | call > "$tmp/bash.json"

docker exec "$executor" sh -ec '
  test ! -e /run/artex-runner/runner.sock
  test ! -e /app/keys/jwt.key
  test -z "$(printenv ARTEX_PG_DSN || true)"
'
docker exec "$broker" sh -ec 'test ! -e /app/data/tasks/split.txt'
python3 - "$tmp" <<'PY'
import json, pathlib, sys
p = pathlib.Path(sys.argv[1])
for name in ("write", "read", "bash"):
    x = json.loads((p / (name + ".json")).read_text())
    assert not x.get("error") and not x["result"].get("IsError"), (name, x)
assert "broker-executor-ok" in (p / "read.json").read_text()
assert "socket-not-mounted" in (p / "bash.json").read_text()
assert "FAKE_CONTROLLER_KEY" not in (p / "read.json").read_text()
print("Split Docker broker executes tools remotely; executor cannot see control Socket")
PY
