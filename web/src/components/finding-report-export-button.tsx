"use client";

import * as React from "react";
import { DownloadIcon, Loader2Icon } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { buildFindingReportMarkdown, collectFindingReport, findingReportFilename } from "@/lib/finding-report-export";

export function FindingReportExportButton({ findingId, contextTask }: { findingId: string; contextTask?: string }) {
  const [loading, setLoading] = React.useState(false);
  const busy = React.useRef(false);

  async function exportMarkdown() {
    if (busy.current) return;
    busy.current = true;
    setLoading(true);
    try {
      const data = await collectFindingReport(api, findingId, contextTask);
      const markdown = buildFindingReportMarkdown(data);
      const url = URL.createObjectURL(new Blob([markdown], { type: "text/markdown;charset=utf-8" }));
      const link = document.createElement("a");
      link.href = url;
      link.download = findingReportFilename(data.finding);
      document.body.append(link);
      link.click();
      link.remove();
      window.setTimeout(() => URL.revokeObjectURL(url), 1000);
      if (data.packets.some((p) => p.request.binary || p.response.binary ||
        p.request.truncated || p.response.truncated)) {
        toast.warning("报告已导出；二进制或截断正文已注明");
      } else {
        toast.success("已导出漏洞 Markdown 报告");
      }
    } catch (error) {
      toast.error("导出失败：" + (error instanceof Error ? error.message : "请重试"));
    } finally {
      busy.current = false;
      setLoading(false);
    }
  }

  return (
    <Button variant="outline" size="sm" onClick={() => void exportMarkdown()} disabled={loading}>
      {loading ? <Loader2Icon className="animate-spin" /> : <DownloadIcon />}
      导出 Markdown
    </Button>
  );
}
