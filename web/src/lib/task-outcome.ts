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
  locale: "zh-CN" | "en" = "en",
): TaskOutcomePresentation {
  const tr = (zh: string, en: string) => locale === "en" ? en : zh;
  if (requestError) {
    return {
      state: "environment_error",
      title: tr("结果数据暂不可用", "Result data unavailable"),
      detail: tr("结果请求失败，评估结论未知；任务状态不受影响。", "The outcome request failed. The assessment result is unknown; this does not change the task status."),
    };
  }
  if (!outcome) return { state: "loading", title: tr("正在加载结果…", "Loading result data…"), detail: tr("尚无可用的评估结论。", "No result conclusion is available yet.") };

  const status = task.queued ? "queued" : task.paused ? "paused" : task.status;
  if (status === "failed") {
    return {
      state: "task_failed",
      title: tr("任务失败 — 评估未完成", "Task failed — assessment incomplete"),
      detail: tr("任务以失败结束。不能仅凭此结果认定环境故障或评估无风险。", "The task ended in failure. The result alone does not identify an environment failure or establish a clean assessment."),
    };
  }
  if (status !== "done" || outcome.status !== "done") {
    return {
      state: "incomplete",
      title: status === "timeout" ? tr("达到时限 — 评估未完成", "Time limit reached — assessment incomplete") : tr("评估未完成", "Assessment incomplete"),
      detail: status === "done"
        ? tr("任务与结果快照尚未一致，请等待一致的完成结果。", "The task and outcome snapshots are still settling. Wait for a consistent completed result.")
        : tr(`任务状态：${({ created: "已创建", queued: "排队中", running: "运行中", paused: "已暂停", timeout: "已超时", done: "已完成" } as Record<string, string>)[status] ?? status}。这些数量仅为当前观察记录，不代表评估已完成。`, `Task state: ${status}. Counts are observations so far, not a completed assessment.`),
    };
  }
  if (outcome.findings.total === 0) {
    return {
      state: "completed_no_findings",
      title: tr("已完成 — 未登记发现", "Completed — no findings recorded"),
      detail: tr("任务已结束，结果中的登记发现为零。这不能证明目标不存在漏洞。", "The task is done and the outcome contains zero recorded findings. This is not proof that the target is vulnerability-free."),
    };
  }
  return {
    state: "completed_with_findings",
    title: tr("已完成 — 已登记发现", "Completed — findings recorded"),
    detail: tr("登记发现含各自的复核状态；总数不等于已确认漏洞数。", "Recorded findings include their individual review dispositions; the total is not a count of confirmed vulnerabilities."),
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
