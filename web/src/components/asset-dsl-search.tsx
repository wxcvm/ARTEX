"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";

import { SearchIcon } from "lucide-react";

import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

// ── DSL autocomplete ──────────────────────────────────────────────────────────
// Shared by the global asset view (/function/assets) and the per-task test assets
// search, so both search boxes behave and look identical.

const DSL_FIELDS: { name: string; desc: string; ops: { op: string; desc: string }[] }[] = [
  {
    name: "domain",
    desc: "Domain (root domain/subdomain/service domain)",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "ip",
    desc: "IPv4/IPv6 address",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "url",
    desc: "Full URL (service/endpoint)",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "root_domain",
    desc: "Root domain",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "page_title",
    desc: "Page title (HTTP service)",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "icp",
    desc: "ICP filing number",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "service_name",
    desc: "Service name (non-HTTP service)",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "app_name",
    desc: "App name",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "bundle_id",
    desc: "App Bundle ID",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "category",
    desc: "App category",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "app_icp",
    desc: "App ICP filing",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "method",
    desc: "HTTP method GET/POST/PUT/…",
    ops: [
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "service_type",
    desc: "Service type: http | other",
    ops: [
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "record_type",
    desc: "DNS record type A/CNAME/MX/…",
    ops: [
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "technology",
    desc: "Technology fingerprint (array field)",
    ops: [
      { op: "=", desc: "fuzzy match" },
      { op: "==", desc: "exact match" },
      { op: "!=", desc: "exclude" },
    ],
  },
  {
    name: "port",
    desc: "Port number (integer)",
    ops: [
      { op: "==", desc: "equals" },
      { op: "!=", desc: "not equals" },
      { op: ">", desc: "greater than" },
      { op: ">=", desc: "greater than or equal" },
      { op: "<", desc: "less than" },
      { op: "<=", desc: "less than or equal" },
    ],
  },
  {
    name: "status_code",
    desc: "HTTP status code (integer)",
    ops: [
      { op: "==", desc: "equals" },
      { op: "!=", desc: "not equals" },
      { op: ">", desc: "greater than" },
      { op: ">=", desc: "greater than or equal" },
      { op: "<", desc: "less than" },
      { op: "<=", desc: "less than or equal" },
    ],
  },
  { name: "company_id", desc: "Owning company ID (integer)", ops: [{ op: "==", desc: "equals" }] },
  { name: "task_id", desc: "Source task ID (integer)", ops: [{ op: "==", desc: "equals" }] },
];

const LOGIC_OPS = [
  { label: "AND", desc: "AND (both conditions must match)" },
  { label: "OR", desc: "OR (either one matches)" },
];

interface DslSuggestion {
  kind: "field" | "operator" | "logic";
  label: string;
  desc: string;
  replaceStart: number;
  replaceEnd: number;
  insertText: string;
}

function getDslSuggestions(text: string, cursor: number): DslSuggestion[] {
  const before = text.slice(0, cursor);
  // Current token: non-whitespace, non-paren run ending at cursor
  const tokenMatch = before.match(/([^\s()]*$)/);
  const currentToken = tokenMatch?.[1] ?? "";
  const tokenStart = cursor - currentToken.length;

  // Token already contains field+operator → typing a value, no suggestions
  if (/^[a-z_]+(==|!=|>=|<=|=|>|<)/.test(currentToken)) return [];

  // Complete known field name → suggest operators for that field
  const exactField = DSL_FIELDS.find((f) => f.name === currentToken.toLowerCase());
  if (exactField) {
    return exactField.ops.map(({ op, desc }) => ({
      kind: "operator",
      label: `${exactField.name}${op}`,
      desc,
      replaceStart: tokenStart,
      replaceEnd: cursor,
      insertText: `${exactField.name}${op}`,
    }));
  }

  // Everything before the current token (trimmed)
  const beforeToken = before.slice(0, tokenStart).trimEnd();
  const afterExpression = beforeToken.length > 0 && !/\b(AND|OR)\s*$/i.test(beforeToken) && !beforeToken.endsWith("(");

  // Current token is a prefix of AND/OR and follows a complete expression
  if (/^(a|an|and|o|or)$/i.test(currentToken) && afterExpression) {
    return LOGIC_OPS.filter((l) => l.label.startsWith(currentToken.toUpperCase())).map(({ label, desc }) => ({
      kind: "logic",
      label,
      desc,
      replaceStart: tokenStart,
      replaceEnd: cursor,
      insertText: `${label} `,
    }));
  }

  // No current token, after a complete expression → suggest AND/OR
  if (!currentToken && afterExpression) {
    return LOGIC_OPS.map(({ label, desc }) => ({
      kind: "logic",
      label,
      desc,
      replaceStart: cursor,
      replaceEnd: cursor,
      insertText: `${label} `,
    }));
  }

  // Default: suggest fields filtered by prefix
  const prefix = currentToken.toLowerCase();
  return DSL_FIELDS.filter((f) => f.name.startsWith(prefix)).map((f) => ({
    kind: "field",
    label: f.name,
    desc: f.desc,
    replaceStart: tokenStart,
    replaceEnd: cursor,
    insertText: f.name,
  }));
}

function applyDslSuggestion(text: string, s: DslSuggestion): { text: string; cursor: number } {
  const newText = text.slice(0, s.replaceStart) + s.insertText + text.slice(s.replaceEnd);
  return { text: newText, cursor: s.replaceStart + s.insertText.length };
}

const KIND_STYLE: Record<string, string> = {
  field: "text-blue-500 dark:text-blue-400",
  operator: "text-amber-500 dark:text-amber-400",
  logic: "text-emerald-500 dark:text-emerald-400",
};

// AssetDslSearch is the shared DSL search box: a monospace input with a
// field/operator/logic autocomplete popover and a status line ("found N" /
// error / loading). Used by both the global asset view and the per-task view.
const dslDescriptionsZh: Record<string, string> = {
  "Domain (root domain/subdomain/service domain)": "域名（根域名 / 子域名 / 服务域名）",
  "IPv4/IPv6 address": "IPv4 / IPv6 地址",
  "Full URL (service/endpoint)": "完整 URL（服务 / 接口）",
  "Root domain": "根域名",
  "Page title (HTTP service)": "页面标题（HTTP 服务）",
  "ICP filing number": "ICP备案号",
  "Service name (non-HTTP service)": "服务名（非 HTTP 服务）",
  "App name": "应用名称",
  "App Bundle ID": "应用包标识",
  "App category": "应用分类",
  "App ICP filing": "应用 ICP 备案",
  "HTTP method GET/POST/PUT/…": "HTTP 方法 GET/POST/PUT/…",
  "Service type: http | other": "服务类型：http | other",
  "DNS record type A/CNAME/MX/…": "DNS 记录类型 A/CNAME/MX/…",
  "Technology fingerprint (array field)": "技术指纹（数组字段）",
  "Port number (integer)": "端口号（整数）",
  "HTTP status code (integer)": "HTTP 状态码（整数）",
  "Owning company ID (integer)": "所属企业 ID（整数）",
  "Source task ID (integer)": "来源任务 ID（整数）",
  "AND (both conditions must match)": "AND（两个条件均需匹配）",
  "OR (either one matches)": "OR（任一条件匹配）",
  "fuzzy match": "模糊匹配",
  "exact match": "精确匹配",
  "exclude": "排除",
  "equals": "等于",
  "not equals": "不等于",
  "greater than": "大于",
  "greater than or equal": "大于等于",
  "less than": "小于",
  "less than or equal": "小于等于"
};

export function AssetDslSearch({
  query,
  onChange,
  loading,
  error,
  count,
}: {
  query: string;
  onChange: (v: string) => void;
  loading: boolean;
  error: string;
  count?: number;
}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const inputRef = React.useRef<HTMLInputElement>(null);
  const [suggestions, setSuggestions] = React.useState<DslSuggestion[]>([]);
  const [selIdx, setSelIdx] = React.useState(0);
  const [open, setOpen] = React.useState(false);

  const refresh = React.useCallback((val: string, pos: number) => {
    const suggs = getDslSuggestions(val, pos);
    setSuggestions(suggs);
    setSelIdx(0);
    setOpen(suggs.length > 0);
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    onChange(val);
    refresh(val, e.target.selectionStart ?? val.length);
  };

  const apply = React.useCallback(
    (s: DslSuggestion) => {
      const cursor = inputRef.current?.selectionStart ?? query.length;
      // use cursor for logic-kind (insert at cursor), replaceStart/End for others
      const adjusted: DslSuggestion =
        s.kind === "logic" && !query.slice(s.replaceStart, s.replaceEnd)
          ? { ...s, replaceStart: cursor, replaceEnd: cursor }
          : s;
      const { text: newText, cursor: newCursor } = applyDslSuggestion(query, adjusted);
      onChange(newText);
      requestAnimationFrame(() => {
        if (!inputRef.current) return;
        inputRef.current.setSelectionRange(newCursor, newCursor);
        inputRef.current.focus();
        refresh(newText, newCursor);
      });
    },
    [query, onChange, refresh],
  );

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (!open || suggestions.length === 0) return;
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setSelIdx((i) => Math.min(i + 1, suggestions.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setSelIdx((i) => Math.max(i - 1, 0));
    } else if (e.key === "Tab" || e.key === "Enter") {
      const s = suggestions[selIdx];
      if (s) {
        e.preventDefault();
        apply(s);
      }
    } else if (e.key === "Escape") {
      setOpen(false);
    }
  };

  const cursorPos = () => inputRef.current?.selectionStart ?? query.length;

  return (
    <div className="flex flex-col gap-1">
      <div className="relative max-w-lg">
        <SearchIcon className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          ref={inputRef}
          placeholder={tr("DSL 搜索：domain=example AND status_code>=400", "DSL search: domain=example AND status_code>=400")}
          value={query}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          onFocus={() => refresh(query, cursorPos())}
          onClick={() => refresh(query, cursorPos())}
          onBlur={() => setTimeout(() => setOpen(false), 120)}
          className="h-8 pl-8 font-mono text-xs"
        />
        {open && suggestions.length > 0 && (
          <div className="absolute top-full left-0 z-50 mt-1 w-max min-w-full max-w-sm rounded-md border bg-popover py-1 shadow-md">
            {suggestions.map((s, i) => (
              <button
                type="button"
                key={i}
                className={cn(
                  "flex w-full cursor-pointer items-center gap-3 px-3 py-1.5 text-left",
                  i === selIdx ? "bg-accent" : "hover:bg-accent/50",
                )}
                onMouseEnter={() => setSelIdx(i)}
                onMouseDown={(e) => {
                  e.preventDefault();
                  apply(s);
                }}
              >
                <span className={cn("shrink-0 font-mono text-xs font-semibold", KIND_STYLE[s.kind])}>{s.label}</span>
                <span className="text-xs text-muted-foreground">{locale === "en" ? s.desc : dslDescriptionsZh[s.desc] ?? s.desc}</span>
              </button>
            ))}
          </div>
        )}
      </div>
      {query.trim() && !open && (
        <p className="pl-1 text-[11px] text-muted-foreground">
          {loading ? tr("正在搜索…", "Searching…") : error ? <span className="text-destructive">{error}</span> : tr(`找到 ${count ?? 0} 项`, `Found ${count ?? 0} items`)}
        </p>
      )}
    </div>
  );
}
