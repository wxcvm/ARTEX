export const mentionKinds = [
  { kind: "finding", label: "finding", alias: "finding" },
  { kind: "asset", label: "asset", alias: "asset" },
  { kind: "company", label: "company", alias: "company" },
  { kind: "endpoint", label: "endpoint", alias: "api" },
  { kind: "ip", label: "IP", alias: "ip" },
  { kind: "app", label: "app", alias: "app" },
  { kind: "root_domain", label: "domain", alias: "domain" },
  { kind: "subdomain", label: "subdomain", alias: "subdomain" },
  { kind: "service", label: "service", alias: "service" },
] as const;

export type MentionKind = (typeof mentionKinds)[number]["kind"];
const mentionLabelsZh: Record<MentionKind, string> = { finding: "发现结果", asset: "资产", company: "企业", endpoint: "接口", ip: "IP", app: "应用", root_domain: "根域名", subdomain: "子域名", service: "服务" };

// Display labels must never change the canonical mention token or API kind.
export function mentionDisplayLabel(kind: MentionKind, locale: "zh-CN" | "en"): string {
  return locale === "en" ? mentionKinds.find((item) => item.kind === kind)?.label ?? kind : mentionLabelsZh[kind] ?? kind;
}

export interface ChatMention {
  kind: MentionKind;
  id: number;
  label: string;
  description: string;
}

export function activeMention(value: string, caret: number) {
  const before = value.slice(0, caret);
  const start = before.lastIndexOf("@");
  if (start < 0 || (start > 0 && /[\w.+/-]/.test(before[start - 1]))) return null;
  const query = before.slice(start + 1);
  if (/[[\]\r\n@]/.test(query) || query.length > 220) return null;
  return { start, end: caret, query };
}

export function mentionSearch(query: string) {
  const text = query.trimStart().toLowerCase();
  for (const item of mentionKinds) {
    for (const alias of [item.label.toLowerCase(), item.alias]) {
      if (text === alias || text.startsWith(`${alias} `) || (/[^a-z]/.test(alias) && text.startsWith(alias))) {
        return { kind: item.kind, query: query.trimStart().slice(alias.length).trim(), categories: [] };
      }
    }
  }
  const categories = mentionKinds.filter(
    (item) => item.label.toLowerCase().startsWith(text) || item.alias.startsWith(text),
  );
  return { kind: "" as const, query: query.trim(), categories };
}

export function mentionToken(item: ChatMention) {
  const kind = mentionKinds.find((entry) => entry.kind === item.kind)?.label ?? "asset";
  const label = item.label
    .replace(/[[\]]/g, (char) => (char === "[" ? "(" : ")"))
    .replace(/\s+/g, " ")
    .slice(0, 100);
  return `@[${kind}#${item.id} ${label}]`;
}

export function selectedMentions(value: string) {
  return [...value.matchAll(/@\[(finding|asset|company|endpoint|IP|app|domain|subdomain|service)#([0-9]+)(?: ([^\]\r\n]*))?\]/g)].map(
    (match) => ({
      token: match[0],
      label: `${match[1]} #${match[2]}${match[3] ? ` · ${match[3]}` : ""}`,
      start: match.index,
    }),
  );
}
