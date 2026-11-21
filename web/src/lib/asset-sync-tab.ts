export function assetSyncTabFromSearch(search: string): "scopesentry" | "purge-intel" {
  return new URLSearchParams(search).get("tab") === "purge-intel" ? "purge-intel" : "scopesentry";
}

/** A configured but unreachable MCP source is not ready for synchronization. */
export function scopeSentryReady(
  status: { exists: boolean; configured: boolean; enabled: boolean; reachable: boolean } | null,
): boolean {
  return Boolean(status?.exists && status.configured && status.enabled && status.reachable);
}
