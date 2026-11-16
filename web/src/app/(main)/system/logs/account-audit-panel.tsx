"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { http } from "@/lib/api";
import { MOCK } from "@/lib/mock/enabled";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

type Event = { id: number; created_at: string; action: string; result: string };
type Page = { items: Event[]; has_more: boolean };

export function AccountAuditPanel() {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  const [open, setOpen] = React.useState(false);
  const [action, setAction] = React.useState("");
  const [result, setResult] = React.useState("");
  const [before, setBefore] = React.useState(0);
  const [page, setPage] = React.useState<Page>({ items: [], has_more: false });
  const [revision, setRevision] = React.useState(0);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  React.useEffect(() => {
    if (!open || MOCK) return;
    let active = true;
    setBusy(true); setError("");
    const q = new URLSearchParams({ action, result, before: String(before), limit: "50" });
    http<Page>("/logs/account-audit?" + q).then(data => { if (active) setPage(data); })
      .catch(e => { if (active) { setError((e as Error).message); setPage({items:[],has_more:false}); } })
      .finally(() => { if (active) setBusy(false); });
    return () => { active = false; };
  }, [open, action, result, before, revision]);
  const actionLabel = (value: string) => value === "login" ? t("登录", "Sign in") : value === "initialize_password" ? t("初始化密码", "Initialize password") : t("修改密码", "Change password");
  const resultLabel = (value: string) => value === "success" ? t("成功", "Success") : value === "failure" ? t("失败", "Failure") : t("错误", "Error");
  return <section className="rounded-xl border bg-card p-4 space-y-3">
    <div className="flex flex-wrap items-center gap-2">
      <h2 className="flex-1 font-semibold">{t("账号操作审计", "Account activity audit")}</h2>
      <Button size="sm" variant="outline" onClick={() => setOpen(v => !v)}>{open ? t("收起", "Collapse") : t("查看审计", "View audit")}</Button>
    </div>
    <p className="text-xs text-muted-foreground">{t("记录登录、密码初始化及修改结果，不保存密码、令牌或请求正文。", "Records sign-in, password initialization and change outcomes without passwords, tokens or request bodies.")}</p>
    {open && (MOCK ? <p className="text-sm text-muted-foreground">{t("演示模式无真实审计记录。", "No real account audit records in demo mode.")}</p> : <>
      <div className="flex flex-wrap gap-2">
        <select aria-label={t("操作类型", "Action")} className="h-9 rounded-md border bg-background px-2 text-sm" value={action} onChange={e => {setAction(e.target.value);setBefore(0);}}>
          <option value="">{t("全部操作", "All actions")}</option>
          {["login","initialize_password","change_password"].map(v => <option key={v} value={v}>{actionLabel(v)}</option>)}
        </select>
        <select aria-label={t("操作结果", "Result")} className="h-9 rounded-md border bg-background px-2 text-sm" value={result} onChange={e => {setResult(e.target.value);setBefore(0);}}>
          <option value="">{t("全部结果", "All results")}</option>
          {["success","failure","error"].map(v => <option key={v} value={v}>{resultLabel(v)}</option>)}
        </select>
        <Button size="sm" variant="outline" disabled={busy} onClick={() => {setBefore(0);setRevision(v=>v+1);}}>{t("刷新", "Refresh")}</Button>
      </div>
      {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
      {busy ? <p className="text-sm text-muted-foreground">{t("加载中…", "Loading…")}</p> : !page.items.length ? <p className="text-sm text-muted-foreground">{t("暂无记录", "No events")}</p> : <div className="divide-y">
        {page.items.map(e => <div key={e.id} className="flex flex-wrap items-center gap-3 py-2 text-sm"><span className="text-muted-foreground tabular-nums">{new Date(e.created_at).toLocaleString()}</span><span className="flex-1">{actionLabel(e.action)}</span><span className={e.result === "success" ? "text-emerald-600" : "text-destructive"}>{resultLabel(e.result)}</span></div>)}
      </div>}
      <Button size="sm" variant="ghost" disabled={busy || !page.has_more || !page.items.length} onClick={() => setBefore(page.items[page.items.length-1].id)}>{t("更早记录", "Earlier events")}</Button>
    </>)}
  </section>;
}
