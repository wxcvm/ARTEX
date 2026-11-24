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
