"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { OnboardingGuide } from "@/components/onboarding-guide";
import { api } from "@/lib/api";
import { auth } from "@/lib/auth";

export default function SetupPage() {
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
    if (password !== confirm) { setError("两次密码不一致，请重新核对。"); return; }
    if (password.length < 8) { setError("登录密码至少需要 8 位。"); return; }
    setLoading(true); setError("");
    try { const { token } = await api.initPassword(password, setupToken.trim()); auth.setToken(token); setSetupToken(""); router.replace("/welcome"); }
    catch (err) { setError(err instanceof Error ? err.message : "设置未完成，请重试。"); }
    finally { setLoading(false); }
  }

  if (checking) return <main className="p-8" role="status">正在确认账号初始化状态…</main>;
  return <main className="mx-auto min-h-dvh w-full max-w-xl space-y-6 px-4 py-10 sm:px-8">
    <header className="space-y-2"><p className="text-sm text-muted-foreground">准备数据库 ✓ → 设置登录密码 → 连接 AI → 开始使用</p>
      <h1 className="text-2xl font-semibold">{unavailable ? "暂时无法确认账号状态" : "第 2 步：设置 ARTEX 登录密码"}</h1>
      <p className="text-sm text-muted-foreground">{unavailable ? "后端或数据库暂时不可用。为保留已有账号，此时不开放初始化。恢复服务后重试。" : "这是以后进入 ARTEX 使用的密码，不是数据库密码或 AI API Key。"}</p></header>
    {unavailable ? <Button onClick={() => window.location.reload()}>重新检查服务</Button> : <form onSubmit={submit} className="space-y-4">
      {automatic && <p role="status" className="rounded-lg border p-3 text-sm">桌面已自动取得一次性初始化令牌，无需查日志或复制。</p>}
      {!automatic && <div className="space-y-3 rounded-lg border p-3"><p className="text-sm">桌面自启实例会自动取得令牌；外部实例或浏览器请使用手动入口。</p><Button type="button" variant="outline" onClick={() => setManual(!manual)}>手动提供初始化令牌</Button></div>}
      {manual && <div className="space-y-2"><Label htmlFor="setup-token">一次性初始化令牌</Label><Input id="setup-token" type="password" value={setupToken} onChange={e => { setSetupToken(e.target.value); setAutomatic(false); }} autoComplete="off" /><p className="text-xs text-muted-foreground">从当前后端启动日志的 first-time setup token 行获取。具体位置见下方详细指南。不要填写数据库密码。</p></div>}
      <div className="space-y-2"><Label htmlFor="password">登录密码（至少 8 位）</Label><Input id="password" type="password" autoComplete="new-password" value={password} onChange={e => setPassword(e.target.value)} required minLength={8} /></div>
      <div className="space-y-2"><Label htmlFor="confirm">再次填写登录密码</Label><Input id="confirm" type="password" autoComplete="new-password" value={confirm} onChange={e => setConfirm(e.target.value)} required minLength={8} /></div>
      <Button type="submit" disabled={loading || !password || !confirm || !setupToken.trim()}>{loading ? "正在设置…" : "设置密码并继续"}</Button>
    </form>}
    {error && <p role="alert" className="break-words text-sm text-destructive">{error}</p>}
    <OnboardingGuide step="account" />
  </main>;
}
