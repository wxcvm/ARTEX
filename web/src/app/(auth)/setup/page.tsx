"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { OnboardingGuide } from "@/components/onboarding-guide";
import { api } from "@/lib/api";
import { auth } from "@/lib/auth";

export default function SetupPage() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const router = useRouter();
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [setupToken, setSetupToken] = useState("");
  const [automatic, setAutomatic] = useState(false);
  const [manual, setManual] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [checking, setChecking] = useState(true);
  const [unavailable, setUnavailable] = useState(false);

  useEffect(() => {
    let alive = true;
    const receive = (event: Event) => {
      const value: unknown = (event as CustomEvent).detail;
      if (alive && typeof value === "string" && /^[a-fA-F0-9]{48}$/.test(value)) { setSetupToken(value); setAutomatic(true); }
    };
    window.addEventListener("artex:setup-token", receive);
    api.authStatus().then(({ initialized }) => {
      if (!alive) return;
      if (initialized) { router.replace("/login"); return; }
      const desktop = (window as unknown as { chrome?: { webview?: { postMessage: (message: string) => void } } }).chrome?.webview;
      if (desktop) desktop.postMessage("artex:request-setup-token");
      else setManual(true);
    }).catch(err => { if (alive) { setError(err instanceof Error ? err.message : "无法连接后端服务"); setUnavailable(true); } })
      .finally(() => { if (alive) setChecking(false); });
    return () => { alive = false; window.removeEventListener("artex:setup-token", receive); };
  }, [router]);

  async function submit(event: React.FormEvent) {
    event.preventDefault();
    if (loading) return;
    if (password !== confirm) { setError(tr("两次密码不一致，请重新核对。", "Passwords do not match.")); return; }
    if (password.length < 8) { setError(tr("登录密码至少需要 8 位。", "Login password must be at least 8 characters.")); return; }
    setLoading(true); setError("");
    try { const { token } = await api.initPassword(password, setupToken.trim()); auth.setToken(token); setSetupToken(""); router.replace("/welcome"); }
    catch (err) { setError(err instanceof Error ? err.message : tr("设置未完成，请重试。", "Setup incomplete; retry.")); }
    finally { setLoading(false); }
  }

  if (checking) return <main className="p-8" role="status">{tr("正在确认账号初始化状态…", "Checking account initialization…")}</main>;
  return <main className="mx-auto min-h-dvh w-full max-w-xl space-y-6 px-4 py-10 sm:px-8">
    <header className="space-y-2"><p className="text-sm text-muted-foreground">{tr("准备数据库 ✓ → 设置登录密码 → 连接 AI → 开始使用", "Prepare database ✓ → Set login password → Connect AI → Get started")}</p>
      <h1 className="text-2xl font-semibold">{unavailable ? tr("暂时无法确认账号状态", "Unable to confirm account status") : tr("第 2 步：设置 ARTEX 登录密码", "Step 2: Set ARTEX login password")}</h1>
      <p className="text-sm text-muted-foreground">{unavailable ? tr("后端或数据库暂时不可用。为保留已有账号，此时不开放初始化。恢复服务后重试。", "Backend or database unavailable. Initialization is disabled to preserve existing accounts. Retry after service recovery.") : tr("这是以后进入 ARTEX 使用的密码，不是数据库密码或 AI API Key。", "This is your ARTEX login password, not the database password or AI API key.")}</p></header>
    {unavailable ? <Button onClick={() => window.location.reload()}>{tr("重新检查服务", "Recheck service")}</Button> : <form onSubmit={submit} className="space-y-4">
      {automatic && <p role="status" className="rounded-lg border p-3 text-sm">{tr("桌面已自动取得一次性初始化令牌，无需查日志或复制。", "Desktop obtained the one-time setup token automatically; no log lookup or copying needed.")}</p>}
      {!automatic && <div className="space-y-3 rounded-lg border p-3"><p className="text-sm">{tr("桌面自启实例会自动取得令牌；外部实例或浏览器请使用手动入口。", "Desktop-managed instances obtain the token automatically; external instances or browsers use the manual entry.")}</p><Button type="button" variant="outline" onClick={() => setManual(!manual)}>{tr("手动提供初始化令牌", "Enter setup token manually")}</Button></div>}
      {manual && <div className="space-y-2"><Label htmlFor="setup-token">{tr("一次性初始化令牌", "One-time setup token")}</Label><Input id="setup-token" type="password" value={setupToken} onChange={e => { setSetupToken(e.target.value); setAutomatic(false); }} autoComplete="off" /><p className="text-xs text-muted-foreground">{tr("从当前后端启动日志的 first-time setup token 行获取。具体位置见下方详细指南。不要填写数据库密码。", "Get the first-time setup token from current backend startup logs. See the guide below. Do not enter the database password.")}</p></div>}
      <div className="space-y-2"><Label htmlFor="password">{tr("登录密码（至少 8 位）", "Login password (at least 8 characters)")}</Label><Input id="password" type="password" autoComplete="new-password" value={password} onChange={e => setPassword(e.target.value)} required minLength={8} /></div>
      <div className="space-y-2"><Label htmlFor="confirm">{tr("再次填写登录密码", "Confirm login password")}</Label><Input id="confirm" type="password" autoComplete="new-password" value={confirm} onChange={e => setConfirm(e.target.value)} required minLength={8} /></div>
      <Button type="submit" disabled={loading || !password || !confirm || !setupToken.trim()}>{loading ? tr("正在设置…", "Setting up…") : tr("设置密码并继续", "Set password and continue")}</Button>
    </form>}
    {error && <p role="alert" className="break-words text-sm text-destructive">{error}</p>}
    <OnboardingGuide step="account" />
  </main>;
}
