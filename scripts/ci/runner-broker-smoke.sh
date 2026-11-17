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
  -e ARTEX_RUNNER_REMOTE_PEER_ONLY=1 \
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

# A model-controlled command must NOT be able to use its own execution
# socket to create a new privileged RPC session bypassing ARTEX's tool audit.
cat > "$tmp/work/attack.json" <<'JSON'
{"session":"forged-executor","name":"Write","working_dir":"/app/data/tasks","input":{"file_path":"bypass.txt","content":"unauthorized-write"}}
JSON
python3 - <<'PY' | call > "$tmp/local-peer.json"
import json
cmd = (
    'test -S /run/artex-executor/executor.sock && command -v curl >/dev/null && '
    '(curl --noproxy "*" --silent --max-time 2 '
    '--unix-socket /run/artex-executor/executor.sock '
    '-H "Content-Type: application/json" '
    '--data-binary @/app/data/tasks/attack.json '
    'http://localhost/call >/tmp/direct-rpc.log 2>&1 || true) && '
    'test ! -e /app/data/tasks/bypass.txt && printf local-rpc-blocked'
)
print(json.dumps({"session":"split", "name":"Bash",
                  "working_dir":"/app/data/tasks", "input":{"command":cmd}}))
PY

# Exercise the real interactive PTY lifecycle through Broker -> Executor,
# not just the synchronous Bash code path.
printf '%s' '{"session":"split","name":"shell_open","input":{"command":"sh","mode":"hands-free"},"working_dir":"/app/data/tasks"}' | call > "$tmp/pty-open.json"
pty_session="$(python3 - "$tmp/pty-open.json" <<'PY'
import json, sys
response = json.load(open(sys.argv[1]))
assert not response.get("error") and not response["result"].get("IsError"), response
blocks = response["result"].get("Content", [])
text = "".join(b.get("text", b.get("Text", "")) for b in blocks)
sid = json.loads(text)["session_id"]
assert sid and all(c.isalnum() or c in "_-" for c in sid), sid
print(sid)
PY
)"
printf '{"session":"split","name":"shell_send","input":{"session_id":"%s","text":"printf pty-%%s-ok smoke","submit":true,"wait":false},"working_dir":"/app/data/tasks"}' "$pty_session" | call > "$tmp/pty-send.json"
for i in $(seq 1 30); do
  printf '{"session":"split","name":"shell_read","input":{"session_id":"%s","since_cursor":0,"drain":true},"working_dir":"/app/data/tasks"}' "$pty_session" | call > "$tmp/pty-read.json"
  if grep -q "pty-smoke-ok" "$tmp/pty-read.json"; then break; fi
  sleep .1
done
printf '{"session":"split","name":"shell_close","input":{"session_id":"%s"},"working_dir":"/app/data/tasks"}' "$pty_session" | call > "$tmp/pty-close.json"

docker exec "$executor" sh -ec '
  test ! -e /run/artex-runner/runner.sock
  test ! -e /app/keys/jwt.key
  test -z "$(printenv ARTEX_PG_DSN || true)"
'
docker exec "$broker" sh -ec 'test ! -e /app/data/tasks/split.txt'
python3 - "$tmp" <<'PY'
import json, pathlib, sys
p = pathlib.Path(sys.argv[1])
for name in ("write", "read", "bash", "local-peer", "pty-open", "pty-send", "pty-read", "pty-close"):
    x = json.loads((p / (name + ".json")).read_text())
    assert not x.get("error") and not x["result"].get("IsError"), (name, x)
assert "broker-executor-ok" in (p / "read.json").read_text()
assert "socket-not-mounted" in (p / "bash.json").read_text()
assert "local-rpc-blocked" in (p / "local-peer.json").read_text()
assert not (p / "work" / "bypass.txt").exists()
assert "pty-smoke-ok" in (p / "pty-read.json").read_text()
assert "FAKE_CONTROLLER_KEY" not in (p / "read.json").read_text()
print("Broker and Executor communicate across separate networks only via Unix Socket")
PY
