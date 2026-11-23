# Mode B: Purge × ARTEX opt-in chained egress (phase 1)

> **Status:** implementation candidate. Static checks are not a runtime leak test.
> Intended topology: ARTEX isolated tool runner → Mihomo → existing airport proxy → VPS1 (VLESS+REALITY or HY2) → [later: VPS2 egress].
> VPS2 is **not yet provisioned or enforced** by this change. The proxy chain only exists after the operator supplies real node values.

## Why this is a separate Compose file

The default Docker and Windows installs stay unchanged. Use **only** the base file plus the Mode B file, **never** combine this with `docker-compose.runner.yml`. The latter uses a recording relay on another network; combining both would create an uncontrolled alternative exit.

The controller/LLM, PostgreSQL, third-party MCP connectors, browser services outside the executor, and native Windows clients are **not** covered by this phase. The isolation protects only `tool-executor` (Norma Read/Write/Bash/PTY/background built-ins). An external subprocess in another service is outside this boundary.

## Set up on an owned Linux Docker host

1. Build the current fork's `artex:local` image per the main README. Keep `.env` local and confidential.
2. Copy `network/mode-b/config.example.yaml` to `network/mode-b/config.yaml` and set your own airport SOCKS5 address and VPS1 VLESS+REALITY credentials. **Do not commit `config.yaml`**; it is gitignored.
3. `airport-entry` is the **first** hop. On Windows Docker Desktop, `host.docker.internal` usually resolves to the host; on Linux this file adds `host-gateway`. The airport SOCKS listener must be deliberately reachable from the gateway, and should not be exposed to your LAN/public internet. A host-local listener may not be reachable through `host.docker.internal`; test this instead of enabling wide-open LAN binding.
4. Keep `vps1.dialer-proxy: airport-entry`, `MATCH,vps1`, and **no** DIRECT/fallback group. VPS1 `server` should be a literal IP to avoid bootstrap DNS outside the chain. VLESS/REALITY is the included example, not a new server: use your existing VPS1 listener. For HY2, replace only the VPS1 node with a valid `hysteria2` node; the airport SOCKS5 proxy must support UDP forwarding for the QUIC handshake. If it cannot, prefer the existing TCP VLESS path rather than direct-fallback.
5. Check and start:

```sh
test -f network/mode-b/config.yaml || { echo 'missing local Mode B config'; exit 1; }
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml config --quiet
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml up -d
```

The proxy image may be pinned through `ARTEX_MODE_B_MIHOMO_IMAGE`; the default is the upstream stable latest tag; pin a tested image digest for production. A checked-in example does not prove node availability or a successful SOCKS UDP bridge.

## Docker Engine prerequisite and host-gateway bypass

