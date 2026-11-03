"use client";

import * as React from "react";

import { CheckIcon, CopyIcon, FileTextIcon, RefreshCwIcon } from "lucide-react";
import { toast } from "sonner";

import { Markdown } from "@/components/markdown";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { copyText } from "@/lib/utils";

export function ReportTab({ taskId }: { taskId: string }) {
  const [report, setReport] = React.useState<string>("");
  const [loading, setLoading] = React.useState(true);
  const [copied, setCopied] = React.useState(false);
  const [error, setError] = React.useState("");
  const [refresh, setRefresh] = React.useState(0);

  React.useEffect(() => {
    let active = true;
    setLoading(true);
    api
      .report(taskId)
      .then((text) => {
        if (active) {
          setReport(text);
          setError("");
        }
      })
      .catch((cause) => {
        if (active) setError(cause instanceof Error ? cause.message : "Report request failed.");
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [taskId, refresh]);

  async function copy() {
    if (!report) return;
    const ok = await copyText(report);
    if (ok) {
      setCopied(true);
      toast.success("Markdown copied");
      setTimeout(() => setCopied(false), 1500);
    } else {
      toast.error("Copy failed; please select the text and copy it manually");
    }
  }

  let content: React.ReactNode;
  if (loading && !report && !error) {
    content = (
      <div className="flex flex-col items-center justify-center gap-2 rounded-md border border-dashed py-16 text-muted-foreground text-sm">
        <FileTextIcon className="size-8 opacity-40" />
        Loading…
      </div>
    );
  } else if (report) {
    content = (
      <div className="max-h-[60vh] overflow-auto rounded-md border bg-muted/20 p-4">
        <Markdown text={report} />
      </div>
    );
  } else if (error) {
    content = null;
  } else {
    content = (
      <div className="flex flex-col items-center justify-center gap-2 rounded-md border border-dashed py-16 text-muted-foreground text-sm">
        <FileTextIcon className="size-8 opacity-40" />
        <span>No report artifact available</span>
        <span className="px-4 text-center text-xs">An absent report does not mean the assessment completed without findings.</span>
      </div>
    );
  }

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between">
        <CardTitle className="flex items-center gap-2 text-sm">
          <FileTextIcon className="size-4" /> Penetration test report (Markdown)
        </CardTitle>
        <div className="flex gap-2">
          <Button size="sm" variant="outline" disabled={loading} onClick={() => setRefresh((n) => n + 1)}>
            <RefreshCwIcon className={loading ? "animate-spin" : ""} /> {error ? "Retry report" : "Refresh report"}
          </Button>
          {report && (
            <Button size="sm" variant="outline" onClick={copy}>
              {copied ? <CheckIcon /> : <CopyIcon />} Copy
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {error && (
          <div role="alert" className="rounded-md border border-destructive/30 p-3 text-sm" data-testid="report-request-error">
            <p className="font-medium text-destructive">Report unavailable — request failed</p>
            <p className="break-words">{error}</p>
            <p className="text-muted-foreground">This is a report-loading error, not evidence of zero findings or a task failure.</p>
            {report && <p className="text-muted-foreground">The last successfully loaded report is shown below and may be stale.</p>}
          </div>
        )}
        {content}
      </CardContent>
    </Card>
  );
}
