"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api, http } from "@/lib/api";

type CompanyOption = { id: number; name: string };
type Snapshot = { id: number; digest: string; received_at: string };
type ImportResult = { id: number; names: number; ips: number; historical_only: boolean };

export function PurgeIntelPanel() {
  const [companies, setCompanies] = React.useState<CompanyOption[]>([]);
  const [company, setCompany] = React.useState("");
  const [raw, setRaw] = React.useState("");
  const [history, setHistory] = React.useState<Snapshot[]>([]);
  const [detail, setDetail] = React.useState<unknown>(null);
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
  const run = async (action: "import" | "history" | "detail", id?: number) => {
    if (!company || busy) return;
    setBusy(true); setError("");
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
            onChange={e => { setCompany(e.target.value); setHistory([]); setDetail(null); setResult(null); setError(""); }}>
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
            <Button variant="outline" size="sm" disabled={busy} onClick={() => void run("detail", row.id)}>查看证据</Button>
          </div>)}
        </div>}
        {detail !== null && <pre className="max-h-96 overflow-auto whitespace-pre-wrap break-all rounded-md border p-3 text-xs">{JSON.stringify(detail, null, 2)}</pre>}
      </CardContent>
    </Card>
  );
}
