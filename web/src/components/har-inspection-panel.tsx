"use client";

import * as React from "react";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { inspectHar, type HarInspection } from "@/lib/har-inspection";

const MAX_HAR_BYTES = 32 * 1024 * 1024;

export function HarInspectionPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [summary, setSummary] = React.useState<HarInspection | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");

  async function inspectFile(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.[0];
    event.currentTarget.value = "";
    if (!file) return;
    setSummary(null);
    setError("");
    if (file.size > MAX_HAR_BYTES) {
      setError(tr("本地检查支持最多 32 MiB 的 HAR 文件。", "Local inspection supports HAR files up to 32 MiB."));
      return;
    }
    setBusy(true);
    try {
      setSummary(inspectHar(JSON.parse(await file.text())));
    } catch {
      // Do not echo parse errors that may contain snippets of sensitive HAR data.
      setError(tr("无法读取 HAR，请确认文件是包含 log.entries 的有效 JSON。", "Unable to read HAR. Use a valid JSON file with log.entries."));
    } finally {
      setBusy(false);
    }
  }

  const metrics: [string, number][] = summary ? [
    [tr("全部条目", "All entries"), summary.entries],
    [tr("HTTP(S) 条目", "HTTP(S) entries"), summary.httpEntries],
    [tr("无效或非 HTTP 条目", "Invalid or non-HTTP entries"), summary.invalidEntries],
    [tr("已观测状态", "Observed status"), summary.observedStatus],
    [tr("未观测状态", "Unobserved status"), summary.unobservedStatus],
    [tr("无效状态", "Invalid status"), summary.invalidStatus],
    [tr("响应正文已内嵌（含空文本）", "Embedded response text (including empty)"), summary.responseTextPresent],
    [tr("响应正文未内嵌", "Response text not embedded"), summary.responseTextNotEmbedded],
    [tr("含编码标记的正文", "Encoded response text"), summary.encodedResponseText],
    [tr("请求 postData 缺少文本", "Request postData without text"), summary.requestPostDataNotEmbedded],
    [tr("含重复头字段的条目", "Entries with repeated header names"), summary.entriesWithDuplicateHeaders],
    [tr("已记录带时区的时间", "Recorded time with timezone"), summary.captureTimePresent],
    [tr("捕获时间未知", "Capture time unknown"), summary.captureTimeUnknown],
  ] : [];

  return (
    <details className="rounded-xl border border-border/70 bg-card/90 p-3 sm:p-4">
      <summary className="cursor-pointer text-sm font-semibold">{tr("HAR 本地检查", "Local HAR inspection")}</summary>
      <p className="mt-2 text-xs text-muted-foreground">
        {tr("检查浏览器或代理导出的 HAR。只在本地读取并统计，不上传、不写入流量库、不重新发出请求。", "Inspect a browser or proxy HAR export locally. No upload, traffic-store writes, or request replay.")}
      </p>
      <input type="file" accept=".har,.json,application/json" disabled={busy}
        aria-label={tr("选择 HAR 文件进行本地检查", "Select HAR file for local inspection")}
        onChange={(event) => void inspectFile(event)}
        className="mt-3 max-w-full text-sm" />
      {busy && <p role="status" className="mt-2 text-xs">{tr("正在检查文件…", "Inspecting file…")}</p>}
      {error && <p role="alert" className="mt-2 text-sm text-destructive">{error}</p>}
      {summary && (
        <div className="mt-3 space-y-2" role="status">
          <dl className="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-4">
            {metrics.map(([name, count]) => (
              <div key={name} className="rounded-md border p-2">
                <dt className="text-xs text-muted-foreground">{name}</dt>
                <dd className="mt-1 text-lg font-semibold tabular-nums">{count}</dd>
              </div>
            ))}
          </dl>
          <p className="text-xs text-muted-foreground">
            {tr("除全部与无效条目外，统计仅覆盖 HTTP(S) 条目。正文未内嵌不等于空响应；编码标记不证明正文完整；重复头字段可能在普通键值映射中丢失。这里仅检查记录，不代表导入或漏洞结论。", "Except for all and invalid entries, counts cover HTTP(S) entries only. Missing embedded text is not an empty response; an encoding marker does not prove completeness; repeated headers may be lost in a flat map. This is an inspection, not an import or vulnerability conclusion.")}
          </p>
        </div>
      )}
    </details>
  );
}
