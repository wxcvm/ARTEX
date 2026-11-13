"use client";

import * as React from "react";

import {
  Loader2Icon,
  PlugZapIcon,
  PlusIcon,
  RefreshCwIcon,
  RotateCcwIcon,
  SaveIcon,
  StarIcon,
  Trash2Icon,
  ZapIcon,
} from "lucide-react";
import { toast } from "sonner";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { api } from "@/lib/api";
import type { LLMPoolMember, LLMPoolStatus, LLMProfile, LLMRetryOverride } from "@/lib/types";
import { cn } from "@/lib/utils";

import { ProfileRetryFields, RetryPolicyPanel, ZERO_OVERRIDE } from "./_components/retry";
import { KeyDiagnosticsPanel } from "./_components/key-diagnostics";
import { OutboundAuditPanel } from "./_components/outbound-audit";

// The thinking switch (thinking.type) and the thinking effort (reasoning_effort) are two [independent] fields,
// each set on its own — some APIs have no thinking field and activate thinking via the effort parameter alone, hence the decoupling.
// An empty string in storage = the field is [not sent]; Radix Select rejects an empty value, so the UI uses "none"
// as the sentinel for not sent, converted to and from "" on load and save (NONE / fromStore / toStore).
const NONE = "none";
const fromStore = (v?: string) => (v ? v : NONE);
const toStore = (v: string) => (v === NONE ? "" : v);
const THINKING_TYPES: { value: string; label: string }[] = [
  { value: NONE, label: "不发送（默认）" },
  { value: "disabled", label: "关闭" },
  { value: "enabled", label: "开启" },
];
// Which request field name carries the output cap (only meaningful for the openai format). NONE ↔ "" shares the same sentinel conversion.
const MAX_TOKENS_FIELDS: { value: string; label: string }[] = [
  { value: NONE, label: "max_tokens（默认）" },
  { value: "max_completion_tokens", label: "max_completion_tokens" },
];
// The other formats hard-code their field names, so the option is meaningless for them; the hint texts spell this out.
const MAX_TOKENS_FIELD_HINTS: Record<string, string> = {
  openai:
    "Which key carries the cap. max_tokens is the default and the vast majority of compatible gateways only accept it; OpenAI's official reasoning models (the o series / GPT-5) conversely only accept max_completion_tokens and fail outright with unsupported_parameter when they receive max_tokens.",
  anthropic: "Only selectable for the openai format. Anthropic's field name is fixed as max_tokens.",
  "openai-responses": "Only selectable for the openai format. The Responses API's field name is fixed as max_output_tokens.",
  ollama: "Only selectable for the openai format. Ollama's field name is fixed as max_tokens.",
};
// LLM profile formats. The default in the create form is "anthropic"; "ollama" is the
// keyless local-server format (base URL defaults to http://127.0.0.1:11434/v1).
type ProfileFormat = "anthropic" | "openai" | "openai-responses" | "ollama";
const OLLAMA_DEFAULT_BASE_URL = "http://127.0.0.1:11434/v1";
const EFFORT_LEVELS: { value: string; label: string }[] = [
  { value: NONE, label: "不发送（默认）" },
  { value: "low", label: "低" },
  { value: "medium", label: "中" },
  { value: "high", label: "高" },
  { value: "xhigh", label: "极高" },
  { value: "max", label: "最大" },
];

function modelTestAdvice(message: string): string {
  if (/401|403|unauthorized|authentication|invalid.*(?:key|token)/i.test(message))
    return "认证未通过：请检查 API 密钥、账户权限和接口地址。";
  if (/404|model.*not found|unknown model|not found.*model/i.test(message))
    return "接口没有找到指定模型：检查模型名称，或使用旁边的按钮获取可用模型列表。";
  if (/429|quota|rate.limit|insufficient.balance/i.test(message))
    return "可能触发了速率限制或额度不足：检查服务端配额后重新测试。";
  if (/timeout|deadline exceeded|timed out|超时/i.test(message))
    return "连接超时：请检查网络、代理设置和模型服务是否正在响应。";
  if (/connection refused|ECONNREFUSED|connectex|fetch failed|dial tcp/i.test(message))
    return "无法连接到服务：请检查服务是否启动、接口地址及端口是否正确。";
  if (/certificate|\bTLS\b|\bSSL\b/i.test(message))
    return "TLS 连接失败：请检查接口 HTTPS 地址及服务端证书配置。";
  return "请检查接口类型、模型名称、地址和密钥；展开技术详情可查看服务端返回的原因。";
}

function cooldownText(secs: number) {
  if (secs <= 0) return "";
  if (secs < 60) return `${secs}s`;
  return `${Math.ceil(secs / 60)}min`;
}

// The "healthy or not" badge a profile shows on its card. A profile without a Key cannot send any request at all — say that before the breaker;
// the exception is Ollama, a keyless local server that IS callable without a key. The other states come from the pool's breaker
// records (with polling off no new records appear, so "healthy" = no known failure).
type Health = { label: string; cls: string; hint?: string };
function healthOf(p: LLMProfile, m?: LLMPoolMember): Health {
  if (!p.api_key_hint && p.format !== "ollama") {
    return {
      label: "未配置密钥",
      cls: "border-muted-foreground/40 text-muted-foreground",
      hint: "未设置 API Key，无法调用",
    };
  }
  if (m?.state === "tripped") {
    return {
      label: m.cooldown_secs > 0 ? `已熔断 · ${cooldownText(m.cooldown_secs)}` : "已熔断",
      cls: "border-destructive/50 text-destructive",
      hint: m.last_error,
    };
  }
  if (m?.state === "degraded") {
    return {
      label: `已降级 · 失败 ${m.fails} 次`,
      cls: "border-amber-500/50 text-amber-600 dark:text-amber-400",
      hint: m.last_error,
    };
  }
  return { label: "正常", cls: "border-emerald-500/50 text-emerald-600 dark:text-emerald-400" };
}

