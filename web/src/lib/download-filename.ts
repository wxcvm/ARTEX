/** Read RFC 6266 download names without treating quoted semicolons as delimiters. */
export function downloadFilename(contentDisposition: string, fallback: string): string {
  const parameters = new Map<string, string>();
  const pattern = /(?:^|;)\s*([^=;\s]+)\s*=\s*("(?:\\.|[^"\\])*"|[^;]*)/g;
  for (const match of contentDisposition.matchAll(pattern)) {
    const key = match[1].toLowerCase();
    const raw = match[2].trim();
    const value = raw.startsWith('"') && raw.endsWith('"')
      ? raw.slice(1, -1).replace(/\\(.)/g, "$1")
      : raw;
    if (!parameters.has(key)) parameters.set(key, value);
  }

  // filename* wins regardless of parameter order. The optional language tag
  // is metadata, not part of the filename. Decode only once; '+' stays '+'.
  const extended = parameters.get("filename*")?.match(/^UTF-8'[^']*'(.*)$/i);
  if (extended) {
    try {
      const name = usableFilename(decodeURIComponent(extended[1]));
      if (name) return name;
    } catch {
      // A broken escape sequence or invalid UTF-8 must not break the download.
    }
  }
  return usableFilename(parameters.get("filename") ?? "") || fallback;
}

function usableFilename(value: string): string {
  if (/[\x00-\x1f\x7f]/.test(value)) return "";
  // Headers supply a suggested name, not a destination path (RFC 6266 §4.3).
  const name = value.split(/[\\/]/).pop()?.trim() ?? "";
  return name === "." || name === ".." ? "" : name;
}
