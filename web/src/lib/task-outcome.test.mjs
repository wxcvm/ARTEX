import assert from "node:assert/strict";
import test from "node:test";
import { isTaskOutcome, taskOutcomePresentation } from "./task-outcome.ts";
import { taskOutcomeScenario, taskOutcomeScenarioReport } from "./mock/task-outcome-scenarios.ts";

const states = {
  9101: "incomplete", 9102: "incomplete", 9103: "incomplete", 9104: "incomplete",
  9105: "completed_no_findings", 9106: "completed_with_findings", 9107: "task_failed",
  9108: "environment_error", 9109: "completed_no_findings", 9110: "incomplete", 9111: "incomplete",
};

for (const [id, expected] of Object.entries(states)) {
  test(`browser fixture ${id}: ${expected}`, () => {
    const fixture = taskOutcomeScenario(id);
    assert.ok(isTaskOutcome(fixture.outcome));
    const actual = taskOutcomePresentation(fixture.task, fixture.outcome, fixture.outcomeError);
    assert.equal(actual.state, expected);
  });
}

test("all non-done lifecycle states keep zero findings incomplete", () => {
  const { outcome } = taskOutcomeScenario("9105");
  for (const status of ["created", "queued", "running", "paused", "timeout", "unknown"]) {
    assert.equal(taskOutcomePresentation({ status }, outcome).state, "incomplete", status);
  }
});

test("paused and queued flags override stale done and running states", () => {
  const { outcome } = taskOutcomeScenario("9105");
  for (const status of ["done", "running"]) {
    assert.equal(taskOutcomePresentation({ status, paused: true }, outcome).state, "incomplete");
    assert.equal(taskOutcomePresentation({ status, queued: true }, outcome).state, "incomplete");
  }
});

test("missing and invalid responses are never treated as zero findings", () => {
  const { task, outcome } = taskOutcomeScenario("9105");
  assert.equal(taskOutcomePresentation(task, null).state, "loading");
  for (const invalid of [null, {}, { ...outcome, findings: {} },
    { ...outcome, findings: { ...outcome.findings, total: -1 } },
    { ...outcome, findings: { ...outcome.findings, confirmed: 1 } },
    { ...outcome, findings: { ...outcome.findings, high: 1 } },
    { ...outcome, skills_used: null },
    { ...outcome, skills_used: [{ skill: "api-recon", calls: -1 }] },
    { ...outcome, lesson_exposure: { shadow: 1, injected: -1 } },
  ]) {
    assert.equal(isTaskOutcome(invalid), false);
  }
});

test("a request failure cannot keep a stale successful zero-result conclusion", () => {
  const { task, outcome } = taskOutcomeScenario("9105");
  assert.equal(taskOutcomePresentation(task, outcome, "503").state, "environment_error");
  assert.equal(outcome.findings.total, 0); // presentation never rewrites the last successful snapshot
});

test("a failed task without request failure is not an environment error", () => {
  const { task, outcome } = taskOutcomeScenario("9107");
  assert.equal(taskOutcomePresentation(task, outcome).state, "task_failed");
});

test("pending or dismissed records cannot become completed_no_findings", () => {
  const { task, outcome } = taskOutcomeScenario("9106");
  assert.equal(outcome.findings.confirmed, 0);
  assert.equal(taskOutcomePresentation(task, outcome).state, "completed_with_findings");
});

test("reports distinguish incomplete empty results and explicit request errors", () => {
  assert.equal(taskOutcomeScenarioReport("9101"), "");
  assert.match(taskOutcomeScenarioReport("9105"), /Recorded findings: 0/);
  assert.match(taskOutcomeScenarioReport("9106"), /Recorded findings: 2/);
  assert.throws(() => taskOutcomeScenarioReport("9109"), /503/);
  assert.equal(taskOutcomeScenarioReport("ordinary-task"), null);
});


test("localized outcome presentation preserves every lifecycle classification and source snapshot", () => {
  for (const [id, expected] of Object.entries(states)) {
    const { task, outcome, outcomeError } = taskOutcomeScenario(id);
    const original = JSON.stringify(outcome);
    const zh = taskOutcomePresentation(task, outcome, outcomeError, "zh-CN");
    const en = taskOutcomePresentation(task, outcome, outcomeError, "en");
    assert.equal(zh.state, expected);
    assert.equal(en.state, expected);
    assert.match(zh.title, /[\u3400-\u9fff]/u);
    assert.notEqual(zh.title, en.title);
    assert.equal(JSON.stringify(outcome), original);
  }
  const { task, outcome } = taskOutcomeScenario("9105");
  assert.match(taskOutcomePresentation(task, outcome, "", "zh-CN").detail, /不能证明/);
});
