import type { Task, TaskOutcome } from "./types";

export type TaskOutcomeState =
  | "loading"
  | "environment_error"
  | "incomplete"
  | "task_failed"
  | "completed_no_findings"
  | "completed_with_findings";

export interface TaskOutcomePresentation {
  state: TaskOutcomeState;
  title: string;
  detail: string;
}

type Lifecycle = Pick<Task, "status" | "paused" | "queued">;

// Keep presentation separate from persisted lifecycle values. Failed/timeout
// tasks are not environment errors, and incomplete zero counts are not success.
export function taskOutcomePresentation(
  task: Lifecycle,
  outcome: TaskOutcome | null,
  requestError = "",
): TaskOutcomePresentation {
  if (requestError) {
    return {
      state: "environment_error",
      title: "Result data unavailable",
      detail: "The outcome request failed. The assessment result is unknown; this does not change the task status.",
    };
  }
  if (!outcome) return { state: "loading", title: "Loading result data…", detail: "No result conclusion is available yet." };

  const status = task.queued ? "queued" : task.paused ? "paused" : task.status;
  if (status === "failed") {
    return {
      state: "task_failed",
      title: "Task failed — assessment incomplete",
      detail: "The task ended in failure. The result alone does not identify an environment failure or establish a clean assessment.",
    };
  }
  if (status !== "done" || outcome.status !== "done") {
    return {
      state: "incomplete",
      title: status === "timeout" ? "Time limit reached — assessment incomplete" : "Assessment incomplete",
      detail: status === "done"
        ? "The task and outcome snapshots are still settling. Wait for a consistent completed result."
        : `Task state: ${status}. Counts are observations so far, not a completed assessment.`,
    };
  }
  if (outcome.findings.total === 0) {
    return {
      state: "completed_no_findings",
      title: "Completed — no findings recorded",
      detail: "The task is done and the outcome contains zero recorded findings. This is not proof that the target is vulnerability-free.",
    };
  }
  return {
    state: "completed_with_findings",
    title: "Completed — findings recorded",
    detail: "Recorded findings include their individual review dispositions; the total is not a count of confirmed vulnerabilities.",
  };
}

// A malformed/partial response must not be silently turned into a zero result.
export function isTaskOutcome(value: unknown): value is TaskOutcome {
  if (!value || typeof value !== "object") return false;
  const v = value as Record<string, unknown>;
  const counts = (object: unknown, keys: string[]) => {
    if (!object || typeof object !== "object") return false;
    const row = object as Record<string, unknown>;
    return keys.every((key) => typeof row[key] === "number" && Number.isSafeInteger(row[key]) && (row[key] as number) >= 0);
  };
  const valid = typeof v.status === "string"
    && counts(v, ["task_id", "category_id", "company_id", "tool_calls", "skill_misses", "wall_seconds"])
    && counts(v.findings, ["total", "confirmed", "open", "false_positive", "duplicate", "critical", "high", "medium", "low"])
    && counts(v.retests, ["reproduced", "fixed", "inconclusive"])
    && counts(v.tokens, ["input_tokens", "output_tokens", "cache_read_tokens", "cache_write_tokens"])
    && counts(v.lesson_exposure, ["shadow", "injected"])
    && Array.isArray(v.skills_used)
    && v.skills_used.every((item: unknown) => {
      if (!item || typeof item !== "object") return false;
      const row = item as Record<string, unknown>;
      return typeof row.skill === "string" && row.skill.length > 0
        && typeof row.calls === "number" && Number.isSafeInteger(row.calls) && row.calls >= 1;
    });
  if (!valid) return false;
  const findings = v.findings as TaskOutcome["findings"];
  return findings.confirmed + findings.open + findings.false_positive + findings.duplicate <= findings.total
    && findings.critical + findings.high + findings.medium + findings.low <= findings.confirmed;
}
