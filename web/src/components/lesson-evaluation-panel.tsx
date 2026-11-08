"use client";

import * as React from "react";

import { BarChart3Icon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import type { LessonEvaluation, LessonEvaluationGroup } from "@/lib/types";

const labels: Record<LessonEvaluationGroup["exposure"], string> = {
  inject: "Lessons injected",
  shadow: "Shadow only",
  mixed: "Both modes",
  none: "No exposure recorded",
};

function per(value: number, count: number): string {
  return count > 0 ? (value / count).toFixed(1) : "—";
}

function rate(value: number, count: number): string {
  return count > 0 ? `${(100 * value / count).toFixed(1)}%` : "—";
}

// This panel only runs its aggregate query when explicitly opened.
// Counts describe observed, non-randomized task cohorts; not a causal effect.
export function LessonEvaluationPanel({ taskId }: { taskId: string }) {
  const [snapshot, setSnapshot] = React.useState<{
    taskId: string;
    data: LessonEvaluation | null;
    error: string;
    loading: boolean;
  }>({ taskId, data: null, error: "", loading: false });

  const current = snapshot.taskId === taskId
    ? snapshot
    : { taskId, data: null, error: "", loading: false };

  const load = async () => {
    setSnapshot({ taskId, data: null, error: "", loading: true });
    try {
      const data = await api.taskLessonEvaluation(taskId);
      setSnapshot((prev) => prev.taskId === taskId
        ? { taskId, data, error: "", loading: false }
        : prev);
    } catch (error) {
      setSnapshot((prev) => prev.taskId === taskId
        ? { taskId, data: null, error: error instanceof Error ? error.message : "Comparison unavailable", loading: false }
        : prev);
    }
  };

  return (
    <div className="space-y-3 border-t pt-3">
      <Button size="sm" variant="outline" disabled={current.loading} onClick={() => void load()}>
        <BarChart3Icon className="size-4" />
        {current.loading ? "Comparing…" : current.data ? "Refresh comparison" : "Compare similar completed tasks"}
      </Button>
      {current.error && <p role="alert" className="text-xs text-destructive">{current.error}</p>}
      {current.data && (
        <div className="space-y-2">
          <p className="text-xs text-muted-foreground">
            Completed, unarchived tasks with the same company and category.
            Observational comparison only; no causal lift is implied.
          </p>
          {current.data.groups.length === 0 ? (
            <p className="text-sm text-muted-foreground">No comparable completed tasks recorded.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[680px] text-left text-xs tabular-nums">
                <thead className="border-b text-muted-foreground">
                  <tr>
                    <th className="py-2 pr-3 font-medium">Lesson exposure</th>
                    <th className="py-2 pr-3 font-medium">Tasks</th>
                    <th className="py-2 pr-3 font-medium">Confirmed / task</th>
                    <th className="py-2 pr-3 font-medium">False positives / reviewed</th>
                    <th className="py-2 pr-3 font-medium">Tools / confirmed</th>
                    <th className="py-2 pr-3 font-medium">Tokens / confirmed</th>
                    <th className="py-2 font-medium">Skills / task</th>
                  </tr>
                </thead>
                <tbody>
                  {current.data.groups.map((group) => (
                    <tr key={group.exposure} className="border-b last:border-0">
                      <th scope="row" className="py-2 pr-3 font-medium">{labels[group.exposure]}</th>
                      <td className="py-2 pr-3">{group.tasks}</td>
                      <td className="py-2 pr-3">{per(group.confirmed, group.tasks)}</td>
                      <td className="py-2 pr-3">{rate(
                        group.false_positive,
                        group.confirmed + group.false_positive + group.duplicate,
                      )}</td>
                      <td className="py-2 pr-3">{per(group.tool_calls, group.confirmed)}</td>
                      <td className="py-2 pr-3">{per(
                        group.input_tokens + group.output_tokens, group.confirmed,
                      )}</td>
                      <td className="py-2">{per(group.skill_loads, group.tasks)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          <p className="text-xs text-muted-foreground">
            Rates use pooled counts; zero confirmed findings show — rather than a made-up efficiency.
            "No exposure" means no recorded usage, not a control group. Mixed-mode tasks are separate.
          </p>
        </div>
      )}
    </div>
  );
}
