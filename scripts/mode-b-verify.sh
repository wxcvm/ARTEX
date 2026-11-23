#!/usr/bin/env bash
# Read-only Linux Docker Mode B smoke: verify a direct IPv4 request is blocked
# and a proxied IPv4 request matches an operator-supplied expected exit.
# This neither reconfigures the chain nor proves the individual intermediate hops.
set -euo pipefail
cd "$(dirname "$0")/.."
if [[ $# -ne 0 || -z "${MODE_B_EXPECTED_EXIT_IP:-}" ]]; then
  echo 'Usage: MODE_B_EXPECTED_EXIT_IP=<your expected VPS exit IPv4> bash scripts/mode-b-verify.sh' >&2
  exit 2
fi
if [[ ! -f network/mode-b/config.yaml ]]; then
  echo 'Mode B local Mihomo config is missing' >&2
  exit 2
fi
python3 - "$MODE_B_EXPECTED_EXIT_IP" <<'PY'
import ipaddress, sys
address = ipaddress.ip_address(sys.argv[1])
if address.version != 4 or not address.is_global:
    raise SystemExit('Expected exit must be a public IPv4 address')
PY

compose=(docker compose -f docker-compose.yml -f docker-compose.mode-b.yml)
"${compose[@]}" config --quiet
for service in tool-executor tool-runner mode-b-egress; do
  cid="$("${compose[@]}" ps -q "$service")"
  if [[ -z "$cid" || "$(docker inspect -f '{{.State.Running}}' "$cid")" != true ]]; then
    echo "Mode B service is not running: $service" >&2
    exit 1
  fi
done

# A tool-controlled command must not escape by explicitly bypassing proxy env.
if "${compose[@]}" exec -T tool-executor sh -ec '
  env -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY \
      -u http_proxy -u https_proxy -u all_proxy \
      curl --noproxy "*" --silent --connect-timeout 3 --max-time 5 http://1.1.1.1/ >/dev/null 2>&1
'; then
  echo 'FAILED: isolated executor has a direct IPv4 exit' >&2
  exit 1
fi
echo 'PASS: executor direct IPv4 connection failed'

# Use a fixed HTTPS endpoint only for this explicit operator-triggered test.
# Successful HTTP response alone does not prove which intermediate hops were used.
actual="$("${compose[@]}" exec -T tool-executor \
  curl --noproxy '' --proxy http://mode-b-egress:7890 -4 \
  --fail --silent --show-error --connect-timeout 8 --max-time 20 \
  https://api.ipify.org)"
if [[ "$actual" != "$MODE_B_EXPECTED_EXIT_IP" ]]; then
  echo 'FAILED: observed public IPv4 does not match expected exit' >&2
  exit 1
fi
echo 'PASS: proxied HTTPS exit matches expected public IPv4'
echo 'NOT TESTED: airport/VPS1 hop evidence, VPS2 readiness, host DNS/IPv6, MCP, model API, failure outages'
