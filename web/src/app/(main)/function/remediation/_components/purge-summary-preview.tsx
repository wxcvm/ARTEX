"use client";

import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { parsePurgeRemediationSummary, purgeSummaryMaxBytes, type PurgeRemediationSummary } from "@/lib/purge-remediation-summary";
import { uiText } from "@/lib/ui-locale";
import { useUiLocale } from "@/lib/ui-locale.client";

export function PurgeRemediationPreview() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [summary, setSummary] = useState<PurgeRemediationSummary | null>(null);
  const [phase, setPhase] = useState<"idle" | "reading" | "ready" | "error">("idle");
  const generation = useRef(0);
  const input = useRef<HTMLInputElement>(null);

  async function readFile(file: File) {
    const request = ++generation.current;
    setSummary(null);
    setPhase("reading");
    try {
      if (file.size > purgeSummaryMaxBytes) throw new Error("oversize");
      const value = parsePurgeRemediationSummary(await file.text());
      if (generation.current !== request) return;
      if (!value) throw new Error("invalid snapshot");
      setSummary(value);
      setPhase("ready");
    } catch {
      if (generation.current === request) setPhase("error");
    }
  }

  function clear() {
    generation.current++;
    setSummary(null);
    setPhase("idle");
    if (input.current) input.current.value = "";
  }

  const segments = summary ? [
    { key: "pending", label: tr("待研判", "Pending"), value: summary.counts.pending, color: "bg-amber-500" },
    { key: "reviewing", label: tr("研判中", "Reviewing"), value: summary.counts.reviewing, color: "bg-sky-500" },
    { key: "confirmed", label: tr("人工确认", "Confirmed"), value: summary.counts.confirmed, color: "bg-emerald-500" },
    { key: "falsePositive", label: tr("误报", "False positive"), value: summary.counts.falsePositive, color: "bg-slate-400" },
    { key: "resolved", label: tr("已处置标记", "Marked resolved"), value: summary.counts.resolved, color: "bg-violet-500" },
  ] : [];
  const reviewed = summary ? summary.counts.total - summary.counts.pending : 0;
  const cards = summary ? [
    [tr("记录总数", "Total"), summary.counts.total],
    [tr("待研判", "Pending"), summary.counts.pending],
    [tr("研判中", "Reviewing"), summary.counts.reviewing],
    [tr("人工确认", "Manually confirmed"), summary.counts.confirmed],
    [tr("误报", "False positive"), summary.counts.falsePositive],
    [tr("已处置标记", "Marked resolved"), summary.counts.resolved],
  ] as const : [];

  return (
    <section className="space-y-3 rounded-xl border bg-card p-4" aria-labelledby="purge-summary-heading">
      <h2 id="purge-summary-heading" className="text-lg font-semibold">{tr("Purge 整改摘要", "Purge remediation summary")}</h2>
      <p className="text-sm text-muted-foreground">
        {tr("从 Purge 敏感资产面板手动导入整改摘要 JSON。仅在本页内存中预览汇总，不上传、不保存、不合并 ARTEX 发现，也不代表自动验收。", "Manually load the aggregate JSON exported by Purge. Preview in this page only: no upload, persistence, merging with ARTEX findings, or automatic acceptance.")}
      </p>
      <div className="flex flex-wrap items-end gap-2">
        <div className="min-w-0 flex-1 space-y-2">
          <Label htmlFor="purge-summary-file">{tr("选择整改摘要 JSON（最多 8 KB）", "Choose summary JSON (up to 8 KB)")}</Label>
          <Input ref={input} id="purge-summary-file" type="file" accept=".json,application/json"
            onChange={event => {
              const file = event.currentTarget.files?.[0];
              event.currentTarget.value = "";
              if (file) void readFile(file);
            }} />
        </div>
        <Button type="button" variant="outline" disabled={phase === "idle"} onClick={clear}>{tr("清除预览", "Clear preview")}</Button>
      </div>
      {phase === "reading" && <p role="status" className="text-sm">{tr("正在检查摘要…", "Checking summary…")}</p>}
      {phase === "error" && <p role="alert" className="text-sm text-destructive">{tr("摘要格式无效、统计不一致或超过限制。未导入任何内容。", "Invalid, inconsistent or oversized summary. Nothing was imported.")}</p>}
      {summary && phase === "ready" && (
        <div className="space-y-3">
          <p className="text-xs text-muted-foreground">{tr("Purge 本地读取时间（非 ARTEX 实时数据）：", "Purge snapshot time (not live ARTEX data): ")}<time dateTime={summary.capturedAt}>{summary.capturedAt}</time></p>
          <div className="space-y-2 rounded-lg border p-3">
            <div className="flex flex-wrap items-center justify-between gap-2 text-sm">
              <span className="font-medium">{tr("复核状态分布", "Review status distribution")}</span>
              <span className="text-muted-foreground tabular-nums">
                {tr("非待研判 " + reviewed + " / " + summary.counts.total, reviewed + " beyond pending / " + summary.counts.total)}
              </span>
            </div>
            <div role="img" aria-label={segments.map(item => item.label + ": " + item.value).join(", ")}
              className="flex h-3 overflow-hidden rounded-full bg-muted">
              {segments.map(item => <div key={item.key} className={item.color}
                style={{ width: summary.counts.total ? (item.value / summary.counts.total * 100) + "%" : "0%" }}
                title={item.label + ": " + item.value} />)}
            </div>
            <div className="flex flex-wrap gap-x-4 gap-y-2">
              {segments.map(item => <div key={item.key} className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <span aria-hidden className={"size-2 rounded-full " + item.color} />
                <span>{item.label}</span><span className="tabular-nums">{item.value}</span>
              </div>)}
            </div>
          </div>
          <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-6">
            {cards.map(([label, value]) => <div key={label} className="min-w-0 rounded-lg border p-3">
              <p className="text-xs text-muted-foreground">{label}</p>
              <p className="mt-1 text-xl font-semibold tabular-nums">{value}</p>
            </div>)}
          </div>
        </div>
      )}
    </section>
  );
}
