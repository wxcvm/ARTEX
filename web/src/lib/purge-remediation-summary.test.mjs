import assert from "node:assert/strict";
import test from "node:test";
import { parsePurgeRemediationSummary, purgeSummaryMaxBytes } from "./purge-remediation-summary.ts";

const fixture = {
  schema: "purge.remediation-summary.v1",
  capturedAt: "2026-11-24T05:00:00.000Z",
  counts: { total: 5, pending: 1, reviewing: 2, confirmed: 1, falsePositive: 0, resolved: 1 },
};
const parse = value => parsePurgeRemediationSummary(JSON.stringify(value));

test("accepts the Purge aggregate contract, including zero totals", () => {
  assert.deepEqual(parse(fixture), fixture);
  assert.equal(parse({ ...fixture, counts: { total: 0, pending: 0, reviewing: 0, confirmed: 0, falsePositive: 0, resolved: 0 } }).counts.total, 0);
});

test("rejects all unexpected fields, including raw or identifying details", () => {
  assert.equal(parse({ ...fixture, engagement: "private" }), null);
  assert.equal(parse({ ...fixture, counts: { ...fixture.counts, records: [] } }), null);
  assert.equal(parse({ ...fixture, schema: "unknown" }), null);
  assert.equal(parse({ ...fixture, capturedAt: "yesterday" }), null);
});

test("rejects inconsistent, negative, fractional and oversized counts", () => {
  assert.equal(parse({ ...fixture, counts: { ...fixture.counts, total: 6 } }), null);
  for (const pending of [-1, 1.5, "1", Number.MAX_SAFE_INTEGER + 1]) {
    assert.equal(parse({ ...fixture, counts: { ...fixture.counts, pending } }), null);
  }
  assert.equal(parsePurgeRemediationSummary(" ".repeat(purgeSummaryMaxBytes + 1)), null);
  assert.equal(parsePurgeRemediationSummary("not json"), null);
});
