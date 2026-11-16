// Local HAR inventory only. No requests, storage, credentials, or raw entry
// data are returned. Missing inline text is not proof of an empty response.
export type HarInspection = {
  entries: number;
  httpEntries: number;
  invalidEntries: number;
  observedStatus: number;
  unobservedStatus: number;
  invalidStatus: number;
  responseTextPresent: number;
  responseTextNotEmbedded: number;
  encodedResponseText: number;
  requestPostDataNotEmbedded: number;
  entriesWithDuplicateHeaders: number;
  captureTimePresent: number;
  captureTimeUnknown: number;
};

function object(value: unknown): Record<string, unknown> | null {
  return typeof value === "object" && value !== null && !Array.isArray(value)
    ? (value as Record<string, unknown>)
    : null;
}

function isHttpRequest(request: Record<string, unknown> | null): boolean {
  if (!request || typeof request.method !== "string" || !request.method.trim() || typeof request.url !== "string") return false;
  try {
    const url = new URL(request.url);
    return (url.protocol === "http:" || url.protocol === "https:") && Boolean(url.hostname);
  } catch {
    return false;
  }
}

function hasDuplicateHeaders(value: unknown): boolean {
  if (!Array.isArray(value)) return false;
  const names = new Set<string>();
  for (const header of value) {
    const raw = object(header)?.name;
    if (typeof raw !== "string") continue;
    const name = raw.trim().toLowerCase();
    if (!name) continue;
    if (names.has(name)) return true;
    names.add(name);
  }
  return false;
}

function hasCaptureTime(value: unknown): boolean {
  // Require a time and an explicit timezone; do not infer the browser's zone.
  return typeof value === "string" &&
    /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/i.test(value) &&
    Number.isFinite(Date.parse(value));
}

export function inspectHar(value: unknown): HarInspection {
  const entries = object(object(value)?.log)?.entries;
  if (!Array.isArray(entries)) throw new Error("Invalid HAR: log.entries must be an array");
  const result: HarInspection = {
    entries: entries.length, httpEntries: 0, invalidEntries: 0,
    observedStatus: 0, unobservedStatus: 0, invalidStatus: 0,
    responseTextPresent: 0, responseTextNotEmbedded: 0, encodedResponseText: 0,
    requestPostDataNotEmbedded: 0, entriesWithDuplicateHeaders: 0,
    captureTimePresent: 0, captureTimeUnknown: 0,
  };

  for (const item of entries) {
    const entry = object(item);
    const request = object(entry?.request);
    if (!isHttpRequest(request)) {
      result.invalidEntries++;
      continue;
    }
    result.httpEntries++;
    const response = object(entry?.response);
    const status = response?.status;
    if (status == null || status === 0) result.unobservedStatus++;
    else if (typeof status === "number" && Number.isInteger(status) && status >= 100 && status <= 599) result.observedStatus++;
    else result.invalidStatus++;

    const content = object(response?.content);
    if (typeof content?.text === "string") {
      result.responseTextPresent++;
      // Empty embedded text is distinct from missing text. An encoding marker
      // is recorded, not decoded or treated as proof of complete binary bytes.
      if (typeof content.encoding === "string" && content.encoding.trim()) result.encodedResponseText++;
    } else {
      result.responseTextNotEmbedded++;
    }
    if (request?.postData != null && typeof object(request.postData)?.text !== "string") result.requestPostDataNotEmbedded++;
    if (hasDuplicateHeaders(request?.headers) || hasDuplicateHeaders(response?.headers)) result.entriesWithDuplicateHeaders++;
    if (hasCaptureTime(entry?.startedDateTime)) result.captureTimePresent++;
    else result.captureTimeUnknown++;
  }
  return result;
}
