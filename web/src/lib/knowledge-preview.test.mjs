import assert from "node:assert/strict";
import test from "node:test";
import { attachmentLabels, sourceClaimLabel } from "./knowledge-preview.ts";

test("material comparison has distinct states", () => {
  assert.equal(new Set(Object.values(attachmentLabels)).size, 5);
  assert.match(attachmentLabels.conflict, /差异/);
  assert.match(attachmentLabels.ambiguous, /不唯一/);
  assert.match(attachmentLabels.missing, /缺失/);
});
test("source claims never become current verification", () => {
  for (const claim of [true, false]) assert.match(sourceClaimLabel(claim), /本项目未复核/);
  assert.match(sourceClaimLabel(true), /来源声称/);
});
