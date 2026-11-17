#!/usr/bin/env bash
# Smoke test the actual Linux Docker runner boundary with only fake credentials.
set -euo pipefail
cd "$(dirname "$0")/../.."
tmp="$(mktemp -d)"
container="artex-runner-ci-$$"
cleanup() {
  docker rm -f "$container" >/dev/null 2>&1 || true
  rm -rf "$tmp"
}
trap cleanup EXIT
mkdir -p "$tmp/socket" "$tmp/work" "$tmp/controller"
# GitHub hosted Docker can remap container root; allow the isolated test
# container to create the socket. Production uses a private named volume.
chmod 0777 "$tmp/socket"
chmod 0777 "$tmp/work" # fixture-only writable workspace for remapped container UID
printf 'CONTROL_SECRET_CANARY_NOT_FOR_RUNNER\n' > "$tmp/controller/jwt.key"
chmod 600 "$tmp/controller/jwt.key"

# Verify the opt-in Compose overlay cannot accidentally place control secrets,
# database network or the whole data directory in the tool runner.
POSTGRES_PASSWORD=ci-placeholder docker compose \
  -f docker-compose.yml -f docker-compose.runner.yml config --format json > "$tmp/compose.json"
python3 - "$tmp/compose.json" <<'PY'
import json, sys
services = json.load(open(sys.argv[1]))["services"]
broker, executor, controller = (services[n] for n in ("tool-runner", "tool-executor", "artex"))
assert "ARTEX_PG_DSN" in controller.get("environment", {}), "fixture needs control credentials"
blocked = {"ARTEX_PG_DSN", "OPENAI_API_KEY", "ANTHROPIC_API_KEY",
           "ARTEX_LLM_SECRET_KEY", "ARTEX_SETUP_TOKEN", "ARTEX_JWT_KEY_DIR"}
assert not (set(broker.get("environment", {})) & blocked), "broker inherited control vars"
assert not (set(executor.get("environment", {})) & blocked), "executor inherited control vars"
assert set(executor["networks"]) == {"runner-execution"}, "executor joined control network"
assert set(broker["networks"]) == {"runner-control"}, "broker joined execution network"
broker_mounts = {m["target"] for m in broker.get("volumes", [])}
executor_mounts = {m["target"] for m in executor.get("volumes", [])}
assert broker_mounts == {"/run/artex-runner", "/run/artex-executor"}, broker_mounts
assert executor_mounts == {"/run/artex-executor", "/app/data/tasks", "/app/data/sessions"}, executor_mounts
assert "/run/artex-runner" not in executor_mounts
for service in (broker, executor):
    assert service.get("read_only") is True
    assert "ALL" in service.get("cap_drop", [])
assert broker.get("environment", {}).get("ARTEX_RUNNER_REMOTE_PEER_ONLY") == "1"
assert not broker.get("pid") and not controller.get("pid"), "PID namespace sharing bypasses SO_PEERCRED"
print("Compose controller/worker sockets, networks and mounts checked")
PY

CGO_ENABLED=0 go build -o "$tmp/artex" ./cmd/artex
cat > "$tmp/Dockerfile" <<'DOCKER'
FROM alpine:3.21
RUN apk add --no-cache bash curl
COPY artex /app/artex
ENTRYPOINT ["/app/artex"]
DOCKER
docker build -q -t artex-runner-ci:local "$tmp" >/dev/null

# Mount workspace and socket only. Controller secrets stay outside this container.
docker run -d --name "$container" --network none --cap-drop ALL \
  -e ARTEX_RUNNER_REMOTE_PEER_ONLY=1 \
  --security-opt no-new-privileges --read-only \
  --tmpfs /tmp:rw,nosuid,size=64m \
  -v "$tmp/socket:/run/artex-runner" \
  -v "$tmp/work:/app/data/tasks" \
  artex-runner-ci:local -runner-socket /run/artex-runner/runner.sock >/dev/null
