#!/usr/bin/env bash
# Verify real Docker isolation without contacting user nodes or relying on proxy credentials.
# This is NOT an end-to-end Mode B transit / VPS exit test.
set -euo pipefail
cd "$(dirname "$0")/../.."

docker image inspect artex-runner-ci:local >/dev/null
tmp="$(mktemp -d)"
project="artexmodeb${GITHUB_RUN_ID:-0}${RANDOM}"
export POSTGRES_PASSWORD=fixture-only ARTEX_IMAGE=artex-runner-ci ARTEX_TAG=local
cleanup() {
  docker rm -f "${mock_gateway:-}" >/dev/null 2>&1 || true
  docker compose -p "$project" -f "$tmp/docker-compose.yml" -f "$tmp/docker-compose.mode-b.yml" \
    down --remove-orphans -v >/dev/null 2>&1 || true
  rm -rf "$tmp"
}
trap cleanup EXIT
cp docker-compose.yml docker-compose.mode-b.yml "$tmp/"
mkdir -p "$tmp/network/mode-b" "$tmp/data/tasks" "$tmp/data/sessions"
# A placeholder file satisfies Compose's local bind path. No gateway is started.
printf 'mode: rule\nrules:\n  - MATCH,REJECT\n' > "$tmp/network/mode-b/config.yaml"
chmod 0777 "$tmp/data/tasks" "$tmp/data/sessions"

compose=(docker compose -p "$project" -f "$tmp/docker-compose.yml" -f "$tmp/docker-compose.mode-b.yml")
"${compose[@]}" config --format json > "$tmp/config.json"
python3 - "$tmp/config.json" <<'PY'
import json, sys
cfg = json.load(open(sys.argv[1]))
services = cfg["services"]
executor, gateway, broker, controller = (services[x] for x in ("tool-executor", "mode-b-egress", "tool-runner", "artex"))
assert cfg["networks"]["runner-execution"]["internal"] is True
assert set(executor["networks"]) == {"runner-execution"}
assert set(gateway["networks"]) == {"runner-execution", "mode-b-uplink"}
assert set(broker["networks"]) == {"runner-control"}
assert set(controller["networks"]) == {"runner-control"}
assert set(services["postgres"]["networks"]) == {"runner-control"}
assert not gateway.get("ports"), "gateway must not publish a proxy port"
assert not executor.get("ports")
assert executor.get("read_only") and gateway.get("read_only")
assert "ALL" in executor.get("cap_drop", [])
assert "NET_RAW" not in executor.get("cap_add", [])
assert "ARTEX_PG_DSN" not in executor.get("environment", {})
assert "ARTEX_PG_DSN" in controller.get("environment", {})
print("Mode B merged Compose: control-plane and execution networks are separated")
PY

# Start only one disposable execution container, *without* the gateway/airport/VPS.
# Check that even stripping all proxy settings cannot restore a direct IPv4/IPv6 route.
"${compose[@]}" run --rm --no-deps --entrypoint sh tool-executor -ec '
  test ! -e /run/artex-runner/runner.sock
  test -z "$(printenv ARTEX_PG_DSN || true)"
  if env -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY \
         -u http_proxy -u https_proxy -u all_proxy \
         curl --noproxy "*" --silent --connect-timeout 3 --max-time 5 \
         http://1.1.1.1/ >/dev/null 2>&1; then
    echo "Mode B ERROR: direct IPv4 escaped the isolated executor" >&2
    exit 1
  fi
  if env -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY \
         -u http_proxy -u https_proxy -u all_proxy \
         curl -6 --noproxy "*" --silent --connect-timeout 3 --max-time 5 \
         "https://[2606:4700:4700::1111]/" >/dev/null 2>&1; then
    echo "Mode B ERROR: direct IPv6 escaped the isolated executor" >&2
    exit 1
  fi
  echo "Mode B executor: direct IPv4/IPv6 blocked with proxy env removed"
'

# Simulate a gateway that is reachable but whose upstream fails. Its reply is
# synthetic HTTP 502; this test never contacts an external site or real nodes.
# The existing CI runner image is Alpine, with busybox nc available.
mock_gateway="${project}-mock-gateway"
network="${project}_runner-execution"
docker network inspect "$network" >/dev/null
docker run --detach --rm --name "$mock_gateway" \
  --network "$network" --network-alias mode-b-egress \
  --read-only --cap-drop ALL --security-opt no-new-privileges \
  --entrypoint sh artex-runner-ci:local -ec '
  while :; do
    printf "HTTP/1.1 502 Bad Gateway\r\nConnection: close\r\nContent-Length: 0\r\n\r\n" \
      | busybox nc -l -p 7890
  done
' >/dev/null

# The executor is allowed to talk to the internal gateway, so the observed
# synthetic 502 proves it reached the proxy but did not bypass it.
reply="$("${compose[@]}" run --rm --no-deps --entrypoint sh tool-executor -ec '
  curl --proxy http://mode-b-egress:7890 --noproxy "" \
    --silent --show-error --connect-timeout 3 --max-time 5 \
    --include http://probe.invalid/
')"
case "$reply" in
  *"HTTP/1.1 502 Bad Gateway"*) ;;
  *) echo "Mode B ERROR: expected synthetic upstream failure from gateway" >&2; exit 1 ;;
esac

# Simulate gateway outage. An explicit proxy request MUST NOT turn into a
# successful direct request, and removing the proxy must not allow egress.
docker stop "$mock_gateway" >/dev/null
if "${compose[@]}" run --rm --no-deps --entrypoint sh tool-executor -ec '
  curl --proxy http://mode-b-egress:7890 --noproxy "" \
    --fail --silent --connect-timeout 3 --max-time 5 \
    http://probe.invalid/ >/dev/null 2>&1
'; then
  echo "Mode B ERROR: request succeeded after gateway was stopped" >&2
  exit 1
fi
if "${compose[@]}" run --rm --no-deps --entrypoint sh tool-executor -ec '
  env -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY \
      -u http_proxy -u https_proxy -u all_proxy \
      curl --noproxy "*" --silent --connect-timeout 3 --max-time 5 \
      http://1.1.1.1/ >/dev/null 2>&1
'; then
  echo "Mode B ERROR: direct egress succeeded after gateway outage" >&2
  exit 1
fi
echo "Mode B: internal gateway HTTP 502 and stopped-gateway outage both fail closed"
