"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

// Shared pieces of the LLM retry settings: "attempts + interval" for each of the five retry layers.
//
// The five layers, innermost to outermost: connect (SDK) → empty response (SDK) → same-provider safe window → pool breaker → intent rerun.
// The first three travel with the endpoint, so every model profile can override the global defaults; the last two are process-level and exist only once globally.
//
// All inputs follow the same "empty = unconfigured" semantics, matching the backend db.RetryRule:
//   attempts   empty/0 = built-in default | -1 = disable this layer's retries | >0 = use this count
//   interval   empty/0 = this layer's original exponential backoff | >0 = use this fixed interval in milliseconds

import * as React from "react";

import { Loader2Icon, SaveIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import type { LLMRetryOverride, LLMRetryPolicy, LLMRetryRule } from "@/lib/types";

export const ZERO_RULE: LLMRetryRule = { attempts: 0, interval_ms: 0 };
export const ZERO_OVERRIDE: LLMRetryOverride = {
  connect: ZERO_RULE,
  empty: ZERO_RULE,
  stream: ZERO_RULE,
};
const ZERO_POLICY: LLMRetryPolicy = {
  ...ZERO_OVERRIDE,
  breaker: ZERO_RULE,
  intent: ZERO_RULE,
};

type LayerMeta = {
  title: string;
  /** Where this retry layer happens and who runs it */
  where: string;
  /** Which errors reach this layer — down to status codes, no guessing */
  trigger: string;
  /** Errors that look similar but do [not] go through this layer, so nobody thinks it is a bug after seeing no effect */
  skips?: string;
  desc: string;
  attemptsLabel: string;
  /** Default attempts when left empty, used for the placeholder */
  defAttempts: number;
  /** Default interval strategy when left empty, used for the placeholder */
  defInterval: string;
  /** Meaning of entering -1 for attempts */
  offHint: string;
};

export const RETRY_LAYERS = {
  connect: {
    title: "Connect retry",
    where: "SDK · before a 200 arrives",
    trigger:
      "Cannot connect or no 200 yet: network-layer errors such as connection reset / read-write timeout / DNS failure, plus HTTP 408, 429, 500, 502, 503, 504.",
    skips: "The remaining status codes (400 / 401 / 403 / 404 / 413 / 422 etc.) are deterministic rejections — resending fails the same way, so they bubble up directly.",
    desc: "Resends the same request verbatim. Once the stream has started (a 200 was received), a mid-stream drop is no longer this layer's business.",
    attemptsLabel: "Retry attempts",
    defAttempts: 3,
    defInterval: "0.5s→1s→2s exponential (capped at 8s)",
    offHint: "-1 = no retries at all; the failure bubbles up immediately",
  },
  empty: {
    title: "Empty-response retry",
    where: "SDK · openai format only",
    trigger:
      "HTTP 200 with a normal stop finish_reason, yet not a single content block in the response — gateway empty frames, dropped thinking-field frames and sampling hiccups all look like this.",
    skips: "Empty because of max_tokens truncation does not count (fix that by raising the output cap; resending just hits it again).",
    desc: "Resends the entire prompt, so it is relatively expensive on long contexts; keep the count small.",
    attemptsLabel: "Retry attempts",
    defAttempts: 2,
    defInterval: "0.5s→1s→2s exponential (capped at 8s)",
    offHint: "-1 = hand the empty response through as-is",
  },
  stream: {
    title: "Same-provider safe-window retry",
    where: "This project · before any output is delivered",
    trigger:
      "Trouble only after the stream is established (a 200 received): mid-stream disconnects, provider overloaded, in-stream 429 / 5xx error events — while not a single token has been handed to the caller yet.",
    skips:
      "Quota exhausted (402 / insufficient_quota, handed to the pool to switch configs), context too long (413 / context length, handed to compression), and deterministic rejections (400 / 401 / 403 / 404 / 422) are never retried.",
    desc: "Replays the same request on the same config. Since no output has been delivered yet, the replay cannot duplicate model output or tool executions.",
    attemptsLabel: "Retry attempts",
    defAttempts: 2,
    defInterval: "0.5s→1s exponential (capped at 4s)",
    offHint: "-1 = a broken stream goes straight to the outer intent rerun",
  },
  breaker: {
    title: "Pool breaker",
    where: "This project · process-level, one global copy",
    trigger:
      "Trips when transient failures (429, 5xx, network errors) accumulate consecutively to the threshold; deterministic failures — insufficient balance (402), dead key (401 / 403), missing model (404) — ignore the threshold and trip on the first occurrence.",
    skips: "One success resets the counter, so an occasionally flaky config is not slowly ground into tripping.",
    desc: "After tripping it enters cooldown, and while cooling the pool skips the config outright. State is persisted, surviving restarts.",
    attemptsLabel: "Consecutive failures to trip",
    defAttempts: 3,
    defInterval: "1min→5min→30min gradient",
    offHint: "-1 = transient failures never trip (deterministic failures still do)",
  },
  intent: {
    title: "Intent rerun",
    where: "This project · process-level, one global copy",
    trigger:
      "The earlier layers could not save it: the worker ends with model_error — inner retries all exhausted, or the stream broke only after output delivery had started (when replay is unsafe and only a full rerun works).",
    skips: "Quota exhaustion is already handled by the pool switching configs and is not rerun here; when the task is paused / aborted / wrapping up it yields immediately, spending no backoff time.",
    desc: "Reruns the whole intent from scratch. It is the outermost layer, so one rerun multiplies the inner layers' attempt counts once more.",
    attemptsLabel: "Rerun attempts",
    defAttempts: 2,
    defInterval: "Fixed 3s",
    offHint: "-1 = no rerun; the intent is marked blocked outright",
  },
} satisfies Record<string, LayerMeta>;

type LayerKey = keyof typeof RETRY_LAYERS;

const RETRY_LAYERS_ZH: Record<LayerKey, LayerMeta> = {
  "connect": {
    "title": "连接重试",
    "where": "SDK · 收到 HTTP 200 之前",
    "trigger": "连接未建立或尚未收到 200：连接重置、读写超时、DNS 故障等网络层错误，以及 HTTP 408、429、500、502、503、504。",
    "skips": "其余状态码（400 / 401 / 403 / 404 / 413 / 422 等）属于确定性拒绝，重发仍会失败，直接向外层返回。",
    "desc": "原样重发同一请求。收到 200 后流式响应已建立，中途断流不由本层处理。",
    "attemptsLabel": "重试次数",
    "defAttempts": 3,
    "defInterval": "0.5 秒→1 秒→2 秒指数退避（最多 8 秒）",
    "offHint": "-1 = 完全不重试，失败直接返回"
  },
  "empty": {
    "title": "空响应重试",
    "where": "SDK · 仅 OpenAI 格式",
    "trigger": "HTTP 200 且正常 stop 结束，但响应无任何内容块；可能是网关空帧、思考字段丢帧或采样异常。",
    "skips": "max_tokens 截断导致的空响应不计入；应提高输出上限，重发仍会截断。",
    "desc": "重新发送完整提示词，长上下文成本较高，建议保持较少次数。",
    "attemptsLabel": "重试次数",
    "defAttempts": 2,
    "defInterval": "0.5 秒→1 秒→2 秒指数退避（最多 8 秒）",
    "offHint": "-1 = 原样返回空响应"
  },
  "stream": {
    "title": "同供应商安全窗口重试",
    "where": "本项目 · 尚未交付任何输出",
    "trigger": "流式响应建立后（收到 200）才发生故障：中途断连、服务过载、流内 429 / 5xx 错误，且尚未向调用方交付任何 Token。",
    "skips": "额度耗尽（402 / insufficient_quota，由模型池切换配置）、上下文过长（413 / context length，由压缩处理）及确定性拒绝（400 / 401 / 403 / 404 / 422）不重试。",
    "desc": "在同一配置重放同一请求。尚未交付输出，因此不会重复模型输出或工具执行。",
    "attemptsLabel": "重试次数",
    "defAttempts": 2,
    "defInterval": "0.5 秒→1 秒指数退避（最多 4 秒）",
    "offHint": "-1 = 断流直接交给外层意图重跑"
  },
  "breaker": {
    "title": "模型池熔断",
    "where": "本项目 · 进程级全局配置",
    "trigger": "429、5xx、网络错误等瞬时故障连续达到阈值时熔断；余额不足（402）、密钥失效（401 / 403）、模型不存在（404）等确定性故障首次即熔断。",
    "skips": "一次成功即清零计数，偶发故障不会逐渐累积至熔断。",
    "desc": "熔断后进入冷却期，模型池跳过该配置。状态持久化，重启后仍保留。",
    "attemptsLabel": "连续失败熔断阈值",
    "defAttempts": 3,
    "defInterval": "1 分钟→5 分钟→30 分钟阶梯",
    "offHint": "-1 = 瞬时故障不熔断（确定性故障仍熔断）"
  },
  "intent": {
    "title": "意图重跑",
    "where": "本项目 · 进程级全局配置",
    "trigger": "前几层均未恢复：执行智能体以 model_error 结束，内部重试耗尽，或已经交付输出后断流（不能安全重放，只能重新运行整个意图）。",
    "skips": "额度耗尽已由模型池切换处理，本层不重跑；任务暂停、中止或收尾时立即退出，不消耗退避时间。",
    "desc": "从头重跑整个意图。作为最外层，一次重跑会再次乘上内部各层尝试次数。",
    "attemptsLabel": "重跑次数",
    "defAttempts": 2,
    "defInterval": "固定 3 秒",
    "offHint": "-1 = 不重跑，意图直接标记为阻塞"
  }
};

/** Milliseconds in human terms, only echoed next to the input so nobody counts zeros. */
function humanMs(ms: number) {
  if (!Number.isFinite(ms) || ms <= 0) return "";
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60_000) return `${Number((ms / 1000).toFixed(2))}s`;
  return `${Number((ms / 60_000).toFixed(2))}min`;
}

