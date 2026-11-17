#!/usr/bin/env bash
# Real two-container broker/executor check. The processes are in different
# Docker networks and communicate through an execution-only Unix socket.
set -euo pipefail
cd "$(dirname "$0")/../.."
tmp="$(mktemp -d)"
control_net="artex-broker-control-$$"
execution_net="artex-broker-execution-$$"
broker="artex-broker-$$"
executor="artex-executor-$$"
cleanup() {
  docker rm -f "$broker" "$executor" >/dev/null 2>&1 || true
  docker network rm "$execution_net" "$control_net" >/dev/null 2>&1 || true
  rm -rf "$tmp"
}
trap cleanup EXIT
mkdir -p "$tmp/controller-socket" "$tmp/executor-socket" "$tmp/work" "$tmp/private"
# Disposable CI fixtures only; production Compose volumes keep native ownership.
chmod 0777 "$tmp/controller-socket" "$tmp/executor-socket" "$tmp/work"
echo FAKE_CONTROLLER_KEY > "$tmp/private/jwt.key"
docker image inspect artex-runner-ci:local >/dev/null
docker network create "$control_net" >/dev/null
docker network create "$execution_net" >/dev/null

docker run -d --name "$executor" --network "$execution_net" \
  --read-only --cap-drop ALL --security-opt no-new-privileges \
  --tmpfs /tmp:rw,nosuid,size=64m \
  -v "$tmp/executor-socket:/run/artex-executor" \
  -v "$tmp/work:/app/data/tasks" artex-runner-ci:local \
  -runner-socket /run/artex-executor/executor.sock >/dev/null

docker run -d --name "$broker" --network "$control_net" \
  --read-only --cap-drop ALL --security-opt no-new-privileges \
  -e ARTEX_RUNNER_REMOTE_PEER_ONLY=1 \
  -v "$tmp/controller-socket:/run/artex-runner" \
  -v "$tmp/executor-socket:/run/artex-executor" artex-runner-ci:local \
  -runner-broker-socket /run/artex-runner/runner.sock \
  -runner-broker-target /run/artex-executor/executor.sock >/dev/null
for i in $(seq 1 100); do
  if test -S "$tmp/controller-socket/runner.sock" && test -S "$tmp/executor-socket/executor.sock"; then break; fi
  if ! docker ps --format '{{.Names}}' | grep -qx "$broker"; then
    docker logs "$broker"
    exit 1
  fi
  sleep .1
done
test -S "$tmp/controller-socket/runner.sock"
test -S "$tmp/executor-socket/executor.sock"

call() {
  sudo curl -fsS --max-time 30 \
    --unix-socket "$tmp/controller-socket/runner.sock" \
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
print("Broker and Executor communicate across separate networks only via Unix Socket")
PY
