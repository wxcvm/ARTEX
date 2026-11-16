import type { Activity } from "./types";

// History and live responses can overlap or arrive out of order. A persisted
// activity's seq is its identity; replay must not duplicate rows or token totals.
export function mergeActivities(current: Activity[], incoming: Activity[]): Activity[] {
  const bySeq = new Map<number, Activity>();
  for (const activity of current) bySeq.set(activity.seq, activity);
  for (const activity of incoming) bySeq.set(activity.seq, activity);
  return [...bySeq.values()].sort((left, right) => left.seq - right.seq);
}

// Earlier turns/rounds cannot still be executing merely because their
// persisted result is missing. Do not invent a failure or synthetic result.
export function latestRunBoundary(activity: Activity[]): number {
  let seq = 0;
  for (const step of activity) {
    if (step.kind === "user" || step.kind === "intent" || step.kind === "round" || step.kind === "result") {
      seq = Math.max(seq, step.seq);
    }
  }
  return seq;
}