/** Controlled numeric input: empty string ↔ 0; intermediate states ("-", "1e") stay local and never bother the parent. */
function NumField({
  id,
  value,
  onChange,
  placeholder,
  min,
}: {
  id: string;
  value: number;
  onChange: (n: number) => void;
  placeholder: string;
  min: number;
}) {
  const [text, setText] = React.useState(value === 0 ? "" : String(value));
  // Keeps up when the parent swaps in a whole new set of values (policy loaded, profile switched); typing never reaches this,
  // because by then value already equals the parsed result of the local text.
  React.useEffect(() => {
    const incoming = value === 0 ? "" : String(value);
    setText((cur) => (Number(cur || 0) === value ? cur : incoming));
  }, [value]);
  return (
    <Input
      id={id}
      type="number"
      min={min}
      className="w-28 shrink-0"
      value={text}
      placeholder={placeholder}
      onChange={(e) => {
        setText(e.target.value);
        const n = Number(e.target.value);
        onChange(e.target.value.trim() === "" || !Number.isFinite(n) ? 0 : Math.trunc(n));
      }}
    />
  );
}

/** The two knobs of one retry layer. idPrefix keeps the label's htmlFor unique when it appears multiple times on one page. */
export function RetryRuleFields({
  layer,
  idPrefix,
  value,
  onChange,
  compact,
}: {
  layer: LayerKey;
  idPrefix: string;
  value: LLMRetryRule;
  onChange: (r: LLMRetryRule) => void;
  /** true = compact version inside the profile drawer: skips the expanded explanations, keeping only the "which errors reach this layer" line */
  compact?: boolean;
}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const meta = locale === "en" ? RETRY_LAYERS[layer] : RETRY_LAYERS_ZH[layer];
  const human = humanMs(value.interval_ms);
  return (
    <div className={compact ? "grid gap-2" : "grid gap-3 rounded-lg border p-3"}>
      <div className="grid gap-0.5">
        <div className="flex flex-wrap items-baseline gap-2">
          <Label className="text-sm">{meta.title}</Label>
          <span className="text-muted-foreground text-xs">{meta.where}</span>
        </div>
        {/* Which errors reach this layer, down to status codes — if the knobs seem to have no effect, the error most likely does not land on this layer at all. */}
        <p className="text-muted-foreground text-xs">
          <span className="font-medium text-foreground">{tr("触发条件", "Triggers")}</span>: {meta.trigger}
        </p>
        {!compact && meta.skips && (
          <p className="text-muted-foreground text-xs">
            <span className="font-medium text-foreground">{tr("跳过条件", "Skipped here")}</span>: {meta.skips}
          </p>
        )}
        {!compact && <p className="text-muted-foreground text-xs">{meta.desc}</p>}
      </div>
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2">
        <div className="flex items-center gap-2">
          <Label htmlFor={`${idPrefix}-${layer}-n`} className="text-muted-foreground text-xs">
            {meta.attemptsLabel}
          </Label>
          <NumField
            id={`${idPrefix}-${layer}-n`}
            min={-1}
            value={value.attempts}
            placeholder={tr(`默认 ${meta.defAttempts}`, `Default ${meta.defAttempts}`)}
            onChange={(n) => onChange({ ...value, attempts: n })}
          />
        </div>
        <div className="flex items-center gap-2">
          <Label htmlFor={`${idPrefix}-${layer}-ms`} className="text-muted-foreground text-xs">
            {tr("间隔（毫秒）", "Interval ms")}
          </Label>
          <NumField
            id={`${idPrefix}-${layer}-ms`}
            min={0}
            value={value.interval_ms}
            placeholder={tr("默认退避", "Default backoff")}
            onChange={(n) => onChange({ ...value, interval_ms: n })}
          />
          <span className="text-muted-foreground text-xs">{human ? tr(`固定 ${human}`, `Fixed ${human}`) : meta.defInterval}</span>
        </div>
      </div>
      {!compact && <p className="text-muted-foreground text-xs">{tr("留空 = 使用默认值；", "Empty = use the default; ")}{meta.offHint}.</p>}
    </div>
  );
}