**Docker Engine 28 or newer is required** for `driver_opts:
com.docker.network.bridge.gateway_mode_ipv4: isolated` on the
`runner-execution` bridge. The network also sets isolated IPv6 gateway
mode and leaves IPv6 disabled. The Docker daemon used by Docker Desktop
or WSL is what matters, not just the version of a CLI installed on the
Windows host. If this network mode is unsupported, do **not** remove
these options merely to make Compose start: ordinary `internal: true`
still allows connections to services bound to the host bridge gateway,
which could include a separate system proxy and bypass the required VPS
chain. See the [Docker Engine 28 release notes](https://docs.docker.com/engine/release-notes/28/)
and [Docker's port publishing and gateway modes](https://docs.docker.com/engine/network/port-publishing/).

The CI test now checks both the requested `isolated` driver options
and the *actual host bridge interface* for an assigned IPv4 address.
This check applies to a native Linux Docker Engine; with Docker Desktop,
the bridge may be inside its separate Linux VM. It does **not**
establish that the host DNS resolver, external MCP services, the gateway
itself, or all IPv6 link-local traffic is confined to the chain.

## Isolation properties (and limits)

- `tool-executor` joins only the Docker network `runner-execution`, explicitly marked `internal: true`. No published ports, no external-network membership and no ordinary default route. Its HTTP(S)/SOCKS-aware command-line clients use `mode-b-egress:7890`.
- `mode-b-egress` is the **only** process attached to both the execution and upstream network. It is a mixed proxy, not a TUN: proxy-unaware raw TCP/UDP/ICMP scanners will be blocked, not magically proxied. Docker DNS, host routing/firewall exceptions and packet-level behavior require separate runtime verification; this is not a host-wide kill switch.
- The previous ARTEX traffic-recording relay is intentionally excluded. `ARTEX_RUNNER_PROXY_ADDR` instead rewrites controller-provided localhost:8788 tool-proxy settings to the chained egress gateway. Traffic capture may not work as before. Do not silently add the relay back.
- No VPS admin secrets, SSH private keys, cloud provider credentials, API keys or host Docker socket are mounted into the executor/gateway.
- A stopped gateway, incorrect airport settings, invalid VPS1 credentials or a broken upstream means HTTP requests should fail. Do not add a DIRECT/fallback rule to restore connectivity.

## Acceptance (requires a real host and both your existing nodes)

Use a dedicated, authorized test target and **do not run live tasks while changing the gateway**.

```sh
# 1. Direct connection must fail even with the proxy environment stripped.
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml exec -T tool-executor \
  sh -c 'env -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY -u http_proxy -u https_proxy -u all_proxy curl --noproxy "*" --connect-timeout 5 -fsS https://example.org/'

# 2. Via gateway: must work after supplying real credentials.
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml exec -T tool-executor \
  curl --proxy http://mode-b-egress:7890 --noproxy "" --connect-timeout 10 -fsS https://api.ipify.org

# 3. Stop gateway, then repeat (2): must fail, without direct fallback.
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml stop mode-b-egress
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml up -d mode-b-egress
```

Check IPv4 and IPv6 separately, monitor DNS from the **host** and gateway, and check the exit from a target you own. Record a positive transit result plus negative airport/VPS1 outage tests. Running a command against `api.ipify.org` shows only that command's observed source IP, not a full application privacy proof. Also verify restart, PTY/background tools and every enabled MCP/browser path before claiming production-ready. Keep local audit logs intact.

### Offline failure-path regression in CI

`scripts/ci/mode-b-network-smoke.sh` now creates a disposable gateway
**emulator**, verifies its injected HTTP 502 response is actually received
by an isolated executor, stops that gateway, and verifies that neither
explicit-proxy nor unproxied direct IPv4 requests can succeed after the
outage. It also checks direct IPv4/IPv6 denial with all proxy
environment variables unset. These are **real Docker network** tests
without accessing VPS nodes, provisioning servers or using credentials.

The emulated gateway is **not Mihomo**. This checks only the executor's
network isolation and absence of an automatic HTTP-client fallback, not
Mihomo's runtime route selection, live airport/VPS2 interoperability,
host-level DNS/IPv6 paths, or arbitrary external MCP processes. For
Mihomo runtime failure and actual exit-IP acceptance, use the manual
checks below with tasks stopped.

## Repeatable, operator-triggered exit check

After Mode B is started and all its local nodes are configured, run on the same
Linux Docker host as ARTEX:

```sh
MODE_B_EXPECTED_EXIT_IP=<your-known-VPS1-or-VPS2-public-IPv4> bash scripts/mode-b-verify.sh
```

The script checks that the executor cannot directly connect to an external
IPv4 address, then makes one explicit proxied HTTPS request to compare the
observed public IPv4 with the provided expected exit. It does not change
Docker networks, stop a gateway, log the observed IP or use real credentials in
CI. A successful check **does not establish the airport hop, VPS2 tunnel or
host-level DNS/IPv6 isolation**. It intentionally rejects placeholder and
private expected IPs. Inspect the proxy logs and the host's DNS/IPv6 paths
separately. A gateway outage test must be performed only while no live tasks
are running; restoring it after the test is an explicit operator action.

## Three-hop proxy-chain variant (opt-in, no VPS1 routing changes)

`network/mode-b/config.vps2.example.yaml` adds an explicit second remote
exit to the current airport → VPS1 sequence. The resulting path for a
**proxy-aware** task request is:

```text
isolated executor → private Mihomo → airport SOCKS → existing VPS1 VLESS/REALITY → VPS2 VLESS/REALITY → target
```

Here, the VPS2 node uses `dialer-proxy: vps1`, VPS1 still uses
`dialer-proxy: airport-entry`, and the only routing rule is
`MATCH,vps2`. DNS over HTTPS is tagged `#vps2`. Thus the user has
one explicit exit choice and no automatic `DIRECT` failover.

To opt in, copy `config.vps2.example.yaml` instead of
`config.example.yaml` to the gitignored `config.yaml`. Supply both
servers' actual (distinct) IPs and credentials locally, then restart the
gateway when tasks are idle. Do not commit your node secrets:

```sh
cp network/mode-b/config.vps2.example.yaml network/mode-b/config.yaml
# Edit network/mode-b/config.yaml locally; set real airport, VPS1 and VPS2 settings.
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml config --quiet
docker compose -f docker-compose.yml -f docker-compose.mode-b.yml up -d
MODE_B_EXPECTED_EXIT_IP=<VPS2-public-IPv4> bash scripts/mode-b-verify.sh
```

**This does not provision VPS2** or require changing VPS1's host-wide
default route. You need a real second remote node (an already-owned node or
an eligible free-tier resource; availability not guaranteed). The two VLESS
+ REALITY hops are an illustrative protocol combination, not a verified
compatibility guarantee. When an intermediate node cannot support this
`dialer-proxy` chain, change the intermediate transport after independent
testing rather than permitting a `DIRECT` fallback. Mihomo's upstream
documentation cautions against using UDP-heavy or TLS-camouflage transport
as a relayed *middle* node; see
https://wiki.metacubex.one/en/config/proxies/dialer-proxy/ .

Keep both topologies separate: the original two-hop config remains supported.
A successful public exit-IP check still **does not prove all intermediate
hops or cover non-Runner traffic**. Operators should correlate node-side
connection logs during an authorized test, and stop tasks before testing
failure/outage behavior.

## Phase 2: VPS2

The planned final path is airport → VPS1 → VPS2, preserving VPS1's existing HY2/VLESS listener. To obtain two independent servers at zero extra cost, investigate an eligible free-tier VM; availability and eligibility are **not guaranteed**. An alternative to the proxy-chain variant is to establish a WireGuard tunnel **only on VPS1 for the service's outbound traffic**, with routing and firewall rules that prevent the VLESS/HY2 handler from failing open when WireGuard is down. Do **not** simply replace VPS1's whole default route: that can sever SSH/admin ingress, leak service traffic or create loops. This phase needs a separately reviewed deployment and tests proving the observed exit really is VPS2.

## Related scope / handoff

Purge's existing VPS SSH management route does **not** route ARTEX or model traffic. Purge's existing proxy settings control only explicit plugin operations. A later read-only Purge network status panel should ingest **verified** ARTEX mode/runner/entry/exit status, and show `not_tested` until real observations are available. Never infer a working chain from saved profile fields alone.

## Purge read-only status (phase 1b)

The existing `/api/privacy/bridge-overview` response now includes a fixed-value
`execution_network` aggregate when the service sees the opt-in Runner environment.
Its `mode=B` and `configuration=configured_only` mean **environment variables
were set**, NOT that the Docker network, airport, VPS1 or VPS2 succeeded. The
`verification` and `final_exit` fields are deliberately `not_tested` until
an explicit, independently verified runtime acceptance check exists. Purge PR
wxcvm/dsh-purge#507 reads these fields through its existing scoped credential
and displays them without exposing proxy URLs or node credentials.
