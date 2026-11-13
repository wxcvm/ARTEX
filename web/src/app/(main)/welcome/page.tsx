"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { OnboardingGuide } from "@/components/onboarding-guide";
import { api } from "@/lib/api";
import { modelErrorGuide, onboardingSteps, testAndEnableModel } from "@/lib/onboarding";
import type { LLMProfile } from "@/lib/types";

export default function WelcomePage() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [profiles, setProfiles] = useState<LLMProfile[]>([]);
  const [configured, setConfigured] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selected, setSelected] = useState("");
  const [format, setFormat] = useState("openai");
  const [base, setBase] = useState("");
  const [model, setModel] = useState("");
  const [key, setKey] = useState("");
  const [proxy, setProxy] = useState("");
  const [busy, setBusy] = useState(false);
  const [result, setResult] = useState<{ reply?: string; latency_ms?: number } | null>(null);
  const [done, setDone] = useState(false);
  const saved = useRef<{ fingerprint: string; id: string } | null>(null);
  const active = profiles.find(p => p.is_default);
  const existing = profiles.find(p => p.id === selected);
  // Credentials stay only in this mounted form and the existing API, never in progress storage.
  const fingerprint = JSON.stringify([format, base, model, key, proxy]);

  async function refresh() {
    setLoading(true); setError("");
    try { const [profiles, state] = await Promise.all([api.llmProfiles(), api.getLLM()]); setProfiles(profiles); setConfigured(state.configured); }
    catch (e) { setConfigured(false); setError((e as Error).message); }
    finally { setLoading(false); }
  }
  useEffect(() => { void refresh(); }, []);
  useEffect(() => { setResult(null); setDone(false); }, [fingerprint, selected]);

  async function connect() {
    if (busy) return;
    setBusy(true); setError(""); setResult(null); setDone(false);
    try {
      const p = existing;
      const r = await testAndEnableModel({
        test: () => api.testLLM(p?.format ?? format, p?.model ?? model.trim(), p?.base_url ?? base.trim(), p ? "" : key,
          p?.proxy ?? proxy.trim(), p?.thinking_type ?? "", p?.reasoning_effort ?? "", p ? Number(p.id) : undefined,
          p?.streaming ?? true, p?.session_header_key ?? "", p?.user_agent ?? "").then(r => { if (r.ok && r.reply?.trim()) setResult(r); return r; }),
        save: async () => {
          if (p) return p.id;
          const response = await api.saveLLMProfile({ ...(saved.current?.fingerprint === fingerprint ? { id: Number(saved.current.id) } : {}),
            name: `${format} · ${model.trim()}`, format, base_url: base.trim(), model: model.trim(), api_key: key, proxy: proxy.trim(), streaming: true });
          const id = String(response.id); saved.current = { fingerprint, id }; return id;
        },
        activate: async (id) => { await api.activateLLMProfile(id); },
      });
      setResult(r); setDone(true); setProfiles(await api.llmProfiles()); setConfigured((await api.getLLM()).configured);
    } catch (e) { setError((e as Error).message); }
    finally { setBusy(false); }
  }

  return <main className="mx-auto w-full max-w-3xl space-y-6 p-4 sm:p-8">
    <header><p className="text-sm text-muted-foreground">{tr("首次使用 · 可随时从菜单重新打开", "First use · reopen from the menu anytime")}</p><h1 className="mt-1 text-2xl font-semibold">{tr("欢迎使用 ARTEX", "Welcome to ARTEX")}</h1>
      <p className="mt-2 text-sm text-muted-foreground">{tr("按四步完成准备，不需要先理解所有参数。", "Complete four steps without needing to understand every parameter first.")}</p></header>
    <ol aria-label={tr("配置进度", "Setup progress")} className="grid grid-cols-2 gap-3 sm:grid-cols-4">{onboardingSteps.map((s, i) => <li key={s} className="rounded-lg border p-3 text-sm">
      <span className="text-muted-foreground">{i + 1}/4</span><p className="mt-1 font-medium">{s}</p><p className="mt-1 text-xs text-muted-foreground">{i < 2 ? tr("已连接并登录", "Connected and signed in") : i === 2 ? done ? tr("已验证并启用", "Verified and enabled") : configured ? tr("已配置 · 可重测", "Configured · retest available") : tr("待配置", "Pending configuration") : tr("手动开始", "Start manually")}</p></li>)}</ol>
    <details className="rounded-lg border p-4"><summary className="cursor-pointer text-sm">{tr("查看已完成的数据库与登录步骤", "View completed database and login steps")}</summary><div className="mt-3 space-y-3"><OnboardingGuide step="database" /><OnboardingGuide step="account" /></div></details>
    <section className="space-y-4 rounded-xl border p-4 sm:p-6"><h2 className="text-lg font-semibold">{tr("第 3 步：连接 AI", "Step 3: Connect AI")}</h2>
      {loading ? <p role="status">{tr("正在读取实际模型配置…", "Reading actual model configuration…")}</p> : <>
        {configured && <div className="rounded-lg bg-muted p-3 text-sm">{tr("当前模型：", "Current model: ")}{active?.name ?? tr("环境配置", "Environment configuration")}（{active?.model ?? tr("已启用", "Enabled")}）。{done ? tr("本次测试已通过。", "Test passed this session.") : tr("已配置，尚未在本次引导重测；不会自动发送请求。", "Configured, not retested in this session; no automatic requests.")}
          <Link className="mt-3 block text-primary underline" href="/function/tasks">{tr("保留当前配置，继续使用（不代表本次已测试）", "Keep current configuration and continue (not tested in this session)")}</Link></div>}
        <fieldset disabled={busy} className="space-y-4">
          {profiles.length > 0 && <div className="space-y-2"><Label htmlFor="existing-model">{tr("选择已有配置或新增", "Select existing or add a profile")}</Label><select id="existing-model" className="w-full rounded-md border bg-background p-2" value={selected} onChange={e => setSelected(e.target.value)}><option value="">{tr("新增模型配置", "New model profile")}</option>{profiles.map(p => <option key={p.id} value={p.id}>{p.name} · {p.model}</option>)}</select></div>}
          {!existing && <>
            <div className="space-y-2"><Label htmlFor="model-format">{tr("协议", "Protocol")}</Label><select id="model-format" className="w-full rounded-md border bg-background p-2" value={format} onChange={e => setFormat(e.target.value)}><option value="openai">{tr("OpenAI 兼容（Chat Completions）", "OpenAI compatible (Chat Completions)")}</option><option value="openai-responses">OpenAI Responses</option><option value="anthropic">Anthropic</option><option value="ollama">{tr("Ollama 本机模型", "Local Ollama model")}</option></select><p className="text-xs text-muted-foreground">{tr("按提供方文档选择；“OpenAI 兼容”不代表只能用 OpenAI 模型。", "Follow provider documentation; OpenAI-compatible does not mean OpenAI models only.")}</p></div>
            <div className="space-y-2"><Label htmlFor="model-base">接口地址</Label><Input id="model-base" value={base} onChange={e => setBase(e.target.value)} placeholder={format === "ollama" ? "示例：http://127.0.0.1:11434" : "示例：https://api.example.com/v1"} /><p className="text-xs text-muted-foreground">填写服务提供方给出的 API 地址，不是网页聊天地址；示例不能直接连接。</p></div>
            <div className="space-y-2"><Label htmlFor="model-name">{tr("模型名", "Model name")}</Label><Input id="model-name" value={model} onChange={e => setModel(e.target.value)} placeholder={tr("填写提供方的准确模型 ID", "Enter the provider's exact model ID")} /></div>
            <div className="space-y-2"><Label htmlFor="model-key">API Key{format === "ollama" ? tr("（本机服务通常可留空）", " (usually optional for local service)") : ""}</Label><Input id="model-key" type="password" autoComplete="off" value={key} onChange={e => setKey(e.target.value)} /><p className="text-xs text-muted-foreground">{tr("从模型平台获取，不是 ARTEX 或数据库密码；不会写入引导进度。", "Obtain from the model provider, not ARTEX or the database. It is not stored in onboarding progress.")}</p></div>
            <details className="rounded-lg border p-3"><summary className="cursor-pointer text-sm">{tr("高级设置（可选）", "Advanced settings (optional)")}</summary><div className="mt-3 space-y-2"><Label htmlFor="model-proxy">{tr("代理地址", "Proxy URL")}</Label><Input id="model-proxy" value={proxy} onChange={e => setProxy(e.target.value)} placeholder={tr("不需要代理时留空", "Leave empty when no proxy is needed")} /><p className="text-xs text-muted-foreground">{tr("限速、上下文和推理参数沿用默认值；保存后可在完整模型设置调整。", "Rate limits, context and reasoning use defaults; adjust in full model settings after saving.")}</p></div></details>
          </>}
          <Button onClick={() => void connect()} disabled={busy || (!existing && (!base.trim() || !model.trim() || (format !== "ollama" && !key.trim())))}>{busy ? tr("正在测试、保存并启用…", "Testing, saving and enabling…") : tr("测试并启用", "Test and enable")}</Button>
          <p className="text-xs text-muted-foreground">{tr("此按钮会发送一次真实模型请求。只有模型返回有效回复后，才保存并启用配置。", "This sends a real model request. The profile is saved and enabled only after a valid response.")}</p>
        </fieldset>
      </>}
      {error && <div role="alert" className="space-y-2 rounded-lg border border-destructive p-3 text-sm"><p className="break-words">{error}</p><p>{modelErrorGuide(error)}</p><Button variant="outline" disabled={busy} onClick={() => void refresh()}>{tr("重新读取配置", "Reload configuration")}</Button></div>}
      {result && <div role="status" className="rounded-lg bg-muted p-3 text-sm"><p>{tr("模型实际回复：", "Actual model response: ")}{result.reply}</p><p>{tr("耗时：", "Latency: ")}{result.latency_ms ?? tr("未知", "Unknown")} ms · {done ? tr("已保存并启用", "Saved and enabled") : tr("模型测试已通过，保存及启用状态见操作提示", "Model test passed; see feedback for save/enable status")}</p></div>}
      <OnboardingGuide step="model" /><Link className="block text-sm text-primary underline" href="/system/llm">{tr("打开完整模型设置", "Open full model settings")}</Link>
    </section>
    <section className="space-y-4 rounded-xl border p-4 sm:p-6"><h2 className="text-lg font-semibold">{tr("第 4 步：开始使用", "Step 4: Get started")}</h2>
      <p className="text-sm text-muted-foreground">{configured || done ? tr("配置已具备开始使用的条件。任务需由你填写并启动。", "Configuration is ready. Enter and start tasks yourself.") : tr("模型尚未配置；可以先浏览知识库，稍后回来连接 AI。", "No model configured. Browse the knowledge base and connect AI later.")}</p>
      <div className="flex flex-wrap gap-3">{(configured || done) && <Link className="rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground" href="/function/tasks">{tr("创建第一个任务", "Create first task")}</Link>}<Link className="rounded-md border px-4 py-2 text-sm" href="/system/knowledge">{tr("浏览知识库", "Browse knowledge base")}</Link></div>
      <OnboardingGuide step="start" /><Link className="block text-sm text-primary underline" href="/help">{tr("完整中文使用指南：结果、复核与导出", "Complete Chinese guide: results, review and export")}</Link>
    </section>
  </main>;
}
