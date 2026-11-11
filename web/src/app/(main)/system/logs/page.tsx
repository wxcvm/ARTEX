"use client";

import * as React from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { sseUrl } from "@/lib/api";
import { MOCK } from "@/lib/mock/enabled";
import type { LogLine } from "@/lib/types";
import { cn } from "@/lib/utils";

// Mock demo: no backend SSE, so stuff in a few sample log lines.
const MOCK_LOGS: LogLine[] = [
  { seq: 1, ts: "2026-07-26T03:55:00Z", level: "info", tag: "engine", text: "ARTEX v0.1.0 backend listening on :8787 (workers=3)" },
  { seq: 2, ts: "2026-07-26T03:55:01Z", level: "info", tag: "config", text: "LLM configured from DB: anthropic / claude-opus-4-8" },
  { seq: 3, ts: "2026-07-26T03:56:10Z", level: "info", tag: "planner", text: "task t-acme-web: planning round 3, generated intent i-4" },
  { seq: 4, ts: "2026-07-26T03:57:00Z", level: "warn", tag: "guard", text: "block bash: target out of bounds, out.evil.example is not in scope" },
  { seq: 5, ts: "2026-07-26T03:57:30Z", level: "info", tag: "work#1", text: "report_finding: Default Credentials (high) saved to database" },
  { seq: 6, ts: "2026-07-26T03:58:20Z", level: "error", tag: "work#3", text: "intercept: mysqldump matched a destructive rule, awaiting manual approval" },
];

const levelTone: Record<LogLine["level"], string> = {
  info: "text-foreground/80",
  warn: "text-amber-600 dark:text-amber-400",
  error: "text-red-600 dark:text-red-400",
};
const levelDot: Record<LogLine["level"], string> = {
  info: "bg-muted-foreground/40",
  warn: "bg-amber-500",
  error: "bg-red-500",
};

function fmtTime(ts: string) {
  const d = new Date(ts);
  if (isNaN(d.getTime())) return "";
  const p = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}