for i in $(seq 1 100); do
  if test -S "$tmp/socket/runner.sock"; then break; fi
  if ! docker ps --format '{{.Names}}' | grep -qx "$container"; then
    docker logs "$container"
    exit 1
  fi
  sleep .1
done
test -S "$tmp/socket/runner.sock"

call() {
  # Deliberately root-owned 0600 socket.
  sudo curl --silent --show-error --fail --max-time 30 \
    --unix-socket "$tmp/socket/runner.sock" \
    -H 'Content-Type: application/json' \
    --data-binary @- "http://localhost/call"
}
printf '%s' '{"session":"ci","name":"Write","input":{"file_path":"fixture.txt","content":"hello-runner"},"working_dir":"/app/data/tasks"}' | call > "$tmp/write.json"
printf '%s' '{"session":"ci","name":"Read","input":{"file_path":"fixture.txt"},"working_dir":"/app/data/tasks"}' | call > "$tmp/read.json"
printf '%s' '{"session":"ci","name":"Bash","input":{"command":"test -z \"$(printenv ARTEX_PG_DSN || true)\" && test ! -e /app/controller/jwt.key && printf isolated-runner"},"working_dir":"/app/data/tasks"}' | call > "$tmp/bash.json"
printf '%s' '{"session":"ci","name":"Read","input":{"file_path":"/app/controller/jwt.key"},"working_dir":"/app/data/tasks"}' | call > "$tmp/secret.json"

# Model-controlled Bash can open the socket file, but must not be able to use
# it as an alternative RPC endpoint to bypass the controller's tool policy.
# The POST is fully valid; if it is accepted the canary file will appear.
python3 - <<'PY' | call > "$tmp/local-peer.json"
import json
attack = json.dumps({
    "session": "ci-bypass", "name": "Write",
    "working_dir": "/app/data/tasks",
    "input": {"file_path": "bypass.txt", "content": "unauthorized"}
})
cmd = (
    "curl --max-time 2 --silent --unix-socket /run/artex-runner/runner.sock "
    "-H 'Content-Type: application/json' --data "
    + __import__("shlex").quote(attack)
    + " http://localhost/call >/tmp/local-peer-response 2>&1 || true; "
      "test ! -e /app/data/tasks/bypass.txt && printf local-peer-blocked"
)
print(json.dumps({"session": "ci", "name": "Bash", "working_dir": "/app/data/tasks",
                  "input": {"command": cmd}}))
PY

# The same per-session manager must support background launch and retrieval.
printf '%s' '{"session":"ci","name":"Bash","input":{"command":"printf background-ok; sleep 1","run_in_background":true},"working_dir":"/app/data/tasks"}' | call > "$tmp/background.json"
task_id="$(python3 - "$tmp/background.json" <<'PY'
import json, re, sys
data = json.load(open(sys.argv[1]))
text = "".join(b.get("text", "") for b in data["result"]["Content"])
match = re.search(r"task_[0-9]+", text)
assert match, text
print(match.group(0))
PY
)"
printf '{"session":"ci","name":"TaskOutput","input":{"task_id":"%s","block":true,"timeout_ms":10000},"working_dir":"/app/data/tasks"}' "$task_id" | call > "$tmp/taskoutput.json"

python3 - "$tmp" <<'PY'
import json, pathlib, sys
p = pathlib.Path(sys.argv[1])
for name in ("write", "read", "bash"):
    x = json.loads((p / (name + ".json")).read_text())
    assert not x.get("error") and not x["result"].get("IsError"), (name, x)
assert "hello-runner" in (p / "read.json").read_text()
assert "isolated-runner" in (p / "bash.json").read_text()
assert "local-peer-blocked" in (p / "local-peer.json").read_text()
assert not (p / "work" / "bypass.txt").exists()
assert "background-ok" in (p / "taskoutput.json").read_text()
secret = json.loads((p / "secret.json").read_text())
assert secret["result"]["IsError"], secret
assert "CONTROL_SECRET_CANARY_NOT_FOR_RUNNER" not in str(secret)
print("Docker runner: file/Bash tools work; controller secret inaccessible")
PY
