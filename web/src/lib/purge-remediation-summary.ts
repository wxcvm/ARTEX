/** Local-only, versioned aggregate exchange; never accept finding or credential details. */
export type PurgeRemediationSummary = {
  schema: "purge.remediation-summary.v1";
  capturedAt: string;
  counts: {
    total: number;
    pending: number;
    reviewing: number;
    confirmed: number;
    falsePositive: number;
    resolved: number;
  };
};

export const purgeSummaryMaxBytes = 8192;
const countKeys = ["total", "pending", "reviewing", "confirmed", "falsePositive", "resolved"] as const;

function exactObject(value: unknown, keys: readonly string[]): value is Record<string, unknown> {
  if (!value || typeof value !== "object" || Array.isArray(value)) return false;
  const own = Object.keys(value);
  return own.length === keys.length && keys.every(key => Object.hasOwn(value, key));
}

export function parsePurgeRemediationSummary(text: string): PurgeRemediationSummary | null {
  if (typeof text !== "string" || text.length > purgeSummaryMaxBytes) return null;
  try {
    const raw: unknown = JSON.parse(text);
    if (!exactObject(raw, ["schema", "capturedAt", "counts"]) || raw.schema !== "purge.remediation-summary.v1") return null;
    if (typeof raw.capturedAt !== "string") return null;
    const date = new Date(raw.capturedAt);
    if (!Number.isFinite(date.getTime()) || date.toISOString() !== raw.capturedAt) return null;
    if (!exactObject(raw.counts, countKeys)) return null;
    const values = raw.counts;
    const counts = {} as PurgeRemediationSummary["counts"];
    for (const key of countKeys) {
      const value = values[key];
      if (typeof value !== "number" || !Number.isSafeInteger(value) || value < 0) return null;
      counts[key] = value;
    }
    const { total, pending, reviewing, confirmed, falsePositive, resolved } = counts;
    if (pending + reviewing + confirmed + falsePositive + resolved !== total) return null;
    return { schema: "purge.remediation-summary.v1", capturedAt: raw.capturedAt, counts };
  } catch {
    return null;
  }
}