// ─────────────────────────────────────────────────────────────────────────────
// Pool settings drawer
// ─────────────────────────────────────────────────────────────────────────────

function PoolSheet({
  open,
  onOpenChange,
  pool,
  onReload,
}: {
  open: boolean;
  onOpenChange: (o: boolean) => void;
  pool: LLMPoolStatus | null;
  onReload: () => Promise<void>;
}) {
  const [busy, setBusy] = React.useState(false);

  // The cooldown countdown is the remaining seconds computed by the backend — fetched on a timer only while the drawer is open and some profile is unhealthy, to keep it ticking.
  React.useEffect(() => {
    if (!open || !pool?.enabled || !pool.chain.some((m) => m.state !== "ok")) return;
    const t = setInterval(() => void onReload(), 10_000);
    return () => clearInterval(t);
  }, [open, pool, onReload]);

  async function toggle(patch: { llm_pool_enabled?: boolean; llm_pool_bind_fallback?: boolean }) {
    if (busy) return;
    setBusy(true);
    try {
      await api.setSettings(patch);
      await onReload();
      if (patch.llm_pool_enabled !== undefined) {
        toast.success(patch.llm_pool_enabled ? "已启用模型池" : "已关闭模型池");
      } else {
        toast.success("故障切换设置已更新");
      }
    } catch (e) {
      toast.error(`应用设置失败：${(e as Error).message}`);
    } finally {
      setBusy(false);
    }
  }

  async function recover(id?: string) {
    try {
      await api.resetLLMPool(id);
      await onReload();
      toast.success(id ? "配置已恢复" : "全部配置已恢复");
    } catch (e) {
      toast.error(`恢复失败：${(e as Error).message}`);
    }
  }

  const enabled = pool?.enabled ?? false;
  const chain = pool?.chain ?? [];
  // Members taking part in the pool (excluding those flagged "excluded from pool"), in the exact order the backend tries them.
  const inChain = chain.filter((m) => m.active || !m.excluded);
  const tripped = chain.filter((m) => m.state === "tripped");

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="flex flex-col gap-0 p-0 data-[side=right]:sm:max-w-lg">
        <SheetHeader className="px-4">
          <SheetTitle className="flex items-center gap-2">
            <ZapIcon className="size-4" /> 模型池 · 故障切换
          </SheetTitle>
          <SheetDescription>
            开启后，<b>未固定模型配置</b>的智能体在当前配置不可用时（余额不足、密钥失效、限流或服务错误），会自动切换到下一项配置。
          </SheetDescription>
        </SheetHeader>

        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-4 pb-6">
          <div className="flex items-center justify-between gap-4 rounded-lg border p-3">
            <div className="grid gap-0.5">
              <Label className="text-sm">启用模型池</Label>
              <p className="text-muted-foreground text-xs">默认关闭。关闭时仅使用当前启用的配置，失败不会切换。</p>
            </div>
            <Switch
              checked={enabled}
              disabled={busy}
              onCheckedChange={(v) => void toggle({ llm_pool_enabled: v })}
              aria-label="模型池开关"
            />
          </div>

          {enabled && (
            <>
              <div className="flex items-center justify-between gap-4 rounded-lg border p-3">
                <div className="grid gap-0.5">
                  <Label className="text-sm">固定配置失败时也允许切换</Label>
                  <p className="text-muted-foreground text-xs">
                    默认关闭：智能体或任务指定配置后，仅使用该配置，失败时不会自动切换其他模型。
                    开启后，指定配置失败时也会按下方模型池顺序切换。
                  </p>
                </div>
                <Switch
                  checked={pool?.bind_fallback ?? false}
                  disabled={busy}
                  onCheckedChange={(v) => void toggle({ llm_pool_bind_fallback: v })}
                  aria-label="固定配置失败切换开关"
                />
              </div>

              <Separator />

              <div className="grid gap-2">
                <div className="flex items-center justify-between">
                  <Label className="text-sm">模型池顺序</Label>
                  {tripped.length > 0 && (
                    <Button size="sm" variant="ghost" onClick={() => void recover()}>
                      <RotateCcwIcon /> Recover all
                    </Button>
                  )}
                </div>
                {inChain.length < 2 && (
                  <p className="text-muted-foreground text-xs">
                    Only {inChain.length} usable profile(s) right now, so the pool has no effect — at least 2 usable profiles (API Key set, or Ollama) with pool participation are
                    required.
                  </p>
                )}
                {chain.map((m) => {
                  const excluded = m.excluded && !m.active;
                  const order = excluded ? null : inChain.findIndex((x) => x.profile_id === m.profile_id) + 1;
                  return (
                    <div
                      key={m.profile_id}
                      className={cn(
                        "grid gap-1 rounded-lg border p-2.5 text-sm",
                        excluded && "opacity-55",
                        m.state === "tripped" && "border-destructive/40",
                      )}
                    >
                      <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
                        <span className="w-5 shrink-0 text-center font-mono text-muted-foreground text-xs">
                          {order ?? "—"}
                        </span>
                        <span className="font-medium">{m.name}</span>
                        {m.active && (
                          <Badge variant="outline" className="border-amber-400/50 text-amber-500">
                            已启用
                          </Badge>
                        )}
                        {excluded && <Badge variant="outline">排除出模型池</Badge>}
                        <div className="ml-auto flex items-center gap-2">
                          {m.state === "tripped" && m.cooldown_secs > 0 && (
                            <span className="text-muted-foreground text-xs">Cooldown {cooldownText(m.cooldown_secs)}</span>
                          )}
                          {m.state === "degraded" && (
                            <span className="text-muted-foreground text-xs">{m.fails} consecutive failures</span>
                          )}
                          {m.state !== "ok" && (
                            <Button
                              size="icon"
                              variant="ghost"
                              className="size-7"
                              aria-label="立即恢复"
                              title="立即恢复：清除熔断状态，下次调用重新尝试此配置"
                              onClick={() => void recover(m.profile_id)}
                            >
                              <RotateCcwIcon className="size-3.5" />
                            </Button>
                          )}
                        </div>
                      </div>
                      <div className="flex flex-wrap items-center gap-x-3 pl-7 text-muted-foreground text-xs">
                        <code className="truncate font-mono">{m.model}</code>
                        {!m.active && <span>Priority {m.priority}</span>}
                      </div>
                      {m.last_error && (
                        <p className="truncate pl-7 font-mono text-muted-foreground text-xs" title={m.last_error}>
                          {m.last_error}
                        </p>
                      )}
                    </div>
                  );
                })}
                {chain.length === 0 && (
                  <div className="rounded-lg border border-dashed p-4 text-center text-muted-foreground text-sm">
                    暂无模型配置
                  </div>
                )}
              </div>

              <div className="rounded-lg border border-dashed p-3 text-muted-foreground text-xs leading-relaxed">
                当前启用配置始终排在首位，其余按优先级从高到低排列。配置失败后进入冷却（60 秒 → 5 分钟 → 30 分钟），冷却期间跳过，恢复后自动重新参与。上下文窗口无法容纳当前请求的配置也会跳过。固定模型的智能体和任务默认不参与模型池。
              </div>
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}

// ─────────────────────────────────────────────────────────────────────────────
// Model profile drawer (create / edit share the same form)
// ─────────────────────────────────────────────────────────────────────────────

function ProfileSheet({
  profile,
  open,
  onOpenChange,
  onSaved,
}: {
  profile: LLMProfile | null; // null = create
  open: boolean;
  onOpenChange: (o: boolean) => void;
  onSaved: (id: string) => void;
}) {
  const isNew = !profile;
  const [name, setName] = React.useState("");
  const [format, setFormat] = React.useState<ProfileFormat>("anthropic");
  const [model, setModel] = React.useState("");
  const [baseUrl, setBaseUrl] = React.useState("");
  const [proxy, setProxy] = React.useState("");
  const [apiKey, setApiKey] = React.useState("");
  const [keyHint, setKeyHint] = React.useState("");
  const [rps, setRps] = React.useState("0");
  const [rpm, setRpm] = React.useState("0");
  const [cw, setCw] = React.useState("0"); // context window (K tokens); 0 = default 200K
  const [thinkingType, setThinkingType] = React.useState(NONE);
  const [effort, setEffort] = React.useState(NONE);
  const [priority, setPriority] = React.useState("0"); // pool order; higher = tried earlier
  const [poolExclude, setPoolExclude] = React.useState(false);
  const [streaming, setStreaming] = React.useState(true); // true = streaming (default); false = non-streaming
  const [maxTokens, setMaxTokens] = React.useState("0"); // per-reply output cap; 0 = not sent
  const [maxTokensField, setMaxTokensField] = React.useState(NONE); // which field name carries the cap; NONE = max_tokens
  const [sessionHeaderKey, setSessionHeaderKey] = React.useState(""); // custom session header name; empty = not sent
  const [userAgent, setUserAgent] = React.useState(""); // custom User-Agent; empty = keep the default UA
  const [retry, setRetry] = React.useState<LLMRetryOverride>(ZERO_OVERRIDE); // retry override for this profile; all 0 = follow the global
  const [testing, setTesting] = React.useState(false);
  const testVersion = React.useRef(0);
  const [testFeedback, setTestFeedback] = React.useState<{
    ok: boolean;
    headline: string;
    hint: string;
    detail?: string;
  } | null>(null);
  const [saving, setSaving] = React.useState(false);
  const [models, setModels] = React.useState<string[]>([]);
  const modelListVersion = React.useRef(0);
  const [modelListFeedback, setModelListFeedback] = React.useState<{ ok: boolean; message: string } | null>(null);
  const [loadingModels, setLoadingModels] = React.useState(false);
  const [modelsOpen, setModelsOpen] = React.useState(false);
  const [advancedOpen, setAdvancedOpen] = React.useState(false);

  // Every time it opens, the form is re-filled from the passed-in profile (reset to defaults when creating). Closing and reopening
  // the drawer is a clean start, leaving no leftovers from the previous profile.
  React.useEffect(() => {
    if (!open) return;
    setName(profile?.name ?? "");
    setFormat(
      profile?.format === "openai" || profile?.format === "openai-responses" || profile?.format === "ollama"
        ? profile.format
        : "anthropic",
    );
    setModel(profile?.model ?? "");
    setBaseUrl(profile?.base_url ?? "");
    setProxy(profile?.proxy ?? "");
    setRps(String(profile?.rate_per_second ?? 0));
    setRpm(String(profile?.rate_per_minute ?? 0));
    setCw(String(profile?.context_window_k ?? 0));
    setThinkingType(fromStore(profile?.thinking_type));
    setEffort(fromStore(profile?.reasoning_effort));
    setPriority(String(profile?.priority ?? 0));
    setPoolExclude(profile?.pool_exclude ?? false);
    setStreaming(profile?.streaming ?? true);
    setMaxTokens(String(profile?.max_tokens ?? 0));
    setMaxTokensField(fromStore(profile?.max_tokens_field));
    setSessionHeaderKey(profile?.session_header_key ?? "");
    setUserAgent(profile?.user_agent ?? "");
    setRetry(profile?.retry ?? ZERO_OVERRIDE);
    setApiKey("");
    setKeyHint(profile?.api_key_hint ?? "");
    setModels([]);
    setModelsOpen(false);
    setModelListFeedback(null);
    testVersion.current += 1;
    setTestFeedback(null);
    setAdvancedOpen(Boolean(profile));
  }, [open, profile]);

  // A connection test only applies to the exact endpoint settings tested.
  React.useEffect(() => {
    testVersion.current += 1;
    setTestFeedback(null);
  }, [format, model, baseUrl, apiKey, proxy, thinkingType, effort, streaming, sessionHeaderKey, userAgent]);

  const profileId = profile ? Number(profile.id) : undefined;

  // Results from a previous endpoint or saved profile must never be offered
  // as selectable models for a newly edited connection.
  React.useEffect(() => {
    modelListVersion.current += 1;
    setModels([]);
    setModelsOpen(false);
    setModelListFeedback(null);
    setLoadingModels(false);
  }, [open, profile?.id, format, baseUrl, apiKey, proxy]);

  async function loadModels() {
    if (loadingModels) return;
    setLoadingModels(true);
    setModelsOpen(false);
    setModels([]);
    setModelListFeedback(null);
    const version = modelListVersion.current;
    try {
      const r = await api.fetchLLMModels(format, baseUrl, apiKey, proxy, profileId);
      if (version !== modelListVersion.current) return;
      if (r.ok && r.models && r.models.length > 0) {
        setModels(r.models);
        setModelsOpen(true);
        setModelListFeedback({ ok: true, message: `已找到 ${r.models.length} 个模型，请从列表中选择。` });
      } else {
        const reason = r.error || "服务端没有返回模型列表";
        setModelListFeedback({ ok: false, message: `无法获取模型列表：${reason}。你仍可手动填写模型名称。` });
      }
    } catch (e) {
      if (version !== modelListVersion.current) return;
      const reason = e instanceof Error ? e.message : "请求失败";
      setModelListFeedback({ ok: false, message: `无法获取模型列表：${reason}。请检查接口地址、网络或密钥，也可以手动输入。` });
    } finally {
      if (version === modelListVersion.current) setLoadingModels(false);
    }
  }

  async function testConnection() {
    if (testing) return;
    setTesting(true);
    setTestFeedback(null);
    const version = testVersion.current;
    try {
      // Use the same options as a normal LLM request; an empty key when editing
      // refers to the existing stored key, rather than accidentally testing blank credentials.
      const r = await api.testLLM(
        format,
        model,
        baseUrl,
        apiKey,
        proxy,
        toStore(thinkingType),
        toStore(effort),
        profileId,
        streaming,
        sessionHeaderKey.trim(),
        userAgent.trim(),
      );
      if (version !== testVersion.current) return;
      if (r.ok) {
        const latency = typeof r.latency_ms === "number" ? ` · ${r.latency_ms} ms` : "";
        setTestFeedback({
          ok: true,
          headline: `模型连接成功${latency}`,
          hint: `已收到来自 ${r.model || model || "模型服务"} 的有效测试响应。`,
          detail: r.reply ? `模型回复：${r.reply}` : undefined,
        });
        toast.success("模型连接测试通过");
      } else {
        const detail = r.error ?? "服务端没有返回详细错误";
        setTestFeedback({ ok: false, headline: "模型连接测试失败", hint: modelTestAdvice(detail), detail });
        toast.error("模型连接测试失败");
      }
    } catch (error) {
      if (version !== testVersion.current) return;
      const detail = error instanceof Error ? error.message : String(error);
      setTestFeedback({ ok: false, headline: "无法完成连接测试", hint: modelTestAdvice(detail), detail });
      toast.error("模型连接测试失败");
    } finally {
      setTesting(false);
    }
  }

  async function save() {
    if (!name.trim() || !model.trim()) {
      toast.error("请填写配置名称和模型");
      return;
    }
    if (saving) return;
    setSaving(true);
    try {
      const { id } = await api.saveLLMProfile({
        ...(profile ? { id: Number(profile.id) } : {}),
        name: name.trim(),
        format,
        model: model.trim(),
        base_url: baseUrl.trim(),
        proxy: proxy.trim(),
        api_key: apiKey,
        rate_per_second: Number(rps) || 0,
        rate_per_minute: Number(rpm) || 0,
        context_window_k: Number(cw) || 0,
        thinking_type: toStore(thinkingType),
        reasoning_effort: toStore(effort),
        priority: Number(priority) || 0,
        pool_exclude: poolExclude,
        streaming,
        max_tokens: Math.max(0, Number(maxTokens) || 0),
        // The field-name switch only matters for openai (Chat Completions); every other format falls back to the default,
        // and the backend re-applies the same normalization; this just keeps the UI from sending contradictory values.
        max_tokens_field: format === "openai" ? toStore(maxTokensField) : "",
        session_header_key: sessionHeaderKey.trim(),
        user_agent: userAgent.trim(),
        retry,
      });
      if (isNew) toast.success(`已创建：${name.trim()}（在卡片中设为启用即可使用）`);
      else toast.success(profile?.is_default ? "已保存；当前启用配置立即生效，无需重启" : "已保存");
      onSaved(String(id));
      onOpenChange(false);
    } catch (e) {
      toast.error(`保存失败：${(e as Error).message}`);
    } finally {
      setSaving(false);
    }
  }

  let keyPlaceholder = "sk-…";
  if (format === "ollama") {
    keyPlaceholder = "optional — a local Ollama server runs keyless";
  } else if (keyHint) {
    keyPlaceholder = `已设置（${keyHint}），留空保留`;
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side="right"
        className="flex min-w-0 flex-col gap-0 overflow-hidden p-0 data-[side=right]:w-full data-[side=right]:max-w-full data-[side=right]:sm:max-w-xl"
      >
        <SheetHeader className="px-4">
          <SheetTitle className="flex items-center gap-2">
            {isNew ? "添加模型连接" : `编辑模型：${profile?.name}`}
            {profile?.is_default && (
              <Badge variant="outline" className="border-amber-400/50 text-amber-500">
                已启用
              </Badge>
            )}
          </SheetTitle>
          <SheetDescription>
            {isNew
              ? "填写模型名称、接口类型与访问凭据；保存后可在卡片上设为默认模型。"
              : "修改后保存即可生效；当前默认模型的修改无需重启服务。"}
          </SheetDescription>
        </SheetHeader>

        <div className="flex min-h-0 min-w-0 flex-1 flex-col gap-4 overflow-y-auto px-4 pb-5 sm:px-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="p-name">配置名称</Label>
              <Input
                id="p-name"
                placeholder="例如：我的日常模型"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="grid gap-2">
              <Label>接口类型</Label>
              <Select value={format} onValueChange={(v) => setFormat(v as ProfileFormat)}>
                <SelectTrigger>
                  <SelectValue placeholder="选择接口类型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="anthropic">Anthropic API</SelectItem>
                  <SelectItem value="openai">OpenAI 兼容（Chat Completions）</SelectItem>
                  <SelectItem value="openai-responses">OpenAI Responses API</SelectItem>
                  <SelectItem value="ollama">本机 Ollama（无需密钥）</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid min-w-0 gap-2">
            <Label htmlFor="p-model">模型名称</Label>
            <div className="flex min-w-0 flex-wrap gap-2">
              <Input
                id="p-model"
                className="min-w-[10rem] flex-1 font-mono"
                placeholder={format === "ollama" ? "llama3.2" : "claude-opus-4-8"}
                value={model}
                onChange={(e) => setModel(e.target.value)}
              />
              {/* modal: this Popover's content is portaled to <body>, outside the Sheet's scroll lock,
                  so without modal the list renders but cannot scroll. modal lets it hold the topmost scroll lock itself. */}
              <Popover open={modelsOpen} onOpenChange={setModelsOpen} modal>
                <PopoverTrigger asChild>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="shrink-0"
                    disabled={loadingModels}
                    onClick={() => void loadModels()}
                    aria-label="从接口获取可用模型列表"
                    title="从当前接口获取可用模型；不支持列举时可手动输入"
                  >
                    {loadingModels ? <Loader2Icon className="size-4 animate-spin" /> : <RefreshCwIcon className="size-4" />}
                    {loadingModels ? "正在获取…" : "获取列表"}
                  </Button>
                </PopoverTrigger>
                {models.length > 0 && (
                  <PopoverContent className="max-h-72 w-72 gap-0 overflow-y-auto overscroll-contain p-1" align="end">
                    {models.map((m) => (
                      <button
                        key={m}
                        type="button"
                        className="w-full shrink-0 rounded-md px-2 py-1.5 text-left font-mono text-xs hover:bg-accent hover:text-accent-foreground"
                        onClick={() => {
                          setModel(m);
                          setModelsOpen(false);
                        }}
                      >
                        {m}
                      </button>
                    ))}
                  </PopoverContent>
                )}
              </Popover>
            </div>
            <p role={modelListFeedback?.ok ? "status" : modelListFeedback ? "alert" : undefined}
              className={cn("text-xs leading-relaxed break-words", modelListFeedback?.ok ? "text-emerald-500" : "text-muted-foreground")}>
              {modelListFeedback?.message ?? "可以直接输入模型名称；获取列表只是可选的辅助操作。"}
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="p-base-url">接口地址（可选）</Label>
            <Input
              id="p-base-url"
              className="font-mono"
              placeholder={format === "ollama" ? OLLAMA_DEFAULT_BASE_URL : "https://api.openai.com/v1"}
              value={baseUrl}
              onChange={(e) => setBaseUrl(e.target.value)}
            />
            <p className="text-xs leading-relaxed text-muted-foreground">
              {format === "ollama"
                ? "本机 Ollama 默认使用 127.0.0.1:11434，可留空。"
                : "官方接口可留空；使用 LM Studio 或其他兼容网关时填写对应 API 地址。"}
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="p-api-key">
              {format === "ollama" ? "API 密钥（Ollama 可留空）" : "API 密钥"}
            </Label>
            <Input
              id="p-api-key"
              type="password"
              placeholder={keyPlaceholder}
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
            />
          </div>

          <section className="min-w-0 rounded-xl border border-border/70 bg-muted/20">
            <button
              type="button"
              className="flex w-full items-center justify-between gap-3 px-4 py-3 text-left text-sm font-semibold hover:bg-muted/30"
              aria-expanded={advancedOpen}
              aria-controls="model-advanced-settings"
              onClick={() => setAdvancedOpen((value) => !value)}
            >
              <span>高级设置 <span className="font-normal text-muted-foreground">· 代理、速率、推理、重试</span></span>
              <span className="shrink-0 text-xs font-normal text-muted-foreground">{advancedOpen ? "收起" : "展开"}</span>
            </button>
            <div id="model-advanced-settings" className={advancedOpen ? "grid min-w-0 gap-4 border-t border-border/70 px-4 py-4" : "hidden"}>
          <div className="grid gap-2">
            <Label htmlFor="p-proxy">代理（可选）</Label>
            <Input
              id="p-proxy"
              className="font-mono"
              placeholder="socks5://user:pass@127.0.0.1:1080 · http://127.0.0.1:8080"
              value={proxy}
              onChange={(e) => setProxy(e.target.value)}
            />
            <p className="text-muted-foreground text-xs">
              仅模型出站请求经过此代理，支持 HTTP、HTTPS、SOCKS5，可配置认证信息（例如 socks5://user:pass@host:port；密码含特殊字符时需 URL 编码）。留空为直连。
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="p-session-header">自定义会话请求头（可选）</Label>
            <Input
              id="p-session-header"
              className="font-mono"
              placeholder="例如 x-session-id（留空不发送）"
              value={sessionHeaderKey}
              onChange={(e) => setSessionHeaderKey(e.target.value)}
            />
            <p className="text-muted-foreground text-xs">
              设置请求头名称后，每次请求会自动携带<b>当前会话 ID</b>（例如对话 conv-12、执行会话 exp3-worker-i87）。适用于按会话 ID 做提示词缓存或粘性路由的网关；同一会话保持一致，不同会话相互独立。留空不发送。
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="p-user-agent">自定义 User-Agent（可选）</Label>
            <Input
              id="p-user-agent"
              className="font-mono"
              placeholder="例如 my-app/1.0（留空使用默认 UA）"
              value={userAgent}
              onChange={(e) => setUserAgent(e.target.value)}
            />
            <p className="text-muted-foreground text-xs">
              覆盖每次模型出站请求的 <b>User-Agent</b> 请求头，适用于依据 User-Agent 路由或匹配白名单的网关。留空保留默认值。
            </p>
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            <div className="grid gap-2">
              <Label htmlFor="p-rps">每秒速率限制</Label>
              <Input id="p-rps" type="number" min={0} value={rps} onChange={(e) => setRps(e.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="p-rpm">每分钟速率限制</Label>
              <Input id="p-rpm" type="number" min={0} value={rpm} onChange={(e) => setRpm(e.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="p-cw">上下文窗口（K）</Label>
              <Input
                id="p-cw"
                type="number"
                min={0}
                max={1000}
                value={cw}
                onChange={(e) => setCw(e.target.value)}
                placeholder="200"
              />
            </div>
          </div>
          <p className="-mt-2 text-muted-foreground text-xs">
            速率限制为 0 表示不限速，由所有智能体共享。上下文窗口单位为 K（千 Token），0 使用默认 200K，最大 1000K（即 1M）。设置过高可能导致上下文压缩无法触发。
          </p>

          <div className="grid gap-3 rounded-lg border p-3">
            <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label htmlFor="p-priority" className="text-sm">
                  模型池优先级
                </Label>
                <p className="text-muted-foreground text-xs">
                  数值越大越优先。当前启用配置始终优先，不受此数值影响。相同优先级的配置轮流排在前面，分摊额度。
                </p>
              </div>
              <Input
                id="p-priority"
                type="number"
                className="w-24 shrink-0"
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
              />
            </div>
            <div className="flex flex-col items-start gap-3 border-t pt-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label className="text-sm">排除出模型池</Label>
                <p className="text-muted-foreground text-xs">
                  开启后，此配置不会作为故障切换目标，但智能体或任务仍可明确指定它。适用于保留给特定智能体的高成本模型，避免被其他智能体的失败消耗。
                </p>
              </div>
              <Switch checked={poolExclude} onCheckedChange={setPoolExclude} aria-label="排除出模型池" />
            </div>
            <div className="flex flex-col items-start gap-3 border-t pt-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label className="text-sm">流式输出</Label>
                <p className="text-muted-foreground text-xs">
                  开启（默认）使用 SSE 流式输出，实时显示进度和 Token 数量。关闭后使用非流式请求（stream:false，一次返回完整响应），可兼容 SSE 实现异常的网关，但不再显示实时生成进度。
                </p>
              </div>
              <Switch checked={streaming} onCheckedChange={setStreaming} aria-label="流式输出" />
            </div>
          </div>

          <div className="grid gap-3 rounded-lg border p-3">
            <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label htmlFor="p-max-tokens" className="text-sm">
                  输出上限 · max tokens
                </Label>
                <p className="text-muted-foreground text-xs">
                  单次回复可生成的最大 Token 数，每次请求均携带。0（默认）表示不发送该字段，使用服务端默认值。此值不同于上方的上下文窗口：后者表示模型总容量，仅在本地用于计算压缩阈值。设置过低可能导致推理模型在思考过程中被截断，无法输出答案。
                </p>
              </div>
              <Input
                id="p-max-tokens"
                type="number"
                min={0}
                className="w-28 shrink-0"
                value={maxTokens}
                onChange={(e) => setMaxTokens(e.target.value)}
                placeholder="0"
              />
            </div>
            <div className="flex flex-col items-start gap-3 border-t pt-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label className="text-sm">上限字段名称</Label>
                <p className="text-muted-foreground text-xs">{MAX_TOKENS_FIELD_HINTS[format]}</p>
              </div>
              <Select
                value={format === "openai" ? maxTokensField : NONE}
                onValueChange={setMaxTokensField}
                disabled={format !== "openai"}
              >
                <SelectTrigger className="w-56 shrink-0">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {MAX_TOKENS_FIELDS.map((o) => (
                    <SelectItem key={o.value} value={o.value}>
                      {o.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid gap-3 rounded-lg border p-3">
            <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label className="text-sm">思考开关 · thinking.type</Label>
                <p className="text-muted-foreground text-xs">
                  控制是否发送 thinking 字段。不发送表示省略该字段（兼容不支持此字段的模型，如 MiniMax）；关闭发送 disabled，开启发送 enabled。与下方思考强度独立。
                </p>
              </div>
              <Select value={thinkingType} onValueChange={setThinkingType}>
                <SelectTrigger className="w-32 shrink-0">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {THINKING_TYPES.map((o) => (
                    <SelectItem key={o.value} value={o.value}>
                      {o.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="flex flex-col items-start gap-3 border-t pt-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="grid gap-0.5">
                <Label className="text-sm">思考强度 · reasoning_effort</Label>
                <p className="text-muted-foreground text-xs">
                  A separate effort knob (OpenAI reasoning_effort / Anthropic output_config.effort). Some APIs have no thinking
                  field and activate thinking via effort alone, so it can be set on its own without sending the thinking switch.
                </p>
              </div>
              <Select value={effort} onValueChange={setEffort}>
                <SelectTrigger className="w-32 shrink-0">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {EFFORT_LEVELS.map((o) => (
                    <SelectItem key={o.value} value={o.value}>
                      {o.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <ProfileRetryFields value={retry} onChange={setRetry} />
            </div>
          </section>
        </div>

        <div className="border-t border-border/70 bg-background">
          {testing && (
            <div role="status" aria-live="polite" className="flex items-center gap-2 px-4 pt-3 text-sm text-muted-foreground sm:px-5">
              <Loader2Icon className="size-4 animate-spin" />
              正在向模型发送测试请求并检查响应…
            </div>
          )}
          {!testing && testFeedback && (
            <div
              role={testFeedback.ok ? "status" : "alert"}
              className={cn(
                "mx-4 mt-3 max-h-44 overflow-y-auto rounded-lg border px-3 py-2 text-sm sm:mx-5",
                testFeedback.ok
                  ? "border-emerald-500/30 bg-emerald-500/10"
                  : "border-amber-500/30 bg-amber-500/10",
              )}
            >
              <p className="font-semibold">{testFeedback.headline}</p>
              <p className="mt-1 leading-relaxed text-muted-foreground">{testFeedback.hint}</p>
              {testFeedback.detail && (
                <details className="mt-2">
                  <summary className="cursor-pointer text-xs font-medium">查看{testFeedback.ok ? "模型回复" : "原始错误"}</summary>
                  <pre className="mt-1 whitespace-pre-wrap break-all font-mono text-xs text-muted-foreground">{testFeedback.detail}</pre>
                </details>
              )}
            </div>
          )}
          <div className="flex flex-wrap gap-2 px-4 py-3 sm:px-5">
          <Button variant="outline" className="min-w-0 flex-1 sm:flex-none" onClick={testConnection} disabled={testing}>
            {testing ? <Loader2Icon className="animate-spin" /> : <PlugZapIcon />}
            {testing ? "正在测试…" : "测试连接"}
          </Button>
          <Button onClick={save} disabled={saving} className="min-w-0 flex-1">
            {saving && <Loader2Icon className="animate-spin" />}
            {!saving && (isNew ? <PlusIcon /> : <SaveIcon />)}
            {isNew ? "保存模型配置" : "保存更改"}
          </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}

// ─────────────────────────────────────────────────────────────────────────────

export default function LLMPage() {
  const [profiles, setProfiles] = React.useState<LLMProfile[]>([]);
  const [pool, setPool] = React.useState<LLMPoolStatus | null>(null);
  const [poolOpen, setPoolOpen] = React.useState(false);
  // The drawer's open flag and its content are stored separately: editing stays put on close, otherwise during the closing animation the title would flash
  // from "Edit X" to "New". editing = null means create.
  const [editOpen, setEditOpen] = React.useState(false);
  const [editing, setEditing] = React.useState<LLMProfile | null>(null);
  const openEditor = React.useCallback((p: LLMProfile | null) => {
    setEditing(p);
    setEditOpen(true);
  }, []);

  const loadPool = React.useCallback(async () => {
    try {
      setPool(await api.llmPool());
    } catch {
      /* ignore */
    }
  }, []);

  const load = React.useCallback(async () => {
    try {
      setProfiles(await api.llmProfiles());
    } catch {
      /* ignore */
    }
    await loadPool();
  }, [loadPool]);

  React.useEffect(() => {
    void load();
  }, [load]);

  // The health badge on each card reads pool state by profile id.
  const health = React.useMemo(() => {
    const m = new Map<string, LLMPoolMember>();
    for (const c of pool?.chain ?? []) m.set(c.profile_id, c);
    return m;
  }, [pool]);

  async function activate(id: string, name: string) {
    try {
      await api.activateLLMProfile(id);
      toast.success(`已启用：${name}`);
      await load();
    } catch (e) {
      toast.error(`启用失败：${(e as Error).message}`);
    }
  }

  async function remove(p: LLMProfile) {
    if (p.is_default) {
      toast.error("无法删除当前启用的配置");
      return;
    }
    try {
      await api.deleteLLMProfile(p.id);
      toast.success(`已删除：${p.name}`);
      await load();
    } catch (e) {
      toast.error(`删除失败：${(e as Error).message}`);
    }
  }

  const poolOn = pool?.enabled ?? false;

  return (
    <div className="flex min-w-0 flex-1 flex-col gap-5 pb-6">
      <header className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-border/70 bg-gradient-to-br from-card to-primary/5 px-5 py-5 shadow-sm sm:px-7">
        <div className="min-w-0">
          <p className="mb-1 text-[11px] font-semibold tracking-[0.16em] text-primary">ARTEX / 模型</p>
          <h1 className="text-2xl font-semibold tracking-tight">模型配置</h1>
          <p className="mt-1 max-w-2xl text-sm leading-relaxed text-muted-foreground">
            管理智能体使用的模型连接、默认配置与故障切换。点击配置卡片可编辑，星标表示当前默认配置。
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button size="sm" variant="outline" onClick={() => setPoolOpen(true)}>
            <ZapIcon /> 模型池设置
            {poolOn && (
              <Badge variant="outline" className="ml-1 border-emerald-500/50 text-emerald-600 dark:text-emerald-400">
                已开启
              </Badge>
            )}
          </Button>
          <Button size="sm" variant="outline" onClick={() => openEditor(null)}>
            <PlusIcon /> 新建模型配置
          </Button>
        </div>
      </header>

      <Tabs defaultValue="profiles" className="min-w-0 flex-1">
        <TabsList className="h-auto max-w-full rounded-lg bg-muted/60 p-1">
          <TabsTrigger value="profiles">模型连接</TabsTrigger>
          <TabsTrigger value="retry">重试与退避</TabsTrigger>
          <TabsTrigger value="privacy">密钥保护</TabsTrigger>
          <TabsTrigger value="outbound">连接审计</TabsTrigger>
        </TabsList>

        <TabsContent value="profiles" className="mt-4">
          <div className="grid min-w-0 grid-cols-1 gap-4 lg:grid-cols-2 2xl:grid-cols-3">
            {profiles.map((p) => {
              const h = healthOf(p, health.get(p.id));
              return (
                // biome-ignore lint/a11y/useSemanticElements: the card contains its own action buttons; a native <button> would nest buttons (invalid HTML)
                <Card
                  key={p.id}
                  role="button"
                  tabIndex={0}
                  onClick={() => openEditor(p)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      e.preventDefault();
                      openEditor(p);
                    }
                  }}
                  className={cn(
                    "min-w-0 cursor-pointer gap-2 rounded-xl border-border/70 bg-card/90 py-4 shadow-sm outline-none transition-colors hover:border-primary/40 focus-visible:ring-2 focus-visible:ring-primary/60",
                    p.is_default && "border-amber-400/50 bg-amber-400/5",
                  )}
                >
                  <CardContent className="grid min-w-0 gap-3 px-5">
                    <div className="flex items-start gap-2">
                      <StarIcon
                        className={cn(
                          "mt-0.5 size-4 shrink-0",
                          p.is_default ? "fill-amber-400 text-amber-400" : "text-muted-foreground",
                        )}
                      />
                      <div className="min-w-0 flex-1">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className="truncate text-base font-semibold">{p.name}</span>
                          <Badge variant="outline" className="uppercase">
                            {p.format}
                          </Badge>
                          <Badge variant="outline" className={cn("ml-auto", h.cls)} title={h.hint}>
                            {h.label}
                          </Badge>
                        </div>
                        <code className="mt-1 block truncate font-mono text-muted-foreground text-xs">{p.model}</code>
                      </div>
                    </div>

                    <div className="flex flex-wrap gap-x-3 gap-y-0.5 pl-6 text-muted-foreground text-xs">
                      {p.api_key_hint && <span>{p.api_key_hint}</span>}
                      <span>
                        {p.rate_per_second}/s · {p.rate_per_minute}/min
                      </span>
                      {p.proxy && <span className="truncate">Proxy {p.proxy}</span>}
                      {p.reasoning_effort && (
                        <span>Thinking {p.reasoning_effort === "off" ? "off" : p.reasoning_effort}</span>
                      )}
                      {/* The two pool-related fields only matter when the pool is on; when it is off they take no space */}
                      {poolOn &&
                        !p.is_default &&
                        (p.pool_exclude ? <span>排除出模型池</span> : <span>优先级 {p.priority ?? 0}</span>)}
                    </div>

                    <div className="mt-1 flex gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        className="flex-1"
                        disabled={p.is_default}
                        onClick={(e) => {
                          e.stopPropagation();
                          void activate(p.id, p.name);
                        }}
                      >
                        {p.is_default ? "当前默认" : "设为默认"}
                      </Button>
                      <Button
                        size="icon"
                        variant="outline"
                        aria-label="删除模型配置"
                        onClick={(e) => {
                          e.stopPropagation();
                          void remove(p);
                        }}
                      >
                        <Trash2Icon className="text-destructive" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
            {profiles.length === 0 && (
              <div className="col-span-full flex min-h-64 flex-col items-center justify-center gap-4 rounded-2xl border border-dashed border-border/70 bg-card/70 px-6 py-12 text-center">
                <PlugZapIcon className="size-9 text-primary" aria-hidden="true" />
                <div>
                  <h2 className="text-lg font-semibold">还没有模型配置</h2>
                  <p className="mt-1 max-w-md text-sm leading-relaxed text-muted-foreground">
                    添加已有模型服务的地址、模型名称和凭据，再选择默认配置，即可供智能体使用。
                  </p>
                </div>
                <Button onClick={() => openEditor(null)}>
                  <PlusIcon /> 添加第一个模型
                </Button>
              </div>
            )}
          </div>
        </TabsContent>

        <TabsContent value="privacy" className="mt-4">
          <KeyDiagnosticsPanel profiles={profiles} />
        </TabsContent>

        <TabsContent value="outbound" className="mt-4">
          <OutboundAuditPanel />
        </TabsContent>

        <TabsContent value="retry" className="mt-4">
          <RetryPolicyPanel />
        </TabsContent>
      </Tabs>

      <ProfileSheet profile={editing} open={editOpen} onOpenChange={setEditOpen} onSaved={() => void load()} />
      <PoolSheet open={poolOpen} onOpenChange={setPoolOpen} pool={pool} onReload={loadPool} />
    </div>
  );
}
