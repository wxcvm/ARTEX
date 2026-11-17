export function assetSyncTabFromSearch(search: string): "scopesentry" | "purge-intel" {
  return new URLSearchParams(search).get("tab") === "purge-intel" ? "purge-intel" : "scopesentry";
}
