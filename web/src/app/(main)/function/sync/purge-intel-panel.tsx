"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";


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
type Pairing = { source: "saved" | "environment"; configured: boolean; company_id: number; token?: string };

export function PurgeIntelPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [companies, setCompanies] = React.useState<CompanyOption[]>([]);
  const [company, setCompany] = React.useState("");
  const [raw, setRaw] = React.useState("");
  const [history, setHistory] = React.useState<Snapshot[]>([]);
  const [detail, setDetail] = React.useState<unknown>(null);
  const [analysis, setAnalysis] = React.useState<IntelAnalysis | null>(null);
  const [result, setResult] = React.useState<ImportResult | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  const [pair, setPair] = React.useState<Pairing | null>(null);
  const [pairToken, setPairToken] = React.useState("");
  const [pairBusy, setPairBusy] = React.useState(false);
  const [pairError, setPairError] = React.useState("");
  const [pairCopy, setPairCopy] = React.useState(false);
  const [pairConfirm, setPairConfirm] = React.useState(false);
  React.useEffect(() => {
    let live = true;
    api.companies().then(rows => { if (live) setCompanies(rows); })
      .catch(e => { if (live) setError(String(e.message || e)); });
    return () => { live = false; };
  }, []);
  React.useEffect(() => {
    let live = true;
    http<Pairing>("/purge/intel-pairing").then(p => {
      if (live) setPair(p);
    }).catch(e => { if (live) setPairError(String(e.message || e)); });
    return () => { live = false; };
  }, []);
  const updatePair = async (action: "create" | "revoke") => {
    if (pairBusy || pair?.source === "environment" || (action === "create" && !company)) return;
    if (pair?.configured && !pairConfirm) return;
    setPairBusy(true); setPairError(""); setPairCopy(false); setPairToken(""); setPairConfirm(false);
    try {
      const p = await http<Pairing>("/purge/intel-pairing", {
        method: action === "create" ? "POST" : "DELETE",
        ...(action === "create" ? { body: JSON.stringify({ company_id: Number(company) }) } : {}),
      });
      setPair({ source: p.source, configured: p.configured, company_id: p.company_id });
      setPairToken(p.token || "");
    } catch (e) { setPairError(e instanceof Error ? e.message : "Pairing failed"); }
    finally { setPairBusy(false); }
  };
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
      <CardHeader><CardTitle className="text-base">{tr("Purge 历史资产情报", "Purge historical asset intelligence")}</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground">{tr("从 Purge 测绘中心复制 purge.asset-intel.v1 JSON，选择已登记授权范围的公司并导入。只保存历史证据，不创建实时资产、漏洞或扫描任务。", "Import purge.asset-intel.v1 JSON from Purge into a company with registered authorized scope. Saves historical evidence only; creates no live assets, findings or scan tasks.")}</p>
        <label className="block space-y-1 text-sm">
          <span>{tr("目标公司", "Target company")}</span>
          <select className="w-full rounded-md border bg-background p-2" value={company} disabled={busy}
            onChange={e => { setCompany(e.target.value); setHistory([]); setDetail(null); setAnalysis(null); setResult(null); setError(""); setPairConfirm(false); }}>
            <option value="">{tr("请选择公司", "Select company")}</option>
            {companies.map(c => <option key={c.id} value={String(c.id)}>{c.name}</option>)}
          </select>
        </label>
        <section className="space-y-3 rounded-md border p-3">
          <h3 className="text-sm font-medium">{tr("Purge 本机配对", "Purge local pairing")}</h3>
          <p className="text-sm text-muted-foreground">{tr("选择上方目标公司，生成独立令牌。将公司 ID 和令牌复制到 Purge → ARTEX → 配对配置，保存并检测。令牌只显示一次，重置会使旧令牌失效。", "Select the target company and generate a dedicated token. Copy the company ID and token to Purge → ARTEX → Pairing, then save and check. The token is shown once; rotation invalidates the old token.")}</p>
          {pair && <p role="status" className="text-sm">{pair.configured ? tr(`已配置 · 公司 #${pair.company_id}`, `Configured · Company #${pair.company_id}`) : tr("尚未配对", "Not paired")} {pair.source === "environment" && tr("（宿主环境变量管理）", "(managed by host environment)")}</p>}
          {pair?.source === "environment" ? <p className="text-sm text-muted-foreground">{tr("请在宿主修改环境变量；页面不会覆盖原有配置。", "Modify the host environment; this page does not override it.")}</p> : <>
            {pair?.configured && <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={pairConfirm} disabled={pairBusy} onChange={e => setPairConfirm(e.target.checked)} />{tr("确认替换或撤销已有配对", "Confirm rotation or revocation of the existing pairing")}</label>}
            <div className="flex flex-wrap gap-2">
              <Button disabled={pairBusy || !pair || !company || (pair.configured && !pairConfirm)} onClick={() => void updatePair("create")}>{pair?.configured ? tr("重新生成令牌", "Rotate token") : tr("生成配对令牌", "Generate pairing token")}</Button>
              {pair?.configured && <Button variant="outline" disabled={pairBusy || !pairConfirm} onClick={() => void updatePair("revoke")}>{tr("撤销配对", "Revoke pairing")}</Button>}
            </div>
          </>}
          {pairToken && <div className="space-y-2">
            <p className="text-sm">{tr("公司 ID", "Company ID")}: {pair?.company_id}</p>
            <input aria-label={tr("新配对令牌", "New pairing token")} readOnly value={pairToken} className="w-full rounded-md border bg-background p-2 font-mono text-xs" />
            <div className="flex flex-wrap gap-2">
              <Button variant="outline" size="sm" onClick={async () => {
                try { await navigator.clipboard.writeText(pairToken); setPairCopy(true); }
                catch { setPairError(tr("无法自动复制，请手动复制令牌。", "Clipboard unavailable; copy the token manually.")); }
              }}>{pairCopy ? tr("已复制", "Copied") : tr("复制令牌", "Copy token")}</Button>
              <Button variant="outline" size="sm" onClick={() => { setPairToken(""); setPairCopy(false); }}>{tr("隐藏令牌", "Hide token")}</Button>
            </div>
          </div>}
          {pairError && <p role="alert" className="text-sm text-destructive">{pairError}</p>}
        </section>
        <label className="block space-y-1 text-sm">
          <span>{tr("Purge JSON 清单", "Purge JSON manifest")}</span>
          <textarea className="min-h-36 w-full rounded-md border bg-background p-3 font-mono text-xs"
            value={raw} disabled={busy} onChange={e => setRaw(e.target.value)} placeholder={tr("粘贴离线关联清单 JSON", "Paste offline association manifest JSON")} />
        </label>
        <div className="flex flex-wrap gap-2">
          <Button disabled={busy || !company || !raw.trim()} onClick={() => void run("import")}>{tr("导入历史情报", "Import historical intelligence")}</Button>
          <Button variant="outline" disabled={busy || !company} onClick={() => void run("history")}>{tr("读取导入历史", "Read import history")}</Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {result && <p role="status" className="text-sm">{tr(`已保存快照 #${result.id}：${result.names} 个域名，${result.ips} 个 IP（历史情报）。重复证据复用原快照。`, `Saved snapshot #${result.id}: ${result.names} domains, ${result.ips} IPs (historical). Duplicate evidence reuses the snapshot.`)}</p>}
        {history.length > 0 && <div className="space-y-2">
          <h3 className="text-sm font-medium">{tr("最近导入", "Recent imports")}</h3>
          {history.map(row => <div key={row.id} className="flex flex-wrap items-center justify-between gap-2 rounded-md border p-2 text-sm">
            <span>#{row.id} · {row.received_at}</span>
            <div className="flex gap-2"><Button variant="outline" size="sm" disabled={busy} onClick={() => void run("detail", row.id)}>{tr("查看证据", "View evidence")}</Button>
            <Button size="sm" disabled={busy} onClick={() => void run("analysis", row.id)}>{tr("分析情报", "Analyze intelligence")}</Button></div>
          </div>)}
        </div>}
        {analysis && <IntelAnalysisView analysis={analysis} />}
        {detail !== null && <pre className="max-h-96 overflow-auto whitespace-pre-wrap break-all rounded-md border p-3 text-xs">{JSON.stringify(detail, null, 2)}</pre>}
      </CardContent>
    </Card>
  );
}

