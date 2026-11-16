"use client";

import * as React from "react";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { inspectHar, type HarInspection } from "@/lib/har-inspection";
import { http } from "@/lib/api";
import { Button } from "@/components/ui/button";

const MAX_HAR_BYTES = 32 * 1024 * 1024;

type HarImportResult = {
  entries: number; imported: number; duplicates: number; skipped: number;
  failed: number; unprocessed: number; reasons: Record<string, number>;
};
export function HarInspectionPanel({ onImported }: { onImported?: () => void }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [summary, setSummary] = React.useState<HarInspection | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  const [selectedFile, setSelectedFile] = React.useState<File | null>(null);
  const [importResult, setImportResult] = React.useState<HarImportResult | null>(null);

  async function inspectFile(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.[0];
    event.currentTarget.value = "";
    if (!file) return;
    setSummary(null);
    setSelectedFile(null);
    setImportResult(null);
    setError("");
    if (file.size > MAX_HAR_BYTES) {
      setError(tr("本地检查支持最多 32 MiB 的 HAR 文件。", "Local inspection supports HAR files up to 32 MiB."));
      return;
    }
    setBusy(true);
    try {
      setSummary(inspectHar(JSON.parse(await file.text())));
      setSelectedFile(file);
    } catch {
      // Do not echo parse errors that may contain snippets of sensitive HAR data.
      setError(tr("无法读取 HAR，请确认文件是包含 log.entries 的有效 JSON。", "Unable to read HAR. Use a valid JSON file with log.entries."));
    } finally {
      setBusy(false);
    }
  }

  async function importFile() {
    if (!selectedFile) return;
    setBusy(true);
    setError("");
    setImportResult(null);
    try {
      const result = await http<HarImportResult>("/traffic/har/import", { method: "POST", body: selectedFile });
      setImportResult(result);
      if (result.imported > 0) onImported?.();
    } catch {
      setError(tr("HAR 导入失败，请检查流量服务或重新尝试同一文件。", "HAR import failed. Check the traffic service or retry the same file."));
    } finally { setBusy(false); }
  }
  const reasonLabels: Record<string, string> = {
    invalid_entry: tr("无效条目", "Invalid entry"),
    invalid_request: tr("无效请求", "Invalid request"),
    unknown_capture_time: tr("捕获时间未知", "Unknown capture time"),
    unobserved_status: tr("状态未观测", "Unobserved status"),
    invalid_status: tr("无效状态", "Invalid status"),
    invalid_headers: tr("头字段缺失或无效", "Missing or invalid headers"),
    request_body_missing: tr("请求正文缺失", "Request body missing"),
    response_body_missing: tr("响应正文缺失", "Response body missing"),
    request_size_mismatch: tr("请求正文长度不符", "Request body length mismatch"),
    response_size_mismatch: tr("响应正文长度不符", "Response body length mismatch"),
    unsupported_encoding: tr("不支持的编码", "Unsupported encoding"),
    invalid_encoding: tr("编码正文无效", "Invalid encoded body"),
    storage_error: tr("存储失败", "Storage error"),
  };
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
    [tr("声明的页面", "Declared pages"), summary.declaredPages],
    [tr("无效页面记录", "Invalid page records"), summary.invalidPages],
    [tr("唯一页面归属", "Unique page attribution"), summary.entriesWithKnownPage],
    [tr("未知页面引用", "Unknown page reference"), summary.entriesWithUnknownPage],
    [tr("歧义页面引用", "Ambiguous page reference"), summary.entriesWithAmbiguousPage],
    [tr("未提供页面引用", "No page reference"), summary.entriesWithoutPage],
  ] : [];

  return (
    <details className="rounded-xl border border-border/70 bg-card/90 p-3 sm:p-4">
      <summary className="cursor-pointer text-sm font-semibold">{tr("HAR 检查与导入", "HAR inspection and import")}</summary>
      <p className="mt-2 text-xs text-muted-foreground">
        {tr("先本地检查浏览器或代理 HAR；仅在点击导入时交给当前 ARTEX 服务保存。不会重新发出请求。", "Inspect a browser or proxy HAR locally first. Only clicking Import sends it to the current ARTEX service for storage. No request replay.")}
      </p>
      <input type="file" accept=".har,.json,application/json" disabled={busy}
        aria-label={tr("选择 HAR 文件进行本地检查", "Select HAR file for local inspection")}
        onChange={(event) => void inspectFile(event)}
        className="mt-3 max-w-full text-sm" />
      {busy && <p role="status" className="mt-2 text-xs">{tr("正在处理文件…", "Processing file…")}</p>}
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
            {tr("除全部与无效条目外，统计仅覆盖 HTTP(S) 条目。正文未内嵌不等于空响应；编码标记不证明正文完整；重复头字段可能在普通键值映射中丢失。页面统计来自 HAR 的 pages/pageref；唯一归属仅表示文件内引用匹配，不证明真实浏览器会话或账号身份。这里仅检查记录，不代表导入或漏洞结论。", "Except for all and invalid entries, counts cover HTTP(S) entries only. Missing embedded text is not an empty response; an encoding marker does not prove completeness; repeated headers may be lost in a flat map. Page counts come from HAR pages/pageref; a unique match only establishes an in-file reference, not an actual browser session or account identity. This is an inspection, not an import or vulnerability conclusion.")}
          </p>
          <div className="flex flex-wrap items-center gap-2">
            <Button size="sm" disabled={busy || !selectedFile || summary.httpEntries === 0} onClick={() => void importFile()}>
              {tr("导入已记录流量", "Import recorded traffic")}
            </Button>
            <span className="text-xs text-muted-foreground">{tr("文件可能包含凭据；请确认保存到当前 ARTEX 实例。", "The file may contain credentials; import only into the intended ARTEX instance.")}</span>
          </div>
        </div>
      )}
      {importResult && (
        <div className="mt-3 space-y-2 border-t pt-3 text-sm" role="status">
          <p>{tr(
            "已导入 " + importResult.imported + " 条 · 重复 " + importResult.duplicates + " 条 · 跳过 " + importResult.skipped + " 条 · 失败 " + importResult.failed + " 条 · 未处理 " + importResult.unprocessed + " 条",
            "Imported " + importResult.imported + " · duplicates " + importResult.duplicates + " · skipped " + importResult.skipped + " · failed " + importResult.failed + " · unprocessed " + importResult.unprocessed
          )}</p>
          {Object.entries(importResult.reasons ?? {}).map(([reason, count]) => (
            <p className="text-xs text-muted-foreground" key={reason}>{reasonLabels[reason] ?? reason}: {count}</p>
          ))}
          <p className="text-xs text-muted-foreground">{tr("跳过的记录不会补造数据；重复导入相同文件不会覆盖已有流量。", "Skipped records are not fabricated. Reimporting the identical file does not overwrite traffic.")}</p>
        </div>
      )}
    </details>
  );
}
