"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { http } from "@/lib/api";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

type NameChanges = {
  available: boolean; added: string[]; removed: string[]; changed: string[];
  added_count: number; removed_count: number; changed_count: number;
};
type BodyComparison = {
  baseline_bytes: number; candidate_bytes: number; relation: string; json_fields: NameChanges;
  form_fields?: NameChanges; xml_children?: NameChanges;
};
type RecordedMeta = { id: string; method: string; ts: number; status: number; content_type: string };
type Comparison = {
  baseline: RecordedMeta; candidate: RecordedMeta; same_url: boolean; status_relation: string;
  request_headers: NameChanges; response_headers: NameChanges;
  request_body: BodyComparison; response_body: BodyComparison; evidence_version?: number;
};

export function TrafficCompareSheet({ ids, onClose, findingId, contextTask, evidenceVersion }: { ids: [string, string] | null; onClose: () => void; findingId?: string; contextTask?: string; evidenceVersion?: number }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [data, setData] = React.useState<Comparison | null>(null);
  const [failed, setFailed] = React.useState(false);
  const [swap, setSwap] = React.useState(false);
  const a = ids ? ids[swap ? 1 : 0] : "";
  const b = ids ? ids[swap ? 0 : 1] : "";

  React.useEffect(() => {
    if (!a || !b) return;
    let alive = true;
    const controller = new AbortController();
    setData(null);
    setFailed(false);
    const query = new URLSearchParams({ baseline: a, candidate: b });
    if (contextTask) query.set("context_task", contextTask);
    if (evidenceVersion !== undefined) query.set("version", String(evidenceVersion));
    const endpoint = findingId ? "/exploration/findings/" + encodeURIComponent(findingId) + "/traffic/compare?" : "/traffic/compare?";
    void http<Comparison>(endpoint + query, { signal: controller.signal })
      .then((result) => { if (alive) setData(result); })
      .catch(() => { if (alive) setFailed(true); });
    return () => { alive = false; controller.abort(); };
  }, [a, b, findingId, contextTask, evidenceVersion]);

  const relation = (value: string) => value === "equal" ? tr("相同", "Equal") : value === "different" ? tr("不同", "Different") : tr("未知", "Unknown");
  const changes = (title: string, c: NameChanges) => (
    <section className="rounded-md border p-3">
      <h3 className="mb-2 text-sm font-semibold">{title}</h3>
      {!c.available ? <p className="text-xs text-muted-foreground">{tr("无法比较", "Comparison unavailable")}</p> : (
        <dl className="space-y-2 text-xs">
          {([
            [tr("新增", "Added"), c.added_count, c.added],
            [tr("移除", "Removed"), c.removed_count, c.removed],
            [tr("变更", "Changed"), c.changed_count, c.changed],
          ] as const).map(([label, count, names]) => (
            <div key={label}>
              <dt className="font-medium">{label}: {count}</dt>
              <dd className="break-all text-muted-foreground">{names.join(", ") || "—"}{count > names.length ? " …" : ""}</dd>
            </div>
          ))}
        </dl>
      )}
    </section>
  );
  const body = (title: string, c: BodyComparison) => (
    <section className="space-y-2 rounded-md border p-3">
      <h3 className="text-sm font-semibold">{title}</h3>
      <p className="text-xs">{tr("存储字节", "Stored bytes")}: {c.baseline_bytes} → {c.candidate_bytes} · {relation(c.relation)}</p>
      {changes(tr("JSON 顶层字段", "Top-level JSON fields"), c.json_fields)}
      {c.form_fields?.available && changes(tr("表单字段", "Form fields"), c.form_fields)}
      {c.xml_children?.available && changes(tr("XML 直接子元素", "Direct XML children"), c.xml_children)}
    </section>
  );
  return (
    <Sheet open={ids !== null} onOpenChange={(open) => { if (!open) { setSwap(false); setData(null); setFailed(false); onClose(); } }}>
      <SheetContent className="w-full! max-w-none! gap-0 p-0 sm:w-[48rem]! sm:max-w-[48rem]!">
        <SheetHeader className="border-b px-5 py-4">
          <SheetTitle>{findingId ? tr("已保存证据对比", "Saved evidence comparison") : tr("录制流量对比", "Recorded traffic comparison")}</SheetTitle>
          <SheetDescription>{tr("比较已保存的报文，不发送新请求。", "Compare stored messages without sending new requests.")}</SheetDescription>
        </SheetHeader>
        <div className="min-h-0 flex-1 space-y-3 overflow-auto p-4">
          <Button size="sm" variant="outline" disabled={!ids} onClick={() => setSwap((v) => !v)}>{tr("交换基准与对照", "Swap baseline and candidate")}</Button>
          {failed ? <p role="alert" className="text-sm text-destructive">{tr("无法比较所选记录，记录可能已删除、证据版本已变化或存储不可用。", "Unable to compare selected records; a record may be deleted, the evidence version changed, or storage unavailable.")}</p> :
            !data ? <p role="status" className="text-sm text-muted-foreground">{tr("正在比较录制数据…", "Comparing recorded data…")}</p> : (
              <>
                {data.evidence_version !== undefined && <p className="text-xs text-muted-foreground">{tr("证据版本", "Evidence version")}: {data.evidence_version}</p>}
                <div className="grid gap-2 sm:grid-cols-2">
                  {([[tr("基准", "Baseline"), data.baseline], [tr("对照", "Candidate"), data.candidate]] as const).map(([label, m]) => (
                    <div className="min-w-0 space-y-1 rounded-md border p-3 text-xs" key={label}>
                      <p className="font-semibold">{label}</p>
                      <p className="break-all font-mono">{m.id}</p>
                      <p>{m.method} · {m.status || tr("未观测", "Unobserved")} · {new Date(m.ts * 1000).toLocaleString(locale)}</p>
                      <p className="break-all text-muted-foreground">{m.content_type || "—"}</p>
                    </div>
                  ))}
                </div>
                <p className="text-sm">{tr("URL", "URL")}: {relation(data.same_url ? "equal" : "different")} · {tr("HTTP 状态", "HTTP status")}: {relation(data.status_relation)}</p>
                {changes(tr("请求头", "Request headers"), data.request_headers)}
                {changes(tr("响应头", "Response headers"), data.response_headers)}
                {body(tr("请求正文", "Request body"), data.request_body)}
                {body(tr("响应正文", "Response body"), data.response_body)}
                <p className="text-xs text-muted-foreground">{tr("仅比较已存储的字节；未知或截断的原始录制不等于完整网络响应。JSON 对象、表单和 XML 结构分析限 256 KiB；字段名最多预览 50 个，不显示值。差异不是漏洞结论。", "Compares stored bytes only; incomplete original captures may not represent full network responses. JSON object, form and XML structure analysis is limited to 256 KiB and previews 50 names per group, without values. Differences are not vulnerability conclusions.")}</p>
              </>
            )}
        </div>
      </SheetContent>
    </Sheet>
  );
}
