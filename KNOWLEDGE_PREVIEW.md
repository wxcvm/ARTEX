# Knowledge material preview

This is a read-only import preview, not a replacement for PostgreSQL or a new Agent runtime. It does not change model keys, update channels, task scopes, lifecycle controls or task-experience tables.

## Implementation

- `scripts/knowledge/knowledge_export.py`: offline stdlib SQLite/tar exporter. Opens the database read-only/immutable; rejects pending WAL/journal snapshots; queries only whitelisted fields from an ordinary `poc` table. No archive extraction, material execution, historical task-table reads or network access.
- `knowledge/`: bounded catalog reader, source/identity/body-hash validation, deterministic search and pagination.
- Authenticated GET routes: `/api/knowledge/status`, `/api/knowledge?q=&conflicts=false&offset=0&limit=20`, `/api/knowledge/entries/{id}`. Existing JWT middleware applies; no write endpoints.
- System -> Knowledge (`/system/knowledge/`): search, pagination, source claims and both material bodies. Failed/disabled catalogs do not appear as empty libraries. Source `verified` never becomes current project/task verification.

## Generate and enable

Requires a closed SQLite snapshot without pending WAL/journal data. Output must be a new file. Keep generated catalogs outside Git and public static asset directories: material bodies can themselves contain sensitive historical text even though historical metadata columns are excluded.

```powershell
python scripts/knowledge/knowledge_export.py --database "[PATH]/knowledge.db" --archive "[PATH]/source.tar.gz" --output "[PATH]/knowledge-preview.json"
$env:ARTEX_KNOWLEDGE_FILE = "[PATH]/knowledge-preview.json"
# Start the normal ARTEX build with its existing database configuration.
```

- Limits: 10,000 entries, 32 MiB catalog, 2 MiB per material; default page size 20, maximum 100; query <=512 UTF-8 bytes.
- Identity derives from source DB SHA-256 + code. Body hashes preserve exact UTF-8 bytes; no newline/BOM normalization.
- Attachment states: `none`, `match`, `conflict`, `missing`, `ambiguous`. Duplicate candidates are never resolved arbitrarily. Conflicts retain both bodies.
- No automatic merge, verified promotion, tool registration, Agent ingestion or writeback.
- Invalid/unavailable configured catalog: 503; missing entry: 404; invalid query: 400. No configuration: status enabled=false and list/detail 503.
- Source URLs and material bodies are plain text, not fetched or rendered as active HTML.

## Local delivery - 2026-11-06

`D:\Codex\ARTEX-preview\knowledge-20261106` contains the generated catalog, an independently usable offline HTML viewer, and a Windows ARTEX binary with the integrated page embedded. The offline viewer is a local review artifact, not a production execution flow; opening it needs no database/model/server. Its generator and browser test are under the untracked validation evidence directory.

Actual input: 448 materials; 382 DB bodies; 227 exact-byte attachment conflicts; 0 missing/ambiguous attachments. Generated data is not committed.

### Validation boundaries

Exact results and hashes: `.git/artex-validation/knowledge-delivery-results.json`.

- Exporter fixtures and Go reader/API regressions are separate from model benchmarks.
- Browser checks cover the real offline dataset and the production-built Next.js page with intercepted API responses. This is not full PostgreSQL-backed browser acceptance.
- Windows PostgreSQL initialization returned permission errors in separate preview directories. The existing deployment was not modified or started; no production DB migration occurred.
- The normal Turbopack invocation made no useful progress and was stopped; the supported webpack production build completed. No dependency versions or lockfile changes.
- The earlier temporary headless attempt returned ACCOUNT_SIGN_IN_REQUIRED. On 2026-11-06 the user-provided codex-review entry completed a real read-only guide review (three specified text files, normal model response and exit 0). This is assistant review evidence, not an ARTEX live-model assessment. Desktop settings and credentials were not modified.
- Linux CI and live-model task acceptance have not been run for this slice.

## Rollback

1. Remove ARTEX_KNOWLEDGE_FILE from the process environment: the preview reports disabled; no schema/data rollback is needed.
2. For a binary rollout, retain and restart the prior binary/config pair. The local delivery did not overwrite D:\deepseek\ARTEX-deploy.
3. Source rollback is limited to knowledge/, scripts/knowledge/, server/knowledge*, the knowledge page/types/test, three knowledge routes and one navigation entry. Preserve pending OpenAPI work and the three pre-existing missing PoCs.
4. Offline HTML/JSON and temporary test outputs can be removed independently; they are not required infrastructure.

## Beginner guide

ARTEX_小白使用指南.md provides a Chinese, click-by-click guide for the offline reader, existing deployment, conditional integrated setup, and the user-provided Harness review entry. Its offline HTML edition and separate guide acceptance record are delivered locally alongside knowledge-preview.html.
