"use client";

import * as React from "react";
import { ArrowDownIcon, ArrowUpIcon, PlusIcon } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/lib/api";
import type { FindingReproductionStep, FindingTrafficBinding } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export function FindingReproductionPanel({ findingId, contextTask, readOnly }: {
  findingId: string;
  contextTask?: string;
  readOnly?: boolean;
}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [steps, setSteps] = React.useState<FindingReproductionStep[] | null>(null);
  const [bindings, setBindings] = React.useState<FindingTrafficBinding[]>([]);
  const [error, setError] = React.useState("");
  const [busy, setBusy] = React.useState(false);
  const [dirty, setDirty] = React.useState(false);
  const [reload, setReload] = React.useState(0);

  React.useEffect(() => {
    let active = true;
    setSteps(null);
    setError("");
    setDirty(false);
    api.findingReproduction(findingId, contextTask)
      .then((data) => { if (active) setSteps(data ?? []); })
      .catch((e) => { if (active) setError(e instanceof Error ? e.message : "Request failed"); });
    api.findingTraffic(findingId, contextTask)
      .then((data) => { if (active) setBindings(data.bindings ?? []); })
      .catch(() => { if (active) setBindings([]); });
    return () => { active = false; };
  }, [findingId, contextTask, reload]);

  function update(next: FindingReproductionStep[]) {
    setSteps(next);
    setDirty(true);
  }

  function change(index: number, field: keyof FindingReproductionStep, value: string) {
    if (!steps) return;
    update(steps.map((step, i) => i === index ? { ...step, [field]: value } : step));
  }

  function move(index: number, delta: number) {
    if (!steps || index + delta < 0 || index + delta >= steps.length) return;
    const next = [...steps];
    [next[index], next[index + delta]] = [next[index + delta], next[index]];
    update(next);
  }

  async function save() {
    if (!steps || busy || readOnly) return;
    setBusy(true);
    try {
      const saved = await api.saveFindingReproduction(findingId, steps, contextTask);
      setSteps(saved);
      setDirty(false);
      toast.success(tr("复现记录已保存", "Reproduction records saved"));
    } catch (e) {
      toast.error(e instanceof Error ? e.message : tr("保存失败", "Save failed"));
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">{tr("人工复现记录", "Recorded reproduction steps")}</CardTitle>
        <CardDescription>{tr("记录实际执行过的步骤与观察结果，可关联已有流量证据；保存后随报告与归档保留。", "Record the steps actually performed and observations. Optionally link existing traffic evidence; saved records are included in reports and archives.")}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-3">
        {error && <div role="alert" className="text-sm text-destructive">{error} <Button size="sm" variant="link" onClick={() => setReload((n) => n + 1)}>{tr("重试", "Retry")}</Button></div>}
        {steps === null && !error && <p className="text-sm text-muted-foreground">{tr("正在加载…", "Loading…")}</p>}
        {steps?.length === 0 && <p className="text-sm text-muted-foreground">{tr("暂无人工复现记录", "No recorded reproduction steps")}</p>}
        {steps?.map((step, index) => (
          <div key={index} className="flex flex-col gap-2 rounded-md border p-3">
            <div className="flex items-center gap-2">
              <span className="flex-1 text-sm font-medium">{tr("步骤", "Step")} {index + 1}</span>
              {!readOnly && <>
                <Button size="sm" variant="ghost" aria-label={tr("上移", "Move up")} disabled={busy || index === 0} onClick={() => move(index, -1)}><ArrowUpIcon /></Button>
                <Button size="sm" variant="ghost" aria-label={tr("下移", "Move down")} disabled={busy || index === steps.length - 1} onClick={() => move(index, 1)}><ArrowDownIcon /></Button>
                <Button size="sm" variant="ghost" disabled={busy} onClick={() => update(steps.filter((_, i) => i !== index))}>{tr("删除", "Delete")}</Button>
              </>}
            </div>
            {readOnly ? <>
              <p className="whitespace-pre-wrap break-words text-sm font-medium">{step.title}</p>
              <p className="whitespace-pre-wrap break-words text-sm">{step.observation}</p>
              {step.binding_id && <p className="text-xs text-muted-foreground">{tr("关联流量证据", "Linked traffic evidence")} #{step.binding_id}</p>}
            </> : <>
              <Input value={step.title} aria-label={tr("步骤标题", "Step title")} placeholder={tr("实际执行的步骤", "Step performed")} maxLength={4000} disabled={busy} onChange={(e) => change(index, "title", e.target.value)} />
              <Textarea value={step.observation} aria-label={tr("观察结果", "Observation")} placeholder={tr("实际观察到的结果", "Observed result")} maxLength={4000} rows={2} disabled={busy} onChange={(e) => change(index, "observation", e.target.value)} />
              <select aria-label={tr("关联流量证据", "Linked traffic evidence")} value={step.binding_id || ""} disabled={busy} onChange={(e) => change(index, "binding_id", e.target.value)} className="min-w-0 rounded-md border bg-background p-2 text-sm">
                <option value="">{tr("不关联流量", "No traffic binding")}</option>
                {bindings.map((b) => <option key={b.id} value={b.id}>#{b.id} · {b.snapshot.method} {b.snapshot.url}</option>)}
                {step.binding_id && !bindings.some((b) => b.id === step.binding_id) && <option value={step.binding_id}>#{step.binding_id}</option>}
              </select>
            </>}
          </div>
        ))}
        {!readOnly && steps && <div className="flex flex-wrap gap-2">
          <Button variant="outline" size="sm" disabled={busy} onClick={() => update([...steps, { title: "", observation: "" }])}><PlusIcon />{tr("添加步骤", "Add step")}</Button>
          <Button size="sm" disabled={busy || !dirty || steps.some((step) => !step.title.trim())} onClick={() => void save()}>{busy ? tr("保存中…", "Saving…") : tr("保存记录", "Save records")}</Button>
          {dirty && <Button variant="ghost" size="sm" disabled={busy} onClick={() => setReload((n) => n + 1)}>{tr("放弃未保存修改", "Discard unsaved changes")}</Button>}
        </div>}
      </CardContent>
    </Card>
  );
}
