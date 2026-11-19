"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import { ClipboardList } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { remediationLimits, remediationStages } from "@/lib/remediation";
import { statusLabel } from "@/lib/status";
import type { RemediationList } from "@/lib/types";
import { uiText } from "@/lib/ui-locale";
import { useUiLocale } from "@/lib/ui-locale.client";
import { RemediationRecord } from "./_components/remediation-record";
import { stageLabel } from "./_components/labels";

const selectClass = "h-8 rounded-lg border bg-background px-2 text-sm";
const pageSize = 25;
const severities = ["critical", "high", "medium", "low", "info"] as const;
const initial = { q: "", stage: "", severity: "", page: 1 };

type ListRead = {
  key: string;
  phase: "loading" | "ready" | "error";
  data?: RemediationList;
  error?: string;
};

function Workbench() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const id = useSearchParams().get("finding") ?? "";
  const [query, setQuery] = useState("");
  const [stage, setStage] = useState("");
  const [severity, setSeverity] = useState("");
  const [filters, setFilters] = useState(initial);
  const [revision, setRevision] = useState(0);
  const [read, setRead] = useState<ListRead | null>(null);
  const key = JSON.stringify(filters);
  const active = read?.key === key ? read : null;
  const list = active?.data;
  const phase = active?.phase ?? "loading";
  const refresh = () => {
    setRead(previous => previous?.key === key ? { ...previous, phase: "loading", error: undefined } : previous);
    setRevision(value => value + 1);
  };

  useEffect(() => {
    let live = true;
    setRead(previous => ({ key, phase: "loading", data: previous?.key === key ? previous.data : undefined }));
    api.remediationList({ ...filters, limit: pageSize })
      .then(data => {
        if (!live) return;
        const last = Math.max(1, Math.ceil(data.total / data.page_size));
        if (filters.page > last) {
          setFilters(current => current === filters ? { ...current, page: last } : current);
          return;
        }
        setRead({ key, phase: "ready", data });
      })
      .catch((e: unknown) => {
        if (live) setRead(previous => ({
          key, phase: "error", data: previous?.key === key ? previous.data : undefined,
          error: e instanceof Error ? e.message : "Request failed",
        }));
      });
    return () => { live = false; };
  }, [filters, key, revision]);

  return (
    <main className="mx-auto w-full max-w-7xl space-y-5 p-4 sm:p-8">
      <header className="space-y-2">
        <h1 className="flex items-center gap-2 text-2xl font-semibold">
          <ClipboardList className="size-6" aria-hidden /> {tr("整改工作台", "Remediation workbench")}
        </h1>
        <p className="text-sm text-muted-foreground">
          {tr("只管理已有发现的人工复核、整改计划和验收记录。这里不会启动检查、调用模型或改写原始证据。", "Manage human review, remediation plans, and acceptance for existing findings. This does not run checks, call models or rewrite evidence.")}
        </p>
      </header>
      <div className="grid items-start gap-4 lg:grid-cols-[minmax(16rem,1fr)_minmax(0,1.7fr)]">
        <section className="min-w-0 space-y-4 rounded-xl border bg-card p-4" aria-labelledby="remediation-list-heading">
          <h2 id="remediation-list-heading" className="text-lg font-semibold">{tr("整改记录", "Remediation records")}</h2>
          <form className="space-y-2" onSubmit={event => {
            event.preventDefault();
            setFilters({ q: query.trim(), stage, severity, page: 1 });
          }}>
            <Label htmlFor="remediation-q">{tr("搜索发现标题或责任人", "Search finding title or owner")}</Label>
            <Input id="remediation-q" value={query} maxLength={remediationLimits.q} onChange={event => setQuery(event.target.value)} />
            <div className="flex flex-wrap gap-2">
              <select aria-label={tr("整改阶段", "Remediation stage")} className={selectClass} value={stage} onChange={event => setStage(event.target.value)}>
                <option value="">{tr("全部阶段", "All stages")}</option>
                {remediationStages.map(value => <option key={value} value={value}>{stageLabel(locale, value)}</option>)}
              </select>
              <select aria-label={tr("严重程度", "Severity")} className={selectClass} value={severity} onChange={event => setSeverity(event.target.value)}>
                <option value="">{tr("全部严重程度", "All severities")}</option>
                {severities.map(value => <option key={value} value={value}>{value === "info" ? tr("信息", "Informational") : statusLabel("severity", value, locale)}</option>)}
              </select>
              <Button type="submit">{tr("筛选", "Filter")}</Button>
              <Button type="button" variant="outline" onClick={refresh} disabled={phase === "loading"}>{tr("刷新", "Refresh")}</Button>
            </div>
          </form>
          {phase === "loading" && <p role="status" className="text-sm text-muted-foreground">{tr("正在读取整改记录…", "Loading remediation records…")}</p>}
          {phase === "error" && <div role="alert" className="space-y-2 text-sm">
            <p className="text-destructive">{tr("本次未能读取整改记录。", "Remediation records could not be loaded.")}</p>
            {active?.error && <p className="break-words text-muted-foreground">{active.error}</p>}
            {list && <p className="text-muted-foreground">{tr("下方仅是上次成功读取的记录，可能已过期。", "The last successful results are shown below and might be out of date.")}</p>}
            <Button size="sm" variant="outline" onClick={refresh}>{tr("重试", "Retry")}</Button>
          </div>}
          {phase === "ready" && list?.items.length === 0 && <p className="text-sm text-muted-foreground">{tr("没有符合条件的记录。", "No matching records.")}</p>}
          {list && list.items.length > 0 && (
            <ul className="space-y-2">
              {list.items.map(item => (
                <li key={item.finding_id}>
                  <Link
                    href={`/function/remediation?finding=${encodeURIComponent(item.finding_id)}#remediation-detail`}
                    aria-current={item.finding_id === id ? "true" : undefined}
                    className="block rounded-lg border p-3 hover:bg-muted/50 aria-[current=true]:border-primary"
                  >
                    <div className="break-words text-sm font-medium">{item.title || `#${item.finding_id}`}</div>
                    <div className="mt-1 flex flex-wrap gap-x-2 gap-y-1 text-xs text-muted-foreground">
                      <span>#{item.finding_id}</span><span>{statusLabel("severity", item.severity, locale)}</span>
                      <span>{stageLabel(locale, item.stage)}</span>{item.read_only && <span>{tr("只读", "Read-only")}</span>}
                      {item.needs_review && <span className="font-medium text-amber-700 dark:text-amber-300">{tr("证据待复核", "Evidence needs review")}</span>}
                    </div>
                    {(item.owner || item.due_date) && <div className="mt-1 break-words text-xs text-muted-foreground">
                      {item.owner && <span>{tr("责任人：", "Owner: ")}{item.owner}</span>}
                      {item.owner && item.due_date && " · "}
                      {item.due_date && <span>{tr("期限：", "Due: ")}{item.due_date}</span>}
                    </div>}
                  </Link>
                </li>
              ))}
            </ul>
          )}
          {list && (
            <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <Button size="sm" variant="outline" disabled={phase !== "ready" || list.page <= 1} onClick={() => setFilters(current => ({ ...current, page: list.page - 1 }))}>
                {tr("上一页", "Previous")}
              </Button>
              <Button size="sm" variant="outline" disabled={phase !== "ready" || list.page * list.page_size >= list.total} onClick={() => setFilters(current => ({ ...current, page: list.page + 1 }))}>
                {tr("下一页", "Next")}
              </Button>
              <span>{tr(`第 ${list.page} 页 · 共 ${list.total} 条`, `Page ${list.page} · ${list.total} records`)}</span>
            </div>
          )}
        </section>
        <div id="remediation-detail" className="min-w-0 scroll-mt-20">
          {id ? <RemediationRecord key={id} id={id} onChanged={refresh} /> :
            <p className="rounded-xl border bg-card p-4 text-sm text-muted-foreground">{tr("从列表选择一条记录，查看阶段、当前证据版本及事件历史。", "Select a record to review its stage, evidence version, and event history.")}</p>}
        </div>
      </div>
    </main>
  );
}

export default function RemediationPage() {
  return <Suspense fallback={null}><Workbench /></Suspense>;
}
