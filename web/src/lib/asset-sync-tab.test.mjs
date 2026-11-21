import test from "node:test";
import assert from "node:assert/strict";
import { assetSyncTabFromSearch, scopeSentryReady } from "./asset-sync-tab.ts";

test("Purge direct link selects its intelligence tab without changing regular asset sync entry", () => {
  assert.equal(assetSyncTabFromSearch("?tab=purge-intel"), "purge-intel");
  assert.equal(assetSyncTabFromSearch("?source=purge&tab=purge-intel"), "purge-intel");
  assert.equal(assetSyncTabFromSearch(""), "scopesentry");
  assert.equal(assetSyncTabFromSearch("?tab=scopesentry"), "scopesentry");
  assert.equal(assetSyncTabFromSearch("?tab=unknown"), "scopesentry");
});

test("ScopeSentry workbench requires a currently reachable configured source", () => {
  const valid = { exists: true, configured: true, enabled: true, reachable: true };
  assert.equal(scopeSentryReady(valid), true);
  assert.equal(scopeSentryReady(null), false);
  for (const field of ["exists", "configured", "enabled", "reachable"]) {
    assert.equal(scopeSentryReady({ ...valid, [field]: false }), false, field);
  }
});
