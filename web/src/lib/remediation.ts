import type { RemediationAction, RemediationEvent, RemediationItem, RemediationMutation, RemediationStage } from "./types";

export const remediationStages = ["pending_review", "pending_fix", "pending_acceptance", "accepted"] as const satisfies readonly RemediationStage[];

export const remediationLimits = { q: 200, note: 2000, reference: 200, plan: 3000, owner: 120 } as const;

export type RemediationDraft = { note: string; reference: string; plan: string; owner: string; dueDate: string };

export function remediationDraft(item?: RemediationItem): RemediationDraft {
  return { note: "", reference: "", plan: item?.plan ?? "", owner: item?.owner ?? "", dueDate: item?.due_date ?? "" };
}

export function remediationActions(item: RemediationItem): RemediationAction[] {
  if (item.read_only) return [];
  switch (item.stage) {
    case "pending_review": return ["review"];
    case "pending_fix": return ["update_plan", "submit"];
    case "pending_acceptance":
      return item.submitted_evidence_version === item.evidence_version ? ["accept", "return"] : ["return"];
    case "accepted": return ["reopen"];
  }
}

type DraftError = "unavailable" | "note" | "reference" | "plan" | "scope" | "length";
type MutationResult = { ok: true; value: RemediationMutation } | { ok: false; error: DraftError };

/** Allow only explicit UI-authorized fields; do not silently switch to updated versions. */
export function remediationMutation(
  item: RemediationItem,
  action: RemediationAction,
  draft: RemediationDraft,
  current = true,
): MutationResult {
  if (!current) return { ok: false, error: "scope" };
  if (!remediationActions(item).includes(action)) return { ok: false, error: "unavailable" };
  const note = draft.note.trim();
  if (action !== "update_plan" && !note) return { ok: false, error: "note" };
  if (Array.from(note).length > remediationLimits.note) return { ok: false, error: "length" };
  const base = { expected_version: item.version, note };
  switch (action) {
    case "review":
    case "return":
    case "reopen":
      return { ok: true, value: { ...base, action } };
    case "accept":
      return { ok: true, value: { ...base, action, expected_evidence_version: item.evidence_version } };
    case "submit": {
      const reference = draft.reference.trim();
      if (!reference) return { ok: false, error: "reference" };
      if (Array.from(reference).length > remediationLimits.reference) return { ok: false, error: "length" };
      return { ok: true, value: { ...base, action, reference, expected_evidence_version: item.evidence_version } };
    }
    case "update_plan": {
      const plan = draft.plan.trim();
      const owner = draft.owner.trim();
      const due_date = draft.dueDate.trim();
      if (!plan) return { ok: false, error: "plan" };
      if (Array.from(plan).length > remediationLimits.plan || Array.from(owner).length > remediationLimits.owner) {
        return { ok: false, error: "length" };
      }
      return { ok: true, value: { ...base, action, plan, owner, due_date } };
    }
  }
}

export function mergeRemediationEvents(current: RemediationEvent[], next: RemediationEvent[]): RemediationEvent[] {
  const bySequence = new Map(current.map(event => [event.id, event]));
  for (const event of next) if (!bySequence.has(event.id)) bySequence.set(event.id, event);
  return [...bySequence.values()].sort((a, b) => b.id - a.id);
}
