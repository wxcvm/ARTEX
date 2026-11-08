import type { Finding, Task, TaskOutcome, TaskStatus } from "../types";

// Deterministic, read-only browser fixtures. Only the NEXT_PUBLIC_MOCK client
// calls these helpers; production task IDs and query parameters are unchanged.
// Open /function/tasks/detail?id=9101 (through 9111), then Overview or Report.
const cases: Record<string, { status: TaskStatus; total?: number; paused?: boolean; queued?: boolean; outcomeError?: string; reportError?: string; outcomeStatus?: string }> = {
  "9101": { status: "running" },
  "9102": { status: "paused", paused: true },
  "9103": { status: "queued", queued: true },
  "9104": { status: "timeout" },
  "9105": { status: "done" },
  "9106": { status: "done", total: 2 },
  "9107": { status: "failed" },
  "9108": { status: "done", outcomeError: "Mock outcome request failed (503)." },
  "9109": { status: "done", reportError: "Mock report request failed (503)." },
  "9110": { status: "running", paused: true },
  "9111": { status: "done", outcomeStatus: "running" },
};

export function mockTaskOutcome(task: Task, findings: Pick<Finding, "status" | "severity">[] = []): TaskOutcome {
  const confirmed = findings.filter((f) => ["confirmed", "resolved", "fixed"].includes(f.status));
  return {
    task_id: Number(task.id) || 0,
    status: task.status,
    category_id: task.category_id ?? 0,
    company_id: task.company_ids?.[0] ?? 0,
    findings: {
      total: findings.length,
      confirmed: confirmed.length,
      open: findings.filter((f) => ["pending", "in_progress"].includes(f.status)).length,
      false_positive: findings.filter((f) => f.status === "false_positive").length,
      duplicate: findings.filter((f) => f.status === "duplicate").length,
      critical: confirmed.filter((f) => f.severity === "critical").length,
      high: confirmed.filter((f) => f.severity === "high").length,
      medium: confirmed.filter((f) => f.severity === "medium").length,
      low: confirmed.filter((f) => f.severity === "low").length,
    },
    retests: { reproduced: 0, fixed: 0, inconclusive: 0 },
    tokens: { input_tokens: 0, output_tokens: 0, cache_read_tokens: 0, cache_write_tokens: 0 },
    tool_calls: 0,
    skill_misses: 0,
    wall_seconds: 0,
    skills_used: [],
    lesson_exposure: { shadow: 0, injected: 0 },
  };
}

export function taskOutcomeScenario(id: string) {
  const example = cases[id];
  if (!example) return null;
  const task: Task = {
    id,
    name: `Result-state fixture ${id}`,
    description: "Deterministic UI result-state validation; no assessment is executed.",
    goal: "Keep lifecycle, recorded results, and request failures distinguishable.",
    status: example.status,
    paused: example.paused ?? false,
    queued: example.queued ?? false,
    created_at: "2026-11-02T00:00:00Z",
    in_flight: 0,
    coverage_enabled: false,
  };
  const findings: Finding[] = Array.from({ length: example.total ?? 0 }, (_, index) => ({
    id: `fixture-${id}-${index + 1}`,
    task_id: id,
    name: `Review item ${index + 1}`,
    vulnclass: "Mock review record",
    severity: "low",
    status: "pending",
    summary: "Synthetic record for result-state validation, not a confirmed vulnerability.",
    evidence: "Mock fixture only.",
    ts: task.created_at,
  }));
  const outcome = mockTaskOutcome(task, findings);
  outcome.status = example.outcomeStatus ?? example.status;
  return { task, outcome, findings, outcomeError: example.outcomeError, reportError: example.reportError };
}

export function taskOutcomeScenarioReport(id: string): string | null {
  const fixture = taskOutcomeScenario(id);
  if (!fixture) return null;
  if (fixture.reportError) throw new Error(fixture.reportError);
  if (fixture.task.status !== "done" || fixture.task.paused || fixture.task.queued || fixture.outcome.status !== "done") {
    return "";
  }
  return `# Result-state validation report\n\nTask status: ${fixture.task.status}\n\nRecorded findings: ${fixture.outcome.findings.total}\n\nThis is a deterministic mock artifact, not an assessment of a target.`;
}
