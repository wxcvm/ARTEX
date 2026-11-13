"use client";

import * as React from "react";

import { CpuIcon, FlaskConicalIcon, KeyboardIcon, RadioTowerIcon, SearchIcon, ShieldAlertIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { api } from "@/lib/api";
import { CHAT_SEND_MODE_OPTIONS, type ChatSendMode, setChatSendMode, useChatSendMode } from "@/lib/chat-send-mode";
import type { Settings } from "@/lib/types";

export default function SystemSettingsPage() {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  const [trafficCapture, setTrafficCapture] = React.useState(false);
  const [agentTrafficBinding, setAgentTrafficBinding] = React.useState(false);
  const [webSearch, setWebSearch] = React.useState(false);
  const [backend, setBackend] = React.useState("ddgs");
  const [braveKeySet, setBraveKeySet] = React.useState(false);
  const [braveKeyInput, setBraveKeyInput] = React.useState("");
  const [tavilyKeySet, setTavilyKeySet] = React.useState(false);
  const [tavilyKeyInput, setTavilyKeyInput] = React.useState("");
  const [savingTavilyKey, setSavingTavilyKey] = React.useState(false);
  const [proxyInput, setProxyInput] = React.useState("");
  const [savingProxy, setSavingProxy] = React.useState(false);
  const [globalProxyInput, setGlobalProxyInput] = React.useState("");
  const [savingGlobalProxy, setSavingGlobalProxy] = React.useState(false);
  const [testing, setTesting] = React.useState(false);
  const [loaded, setLoaded] = React.useState(false);
  const [saving, setSaving] = React.useState(false);
  const [savingKey, setSavingKey] = React.useState(false);
  const [pyInterp, setPyInterp] = React.useState("");
  const [workers, setWorkers] = React.useState("3");
  const [savingWorkers, setSavingWorkers] = React.useState(false);
  // Operation constraint injection scope (both on by default).
  const [injectPlanner, setInjectPlanner] = React.useState(true);
  const [injectWorker, setInjectWorker] = React.useState(true);
  // Experimental: noa 上下文压缩 (off by default).
  const [noaCompaction, setNoaCompaction] = React.useState(false);
  // Pure frontend preference: bypasses /api/settings, reads and writes localStorage directly.
  const sendMode = useChatSendMode();

  const apply = React.useCallback((s: Settings) => {
    setTrafficCapture(!!s.traffic_capture);
    setAgentTrafficBinding(!!s.agent_traffic_binding);
    setWebSearch(!!s.web_search_enabled);
    setBackend(s.web_search_backend || "ddgs");
    setBraveKeySet(!!s.brave_key_set);
    setTavilyKeySet(!!s.tavily_key_set);
    setProxyInput(s.web_search_proxy ?? "");
    setGlobalProxyInput(s.global_proxy ?? "");
    setPyInterp(s.python_interpreter ?? "");
    setWorkers(String(s.workers ?? 3));
    setInjectPlanner(s.constraints_inject_planner !== false);
    setInjectWorker(s.constraints_inject_worker !== false);
    setNoaCompaction(!!s.noa_compaction);
  }, []);

  const saveWorkers = () => {
    const n = Number(workers);
    if (!Number.isInteger(n) || n <= 0) {
      toast.error(t("并发数必须为大于 0 的整数", "Concurrency must be a positive integer"));
      return;
    }
    setSavingWorkers(true);
    api
      .setSettings({ workers: n })
      .then((s) => {
        apply(s);
        toast.success(t("已保存并发执行智能体数量（对后续启动的任务生效）", "Saved worker concurrency (applies to new tasks)"));
      })
      .catch((e) => toast.error(t("保存失败：", "Save failed: ") + (e as Error).message))
      .finally(() => setSavingWorkers(false));
  };

  const savePython = () => {
    setSaving(true);
    api
      .setSettings({ python_interpreter: pyInterp.trim() })
      .then((s) => {
        apply(s);
        toast.success(t("已保存 Python 解释器配置", "Saved Python interpreter"));
      })
      .catch((e) => toast.error(t("保存失败：", "Save failed: ") + (e as Error).message))
      .finally(() => setSaving(false));
  };
  const detectPython = () => {
    setSaving(true);
    api
      .detectPython()
      .then((r) => setPyInterp(r.python_interpreter))
      .catch(() => undefined)
      .finally(() => setSaving(false));
  };

  React.useEffect(() => {
    api
      .settings()
      .then(apply)
      .catch(() => undefined)
      .finally(() => setLoaded(true));
  }, [apply]);

  const toggleTraffic = (v: boolean) => {
    setTrafficCapture(v); // optimistic
    setSaving(true);
    api
      .setSettings({ traffic_capture: v })
      .then(apply)
      .catch(() => setTrafficCapture(!v)) // revert on failure
      .finally(() => setSaving(false));
  };

  const toggleInjectPlanner = (v: boolean) => {
    setInjectPlanner(v); // optimistic
    api
      .setSettings({ constraints_inject_planner: v })
      .then(apply)
      .catch(() => setInjectPlanner(!v)); // revert on failure
  };

  const toggleAgentTrafficBinding = (v: boolean) => {
    setAgentTrafficBinding(v);
    setSaving(true);
    api
      .setSettings({ agent_traffic_binding: v })
      .then((s) => {
        apply(s);
        toast.success(v ? t("已启用智能体流量自动关联", "Enabled automatic agent traffic binding") : t("已关闭智能体流量自动关联", "Disabled automatic agent traffic binding"));
      })
      .catch((e) => {
        setAgentTrafficBinding(!v);
        toast.error(`${t("保存失败：", "Save failed: ")}${(e as Error).message}`);
      })
      .finally(() => setSaving(false));
  };

  const toggleInjectWorker = (v: boolean) => {
    setInjectWorker(v); // optimistic
    api
      .setSettings({ constraints_inject_worker: v })
      .then(apply)
      .catch(() => setInjectWorker(!v)); // revert on failure
  };

  const toggleNoaCompaction = (v: boolean) => {
    setNoaCompaction(v); // optimistic
    api
      .setSettings({ noa_compaction: v })
      .then((s) => {
        apply(s);
        toast.success(v ? t("已启用 noa 上下文压缩（对后续运行生效）", "Enabled noa context compaction (applies to new runs)") : t("已关闭 noa 上下文压缩（恢复内置压缩）", "Disabled noa context compaction (restored built-in compaction)"));
      })
      .catch((e) => {
        setNoaCompaction(!v); // revert on failure
        toast.error(`${t("保存失败：", "Save failed: ")}${(e as Error).message}`);
      });
  };

  // Persist a web-search patch (enable and/or backend). Optimistic with refetch.
  const saveWebSearch = (patch: Partial<Settings>) => {
    setSaving(true);
    api
      .setSettings(patch)
      .then((s) => {
        apply(s);
        toast.success(t("已保存网页搜索配置", "Saved web search config"));
      })
      .catch((e) => {
        toast.error(t("保存失败：", "Save failed: ") + (e as Error).message);
        api
          .settings()
          .then(apply)
          .catch(() => undefined);
      })
      .finally(() => setSaving(false));
  };

  const saveBraveKey = () => {
    setSavingKey(true);
    api
      .setSettings({ brave_search_api_key: braveKeyInput })
      .then((s) => {
        apply(s);
        setBraveKeyInput("");
        toast.success(t("已保存 Brave API Key", "Saved Brave API Key"));
      })
      .catch((e) => toast.error(t("保存失败：", "Save failed: ") + (e as Error).message))
      .finally(() => setSavingKey(false));
  };

  const saveTavilyKey = () => {
    setSavingTavilyKey(true);
    api
      .setSettings({ tavily_search_api_key: tavilyKeyInput })
      .then((s) => {
        apply(s);
        setTavilyKeyInput("");
        toast.success(t("已保存 Tavily API Key", "Saved Tavily API Key"));
      })
      .catch((e) => toast.error(t("保存失败：", "Save failed: ") + (e as Error).message))
      .finally(() => setSavingTavilyKey(false));
  };

  const saveProxy = () => {
    setSavingProxy(true);
    api
      .setSettings({ web_search_proxy: proxyInput.trim() })
      .then((s) => {
        apply(s);
        toast.success(proxyInput.trim() ? t("已保存出站代理", "Saved egress proxy") : t("已清除出站代理（直连）", "Cleared egress proxy (direct)"));
      })
      .catch((e) => toast.error(t("保存失败：", "Save failed: ") + (e as Error).message))
      .finally(() => setSavingProxy(false));
  };

  const saveGlobalProxy = () => {
    setSavingGlobalProxy(true);
    api
      .setSettings({ global_proxy: globalProxyInput.trim() })
      .then((s) => {
        apply(s);
        toast.success(globalProxyInput.trim() ? t("已保存全局代理", "Saved global proxy") : t("已清除全局代理（直连）", "Cleared global proxy (direct)"));
      })
      .catch((e) => toast.error(t("保存失败：", "Save failed: ") + (e as Error).message))
      .finally(() => setSavingGlobalProxy(false));
  };

  // Run a real "test" search ("test") against the CURRENT form values (backend +
  // proxy + entered key), falling back to saved values server-side. Toasts result.
  const runTest = () => {
    setTesting(true);
    api
      .testWebSearch({
        web_search_backend: backend,
        web_search_proxy: proxyInput.trim(),
        brave_search_api_key: braveKeyInput,
        tavily_search_api_key: tavilyKeyInput,
      })
      .then((r) => {
        if (r.ok) toast.success(t(`搜索测试成功 · ${r.backend} 返回 ${r.count} 条结果`, `Search test succeeded · ${r.backend} returned ${r.count} results`));
        else toast.error(t("搜索测试失败：", "Search test failed: ") + (r.error || t("未知错误", "Unknown error")));
      })
      .catch((e) => toast.error(t("搜索测试失败：", "Search test failed: ") + (e as Error).message))
      .finally(() => setTesting(false));
  };

  // brave-free selected but no key stored and none being entered → tool stays off.
  const braveNeedsKey = webSearch && backend === "brave-free" && !braveKeySet;

  return (
    <div className="flex min-w-0 flex-1 flex-col gap-5 pb-6">
      <header className="rounded-2xl border border-border/70 bg-gradient-to-br from-card to-primary/5 px-5 py-5 shadow-sm sm:px-7">
        <p className="mb-1 text-[11px] font-semibold tracking-[0.16em] text-primary">{t("ARTEX / 系统设置", "ARTEX / SETTINGS")}</p>
        <h1 className="text-2xl font-semibold tracking-tight">{t("系统设置", "System settings")}</h1>
        <p className="mt-1 max-w-3xl text-sm leading-relaxed text-muted-foreground">
          {t("集中管理网络、模型工具与工作台行为。每项设置的生效范围在对应卡片中说明。", "Manage networking, model tools and workspace behavior. Each card explains when its settings apply.")}
        </p>
      </header>

      {/* Multi-column instead of grid: the web search card is several times taller than the others and
          its height varies with the selected backend (brave/tavily key inputs are conditional). A grid
          would stretch the row to the tallest card and leave blank space beside the others, while columns
          balance-fill by content height. Card spacing uses mb rather than gap — in a multi-column layout
          column-gap only covers the gaps between columns; row spacing must come from the children. */}
      <div className="columns-1 gap-4 xl:columns-2 xl:gap-5">
        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              {t("流量采集", "Traffic capture")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>When enabled, every agent's HTTP traffic is stored in full through the recording proxy, and agents are injected with the traffic_search / traffic_get
              tools and proxy configuration (the prompt includes proxy notes).
              <br />
              When off (default) no traffic is recorded: the agent
              does <b>not</b> receive proxy configuration or traffic tools, and the prompt contains <b>no</b> proxy-related content. Agents are rebuilt immediately after switching.</> : <>开启后，每个智能体的 HTTP 流量通过记录代理完整保存，智能体获得 traffic_search / traffic_get 工具及代理配置。<br />关闭时（默认）不记录流量，不向智能体注入代理配置或流量工具，提示词也不包含代理相关内容。切换后立即重建智能体。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="traffic-capture" className="text-sm font-normal text-muted-foreground">
              {trafficCapture ? t("已开启 · 记录流量并注入代理", "On · records traffic and injects proxy") : t("已关闭 · 不记录流量，不注入代理", "Off · no recording or proxy injection")}
            </Label>
            <Switch
              id="traffic-capture"
              checked={trafficCapture}
              disabled={!loaded || saving}
              onCheckedChange={toggleTraffic}
            />
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              {t("自动关联智能体流量", "Automatic agent traffic binding")}
            </CardTitle>
            <CardDescription id="agent-traffic-binding-description">
              {locale === "en" ? <>Off by default. When enabled, the report agent triggered on vulnerability creation reviews existing HTTP requests/responses and links the matching traffic before writing the report.
              <b>Reading packets and the extra tool calls increase token consumption.</b>
              <br />
              Reporting still works for TCP, uncaptured traffic, or when no traffic matches. This switch does not affect traffic capture, manual binding, or viewing saved evidence. It takes effect for the
              next round of agents; once disabled, new automatic bindings are rejected immediately.</> : <>默认关闭。开启后，报告智能体在创建漏洞时检查已有 HTTP 请求和响应，并在撰写报告前关联匹配流量。<b>读取报文和额外工具调用会增加 Token 消耗。</b><br />TCP、未采集流量或无匹配流量时仍可生成报告。不影响流量采集、手动关联或查看已有证据；对下一轮智能体生效，关闭后立即拒绝新的自动关联。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="agent-traffic-binding" className="text-sm font-normal text-muted-foreground">
              {agentTrafficBinding ? t("已开启 · 增加 Token 消耗", "On · increases token usage") : t("已关闭 · 仍可手动关联", "Off · manual binding still available")}
            </Label>
            <Switch
              id="agent-traffic-binding"
              aria-describedby="agent-traffic-binding-description"
              checked={agentTrafficBinding}
              disabled={!loaded || saving}
              onCheckedChange={toggleAgentTrafficBinding}
            />
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              {t("全局代理", "Global proxy")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>All agents' <b>target traffic</b> egresses through this proxy (hide the source IP / use a jump host). Supports <b>http / https / socks5</b>, optionally with{" "}
              <code>user:pass</code> authentication. Empty = direct connection.
              <br />
              With <b>traffic capture</b> on, it acts as the recording proxy's <b>upstream</b> (traffic is still stored in full, then egresses through this proxy); with capture off, it is injected directly into
              the agent's bash / WebFetch egress. Independent of the web search proxy and the LLM proxy.
              <br />
              <b>Note</b>: with <b>capture off</b>, socks5 relies on each CLI tool honoring <code>ALL_PROXY</code> (curl
              works; some tools may ignore it). If you mainly use socks5, enable traffic capture — on that path the MITM
              does the dialing itself, invisible to tools and reliably effective.</> : <>所有智能体的<b>目标流量</b>通过此代理出站。支持 HTTP、HTTPS 和 SOCKS5，可选 <code>user:pass</code> 认证，留空直连。<br />开启流量采集时，此代理作为记录代理的上游；关闭采集时，直接用于智能体 Bash / WebFetch 出站。与网页搜索代理及模型代理相互独立。<br />注意：关闭采集时，SOCKS5 依赖各命令行工具支持 <code>ALL_PROXY</code>；部分工具可能忽略该变量。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-2">
            <Label htmlFor="global-proxy" className="text-sm font-normal text-muted-foreground">
              {t("代理地址", "Proxy address")}
            </Label>
            <div className="flex min-w-0 flex-wrap items-center gap-2">
              <Input
                className="min-w-56 flex-1"
                id="global-proxy"
                autoComplete="off"
                placeholder={t("socks5://user:pass@host:1080 或 http://host:port（留空直连）", "socks5://user:pass@host:1080 or http://host:port (empty = direct)")}
                value={globalProxyInput}
                disabled={!loaded || savingGlobalProxy}
                onChange={(e) => setGlobalProxyInput(e.target.value)}
              />
              <Button type="button" onClick={saveGlobalProxy} disabled={!loaded || savingGlobalProxy}>
                {t("保存", "Save")}
              </Button>
            </div>
            <p className="text-muted-foreground text-xs">
              {globalProxyInput.trim() ? t("已配置 · 目标流量经此代理出站", "Configured · target traffic uses this proxy") : t("未配置 · 目标流量直接连接", "Not configured · target traffic connects directly")}
            </p>
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <ShieldAlertIcon className="size-4" />
              {t("任务操作约束", "Task operation constraints")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>When enabled, each task's <b>operation constraints</b> (the allow/deny entries maintained under "Operation constraints" in the task overview) are spliced into the corresponding agent's
              system prompt to bound the exploration scope (e.g. "only test the current port", "no brute forcing").
              <br />
              You can control injection into the <b>planner</b> and the <b>worker</b> separately
              ; both are on by default. Switching takes effect immediately (read on the next turn) without rebuilding agents. When off, the agent no longer sees the constraints.</> : <>开启后，任务概览中的允许/拒绝操作约束会加入对应智能体的系统提示词，以限定探索范围。<br />可分别控制<b>规划器</b>和<b>执行智能体</b>，默认均开启。切换后下一轮读取，无需重建智能体；关闭后相应智能体不再看到这些约束。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4">
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="inject-planner" className="text-sm font-normal text-muted-foreground">
                {t("注入规划器", "Inject into planner")}{injectPlanner ? t(" · 开启", " · On") : t(" · 关闭", " · Off")}
              </Label>
              <Switch
                id="inject-planner"
                checked={injectPlanner}
                disabled={!loaded}
                onCheckedChange={toggleInjectPlanner}
              />
            </div>
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="inject-worker" className="text-sm font-normal text-muted-foreground">
                {t("注入执行智能体", "Inject into worker")}{injectWorker ? t(" · 开启", " · On") : t(" · 关闭", " · Off")}
              </Label>
              <Switch
                id="inject-worker"
                checked={injectWorker}
                disabled={!loaded}
                onCheckedChange={toggleInjectWorker}
              />
            </div>
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <FlaskConicalIcon className="size-4" />
              {t("实验功能", "Experimental features")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>Mechanisms still being validated, off by default. They may change agent behavior or affect stability; enable only after understanding the impact.
              <br />
              <b>noa context compaction</b>: the model actively compacts long conversation history (norma v0.4.0). When enabled, the four agent types integrated into the platform (
              <b>planner / worker / main agent / chat</b>) switch to noa-managed context in place of the built-in compaction,
              and the pre-compaction originals are archived under the task working directory for later review. Switching takes effect immediately (for runs started afterwards) without rebuilding agents;
              when disabled, the built-in compaction is restored immediately.</> : <>仍在验证的机制，默认关闭。可能影响智能体行为和稳定性，了解影响后再启用。<br /><b>noa 上下文压缩</b>：模型主动压缩长对话历史（norma v0.4.0）。开启后，规划器、执行智能体、主智能体及对话智能体使用 noa 管理上下文，压缩前原文归档到任务工作目录供查阅。切换对后续运行生效，无需重建智能体；关闭恢复内置压缩。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="noa-compaction" className="text-sm font-normal text-muted-foreground">
              {t("noa 上下文压缩", "noa context compaction")}{noaCompaction ? t(" · 开启", " · On") : t(" · 关闭", " · Off")}
            </Label>
            <Switch
              id="noa-compaction"
              checked={noaCompaction}
              disabled={!loaded}
              onCheckedChange={toggleNoaCompaction}
            />
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <SearchIcon className="size-4" />
              {t("网页搜索", "Web search")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>This is the web search <b>master switch + source config</b>. Only when it is enabled can <b>each agent's config</b> separately choose whether to enable
              <b>web_search</b> (returns only titles/links/snippets without fetching bodies; fetching is handled by WebFetch). Web search does <b>not</b> go through
              the recording proxy and is independent of traffic capture.
              <br />
              Sources: <b>DuckDuckGo (ddgs)</b> (no key needed), <b>Brave (free tier)</b> (requires a Brave API Key),{" "}
              <b>Tavily</b> (requires a Tavily API Key) or <b>DeepSeek</b> (reuses the current LLM config). When the master switch is off, the
              per-agent web search switches are unavailable.</> : <>{t("网页搜索", "Web search")}<b>总开关与来源配置</b>。开启后才可在各智能体配置中单独启用 web_search（返回标题、链接和摘要；正文由 WebFetch 获取）。网页搜索独立于流量采集，不经过记录代理。<br />支持 DuckDuckGo（无需密钥）、Brave 免费层、Tavily（需密钥）以及复用当前模型配置的 DeepSeek。总开关关闭时，各智能体的网页搜索开关不可用。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4">
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="web-search" className="text-sm font-normal text-muted-foreground">
                {webSearch ? t("总开关已开启 · 可按智能体单独启用", "On · can be enabled per agent") : t("已关闭 · 智能体无法启用网页搜索", "Off · agents cannot enable search")}
              </Label>
              <Switch
                id="web-search"
                checked={webSearch}
                disabled={!loaded || saving}
                onCheckedChange={(v) => {
                  setWebSearch(v); // optimistic
                  saveWebSearch({ web_search_enabled: v });
                }}
              />
            </div>

            {webSearch && (
              <div className="flex items-center justify-between gap-4">
                <Label className="text-sm font-normal text-muted-foreground">{t("搜索来源", "Search source")}</Label>
                <Select
                  value={backend}
                  disabled={!loaded || saving}
                  onValueChange={(v) => {
                    setBackend(v); // optimistic
                    saveWebSearch({ web_search_backend: v });
                  }}
                >
                  <SelectTrigger className="w-48 shrink-0">
                    <SelectValue placeholder={t("选择来源", "Choose a source")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="ddgs">{t("DuckDuckGo（ddgs · 免费，无需密钥）", "DuckDuckGo (ddgs · free, no key)")}</SelectItem>
                    <SelectItem value="brave-free">{t("Brave（免费层 · 需要密钥）", "Brave (free tier · requires key)")}</SelectItem>
                    <SelectItem value="tavily">{t("Tavily（需要密钥）", "Tavily (requires key)")}</SelectItem>
                    <SelectItem value="deepseek">{t("DeepSeek（官方）", "DeepSeek (official)")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            )}

            {webSearch && backend === "deepseek" && (
              <div className="border-border/60 bg-muted/30 flex flex-col gap-2 rounded-md border p-3">
                <p className="text-sm font-medium">{t("DeepSeek 官方网页搜索", "DeepSeek official web search")}</p>
                <p className="text-muted-foreground text-xs leading-relaxed">
                  {t("此来源直接复用当前启用的模型配置，仅支持 DeepSeek 官方模型，并且必须使用 Anthropic 协议。OpenAI 协议端点不支持服务端搜索；切换模型配置后可能失效。", "Reuses the active model configuration. Only official DeepSeek models using the Anthropic protocol are supported. OpenAI-protocol endpoints do not support server-side search; changing models may break it.")}
                </p>
                <p className="text-muted-foreground text-xs leading-relaxed">
                  {t("搜索由 DeepSeek 服务端执行，每次搜索会额外调用模型并消耗 Token。请求不经过上述出站代理，也不记录为流量。结果仅含标题及链接；需要正文时请使用 WebFetch。", "Search runs on DeepSeek servers and consumes an extra model call and tokens. It bypasses the egress proxy and is not recorded as traffic. Results include titles and links only; use WebFetch for page bodies.")}
                </p>
                <p className="text-muted-foreground text-xs leading-relaxed">
                  {t("请自行确认上述条件；系统不会阻止配置。使用下方「测试搜索」进行实际验证。", "Verify these conditions yourself; configuration is not blocked. Use Test search below to check.")}
                </p>
              </div>
            )}

            {webSearch && backend === "brave-free" && (
              <div className="flex flex-col gap-2">
                <Label htmlFor="brave-key" className="text-sm font-normal text-muted-foreground">
                  Brave Search API Key
                  {braveKeySet && <span className="ml-2 text-xs text-emerald-500">{t("已配置", "Configured")}</span>}
                </Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="brave-key"
                    type="password"
                    autoComplete="off"
                    placeholder={braveKeySet ? t("已配置（留空保留）", "Configured (leave empty to keep)") : t("输入 Brave API Key", "Enter Brave API Key")}
                    value={braveKeyInput}
                    disabled={!loaded || savingKey}
                    onChange={(e) => setBraveKeyInput(e.target.value)}
                  />
                  <Button
                    type="button"
                    onClick={saveBraveKey}
                    disabled={!loaded || savingKey || braveKeyInput.trim() === ""}
                  >
                    {t("保存", "Save")}
                  </Button>
                </div>
                {braveNeedsKey && (
                  <p className="text-xs text-amber-500">
                    {t("已选择 Brave，但未配置密钥。保存密钥后才能启用搜索工具。", "Brave is selected but no key is configured. Save a key to enable the search tool.")}
                  </p>
                )}
                <p className="text-muted-foreground text-xs">
                  {t("免费层每月约 2,000 次查询。可在 https://brave.com/search/api/ 获取密钥。", "The free tier allows about 2,000 queries per month. Get a key at https://brave.com/search/api/.")}
                </p>
              </div>
            )}

            {webSearch && backend === "tavily" && (
              <div className="flex flex-col gap-2">
                <Label htmlFor="tavily-key" className="text-sm font-normal text-muted-foreground">
                  Tavily Search API Key
                  {tavilyKeySet && <span className="ml-2 text-xs text-emerald-500">{t("已配置", "Configured")}</span>}
                </Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="tavily-key"
                    type="password"
                    autoComplete="off"
                    placeholder={tavilyKeySet ? t("已配置（留空保留）", "Configured (leave empty to keep)") : t("输入 Tavily API Key（tvly-…）", "Enter Tavily API Key (tvly-…)")}
                    value={tavilyKeyInput}
                    disabled={!loaded || savingTavilyKey}
                    onChange={(e) => setTavilyKeyInput(e.target.value)}
                  />
                  <Button
                    type="button"
                    onClick={saveTavilyKey}
                    disabled={!loaded || savingTavilyKey || tavilyKeyInput.trim() === ""}
                  >
                    {t("保存", "Save")}
                  </Button>
                </div>
                {webSearch && backend === "tavily" && !tavilyKeySet && (
                  <p className="text-xs text-amber-500">
                    {t("已选择 Tavily，但未配置密钥。保存密钥后才能启用搜索工具。", "Tavily is selected but no key is configured. Save a key to enable the search tool.")}
                  </p>
                )}
                <p className="text-muted-foreground text-xs">{t("在 https://tavily.com 注册并获取 API 密钥。", "Register and obtain an API key at https://tavily.com.")}</p>
              </div>
            )}

            {webSearch && (
              <div className="flex flex-col gap-2">
                <Label htmlFor="ws-proxy" className="text-sm font-normal text-muted-foreground">
                  {t("出站代理（可选）", "Egress proxy (optional)")}
                </Label>
                <div className="flex min-w-0 flex-wrap items-center gap-2">
                  <Input
                    className="min-w-56 flex-1"
                    id="ws-proxy"
                    autoComplete="off"
                    placeholder={t("http://host:port 或 socks5://host:port（留空直连）", "http://host:port or socks5://host:port (empty = direct)")}
                    value={proxyInput}
                    disabled={!loaded || savingProxy}
                    onChange={(e) => setProxyInput(e.target.value)}
                  />
                  <Button type="button" onClick={saveProxy} disabled={!loaded || savingProxy}>
                    {t("保存", "Save")}
                  </Button>
                </div>
                <p className="text-muted-foreground text-xs">
                  {t("仅用于访问搜索服务的独立出站代理（VPN/SOCKS 等），与流量记录代理无关；可在网络无法直连时使用。", "Dedicated egress proxy for search services (VPN/SOCKS etc.), independent of the traffic recording proxy. Use when direct access is unavailable.")}
                </p>
              </div>
            )}

            {webSearch && (
              <div className="flex items-center justify-between gap-4 border-t pt-4">
                <p className="text-muted-foreground text-xs">
                  {t("使用当前来源、代理及密钥执行一次实际测试搜索，验证配置是否可用。", "Run a real search using the current source, proxy and key to verify the configuration.")}
                </p>
                <Button
                  type="button"
                  variant="outline"
                  onClick={runTest}
                  disabled={!loaded || testing}
                  className="shrink-0"
                >
                  {testing ? t("测试中…", "Testing…") : t("测试搜索", "Test search")}
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              {t("自定义脚本 · Python 环境", "Custom scripts · Python environment")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>Custom <b>script</b> tools run Python with it. Auto-detected at startup (python3 first); you can enter an absolute path to a venv /
              specific version here; leave empty to auto-detect at runtime.</> : <>自定义 script 工具使用此 Python 解释器。启动时自动检测（优先 python3），也可输入虚拟环境或指定版本的绝对路径；留空则运行时自动检测。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            <div className="flex items-center gap-2">
              <Input
                className="font-mono text-sm"
                placeholder={t("/usr/bin/python3（留空自动检测）", "/usr/bin/python3 (empty = auto-detect)")}
                value={pyInterp}
                disabled={!loaded || saving}
                onChange={(e) => setPyInterp(e.target.value)}
              />
              <Button variant="outline" onClick={detectPython} disabled={!loaded || saving}>
                {t("重新检测", "Re-detect")}
              </Button>
              <Button onClick={savePython} disabled={!loaded || saving}>
                {t("保存", "Save")}
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <CpuIcon className="size-4" />
              {t("任务并发与工作智能体", "Task concurrency & workers")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>The number of worker agents running concurrently per task (default 3). Higher values mean more parallel probing and higher consumption. Changes
              <b>apply to tasks started afterwards</b>; running tasks are unaffected.</> : <>每个任务并发运行的执行智能体数量（默认 3）。数值越大，并行操作与资源消耗越多。<b>仅对后续启动的任务生效</b>，不影响正在运行的任务。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            <div className="flex items-center gap-2">
              <Input
                type="number"
                min={1}
                className="w-32 font-mono text-sm"
                placeholder="3"
                value={workers}
                disabled={!loaded || savingWorkers}
                onChange={(e) => setWorkers(e.target.value)}
              />
              <Button onClick={saveWorkers} disabled={!loaded || savingWorkers}>
                {t("保存", "Save")}
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <KeyboardIcon className="size-4" />
              {t("聊天发送快捷键", "Chat send shortcut")}
            </CardTitle>
            <CardDescription>
              {locale === "en" ? <>Shared by the main-agent chat inputs on the chat page and task details; takes effect immediately on selection, no saving needed.
              <br />
              This preference <b>lives only in this browser</b>, does not sync with the account, and must be set again after switching browsers or clearing site data.</> : <>对话页和任务详情的主智能体输入框共用此设置；选择后立即生效，无需保存。
              <br />
              此偏好<b>仅保存在当前浏览器</b>，不会随账户同步；更换浏览器或清除站点数据后需重新设置。</>}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="chat-send-mode" className="text-sm font-normal text-muted-foreground">
              {t("发送方式", "Send mode")}
            </Label>
            <Select value={sendMode} onValueChange={(v) => setChatSendMode(v as ChatSendMode)}>
              <SelectTrigger id="chat-send-mode" className="w-72">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {CHAT_SEND_MODE_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.value === "enter" ? t("Enter 发送，Shift+Enter 换行", "Enter to send, Shift+Enter for newline") : t("Ctrl+Enter 发送，Enter 换行", "Ctrl+Enter to send, Enter for newline")}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
