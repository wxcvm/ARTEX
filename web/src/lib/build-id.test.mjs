import assert from "node:assert/strict";
import test from "node:test";

import nextConfig from "../../next.config.mjs";

test("Next build ID is reproducible and can be pinned per release", async () => {
  const original = process.env.ARTEX_BUILD_ID;
  try {
    delete process.env.ARTEX_BUILD_ID;
    assert.equal(await nextConfig.generateBuildId(), "artex");
    assert.equal(await nextConfig.generateBuildId(), "artex");
    process.env.ARTEX_BUILD_ID = "commit-123456";
    assert.equal(await nextConfig.generateBuildId(), "commit-123456");
  } finally {
    if (original === undefined) delete process.env.ARTEX_BUILD_ID;
    else process.env.ARTEX_BUILD_ID = original;
  }
});
