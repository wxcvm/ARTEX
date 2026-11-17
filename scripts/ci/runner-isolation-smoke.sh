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
runner = services["tool-runner"]
controller = services["artex"]
assert "ARTEX_PG_DSN" in controller.get("environment", {}), "fixture needs controller credentials"
assert not (set(runner.get("environment", {})) & {
    "ARTEX_PG_DSN", "OPENAI_API_KEY", "ANTHROPIC_API_KEY",
    "ARTEX_LLM_SECRET_KEY", "ARTEX_SETUP_TOKEN", "ARTEX_JWT_KEY_DIR"
}), "runner inherited control variables"
assert set(runner["networks"]) == {"runner-execution"}, "runner joined controller network"
destinations = {m["target"] for m in runner.get("volumes", [])}
assert destinations == {"/run/artex-runner", "/app/data/tasks", "/app/data/sessions"}, destinations
assert runner.get("read_only") is True
assert "ALL" in runner.get("cap_drop", [])
print("Compose runner environment, networks and mounts checked")
PY

CGO_ENABLED=0 go build -o "$tmp/artex" ./cmd/artex
cat > "$tmp/Dockerfile" <<'DOCKER'
FROM alpine:3.21
RUN apk add --no-cache bash
COPY artex /app/artex
ENTRYPOINT ["/app/artex"]
DOCKER
docker build -q -t artex-runner-ci:local "$tmp" >/dev/null

# Mount workspace and socket only. Controller secrets stay outside this container.
docker run -d --name "$container" --network none --cap-drop ALL \
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
assert "background-ok" in (p / "taskoutput.json").read_text()
secret = json.loads((p / "secret.json").read_text())
assert secret["result"]["IsError"], secret
assert "CONTROL_SECRET_CANARY_NOT_FOR_RUNNER" not in str(secret)
print("Docker runner: file/Bash tools work; controller secret inaccessible")
PY
