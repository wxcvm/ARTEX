import assert from "node:assert/strict";
import test from "node:test";
import { statusLabel, statusMeta } from "./status.ts";
test("status presentation follows the locale without changing existing semantics", () => {
  assert.equal(statusLabel("task", "running"), "运行中");
  assert.equal(statusLabel("task", "running", "en"), "Running");
  assert.equal(statusLabel("severity", "critical"), "严重");
  assert.equal(statusLabel("finding", "fixed", "en"), "Fixed");
  assert.equal(statusLabel("task", "custom_status"), "custom_status");
  assert.equal(statusLabel("task", "constructor"), "constructor");
  assert.equal(statusMeta("task", "running").label, "Running");
  assert.equal(statusMeta("task", "running").tone, "blue");
});