export default function LogsPage() {
  const [lines, setLines] = React.useState<LogLine[]>([]);
  const [q, setQ] = React.useState("");
  const [level, setLevel] = React.useState<"all" | LogLine["level"]>("all");
  const [paused, setPaused] = React.useState(false);
  const [loadingHistory, setLoadingHistory] = React.useState(false);
  const [hasMore, setHasMore] = React.useState(false);
  const [connection, setConnection] = React.useState<"connecting" | "connected" | "disconnected">("connecting");
  const pendingWhilePaused = React.useRef<LogLine[]>([]);
  const bottom = React.useRef<HTMLDivElement>(null);
  const stick = React.useRef(true);
  const pausedRef = React.useRef(false);
  pausedRef.current = paused;

  // Minimum db_id seen — used as the cursor for loading older history.
  const minDbId = React.useMemo(() => {
    let min = 0;
    for (const l of lines) {
      if (l.db_id && (min === 0 || l.db_id < min)) min = l.db_id;
    }
    return min;
  }, [lines]);

  // Live tail via SSE. on connect, replays whatever is in the ring (includes
  // DB-restored history from the last 100 rows written before restart).
  React.useEffect(() => {
    if (MOCK) {
      setLines(MOCK_LOGS);
      return;
    }
    const es = new EventSource(sseUrl("/api/logs/stream?since=0"));
    es.onopen = () => setConnection("connected");
    es.onerror = () => setConnection("disconnected");
    es.onmessage = (e) => {
      try {
        const l = JSON.parse(e.data) as LogLine;
        if (pausedRef.current) {
          // Pause the display, not the incoming stream; restore missed lines on resume.
          pendingWhilePaused.current = [...pendingWhilePaused.current, l].slice(-3000);
          return;
        }
        setLines((prev) => (prev.some((x) => x.seq === l.seq) ? prev : [...prev, l].slice(-3000)));
      } catch {
        /* ignore malformed frame */
      }
    };
    return () => es.close();
  }, []);

  React.useEffect(() => {
    if (paused || pendingWhilePaused.current.length === 0) return;
    const pending = pendingWhilePaused.current;
    pendingWhilePaused.current = [];
    setLines((prev) => {
      const seen = new Set(prev.map((line) => line.seq));
      const unseen = pending.filter((line) => !seen.has(line.seq));
      return [...prev, ...unseen].slice(-3000);
    });
  }, [paused]);

  // Mark "has more" once we have any db_id in view.
  React.useEffect(() => {
    if (minDbId > 1) setHasMore(true);
  }, [minDbId]);

  async function loadOlderHistory() {
    if (loadingHistory) return;
    setLoadingHistory(true);
    try {
      const params = minDbId > 0 ? `?before=${minDbId}&limit=200` : `?limit=200`;
      const res = await fetch(`/api/logs/history${params}`);
      if (!res.ok) return;
      const data = await res.json() as { items: LogLine[]; has_more: boolean };
      if (data.items?.length) {
        // Assign synthetic seq numbers below current minimum to keep dedup working.
        setLines((prev) => {
          const minSeq = prev.reduce((m, l) => Math.min(m, l.seq ?? 0), 0);
          const older = data.items.map((l, i) => ({
            ...l,
            seq: minSeq - data.items.length + i,
          }));
          // Prepend, dedup by db_id, cap at 5000.
          const merged = [...older, ...prev];
          const seen = new Set<number>();
          const deduped = merged.filter((l) => {
            if (!l.db_id) return true;
            if (seen.has(l.db_id)) return false;
            seen.add(l.db_id);
            return true;
          });
          return deduped.slice(-5000);
        });
      }
      setHasMore(data.has_more ?? false);
    } finally {
      setLoadingHistory(false);
    }
  }

  const filtered = React.useMemo(() => {
    const needle = q.trim().toLowerCase();
    return lines.filter(
      (l) =>
        (level === "all" || l.level === level) &&
        (!needle || l.text.toLowerCase().includes(needle) || l.tag.toLowerCase().includes(needle)),
    );
  }, [lines, q, level]);

  React.useEffect(() => {
    if (stick.current && !paused) bottom.current?.scrollIntoView();
  }, [filtered, paused]);

  const counts = React.useMemo(() => {
    let warn = 0,
      error = 0;
    for (const l of lines) {
      if (l.level === "warn") warn++;
      else if (l.level === "error") error++;
    }
    return { warn, error, total: lines.length };
  }, [lines]);

  return (
    <div className="flex min-w-0 flex-1 flex-col gap-4 pb-5">
      <header className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-border/70 bg-gradient-to-br from-card to-primary/5 px-5 py-5 shadow-sm sm:px-7">
        <div>
          <p className="mb-1 text-[11px] font-semibold tracking-[0.16em] text-primary">ARTEX / LOGS</p>
          <h1 className="text-2xl font-semibold tracking-tight">系统日志</h1>
          <p className="mt-1 text-sm leading-relaxed text-muted-foreground">实时查看后端、任务规划与工作进程的运行记录。</p>
        </div>
        <span role="status" className="inline-flex items-center gap-2 rounded-lg border border-border/70 bg-background/60 px-3 py-2 text-xs text-muted-foreground">
          <span className={cn("size-2 rounded-full", MOCK ? "bg-sky-400" : connection === "connected" ? "bg-emerald-400" : connection === "connecting" ? "bg-amber-400" : "bg-rose-400")} />
          {MOCK ? "演示数据" : connection === "connected" ? "实时连接" : connection === "connecting" ? "正在连接" : "连接中断 · 自动重试"}
        </span>
      </header>
      <div className="flex min-w-0 flex-1 flex-col overflow-hidden rounded-xl border border-border/70 bg-card/90 shadow-sm">
        <div className="flex flex-wrap items-center gap-2 border-b border-border/60 px-4 py-4 sm:px-5">
          <Input
            placeholder="筛选日志内容或标签…"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            className="h-9 min-w-[180px] flex-1 sm:max-w-sm"
          />
          <div className="flex flex-wrap gap-1">
            {(["all", "info", "warn", "error"] as const).map((lv) => (
              <Button
                key={lv}
                size="sm"
                variant={level === lv ? "default" : "outline"}
                className="h-8"
                onClick={() => setLevel(lv)}
              >
                {lv === "all" ? "全部" : lv === "info" ? "信息" : lv === "warn" ? "警告" : "错误"}
              </Button>
            ))}
          </div>
          <Button
            size="sm"
            variant={paused ? "default" : "outline"}
            className="h-8"
            onClick={() => setPaused((p) => !p)}
          >
            {paused ? "继续显示" : "暂停显示"}
          </Button>
          <Button size="sm" variant="outline" className="h-9" onClick={() => { pendingWhilePaused.current = []; setLines([]); }}>
            清空显示
          </Button>
          <span className="ml-auto text-xs text-muted-foreground">
            共 {counts.total} 条 ·{" "}
            <span className="text-amber-600 dark:text-amber-400">{counts.warn} 警告</span> ·{" "}
            <span className="text-red-600 dark:text-red-400">{counts.error} 错误</span>
          </span>
        </div>

        <div
          onScroll={(e) => {
            const el = e.currentTarget;
            stick.current = el.scrollHeight - el.scrollTop - el.clientHeight < 40;
          }}
          className="h-[calc(100dvh-19rem)] min-h-[22rem] min-w-0 overflow-auto bg-card/50 p-3 font-mono text-[13px] leading-6 sm:p-4"
        >
          {hasMore && (
            <div className="flex justify-center py-2">
              <Button
                size="sm"
                variant="outline"
                className="h-7 text-xs"
                disabled={loadingHistory}
                onClick={loadOlderHistory}
              >
                {loadingHistory ? "加载中…" : "加载更早的日志"}
              </Button>
            </div>
          )}
          {filtered.length === 0 ? (
            <p className="px-4 py-16 text-center text-sm text-muted-foreground">
              {q || level !== "all" ? "没有符合筛选条件的日志。" : "暂无日志记录，连接后会自动显示新内容。"}
            </p>
          ) : (
            filtered.map((l) => {
              const body =
                l.tag && l.text.startsWith("[" + l.tag + "]")
                  ? l.text.slice(l.tag.length + 2).trimStart()
                  : l.text;
              return (
                <div key={l.seq} className="flex flex-wrap items-start gap-x-3 gap-y-1 rounded-md border-b border-border/35 px-2 py-2 last:border-b-0 hover:bg-muted/35">
                  <span className={cn("mt-1.5 size-1.5 shrink-0 rounded-full", levelDot[l.level])} />
                  <span className="shrink-0 tabular-nums text-muted-foreground">{fmtTime(l.ts)}</span>
                  {l.tag && (
                    <span className="shrink-0 rounded-md bg-muted px-2 py-0.5 text-[11px] text-foreground/80">{l.tag}</span>
                  )}
                  <span className={cn("min-w-[min(100%,20rem)] flex-1 whitespace-pre-wrap break-words", levelTone[l.level])}>{body}</span>
                </div>
              );
            })
          )}
          <div ref={bottom} />
        </div>
      </div>
    </div>
  );
}
