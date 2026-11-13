"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";


import * as React from "react";

import { CheckIcon, CopyIcon, FileTextIcon, RefreshCwIcon } from "lucide-react";
import { toast } from "sonner";

import { Markdown } from "@/components/markdown";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { copyText } from "@/lib/utils";

export function ReportTab({ taskId }: { taskId: string }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [report, setReport] = React.useState<string>("");
  const [reportTaskId, setReportTaskId] = React.useState(taskId);
  // When the task changes, never display the previous task's report.
  const currentReport = reportTaskId === taskId ? report : "";
  const [loading, setLoading] = React.useState(true);
  const [copied, setCopied] = React.useState(false);
  const [error, setError] = React.useState("");
  const [errorTaskId, setErrorTaskId] = React.useState(taskId);
  const currentError = errorTaskId === taskId ? error : "";
  const [refresh, setRefresh] = React.useState(0);

  React.useEffect(() => {
    let active = true;
    setLoading(true);
    api
      .report(taskId)
      .then((text) => {
        if (active) {
          setReport(text);
          setReportTaskId(taskId);
          setErrorTaskId(taskId);
          setError("");
        }
      })
      .catch((cause) => {
        if (active) {
          setErrorTaskId(taskId);
          setError(cause instanceof Error ? cause.message : "报告请求失败");
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [taskId, refresh]);

  React.useEffect(() => {
    setCopied(false);
  }, [taskId]);

  async function copy() {
    if (!currentReport) return;
    const ok = await copyText(currentReport);
    if (ok) {
      setCopied(true);
      toast.success(tr("报告 Markdown 已复制", "Report Markdown copied"));
      setTimeout(() => setCopied(false), 1500);
    } else {
      toast.error(tr("复制失败，请选中报告内容手动复制", "Copy failed; select the report text and copy manually"));
    }
  }

  let content: React.ReactNode;
  if (loading && !currentReport && !currentError) {
    content = (
      <div className="flex min-h-64 flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border/70 bg-muted/10 px-4 py-14 text-center text-sm text-muted-foreground">
        <FileTextIcon className="size-8 opacity-50" />
        {tr("正在加载任务报告…", "Loading task report…")}
      </div>
    );
  } else if (currentReport) {
    content = (
      <article className="max-h-[65dvh] min-w-0 overflow-auto rounded-xl border border-border/70 bg-background/60 p-4 text-sm leading-7 sm:p-6">
        <Markdown text={currentReport} />
      </article>
    );
  } else if (currentError) {
    content = null;
  } else {
    content = (
      <div className="flex min-h-64 flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-border/70 bg-muted/10 px-4 py-14 text-center text-sm text-muted-foreground">
        <FileTextIcon className="size-8 opacity-50" />
        <strong className="text-foreground">{tr("当前没有可用的任务报告", "No task report available")}</strong>
        <span className="max-w-lg leading-relaxed">{tr("没有生成报告不等于任务没有发现结果。可以稍后刷新，或到“发现结果”标签查看已有数据。", "A missing report does not mean there are no findings. Refresh later or view existing results in Findings.")}</span>
      </div>
    );
  }

  return (
    <Card className="min-w-0 gap-0 overflow-hidden rounded-xl border-border/70 bg-card/90 py-0 shadow-sm">
      <CardHeader className="flex flex-wrap items-center justify-between gap-3 border-b border-border/70 px-4 py-4 sm:px-6">
        <div className="min-w-0">
          <CardTitle className="flex items-center gap-2 text-base font-semibold">
            <FileTextIcon className="size-4" /> {tr("任务报告", "Task report")}
          </CardTitle>
          <p className="mt-1 text-xs text-muted-foreground">{tr("Markdown 报告 · 可滚动阅读、复制原文", "Markdown report · scroll to read or copy the original")}</p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button size="sm" variant="outline" disabled={loading} onClick={() => setRefresh((n) => n + 1)}>
            <RefreshCwIcon className={loading ? "animate-spin" : ""} /> {currentError ? tr("重试加载", "Retry") : tr("刷新报告", "Refresh report")}
          </Button>
          {currentReport && (
            <Button size="sm" variant="outline" onClick={copy}>
              {copied ? <CheckIcon /> : <CopyIcon />} {copied ? tr("已复制", "Copied") : tr("复制 Markdown", "Copy Markdown")}
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="min-w-0 space-y-3 px-4 py-4 sm:px-6">
        {currentError && (
          <div role="alert" className="rounded-md border border-destructive/30 p-3 text-sm" data-testid="report-request-error">
            <p className="font-medium text-destructive">{tr("无法加载任务报告", "Unable to load task report")}</p>
            <p className="break-words">{currentError}</p>
            <p className="text-muted-foreground">{tr("这只是报告请求失败，不代表任务失败或没有发现结果。", "The report request failed; this does not mean the task failed or has no findings.")}</p>
            {currentReport && <p className="text-muted-foreground">{tr("下方保留上次成功加载的报告，内容可能不是最新版本。", "The last successfully loaded report is retained below; it may be out of date.")}</p>}
          </div>
        )}
        {content}
      </CardContent>
    </Card>
  );
}
