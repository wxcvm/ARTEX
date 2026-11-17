# Browser capture sessions and evidence lifecycle

ARTEX already ships a `playwright-cli` skill with **named browser sessions**. Use the upstream CLI rather than creating a second profile manager:

```bash
playwright-cli -s=review-a open https://example.test --persistent
playwright-cli -s=review-b open https://example.test
playwright-cli list
playwright-cli -s=review-b close
playwright-cli -s=review-a close
playwright-cli -s=review-a delete-data
```

A Playwright session name and browser profile **do not yet automatically attach to ARTEX traffic evidence**. Session isolation belongs to the CLI: the ARTEX Chrome extension captures the currently selected tab through `chrome.debugger`, independently of Playwright's session manager.

## Manual tab capture

The extension's Start button attaches to the active tab and enables the DevTools Network domain. Starting while a different tab is already being captured reports the active tab instead of pretending that the requested tab changed. If Network.enable fails after a successful attach, the extension detaches immediately.

Stop freezes new network events, waits for in-flight response-body reads, disables/detaches the debugger and flushes **all** queued batches. If a server upload fails, pending exchanges remain **in memory** for a later Flush now or Stop; a subsequent Start does not erase them. The popup shows the queue count and last upload error. An unexpected debugger detach stops capture, retains unfinished requests as attempted evidence, uploads queued exchanges where possible and reports the reason. Stop and detach retain requests that have no terminal Network event as status **0** unless an HTTP response status was actually observed; they do not invent a response body or claim the request completed.

Do not treat a successful browser event or a queued exchange as persisted evidence until the ARTEX server confirms ingestion. The extension deliberately does not write request/response bodies into browser persistent storage; queued exchanges are lost if the service worker or browser process restarts. The evidence snapshot flow only covers records already ingested by ARTEX.

## Scope

This stage fixes capture **lifecycle reliability** and provides regression coverage through the existing web `npm test` command. It does not claim Playwright-to-ARTEX session provenance, durable offline queuing, automatic browser restart recovery, or independently managed browser profiles in the ARTEX UI. These require a separate evidence metadata contract and explicit lifecycle integration.

## Capture group provenance

Each successful Start in the Chrome extension creates a new random capture group UUID. Every request is tagged when the Network request event arrives, so in-flight body reads and retries keep the original group even if the next capture starts. The extension displays the latest group in its popup; the ingest endpoint persists it as optional **client-provided** metadata, independent of the request URL or headers.

The traffic index stores capture IDs in an optional auxiliary row alongside HAR page attribution. Traffic archive export/import preserves the capture group; normal traffic deletion removes it. When an ingested exchange is linked to a finding, the immutable snapshot copies the group and includes it in its identity hash. Old snapshots without a group retain their previous IDs. The finding card shows a short group prefix, and the detail viewer shows the full value.

This is **capture-batch correlation**, not the browser's authenticated identity, user identity, cookie state or Playwright's `-s` session name. The API accepts metadata from authorized clients, so the value is not independently attested. Captured bodies remain in memory until successful ingestion; no new persistent browser storage or cross-browser session control is introduced.

## Find all requests from a capture group

Open **Traffic → Advanced filters** and paste the exact **Capture group ID** shown in the extension popup or the immutable finding evidence detail. This filters the traffic list by the optional capture group while preserving the existing host, method, response-status and content filters. The same filter is available as `GET /api/traffic?capture_id=<id>`. An empty filter continues to show legacy traffic records without a capture ID.

This is a lookup of recorded traffic, not browser control or an authenticated assertion about which Playwright session produced the traffic. If the original traffic has been deleted, the immutable finding evidence may still show the group, but the traffic-list filter will correctly return no hot rows.

The finding's evidence detail also offers **View group traffic** to open the pre-filled filter directly. This may return no rows if the original hot traffic was deleted; the immutable finding evidence remains available.


## Locate capture groups in traffic

`Page`, `Search`, and agent-side traffic queries return an optional `capture_id` per exchange. The **Capture group** column shows a shortened ID; hover for the full value or click it to apply the existing exact-group filter while preserving other active filters. Unattributed legacy rows display a dash rather than an inferred identity. Groups are extension-supplied correlation hints, not a verified Playwright browser session or user identity.

## Browse recorded capture groups

The **Traffic → Capture groups** picker fetches records only when opened. `GET /api/traffic/groups?page=0&size=20` returns groups ordered by latest recorded exchange, including request count, distinct hosts, first/last timestamp, HTTP 4xx/5xx count and status-0 (unobserved) count. The `has_more` field supports paging without a second global distinct-count query. Selecting a group applies the existing exact filter without clearing other filters.

Summaries include only persisted traffic. A group disappears when its last recorded exchange is deleted; immutable finding evidence can still preserve its old group ID. Status 0 is never counted as success. Client-supplied group IDs are not authenticated browser or Playwright identities.


## Unobserved browser failures and retry receipts

The extension also records `Network.loadingFailed` requests and completed requests
that lack a `Network.responseReceived` event. They use HTTP status **0**
(unobserved), not an invented 2xx response; any actually observed status is retained.
This preserves evidence that a request was attempted without claiming it succeeded.

Each extension exchange now carries its debugger request ID under `exchangeId`,
paired with the existing random `captureId`. ARTEX derives a stable, namespaced
storage key from that pair. If an HTTP upload succeeded but its response was lost,
replaying the same batch acknowledges an already stored record rather than replacing
the original response body or inserting a duplicate. `POST /api/traffic/ingest`
responds with `ingested`, `duplicates`, and `failed` counts. Older clients
without `exchangeId` retain their prior write behavior. The popup counts new and
already-ingested records as confirmed, not as separate requests.

Failed network uploads keep unsent entries in the session's in-memory queue instead
of silently discarding entries above a fixed restore slice. This is not persistent
offline storage: Chrome terminating the worker can still lose unsent data. This
does not establish a Playwright authenticated session identity or guarantee that
the debugger saw every attempted network call.