/** The three-layer override in the model profile drawer (the three layers that travel with the endpoint). */
export function ProfileRetryFields({
  value,
  onChange,
}: {
  value: LLMRetryOverride;
  onChange: (o: LLMRetryOverride) => void;
}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  return (
    <div className="grid gap-3 rounded-lg border p-3">
      <div className="grid gap-0.5">
        <Label className="text-sm">{tr("重试覆盖配置", "Retry overrides")}</Label>
        <p className="text-muted-foreground text-xs">
          {tr("仅适用于当前模型配置，覆盖「重试与退避」全局默认值。留空跟随全局；次数 -1 禁用该层；填写间隔将使用固定退避。熔断和意图重跑属于进程级配置，只能在全局页面调整。", "Applies only to this profile, overriding global Retry & backoff defaults. Empty follows global; attempts -1 disables the layer; a filled interval uses fixed backoff. Breaker and intent rerun are process-level and configured globally.")}
        </p>
      </div>
      {(["connect", "empty", "stream"] as const).map((k) => (
        <div key={k} className="border-t pt-3 first:border-t-0 first:pt-0">
          <RetryRuleFields
            compact
            layer={k}
            idPrefix="pf"
            value={value[k]}
            onChange={(r) => onChange({ ...value, [k]: r })}
          />
        </div>
      ))}
    </div>
  );
}

