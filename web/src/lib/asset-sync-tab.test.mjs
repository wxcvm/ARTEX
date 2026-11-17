import test from "node:test";
import assert from "node:assert/strict";
import { assetSyncTabFromSearch } from "./asset-sync-tab.ts";

test("Purge direct link selects its intelligence tab without changing regular asset sync entry", () => {
  assert.equal(assetSyncTabFromSearch("?tab=purge-intel"), "purge-intel");
  assert.equal(assetSyncTabFromSearch("?source=purge&tab=purge-intel"), "purge-intel");
  assert.equal(assetSyncTabFromSearch(""), "scopesentry");
  assert.equal(assetSyncTabFromSearch("?tab=scopesentry"), "scopesentry");
  assert.equal(assetSyncTabFromSearch("?tab=unknown"), "scopesentry");
});