function IntelAnalysisView({ analysis: a }: { analysis: IntelAnalysis }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const metrics = [
    [tr("域名", "Domains"), a.names], ["IP", a.ips], [tr("DNS 关系", "DNS relations"), a.edges], [tr("历史端口", "Historical ports"), a.port_observations],
    [tr("共享 IP", "Shared IPs"), a.shared_addresses], [tr("未知传输协议", "Unknown transport"), a.unknown_transport_ports],
    [tr("无原始观察时间的端口", "Ports without observation time"), a.ports_without_observation_time], [tr("范围外 DNS 排除", "Out-of-scope DNS excluded"), a.excluded_dns],
  ] as const;
  const groups = [
    [tr("多来源域名", "Multi-source domains"), a.names_with_multiple_sources],
    [tr("未保存 DNS 快照", "No DNS snapshot"), a.names_without_dns],
    [tr("DNS 查询异常", "DNS failures"), a.names_with_dns_failures],
    [tr("无保留的 A / AAAA", "No retained A/AAAA"), a.names_without_retained_address],
    [tr("无域名关联的 IP", "IPs without domain link"), a.ips_without_domain_link],
    [tr("无历史端口证据的 IP", "IPs without historical port evidence"), a.ips_without_port_evidence],
  ] as const;
  return <section className="space-y-3 rounded-lg border p-4">
    <div>
      <h3 className="font-medium">{tr("历史证据分析", "Historical evidence analysis")} · {a.target}</h3>
      <p className="text-xs text-muted-foreground">{tr("只分析保存的证据，不代表实时存活、资产完整性或漏洞确认。", "Analyzes stored evidence only; does not establish live availability, asset completeness or confirmed findings.")}{a.truncated ? tr("该清单已截断。", "Manifest truncated.") : ""}</p>
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
      <div className="mt-2 max-h-48 overflow-auto text-xs">{values.length ? values.map(value => <div key={value} className="break-all py-0.5">{value}</div>)  : tr("暂无记录", "No records")}</div>
    </details>)}
    <details className="rounded-md border p-3 text-sm"><summary className="cursor-pointer">{tr("结构化分析 JSON（可复制）", "Structured analysis JSON (copyable)")}</summary>
      <textarea readOnly className="mt-2 min-h-32 w-full rounded-md border bg-background p-2 font-mono text-xs" value={JSON.stringify(a, null, 2)} />
    </details>
  </section>;
}