/** The "Retry & backoff" tab: global defaults for the five layers. */
export function RetryPolicyPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [policy, setPolicy] = React.useState<LLMRetryPolicy>(ZERO_POLICY);
  const [loading, setLoading] = React.useState(true);
  const [saving, setSaving] = React.useState(false);

  const load = React.useCallback(async () => {
    setLoading(true);
    try {
      const p = await api.llmRetryPolicy();
      setPolicy({ ...ZERO_POLICY, ...p });
    } catch (e) {
      toast.error(`Failed to load retry policy: ${(e as Error).message}`);
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    void load();
  }, [load]);

  async function save() {
    if (saving) return;
    setSaving(true);
    try {
      // The backend clamps out-of-range values back into the interval and returns them; refresh straight from the reply so what you see is what is stored.
      const saved = await api.saveLLMRetryPolicy(policy);
      setPolicy({ ...ZERO_POLICY, ...saved });
      toast.success(tr("已保存，立即生效（正在进行的调用仍使用旧参数）", "Saved; takes effect immediately (in-flight calls retain old parameters)"));
    } catch (e) {
      toast.error(tr(`保存失败：${(e as Error).message}`, `Save failed: ${(e as Error).message}`));
    } finally {
      setSaving(false);
    }
  }

  const set = (k: LayerKey) => (r: LLMRetryRule) => setPolicy((p) => ({ ...p, [k]: r }));

  if (loading) {
    return (
      <div className="flex items-center gap-2 rounded-lg border border-dashed p-10 text-muted-foreground text-sm">
        <Loader2Icon className="size-4 animate-spin" /> {tr("正在加载重试策略…", "Loading retry policy…")}
      </div>
    );
  }

  return (
    <div className="grid gap-4">
      <div className="rounded-lg border bg-muted/30 p-3 text-muted-foreground text-xs leading-relaxed">
        {tr("模型调用失败依次经过五层重试：连接 → 空响应 → 同供应商安全窗口 → 模型池熔断 → 意图重跑。外层仅处理内层已耗尽的故障，因此尝试次数相乘；各层次数都过高时，一次故障可能消耗数十次请求。全部留空保留现有默认行为，前三层可按模型配置覆盖。", "A failed model call passes through five layers: connect → empty response → same-provider safe window → pool breaker → intent rerun. Outer layers handle failures exhausted by inner layers, so attempt counts multiply; high counts can consume dozens of requests. Empty fields retain current defaults; the first three layers can be overridden per profile.")}
      </div>

      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        {(Object.keys(RETRY_LAYERS) as LayerKey[]).map((k) => (
          <RetryRuleFields key={k} layer={k} idPrefix="gl" value={policy[k]} onChange={set(k)} />
        ))}
      </div>

      <div className="flex gap-2">
        <Button onClick={save} disabled={saving}>
          {saving ? <Loader2Icon className="animate-spin" /> : <SaveIcon />}
          {tr("保存", "Save")}
        </Button>
        <Button variant="outline" onClick={() => setPolicy(ZERO_POLICY)} disabled={saving}>
          {tr("全部恢复默认值", "Reset all to defaults")}
        </Button>
      </div>
    </div>
  );
}
