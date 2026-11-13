"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { http } from "@/lib/api";

type DSHReview = {
  finding_id: number;
  engagement: string;
  source_task_id: number | null;
  evidence_version: number | null;
  status: "pending" | "in_progress" | "needs_evidence" | "reviewed";
  note: string;
  source_updated_at: string;
  artex_status: string;
  current_evidence_version: number;
  evidence_changed: boolean;
};

type ImportResult = { imported: number; missing: number; task_mismatch: number };
const labels: Record<DSHReview["status"], string> = {
  pending: "待核查",
  in_progress: "核查中",
  needs_evidence: "待补证据",
  reviewed: "已人工复核",
};

export function DSHReviewFeedbackPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [rows, setRows] = React.useState<DSHReview[]>([]);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  const [message, setMessage] = React.useState("");

  const reload = React.useCallback(async () => {
    const response = await http<{ items: DSHReview[] }>("/exploration/findings/dsh-reviews");
    setRows(response.items || []);
  }, []);

  React.useEffect(() => {
    void reload().catch((err: unknown) =>
      setError(err instanceof Error ? err.message : tr("无法读取 DSH 核查记录", "Unable to read DSH reviews")));
  }, [reload]);

  const importFile = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.currentTarget.files?.[0];
    event.currentTarget.value = "";
    if (!file) return;
    setBusy(true);
    setError("");
    setMessage("");
    try {
      if (file.size > 1024 * 1024) throw new Error(tr("回执文件不得超过 1 MiB", "Feedback file must not exceed 1 MiB"));
      const feedback: unknown = JSON.parse(await file.text());
      const outcome = await http<ImportResult>("/exploration/findings/dsh-reviews/import", {
        method: "POST",
        body: JSON.stringify(feedback),
      });
      setMessage(tr(`已接收 ${outcome.imported} 条；来源漏洞不存在 ${outcome.missing} 条；任务关联不一致 ${outcome.task_mismatch} 条。`, `Received ${outcome.imported}; missing source findings ${outcome.missing}; task mismatches ${outcome.task_mismatch}.`));
      await reload();
    } catch (err) {
      setError(err instanceof Error ? err.message : tr("无法导入 DSH 回执", "Unable to import DSH feedback"));
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{tr("DSH 核查回执", "DSH review feedback")}</CardTitle>
        <CardDescription>
          {tr("从 dsh-purge 的「资源 → ARTEX」导出核查回执，再在这里选择文件。这些是 DSH 的人工核查进度，不会改写 ARTEX 漏洞确认、修复或任务状态。", "Export review feedback from dsh-purge Resources → ARTEX, then select the file here. This is DSH human-review progress; it does not overwrite ARTEX confirmation, remediation or task state.")}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <input type="file" accept=".json,application/json"
            aria-label={tr("导入 DSH 核查回执", "Import DSH review feedback")} disabled={busy}
            onChange={(event) => void importFile(event)}
            className="max-w-full text-sm" />
          <Button variant="outline" size="sm" disabled={busy}
            onClick={() => void reload().catch((err: unknown) =>
              setError(err instanceof Error ? err.message : tr("刷新失败", "Refresh failed")))}>{tr("刷新核查记录", "Refresh reviews")}</Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {message && <p role="status" className="text-sm">{message}</p>}
        {rows.length > 0 ? (
          <div className="space-y-2">
            <p className="text-muted-foreground text-xs">{tr(`最近导入的 ${rows.length} 条外部核查记录`, `${rows.length} recently imported external reviews`)}</p>
            {rows.map((row) => (
              <div key={`${row.engagement}:${row.finding_id}`}
                className="rounded-md border p-3 text-sm">
                <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                  <strong>{tr("ARTEX 漏洞", "ARTEX finding")} #{row.finding_id}</strong>
                  <span>{tr("DSH 靶标：", "DSH engagement: ")}{row.engagement}</span>
                  <span>DSH：{(locale === "en" ? { pending: "Pending", in_progress: "In progress", needs_evidence: "Needs evidence", reviewed: "Human reviewed" } : labels)[row.status] || row.status}</span>
                  <span>ARTEX：{row.artex_status}</span>
                </div>
                <p className="mt-1 text-muted-foreground">
                  {tr("DSH 核查证据版本：", "DSH review evidence version: ")}{row.evidence_version ?? tr("未知", "Unknown")}
                  {" · "}{tr("ARTEX 当前证据版本：", "ARTEX current evidence version: ")}{row.current_evidence_version}
                  {row.evidence_changed ? tr(" · 证据已变化，需重新核查", " · Evidence changed; review again") : ""}
                </p>
                {row.note && <p className="mt-1 whitespace-pre-wrap break-words">{row.note}</p>}
              </div>
            ))}
          </div>
        ) : <p className="text-muted-foreground text-sm">{tr("还没有导入 DSH 核查记录。", "No DSH reviews imported yet.")}</p>}
      </CardContent>
    </Card>
  );
}
