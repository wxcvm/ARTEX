# MCP discovery consistency

Remote HTTP and legacy SSE discovery follow every tools/list cursor. A later page failure, repeated cursor, missing tools array or ambiguous tool name is an error, not a successful partial catalog. Explicit discovery replaces the cached catalog and records success in one database transaction; on failure the previous catalog is retained and the attempt is reported as failed. A successful empty list clears the cache. The MCP page labels cached tools and refreshes health after failed refreshes; disabling a server is shown separately from connection status. Cached results are not proof of current availability. Stdio continues using the existing SDK; active tool execution and visibility rules are unchanged.

Upstream design reference: CyberStrike `packages/cyberstrike/src/mcp/index.ts` connection status and tool-change lifecycle. Independent implementation; no upstream source copied.
