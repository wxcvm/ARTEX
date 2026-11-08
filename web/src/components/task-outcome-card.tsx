"use client";

import * as React from "react";

import { RefreshCwIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { isTaskOutcome, taskOutcomePresentation } from "@/lib/task-outcome";
import type { Task, TaskOutcome } from "@/lib/types";

export function TaskOutcomeCard({ task }: { task: Task }) {
  const [snapshot, setSnapshot] = React.useState<{ taskId: string; data: TaskOutcome | null; error: string }>({
    taskId: task.id, data: null, error: "",
  });
  const [refresh, setRefresh] = React.useState(0);
  const [refreshing, setRefreshing] = React.useState(false);

  React.useEffect(() => {
    let active = true;
    let inFlight = false;
    const load = async () => {
      if (inFlight) return;
      inFlight = true;
      setRefreshing(true);
      try {
        const data = await api.taskOutcome(task.id);
        if (!isTaskOutcome(data)) throw new Error("The outcome response is incomplete or invalid.");
        if (active) setSnapshot({ taskId: task.id, data, error: "" });
      } catch (error) {
        if (active) {
          setSnapshot((old) => ({
            taskId: task.id,
            data: old.taskId === task.id ? old.data : null,
            error: error instanceof Error ? error.message : "Outcome request failed.",
          }));
        }
      } finally {
        inFlight = false;
        if (active) setRefreshing(false);
      }
    };
    void load();
    // Completed tasks can be revived; continue observing rather than freezing a
    // zero-result card after the first terminal snapshot.
    const timer = setInterval(() => void load(), 5000);
    return () => { active = false; clearInterval(timer); };
  }, [task.id, refresh]);

  const current = snapshot.taskId === task.id ? snapshot : { data: null, error: "" };
  const presentation = taskOutcomePresentation(task, current.data, current.error);

  return (
    <Card className="mb-4 gap-3" data-testid="task-outcome-card" data-outcome={presentation.state}>
      <CardHeader className="flex-row items-start justify-between gap-3">
        <div className="space-y-1.5" aria-live="polite">
          <CardDescription>Assessment result</CardDescription>
          <CardTitle className="text-base">{presentation.title}</CardTitle>
          <p className="text-sm text-muted-foreground">{presentation.detail}</p>
        </div>
        <Button variant="outline" size="sm" disabled={refreshing} onClick={() => setRefresh((n) => n + 1)}>
          <RefreshCwIcon className={refreshing ? "animate-spin" : ""} /> {current.error ? "Retry result" : "Refresh result"}
        </Button>
      </CardHeader>
      {(current.data || current.error) && (
        <CardContent className="space-y-2 text-sm">
          {current.error && <p role="alert" className="break-words text-destructive">Outcome request error: {current.error}</p>}
          {current.data && (
            <>
              {current.error && <p className="text-muted-foreground">Last successful snapshot — may be stale:</p>}
              <dl className="flex flex-wrap gap-x-6 gap-y-2 tabular-nums">
                <div><dt className="text-muted-foreground">Recorded findings</dt><dd>{current.data.findings.total}</dd></div>
                <div><dt className="text-muted-foreground">Confirmed / resolved / fixed</dt><dd>{current.data.findings.confirmed}</dd></div>
                <div><dt className="text-muted-foreground">Pending review</dt><dd>{current.data.findings.open}</dd></div>
                <div><dt className="text-muted-foreground">Tool calls</dt><dd>{current.data.tool_calls}</dd></div>
                <div><dt className="text-muted-foreground">Skills loaded</dt><dd>{current.data.skills_used.length}</dd></div>
                <div><dt className="text-muted-foreground">Prior lessons injected</dt><dd>{current.data.lesson_exposure.injected}</dd></div>
              </dl>
              <div className="space-y-1 border-t pt-2 text-xs text-muted-foreground">
                <p>Learning trace · observed usage, not a measure of effectiveness</p>
                <p>
                  Skills: {current.data.skills_used.length
                    ? current.data.skills_used.map((skill) => `${skill.skill} (${skill.calls})`).join(", ")
                    : "none recorded"}
                </p>
                <p>
                  Prior lessons: {current.data.lesson_exposure.shadow} shadow-only,
                  {" "}{current.data.lesson_exposure.injected} injected.
                </p>
              </div>
            </>
          )}
        </CardContent>
      )}
    </Card>
  );
}
