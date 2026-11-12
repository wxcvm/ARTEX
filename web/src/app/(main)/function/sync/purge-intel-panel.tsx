"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api, http } from "@/lib/api";

type CompanyOption = { id: number; name: string };
type Snapshot = { id: number; digest: string; received_at: string };
type IntelAnalysis = {
  target: string; historical_only: boolean; truncated: boolean; excluded_dns: number;
  names: number; ips: number; edges: number; port_observations: number;
  source_counts: Record<string, number>;
  names_with_multiple_sources: string[]; names_without_dns: string[];
  names_with_dns_failures: string[]; names_without_retained_address: string[];
  ips_without_domain_link: string[]; ips_without_port_evidence: string[];
  unknown_transport_ports: number; ports_without_observation_time: number; shared_addresses: number;
};
type ImportResult = { id: number; names: number; ips: number; historical_only: boolean };

export function PurgeIntelPanel() {
  const [companies, setCompanies] = React.useState<CompanyOption[]>([]);
  const [company, setCompany] = React.useState("");
  const [raw, setRaw] = React.useState("");
  const [history, setHistory] = React.useState<Snapshot[]>([]);
  const [detail, setDetail] = React.useState<unknown>(null);
  const [analysis, setAnalysis] = React.useState<IntelAnalysis | null>(null);
  const [result, setResult] = React.useState<ImportResult | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  React.useEffect(() => {
    let live = true;
    api.companies().then(rows => { if (live) setCompanies(rows); })
      .catch(e => { if (live) setError(String(e.message || e)); });
    return () => { live = false; };
  }, []);
  const path = "/companies/" + company + "/purge-intel";
  const run = async (action: "import" | "history" | "detail" | "analysis", id?: number) => {
    if (!company || busy) return;
    setBusy(true); setError(""); setAnalysis(null);
    try {
      if (action === "import") {
        setResult(null); setDetail(null);
        const body = JSON.parse(raw);
        const imported = await http<ImportResult>(path, { method: "POST", body: JSON.stringify(body) });
        setResult(imported);
        setRaw("");
        setHistory((await http<{ items: Snapshot[] }>(path)).items);
      } else if (action === "history") {
        setHistory((await http<{ items: Snapshot[] }>(path)).items);
        setDetail(null);
      } else if (action === "analysis") {
        setDetail(null);
        setAnalysis(await http<IntelAnalysis>(path + "/" + id + "/analysis"));
      } else {
        setDetail(await http<unknown>(path + "/" + id));
      }
    } catch (e) { setError(e instanceof Error ? e.message : "操作失败"); }
    finally { setBusy(false); }
  };
  return (
    <Card>
      <CardHeader><CardTitle className="text-base">Purge 历史资产情报</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground">从 Purge 测绘中心复制 purge.asset-intel.v1 JSON，选择已登记授权范围的公司并导入。只保存历史证据，不创建实时资产、漏洞或扫描任务。</p>
        <label className="block space-y-1 text-sm">
          <span>目标公司</span>
          <select className="w-full rounded-md border bg-background p-2" value={company} disabled={busy}
            onChange={e => { setCompany(e.target.value); setHistory([]); setDetail(null); setAnalysis(null); setResult(null); setError(""); }}>
            <option value="">请选择公司</option>
            {companies.map(c => <option key={c.id} value={String(c.id)}>{c.name}</option>)}
          </select>
        </label>
        <label className="block space-y-1 text-sm">
          <span>Purge JSON 清单</span>
          <textarea className="min-h-36 w-full rounded-md border bg-background p-3 font-mono text-xs"
            value={raw} disabled={busy} onChange={e => setRaw(e.target.value)} placeholder="粘贴离线关联清单 JSON" />
        </label>
        <div className="flex flex-wrap gap-2">
          <Button disabled={busy || !company || !raw.trim()} onClick={() => void run("import")}>导入历史情报</Button>
          <Button variant="outline" disabled={busy || !company} onClick={() => void run("history")}>读取导入历史</Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {result && <p role="status" className="text-sm">已保存快照 #{result.id}：{result.names} 个域名，{result.ips} 个 IP（历史情报）。重复证据复用原快照。</p>}
        {history.length > 0 && <div className="space-y-2">
          <h3 className="text-sm font-medium">最近导入</h3>
          {history.map(row => <div key={row.id} className="flex flex-wrap items-center justify-between gap-2 rounded-md border p-2 text-sm">
            <span>#{row.id} · {row.received_at}</span>
            <div className="flex gap-2"><Button variant="outline" size="sm" disabled={busy} onClick={() => void run("detail", row.id)}>查看证据</Button>
            <Button size="sm" disabled={busy} onClick={() => void run("analysis", row.id)}>分析情报</Button></div>
          </div>)}
        </div>}
        {analysis && <IntelAnalysisView analysis={analysis} />}
        {detail !== null && <pre className="max-h-96 overflow-auto whitespace-pre-wrap break-all rounded-md border p-3 text-xs">{JSON.stringify(detail, null, 2)}</pre>}
      </CardContent>
    </Card>
  );
}

function IntelAnalysisView({ analysis: a }: { analysis: IntelAnalysis }) {
  const metrics = [
    ["域名", a.names], ["IP", a.ips], ["DNS 关系", a.edges], ["历史端口", a.port_observations],
    ["共享 IP", a.shared_addresses], ["未知传输协议", a.unknown_transport_ports],
    ["无原始观察时间的端口", a.ports_without_observation_time], ["范围外 DNS 排除", a.excluded_dns],
  ] as const;
  const groups = [
    ["多来源域名", a.names_with_multiple_sources],
    ["未保存 DNS 快照", a.names_without_dns],
    ["DNS 查询异常", a.names_with_dns_failures],
    ["无保留的 A / AAAA", a.names_without_retained_address],
    ["无域名关联的 IP", a.ips_without_domain_link],
    ["无历史端口证据的 IP", a.ips_without_port_evidence],
  ] as const;
  return <section className="space-y-3 rounded-lg border p-4">
    <div>
      <h3 className="font-medium">历史证据分析 · {a.target}</h3>
      <p className="text-xs text-muted-foreground">只分析保存的证据，不代表实时存活、资产完整性或漏洞确认。{a.truncated ? "该清单已截断。" : ""}</p>
    </div>
    <div className="grid grid-cols-2 gap-2 md:grid-cols-4">
      {metrics.map(([label, value]) => <div key={label} className="rounded-md bg-muted/50 p-3">
        <div className="text-xs text-muted-foreground">{label}</div>
        <div className="text-xl font-semibold tabular-nums">{value}</div>
      </div>)}
    </div>
    <div className="flex flex-wrap gap-2 text-xs">
      {Object.entries(a.source_counts).sort(([x], [y]) => x.localeCompare(y)).map(([source, count]) =>
        <span key={source} className="rounded-md border px-2 py-1">{source} · {count}</span>)}
    </div>
    {groups.map(([title, values]) => <details key={title} className="rounded-md border p-3 text-sm">
      <summary className="cursor-pointer">{title} · {values.length}</summary>
      <div className="mt-2 max-h-48 overflow-auto text-xs">{values.length ? values.map(value => <div key={value} className="break-all py-0.5">{value}</div>) : "暂无记录"}</div>
    </details>)}
    <details className="rounded-md border p-3 text-sm"><summary className="cursor-pointer">结构化分析 JSON（可复制）</summary>
      <textarea readOnly className="mt-2 min-h-32 w-full rounded-md border bg-background p-2 font-mono text-xs" value={JSON.stringify(a, null, 2)} />
    </details>
  </section>;
}
