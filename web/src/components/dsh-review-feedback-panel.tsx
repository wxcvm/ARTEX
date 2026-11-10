"use client";

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
      setError(err instanceof Error ? err.message : "无法读取 DSH 核查记录"));
  }, [reload]);

  const importFile = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.currentTarget.files?.[0];
    event.currentTarget.value = "";
    if (!file) return;
    setBusy(true);
    setError("");
    setMessage("");
    try {
      if (file.size > 1024 * 1024) throw new Error("回执文件不得超过 1 MiB");
      const feedback: unknown = JSON.parse(await file.text());
      const outcome = await http<ImportResult>("/exploration/findings/dsh-reviews/import", {
        method: "POST",
        body: JSON.stringify(feedback),
      });
      setMessage(`已接收 ${outcome.imported} 条；来源漏洞不存在 ${outcome.missing} 条；任务关联不一致 ${outcome.task_mismatch} 条。`);
      await reload();
    } catch (err) {
      setError(err instanceof Error ? err.message : "无法导入 DSH 回执");
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>DSH 核查回执</CardTitle>
        <CardDescription>
          从 dsh-purge 的「资源 → ARTEX」导出核查回执，再在这里选择文件。
          这些是 DSH 的人工核查进度，不会改写 ARTEX 漏洞确认、修复或任务状态。
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex flex-wrap items-center gap-3">
          <input type="file" accept=".json,application/json"
            aria-label="导入 DSH 核查回执" disabled={busy}
            onChange={(event) => void importFile(event)}
            className="max-w-full text-sm" />
          <Button variant="outline" size="sm" disabled={busy}
            onClick={() => void reload().catch((err: unknown) =>
              setError(err instanceof Error ? err.message : "刷新失败"))}>刷新核查记录</Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {message && <p role="status" className="text-sm">{message}</p>}
        {rows.length > 0 ? (
          <div className="space-y-2">
            <p className="text-muted-foreground text-xs">最近导入的 {rows.length} 条外部核查记录</p>
            {rows.map((row) => (
              <div key={`${row.engagement}:${row.finding_id}`}
                className="rounded-md border p-3 text-sm">
                <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                  <strong>ARTEX 漏洞 #{row.finding_id}</strong>
                  <span>DSH 靶标：{row.engagement}</span>
                  <span>DSH：{labels[row.status] || row.status}</span>
                  <span>ARTEX：{row.artex_status}</span>
                </div>
                <p className="mt-1 text-muted-foreground">
                  DSH 核查证据版本：{row.evidence_version ?? "未知"}
                  {" · "}ARTEX 当前证据版本：{row.current_evidence_version}
                  {row.evidence_changed ? " · 证据已变化，需重新核查" : ""}
                </p>
                {row.note && <p className="mt-1 whitespace-pre-wrap break-words">{row.note}</p>}
              </div>
            ))}
          </div>
        ) : <p className="text-muted-foreground text-sm">还没有导入 DSH 核查记录。</p>}
      </CardContent>
    </Card>
  );
}
