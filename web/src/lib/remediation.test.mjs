import assert from "node:assert/strict";
import test from "node:test";
import { mergeRemediationEvents, remediationActions, remediationDraft, remediationMutation } from "./remediation.ts";

const item = { stage: "pending_review", version: 0, evidence_version: 0, read_only: false };
const draft = { ...remediationDraft(), note: "Reviewed manually", reference: "QA-01" };

test("the first write sends version zero with no unused fields", () => {
  assert.deepEqual(remediationMutation(item, "review", draft), {
    ok: true, value: { action: "review", expected_version: 0, note: draft.note },
  });
});

test("an outdated submission can be returned but never accepted", () => {
  const waiting = { ...item, stage: "pending_acceptance", submitted_evidence_version: 0, evidence_version: 1 };
  assert.deepEqual(remediationActions(waiting), ["return"]);
  assert.deepEqual(remediationMutation(waiting, "accept", draft), { ok: false, error: "unavailable" });
  assert.deepEqual(remediationMutation({ ...waiting, submitted_evidence_version: 1 }, "accept", draft).value, {
    action: "accept", note: draft.note, expected_version: 0, expected_evidence_version: 1,
  });
});

test("conflict recovery and read-only state block writes without changing the draft", () => {
  const fixing = { ...item, stage: "pending_fix" };
  assert.deepEqual(remediationMutation(fixing, "submit", draft, false), { ok: false, error: "scope" });
  assert.deepEqual(remediationMutation({ ...fixing, read_only: true }, "submit", draft), { ok: false, error: "unavailable" });
  assert.equal(draft.reference, "QA-01");
});

test("overlapping history pages cannot replace newer copies of an event", () => {
  const newest = { id: 4, note: "read first" };
  assert.deepEqual(mergeRemediationEvents([newest], [{ id: 3 }, { id: 4, note: "duplicate" }]), [newest, { id: 3 }]);
});