"use client";

import * as React from "react";
import { toast } from "sonner";
import { EyeIcon, GitCompareIcon, InfoIcon, PencilIcon, RotateCcwIcon, SaveIcon, Trash2Icon, XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Markdown } from "@/components/markdown";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { api } from "@/lib/api";
import { cn } from "@/lib/utils";
import { skillLabel, skillSummary, toolLabel } from "@/lib/resource-labels";
import type { Agent, AgentDetail, AgentTrigger, MCPServer, PromptVar, PromptVersion, Settings, SkillItem, Tool } from "@/lib/types";

// Traffic tools are host tools gated by the global traffic capture switch: bindable, but
// only usable when capture is on. Keep this list in sync with traffic.SeedToolMetas.
const TRAFFIC_TOOL_KEYS = new Set(["traffic_search", "traffic_get"]);

// AgentEditor is the tabbed editor for one agent, used inside the agents-page
// drawer (and reused full-page for deep links). Tabs: Config & Prompt / MCP / Skill /
// Tools. Config + prompt save as before; visibility + tool bindings toggle live.
export function AgentEditor({ agentKey, onSaved }: { agentKey: string; onSaved?: () => void }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [detail, setDetail] = React.useState<AgentDetail | null>(null);
  const [versions, setVersions] = React.useState<PromptVersion[]>([]);
  const [variables, setVariables] = React.useState<PromptVar[]>([]);
  const [mcp, setMcp] = React.useState<MCPServer[]>([]);
  const [skills, setSkills] = React.useState<SkillItem[]>([]);
  const [tools, setTools] = React.useState<Tool[]>([]);
  const [loaded, setLoaded] = React.useState(false);
  const [viewVer, setViewVer] = React.useState<PromptVersion | null>(null);
  const [diffVer, setDiffVer] = React.useState<PromptVersion | null>(null);

  const [prompt, setPrompt] = React.useState("");
  const [mcpVisible, setMcpVisible] = React.useState<number[]>([]);
  const [skillVisible, setSkillVisible] = React.useState<string[]>([]);
  const [preview, setPreview] = React.useState("");
  const [maxTurns, setMaxTurns] = React.useState("0");
  const [runSecs, setRunSecs] = React.useState("600");
  // "" = follow (unbound); otherwise a profile id string
  const [llmProfileId, setLlmProfileId] = React.useState("");
  const [llmProfiles, setLlmProfiles] = React.useState<NonNullable<AgentDetail["llm_profiles"]>>([]);
  const [webSearch, setWebSearch] = React.useState(false);
  const [interactiveShell, setInteractiveShell] = React.useState(false);
  const [wrapup, setWrapup] = React.useState("");
  const [wrapupDefault, setWrapupDefault] = React.useState("");
  const [wrapupTurns, setWrapupTurns] = React.useState("0");
  const [wrapupTurnsDefault, setWrapupTurnsDefault] = React.useState(5);
  // task-level timeout wrap-up prompt (worker/planner only)
  const [ttSupported, setTtSupported] = React.useState(false);
  const [ttWrapup, setTtWrapup] = React.useState("");
  const [ttWrapupDefault, setTtWrapupDefault] = React.useState("");
  const [ttTurns, setTtTurns] = React.useState("0");
  const [ttTurnsDefault, setTtTurnsDefault] = React.useState(5);
  const [settings, setSettings] = React.useState<Settings | null>(null);

  React.useEffect(() => {
    api.mcpServers().then(setMcp).catch(() => {});
    api.skills().then(setSkills).catch(() => {});
    api.tools().then(setTools).catch(() => {});
    api.settings().then(setSettings).catch(() => {});
  }, []);
  // global gates: traffic tools need traffic capture, web search needs the master switch.
  const captureOn = !!settings?.traffic_capture;
  const webSearchGlobalOn = !!settings?.web_search_enabled;

  const reload = React.useCallback(() => {
    api
      .getAgent(agentKey)
      .then((d) => {
        setDetail(d);
        setPrompt(d.prompt ?? "");
        setVariables(d.variables ?? []);
        setVersions(d.versions ?? []);
        setMcpVisible(d.visibility?.mcp ?? []);
        setSkillVisible(d.visibility?.skill ?? []);
        setMaxTurns(String(d.agent?.max_turns ?? 0));
        setRunSecs(String(d.agent?.run_seconds ?? 600));
        setLlmProfileId(d.agent?.llm_profile_id != null ? String(d.agent.llm_profile_id) : "");
        setLlmProfiles(d.llm_profiles ?? []);
        setWebSearch(!!d.agent?.web_search);
        setInteractiveShell(!!d.agent?.interactive_shell);
        setWrapup(d.wrapup_prompt ?? "");
        setWrapupDefault(d.wrapup_default ?? "");
        setWrapupTurns(String(d.wrapup_max_turns ?? 0));
        setWrapupTurnsDefault(d.wrapup_max_turns_default ?? 5);
        setTtSupported(!!d.task_timeout_wrapup_supported);
        setTtWrapup(d.task_timeout_wrapup_prompt ?? "");
        setTtWrapupDefault(d.task_timeout_wrapup_default ?? "");
        setTtTurns(String(d.task_timeout_wrapup_max_turns ?? 0));
        setTtTurnsDefault(d.task_timeout_wrapup_max_turns_default ?? 5);
      })
      .catch(() => setDetail(null))
      .finally(() => setLoaded(true));
  }, [agentKey]);
  React.useEffect(() => {
    reload();
  }, [reload]);

  async function doPreview() {
    try {
      const r = await api.previewAgentPrompt(agentKey, prompt);
      setPreview(r.error ? tr("渲染错误：", "Render error: ") + r.error : r.rendered);
    } catch (e) {
      setPreview(tr("预览失败：", "Preview failed: ") + (e as Error).message);
    }
  }
  async function savePrompt() {
    try {
      const r = await api.saveAgentPrompt(agentKey, prompt);
      toast.success(tr(`已保存为版本 v${r.version}`, `Saved as version v${r.version}`));
      reload();
      onSaved?.();
    } catch (e) {
      toast.error(tr("保存失败：", "Save failed: ") + (e as Error).message);
    }
  }
  async function resetPrompt() {
    try {
      const r = await api.resetAgentPrompt(agentKey);
      toast.success(tr(`已恢复内置默认配置（v${r.version}）`, `Restored built-in default (v${r.version})`));
      reload();
    } catch (e) {
      toast.error(tr("恢复失败：", "Restore failed: ") + (e as Error).message);
    }
  }
  async function saveWrapup() {
    try {
      const turns = Math.max(0, Math.floor(Number(wrapupTurns) || 0));
      await api.saveAgentWrapup(agentKey, wrapup, turns);
      toast.success(wrapup.trim() || turns > 0 ? tr("收尾配置已保存（下次运行生效）", "Wrap-up saved (next run)") : tr("已清空，将使用内置默认配置", "Cleared; using built-in default"));
      reload();
    } catch (e) {
      toast.error(tr("保存失败：", "Save failed: ") + (e as Error).message);
    }
  }
  async function resetWrapup() {
    try {
      await api.resetAgentWrapup(agentKey);
      toast.success(tr("已恢复内置默认配置", "Built-in default restored"));
      reload();
    } catch (e) {
      toast.error(tr("恢复失败：", "Restore failed: ") + (e as Error).message);
    }
  }
  async function saveTaskTimeoutWrapup() {
    try {
      const turns = Math.max(0, Math.floor(Number(ttTurns) || 0));
      await api.saveAgentTaskTimeoutWrapup(agentKey, ttWrapup, turns);
      toast.success(tr("任务超时收尾配置已保存（下次运行生效）", "Task timeout wrap-up saved (next run)"));
      reload();
    } catch (e) {
      toast.error(tr("保存失败：", "Save failed: ") + (e as Error).message);
    }
  }
  async function resetTaskTimeoutWrapup() {
    try {
      await api.resetAgentTaskTimeoutWrapup(agentKey);
      toast.success(tr("已恢复内置默认配置", "Built-in default restored"));
      reload();
    } catch (e) {
      toast.error(tr("恢复失败：", "Restore failed: ") + (e as Error).message);
    }
  }
  async function saveConfig() {
    try {
      // Only submit the fields this agent actually displays, to avoid overwriting hidden items (such as max_turns for goals) with defaults.
      const patch: Parameters<typeof api.saveAgentConfig>[1] = {
        llm_profile_id: llmProfileId === "" ? null : Number(llmProfileId),
      };
      if (showConfig) {
        patch.max_turns = Math.max(0, Math.floor(Number(maxTurns) || 0));
        patch.run_seconds = Math.max(0, Math.floor(Number(runSecs) || 0));
      }
      if (showWebSearch) patch.web_search = webSearch;
      if (showInteractiveShell) patch.interactive_shell = interactiveShell;
      await api.saveAgentConfig(agentKey, patch);
      toast.success(tr("运行配置已保存（立即生效）", "Run settings saved (effective immediately)"));
      reload();
    } catch (e) {
      toast.error(tr("保存失败：", "Save failed: ") + (e as Error).message);
    }
  }
  // applyVis optimistically updates, persists, and toasts success/failure. On
  // failure it reverts to the prior selection so the UI never lies about state.
  async function applyVis(nextMcp: number[], nextSkill: string[], okMsg: string) {
    const prevMcp = mcpVisible;
    const prevSkill = skillVisible;
    setMcpVisible(nextMcp);
    setSkillVisible(nextSkill);
    try {
      await api.setAgentVisibility(agentKey, nextMcp, nextSkill);
      toast.success(okMsg);
      onSaved?.(); // refresh the list so the card's MCP/Skill counts stay in sync
    } catch (e) {
      setMcpVisible(prevMcp);
      setSkillVisible(prevSkill);
      toast.error(tr("保存失败：", "Save failed: ") + (e as Error).message);
    }
  }
  function toggleMcp(id: number) {
    const on = mcpVisible.includes(id);
    const name = mcp.find((m) => m.id === id)?.name ?? String(id);
    applyVis(
      on ? mcpVisible.filter((x) => x !== id) : [...mcpVisible, id],
      skillVisible,
      `${on ? "Revoked" : "Granted"} MCP "${name}" visibility`,
    );
  }
  function toggleSkill(name: string) {
    const on = skillVisible.includes(name);
    applyVis(
      mcpVisible,
      on ? skillVisible.filter((x) => x !== name) : [...skillVisible, name],
      `${on ? "Revoked" : "Granted"} Skill "${name}" visibility`,
    );
  }
  async function toggleTool(t: Tool) {
    const on = t.agents.includes(agentKey);
    const nextAgents = on ? t.agents.filter((k) => k !== agentKey) : [...t.agents, agentKey];
    // optimistic update
    setTools((ts) => ts.map((x) => (x.key === t.key ? { ...x, agents: nextAgents } : x)));
    try {
      await api.saveTool(t.key, {
        description: t.description,
        schema: t.schema,
        agents: nextAgents,
        enabled: t.enabled,
      });
      toast.success(`${on ? "已解绑" : "已绑定"}工具「${t.key}」`);
      onSaved?.(); // refresh the list so the card's tool count stays in sync
    } catch (e) {
      toast.error(tr("保存工具绑定失败：", "Failed to save tool binding: ") + (e as Error).message);
      reload();
      api.tools().then(setTools).catch(() => {});
    }
  }

  if (loaded && !detail) {
    return <div className="text-muted-foreground p-6 text-center text-sm">{tr("未找到智能体：", "Agent not found: ")}{agentKey}</div>;
  }
  // config is meaningless for the conversational main agent and the fixed-budget
  // goals decomposer; every other agent (workers, custom assistants) honors it.
  const showConfig = agentKey !== "mainagent" && agentKey !== "goals";
  // web search applies to every conversational/executing agent except the one-shot
  // goals decomposer; it's gated by the global master switch.
  const showWebSearch = agentKey !== "goals";
  // interactive shell (persistent PTY session tool family) is likewise open to all agents except goals; no global gate.
  const showInteractiveShell = agentKey !== "goals";
  // Every agent (including goals/mainagent) runs on some LLM, so the "default model" binding is open to all agents.
  const showLLM = true;
  // triggers (P3) only attach to custom agents.
  const isCustom = !!detail && !detail.agent?.builtin;

  return (
    <Tabs defaultValue="prompt" className="flex min-h-0 flex-1 flex-col">
      <TabsList className="mx-4 mt-2 w-fit">
        <TabsTrigger value="prompt">{tr("配置与提示词", "Config & prompts")}</TabsTrigger>
        <TabsTrigger value="wrapup">{tr("收尾提示词", "Wrap-up prompt")}</TabsTrigger>
        <TabsTrigger value="mcp">MCP</TabsTrigger>
        <TabsTrigger value="skill">{tr("技能", "Skills")}</TabsTrigger>
        <TabsTrigger value="tools">{tr("工具", "Tools")}</TabsTrigger>
        {isCustom && <TabsTrigger value="triggers">{tr("触发器", "Triggers")}</TabsTrigger>}
      </TabsList>

      {/* config + prompt */}
      <TabsContent value="prompt" className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
        <div className="grid gap-4">
          {(showLLM || showConfig || showWebSearch || showInteractiveShell) && (
            <div className="grid gap-3 rounded-md border p-3">
              {showLLM && (
                <div className="grid gap-1.5">
                  <Label htmlFor="llm-profile" className="text-xs">{tr("默认模型（模型配置）", "Default model (profile)")}</Label>
                  <div className="flex flex-wrap items-center gap-3">
                    <Select value={llmProfileId || "__follow__"} onValueChange={(v) => setLlmProfileId(v === "__follow__" ? "" : v)}>
                      <SelectTrigger id="llm-profile" className="h-8 w-72">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="__follow__">{tr("跟随任务或全局启用模型", "Follow task or global active model")}</SelectItem>
                        {llmProfiles.map((p) => (
                          <SelectItem key={p.id} value={String(p.id)}>
                            {p.name} ({p.model}){p.is_default ? tr(" · 默认", " · default") : ""}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <span className="text-muted-foreground max-w-md text-xs">
                      {tr("为此智能体绑定固定模型配置（点击「保存配置」生效）。优先级：智能体绑定 > 任务/会话覆盖 > 全局启用模型。", "Pin a model profile for this agent (save to apply). Priority: agent binding > task/session override > global active model.")}
                    </span>
                  </div>
                </div>
              )}
              <div className="flex flex-wrap items-end gap-3">
                {showConfig && (
                  <>
                    <div className="grid gap-1.5">
                      <Label htmlFor="max-turns" className="text-xs">{tr("最大循环轮数（0 = 不限）", "Maximum turns (0 = unlimited)")}</Label>
                      <Input id="max-turns" type="number" min={0} className="h-8 w-32"
                        value={maxTurns} onChange={(e) => setMaxTurns(e.target.value)} />
                    </div>
                    <div className="grid gap-1.5">
                      <Label htmlFor="run-seconds" className="text-xs">{tr("运行时长（秒，0 = 不限）", "Run duration (seconds, 0 = unlimited)")}</Label>
                      <Input id="run-seconds" type="number" min={0} className="h-8 w-32"
                        value={runSecs} onChange={(e) => setRunSecs(e.target.value)} />
                    </div>
                  </>
                )}
                <Button size="sm" variant="outline" onClick={saveConfig}>
                  <SaveIcon /> {tr("保存配置", "Save config")}
                </Button>
              </div>
              {showWebSearch && (
                <div className="flex items-center gap-3 border-t pt-3">
                  <Switch
                    id="web-search"
                    checked={webSearch}
                    disabled={!webSearchGlobalOn}
                    onCheckedChange={setWebSearch}
                  />
                  <div className="grid gap-0.5">
                    <Label htmlFor="web-search" className="text-sm">{tr("网页搜索", "Web search")}</Label>
                    <span className="text-muted-foreground text-xs">
                      {webSearchGlobalOn
                        ? tr("为此智能体启用并保存后，web_search 可搜索互联网", "Enable and save for this agent to use web_search")
                        : tr("请先在系统设置中启用网页搜索并配置后端", "Enable web search and configure a backend in system settings first")}
                    </span>
                  </div>
                </div>
              )}
              {showInteractiveShell && (
                <div className="flex items-center gap-3 border-t pt-3">
                  <Switch
                    id="interactive-shell"
                    checked={interactiveShell}
                    onCheckedChange={setInteractiveShell}
                  />
                  <div className="grid gap-0.5">
                    <Label htmlFor="interactive-shell" className="text-sm">{tr("交互式 Shell", "Interactive shell")}</Label>
                    <span className="text-muted-foreground text-xs">
                      {tr("为此智能体启用并保存后，可使用持久 PTY 会话工具（shell_open/send/read/close/list）操作交互式程序，例如 SSH 或 REPL。", "Enable and save to use persistent PTY session tools (shell_open/send/read/close/list) with interactive programs such as SSH or REPL.")}
                    </span>
                  </div>
                </div>
              )}
            </div>
          )}

          <div className="grid gap-2">
            <Label className="text-muted-foreground text-xs">{tr("变量（点击插入占位符，渲染时替换为运行数据）", "Variables (click to insert placeholders, replaced by runtime data)")}</Label>
            <div className="flex flex-wrap gap-2">
              {variables.map((v) => (
                <Tooltip key={v.name}>
                  <TooltipTrigger asChild>
                    <button type="button" onClick={() => setPrompt((p) => `${p}{{.${v.name}}}`)}
                      className="hover:bg-muted inline-flex items-center gap-1 rounded-md border bg-muted/40 px-2 py-1 font-mono text-xs">
                      {`{{.${v.name}}}`}
                      <Badge variant="secondary" className="px-1 py-0 text-[10px]">{v.source}</Badge>
                    </button>
                  </TooltipTrigger>
                  <TooltipContent className="max-w-xs">
                    <p className="font-medium">{v.description}</p>
                    <p className="text-muted-foreground mt-1">{tr("示例：", "Example: ")}{v.example}</p>
                  </TooltipContent>
                </Tooltip>
              ))}
              {variables.length === 0 && <span className="text-muted-foreground text-xs">{tr("（暂无变量）", "(No variables)")}</span>}
            </div>
          </div>

          <Textarea className="font-mono text-xs" rows={16} value={prompt}
            placeholder={tr("留空使用内置默认提示词", "Leave empty to use built-in prompt")} onChange={(e) => setPrompt(e.target.value)} />

          <div className="flex flex-wrap gap-2">
            <Dialog>
              <DialogTrigger asChild>
                <Button variant="outline" size="sm" onClick={doPreview}>
                  <EyeIcon /> {tr("渲染预览", "Render preview")}
                </Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-2xl">
                <DialogHeader>
                  <DialogTitle>{tr("渲染预览", "Render preview")}</DialogTitle>
                  <DialogDescription>{tr("所有 {{.Var}} 已由后端替换为示例值。", "All {{.Var}} placeholders are replaced with example values by the backend.")}</DialogDescription>
                </DialogHeader>
                <div className="max-h-[60vh] overflow-auto rounded-md border bg-muted/30 p-3">
                  <Markdown text={preview} />
                </div>
              </DialogContent>
            </Dialog>
            <Button size="sm" onClick={savePrompt}>
              <SaveIcon /> {tr("保存为新版本", "Save as new version")}
            </Button>
            <Button variant="outline" size="sm" onClick={resetPrompt}>
              <RotateCcwIcon /> {tr("恢复默认", "Restore default")}
            </Button>
          </div>


          <Separator />
          <div className="grid gap-2">
            <Label className="text-muted-foreground text-xs">{tr("版本历史", "Version history")}</Label>
            <ul className="grid gap-1">
              {versions.map((ver, i) => (
                <li key={ver.version} className="flex items-center gap-2 rounded-md px-1 py-0.5 text-xs hover:bg-muted/50">
                  <span className="font-mono shrink-0">v{ver.version}</span>
                  {i === 0 && <Badge variant="secondary" className="px-1.5 py-0 shrink-0">{tr("当前", "Current")}</Badge>}
                  <span className="text-muted-foreground truncate flex-1">{ver.note}</span>
                  {ver.ts && (
                    <span className="text-muted-foreground/60 shrink-0 tabular-nums">
                      {new Date(ver.ts).toLocaleDateString(locale, { month: "2-digit", day: "2-digit", hour: "2-digit", minute: "2-digit" })}
                    </span>
                  )}
                  <Button variant="ghost" size="icon-sm" className="size-6 shrink-0" onClick={() => setViewVer(ver)}>
                    <EyeIcon className="size-3" />
                  </Button>
                  {i > 0 && versions[0] && (
                    <Button variant="ghost" size="icon-sm" className="size-6 shrink-0" onClick={() => setDiffVer(ver)}>
                      <GitCompareIcon className="size-3" />
                    </Button>
                  )}
                </li>
              ))}
              {versions.length === 0 && (
                <li className="text-muted-foreground text-xs">{tr("（暂无保存版本，正在使用内置默认配置）", "(No saved versions; using built-in defaults)")}</li>
              )}
            </ul>
          </div>

          {/* view-version dialog */}
          <Dialog open={!!viewVer} onOpenChange={(o) => { if (!o) setViewVer(null); }}>
            <DialogContent className="sm:max-w-2xl">
              <DialogHeader>
                <DialogTitle>
                  v{viewVer?.version}
                  {viewVer?.version === versions[0]?.version && (
                    <Badge variant="secondary" className="ml-2 px-1.5 py-0 align-middle">{tr("当前", "Current")}</Badge>
                  )}
                </DialogTitle>
                <DialogDescription>
                  {viewVer?.note || tr("（无备注）", "(no note)")}
                  {viewVer?.ts && (
                    <span className="ml-2 text-muted-foreground/60">
                      {new Date(viewVer.ts).toLocaleString(locale)}
                    </span>
                  )}
                </DialogDescription>
              </DialogHeader>
              <pre className="bg-muted max-h-[55vh] overflow-auto whitespace-pre-wrap rounded-md p-3 font-mono text-xs">
                {viewVer?.template_text || tr("（空）", "(empty)")}
              </pre>
              <div className="flex gap-2 justify-end">
                {viewVer && viewVer.version !== versions[0]?.version && (
                  <Button variant="outline" size="sm" onClick={() => {
                    if (viewVer) { setDiffVer(viewVer); setViewVer(null); }
                  }}>
                    <GitCompareIcon className="mr-1 size-3.5" /> {tr("与当前版本比较", "Compare with current version")}
                  </Button>
                )}
                <Button size="sm" onClick={() => {
                  if (viewVer) { setPrompt(viewVer.template_text); setViewVer(null); toast.success(`已将 v${viewVer.version} 加载到编辑器，点击「保存为新版本」确认`); }
                }}>
                  {tr("加载到编辑器", "Load into editor")}
                </Button>
              </div>
            </DialogContent>
          </Dialog>

          {/* version-diff dialog */}
          <Dialog open={!!diffVer} onOpenChange={(o) => { if (!o) setDiffVer(null); }}>
            <DialogContent className="sm:max-w-3xl">
              <DialogHeader>
                <DialogTitle>{tr(`版本差异：v${diffVer?.version} → v${versions[0]?.version}（当前）`, `Version diff: v${diffVer?.version} → v${versions[0]?.version} (current)`)}</DialogTitle>
                <DialogDescription>
                  <span className="inline-flex items-center gap-3 text-xs">
                    <span className="rounded bg-red-500/15 px-1.5 py-0.5 text-red-600 dark:text-red-400">{tr("- 删除", "- Removed")}</span>
                    <span className="rounded bg-green-500/15 px-1.5 py-0.5 text-green-600 dark:text-green-400">{tr("+ 新增", "+ Added")}</span>
                  </span>
                </DialogDescription>
              </DialogHeader>
              <DiffView oldText={diffVer?.template_text ?? ""} newText={versions[0]?.template_text ?? ""} />
            </DialogContent>
          </Dialog>
        </div>
      </TabsContent>

      {/* wrap-up prompt */}
      <TabsContent value="wrapup" className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
        <div className="grid gap-3">
          <p className="text-muted-foreground text-xs leading-relaxed">
            {tr("当智能体因超时或轮数耗尽而终止时，系统将注入此收尾提示词执行一轮收尾：", "When an agent stops due to timeout or exhausted turns, this prompt runs one wrap-up round:")}
            {tr("先回写已识别但尚未保存的内容，再输出一句总结，避免突然结束。留空使用内置默认配置。", "Save identified but unsaved information first, then summarize before ending. Leave empty for the built-in default.")}
          </p>
          <div className="flex items-center gap-2">
            <Label className="text-xs">{tr("收尾提示词正文", "Wrap-up prompt text")}</Label>
            {wrapup.trim() ? (
              <Badge variant="secondary" className="px-1.5 py-0">{tr("自定义", "Custom")}</Badge>
            ) : (
              <Badge variant="outline" className="px-1.5 py-0">{tr("使用内置默认配置", "Using built-in default")}</Badge>
            )}
          </div>
          <Textarea
            className="font-mono text-xs"
            rows={10}
            value={wrapup}
            placeholder={wrapupDefault || tr("留空使用内置默认收尾提示词", "Leave empty to use default wrap-up prompt")}
            onChange={(e) => setWrapup(e.target.value)}
          />
          <div className="grid gap-1.5">
            <Label htmlFor="wrapup-turns" className="text-xs">
              Wrap-up turns (max turns for the wrap-up phase itself; 0 = built-in default of {wrapupTurnsDefault} turns)
            </Label>
            <Input id="wrapup-turns" type="number" min={0} className="h-8 w-32"
              value={wrapupTurns} onChange={(e) => setWrapupTurns(e.target.value)} />
          </div>
          <div className="flex gap-2">
            <Button size="sm" onClick={saveWrapup}>{tr("保存", "Save")}</Button>
            <Button size="sm" variant="outline" onClick={resetWrapup}>{tr("恢复默认", "Restore default")}</Button>
          </div>
          {wrapupDefault && (
            <>
              <Separator />
              <div className="grid gap-1.5">
                <Label className="text-muted-foreground text-xs">{tr("内置默认值（只读参考）", "Built-in default (read-only reference)")}</Label>
                <pre className="text-muted-foreground max-h-40 overflow-y-auto rounded-md border bg-muted/30 p-2 text-xs whitespace-pre-wrap">
                  {wrapupDefault}
                </pre>
              </div>
            </>
          )}

          {ttSupported && (
            <>
              <Separator className="my-2" />
              <p className="text-muted-foreground text-xs leading-relaxed">
                {tr("任务超时收尾（独立于单次运行收尾）：当整个任务到达超时上限即将结束时注入。", "Task timeout wrap-up (independent of run wrap-up): injected when the entire task reaches its deadline.")}
                {tr("其含义通常与单次运行收尾不同，例如规划器的单次运行要求继续规划，而任务超时要求停止并作出最终判断。留空使用内置默认配置。", "This differs from run wrap-up: the planner may continue planning after a run, but at task timeout it must stop and make a final judgment. Leave empty for default.")}
              </p>
              <div className="flex items-center gap-2">
                <Label className="text-xs">{tr("任务超时收尾提示词正文", "Task timeout wrap-up prompt")}</Label>
                {ttWrapup.trim() ? (
                  <Badge variant="secondary" className="px-1.5 py-0">{tr("自定义", "Custom")}</Badge>
                ) : (
                  <Badge variant="outline" className="px-1.5 py-0">{tr("使用内置默认配置", "Using built-in default")}</Badge>
                )}
              </div>
              <Textarea
                className="font-mono text-xs"
                rows={10}
                value={ttWrapup}
                placeholder={ttWrapupDefault || tr("留空使用内置默认任务超时收尾提示词", "Leave empty to use default task timeout wrap-up")}
                onChange={(e) => setTtWrapup(e.target.value)}
              />
              <div className="grid gap-1.5">
                <Label htmlFor="tt-turns" className="text-xs">
                  Wrap-up turns (0 = built-in default of {ttTurnsDefault} turns)
                </Label>
                <Input id="tt-turns" type="number" min={0} className="h-8 w-32"
                  value={ttTurns} onChange={(e) => setTtTurns(e.target.value)} />
              </div>
              <div className="flex gap-2">
                <Button size="sm" onClick={saveTaskTimeoutWrapup}>{tr("保存", "Save")}</Button>
                <Button size="sm" variant="outline" onClick={resetTaskTimeoutWrapup}>{tr("恢复默认", "Restore default")}</Button>
              </div>
              {ttWrapupDefault && (
                <div className="grid gap-1.5">
                  <Label className="text-muted-foreground text-xs">{tr("内置默认值（只读参考）", "Built-in default (read-only reference)")}</Label>
                  <pre className="text-muted-foreground max-h-40 overflow-y-auto rounded-md border bg-muted/30 p-2 text-xs whitespace-pre-wrap">
                    {ttWrapupDefault}
                  </pre>
                </div>
              )}
            </>
          )}
        </div>
      </TabsContent>

      {/* MCP visibility */}
      <TabsContent value="mcp" className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
        <p className="text-muted-foreground mb-3 text-xs">{tr("查看此智能体可见的 MCP 服务。", "View MCP servers visible to this agent.")}</p>
        <div className="grid gap-2">
          {mcp.map((m) => (
            <label key={m.id} className="flex items-center gap-2 rounded-md border p-2 text-sm">
              <Checkbox checked={mcpVisible.includes(m.id)} onCheckedChange={() => toggleMcp(m.id)} />
              {m.name}
              <span className="text-muted-foreground ml-auto text-xs">{m.transport}</span>
            </label>
          ))}
          {mcp.length === 0 && <span className="text-muted-foreground text-xs">{tr("（暂无 MCP 服务）", "(No MCP servers)")}</span>}
        </div>
      </TabsContent>

      {/* Skill visibility */}
      <TabsContent value="skill" className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
        <p className="text-muted-foreground mb-3 text-xs">{tr("查看此智能体可见的技能。", "View skills visible to this agent.")}</p>
        <div className="grid gap-2">
          {skills.map((s) => (
            <label key={s.name} className="flex items-center gap-2 rounded-md border p-2 text-sm">
              <Checkbox checked={skillVisible.includes(s.name)} onCheckedChange={() => toggleSkill(s.name)} />
              <span className="text-xs" title={s.name}>{skillLabel(s.name)}</span>
              {(skillSummary(s.name) || s.description) && <span className="text-muted-foreground ml-auto truncate text-xs" title={s.description}>{skillSummary(s.name) || s.description}</span>}
            </label>
          ))}
          {skills.length === 0 && <span className="text-muted-foreground text-xs">{tr("（暂无技能）", "(No skills)")}</span>}
        </div>
      </TabsContent>

      {/* Tools bindings */}
      <TabsContent value="tools" className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
        <p className="text-muted-foreground mb-3 text-xs">{tr("查看绑定到此智能体的内置工具。", "View built-in tools bound to this agent.")}</p>
        <div className="grid gap-2">
          {tools.map((t) => {
            const isTraffic = TRAFFIC_TOOL_KEYS.has(t.key);
            const gated = isTraffic && !captureOn; // traffic tools need traffic capture on
            return (
              <label
                key={t.key}
                className={cn(
                  "flex items-center gap-2 rounded-md border p-2 text-sm",
                  gated && "opacity-60",
                )}
              >
                <Checkbox
                  checked={t.agents.includes(agentKey)}
                  disabled={gated}
                  onCheckedChange={() => toggleTool(t)}
                />
                <span className="text-xs" title={t.key}>{toolLabel(t.key, t.system)}</span>
                {isTraffic && (
                  <Badge variant="secondary" className="px-1 py-0 text-[9px]">{tr("流量", "Traffic")}</Badge>
                )}
                {!t.enabled && (
                  <Badge variant="outline" className="text-destructive px-1 py-0 text-[9px]">{tr("已停用", "Disabled")}</Badge>
                )}
                {gated ? (
                  <span className="text-muted-foreground ml-auto text-xs">{tr("需要流量采集", "Requires traffic capture")}</span>
                ) : (
                  t.description && (
                    <span className="text-muted-foreground ml-auto line-clamp-1 max-w-[55%] text-xs">
                      {t.description}
                    </span>
                  )
                )}
              </label>
            );
          })}
          {tools.length === 0 && <span className="text-muted-foreground text-xs">{tr("（暂无工具）", "(No tools)")}</span>}
        </div>
      </TabsContent>

      {/* Triggers (P3, custom agents only) */}
      {isCustom && (
        <TabsContent value="triggers" className="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
          <AgentTriggersTab agentKey={agentKey} agent={detail?.agent} />
        </TabsContent>
      )}
    </Tabs>
  );
}

// ---------- Diff helpers ----------

type DiffLine = { type: "same" | "add" | "del"; text: string };

function computeDiff(oldText: string, newText: string): DiffLine[] {
  const a = oldText.split("\n");
  const b = newText.split("\n");
  const m = a.length;
  const n = b.length;
  const dp: number[][] = Array.from({ length: m + 1 }, () => new Array(n + 1).fill(0));
  for (let i = 1; i <= m; i++)
    for (let j = 1; j <= n; j++)
      dp[i][j] = a[i - 1] === b[j - 1] ? dp[i - 1][j - 1] + 1 : Math.max(dp[i - 1][j], dp[i][j - 1]);
  const result: DiffLine[] = [];
  let i = m;
  let j = n;
  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && a[i - 1] === b[j - 1]) {
      result.unshift({ type: "same", text: a[i - 1] });
      i--;
      j--;
    } else if (j > 0 && (i === 0 || dp[i][j - 1] >= dp[i - 1][j])) {
      result.unshift({ type: "add", text: b[j - 1] });
      j--;
    } else {
      result.unshift({ type: "del", text: a[i - 1] });
      i--;
    }
  }
  return result;
}

function DiffView({ oldText, newText }: { oldText: string; newText: string }) {
  const lines = React.useMemo(() => computeDiff(oldText, newText), [oldText, newText]);
  return (
    <pre className="max-h-[60vh] overflow-auto rounded-md border bg-muted/30 p-2 font-mono text-xs leading-5">
      {lines.map((l, idx) => (
        <div
          key={idx}
          className={cn(
            "whitespace-pre-wrap px-1",
            l.type === "del" && "bg-red-500/15 text-red-700 dark:text-red-400",
            l.type === "add" && "bg-green-500/15 text-green-700 dark:text-green-400",
            l.type === "same" && "text-muted-foreground",
          )}
        >
          <span className="select-none mr-1 opacity-50">{l.type === "del" ? "-" : l.type === "add" ? "+" : " "}</span>
          {l.text}
        </div>
      ))}
    </pre>
  );
}

// AgentTriggersTab manages a custom agent's P3 triggers: list + add + delete.
// Each trigger fires (schedule/finding found/goal met/task timeout/tool call, multi-select) → a new conversation runs
// in parallel with the base user message + auto context appended by the backend.
function AgentTriggersTab({ agentKey, agent }: { agentKey: string; agent?: Agent }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [triggers, setTriggers] = React.useState<AgentTrigger[]>([]);
  const [tools, setTools] = React.useState<Tool[]>([]);
  // post-trigger handling policy (per agent); initial value comes from agent detail, saved on change.
  const [runMode, setRunMode] = React.useState<"serial" | "parallel">(agent?.trigger_run_mode ?? "serial");
  const [mergeMode, setMergeMode] = React.useState<"by_task" | "all" | "none">(agent?.trigger_merge_mode ?? "all");
  const [maxParallel, setMaxParallel] = React.useState(String(agent?.trigger_max_parallel ?? 5));
  React.useEffect(() => {
    setRunMode(agent?.trigger_run_mode ?? "serial");
    setMergeMode(agent?.trigger_merge_mode ?? "all");
    setMaxParallel(String(agent?.trigger_max_parallel ?? 5));
  }, [agent?.trigger_run_mode, agent?.trigger_merge_mode, agent?.trigger_max_parallel]);

  async function saveBehavior(patch: {
    trigger_run_mode?: "serial" | "parallel";
    trigger_merge_mode?: "by_task" | "all" | "none";
    trigger_max_parallel?: number;
  }) {
    try {
      await api.saveAgentConfig(agentKey, patch);
    } catch (e) {
      toast.error(tr("保存策略失败：", "Failed to save policy: ") + (e as Error).message);
    }
  }
  const [onInterval, setOnInterval] = React.useState(false);
  const [intervalSec, setIntervalSec] = React.useState("60");
  const [onFinding, setOnFinding] = React.useState(false);
  const [onGoalMet, setOnGoalMet] = React.useState(false);
  const [onTaskTimeout, setOnTaskTimeout] = React.useState(false);
  const [onToolCall, setOnToolCall] = React.useState(false);
  const [onTaskCreate, setOnTaskCreate] = React.useState(false);
  const [intervalMsg, setIntervalMsg] = React.useState("");
  const [findingMsg, setFindingMsg] = React.useState("");
  const [goalMsg, setGoalMsg] = React.useState("");
  const [taskTimeoutMsg, setTaskTimeoutMsg] = React.useState("");
  const [toolCallMsg, setToolCallMsg] = React.useState("");
  const [taskCreateMsg, setTaskCreateMsg] = React.useState("");
  const [toolNames, setToolNames] = React.useState<string[]>([]);
  const [saving, setSaving] = React.useState(false);
  // null = create mode; non-null = editing the trigger with that id.
  const [editingId, setEditingId] = React.useState<number | null>(null);

  const reload = React.useCallback(() => {
    api.agentTriggers(agentKey).then(setTriggers).catch(() => setTriggers([]));
  }, [agentKey]);
  React.useEffect(() => {
    reload();
  }, [reload]);
  React.useEffect(() => {
    api.tools().then(setTools).catch(() => setTools([]));
  }, []);

  function toggleTool(key: string) {
    setToolNames((prev) => (prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key]));
  }

  // resetForm clears the form and returns to "create" mode.
  function resetForm() {
    setEditingId(null);
    setOnInterval(false);
    setIntervalSec("60");
    setOnFinding(false);
    setOnGoalMet(false);
    setOnTaskTimeout(false);
    setOnToolCall(false);
    setOnTaskCreate(false);
    setIntervalMsg("");
    setFindingMsg("");
    setGoalMsg("");
    setTaskTimeoutMsg("");
    setToolCallMsg("");
    setTaskCreateMsg("");
    setToolNames([]);
  }

  // startEdit fills the form with an existing trigger and enters "edit" mode.
  function startEdit(t: AgentTrigger) {
    setEditingId(t.id);
    setOnInterval(t.interval_sec > 0);
    setIntervalSec(t.interval_sec > 0 ? String(t.interval_sec) : "60");
    setOnFinding(t.on_finding);
    setOnGoalMet(t.on_goal_met);
    setOnTaskTimeout(t.on_task_timeout);
    setOnToolCall(t.on_tool_call);
    setOnTaskCreate(t.on_task_create);
    setIntervalMsg(t.interval_message);
    setFindingMsg(t.finding_message);
    setGoalMsg(t.goal_message);
    setTaskTimeoutMsg(t.task_timeout_message);
    setToolCallMsg(t.tool_call_message);
    setTaskCreateMsg(t.task_create_message);
    setToolNames(t.tool_names ?? []);
  }

  // submit creates or saves changes based on editingId; editing keeps the trigger's enabled state.
  async function submit() {
    const n = onInterval ? Math.max(1, Math.floor(Number(intervalSec) || 0)) : 0;
    if (n === 0 && !onFinding && !onGoalMet && !onTaskTimeout && !onToolCall && !onTaskCreate) {
      toast.error(tr("至少选择一个触发条件", "Select at least one trigger condition"));
      return;
    }
    if (onToolCall && toolNames.length === 0) {
      toast.error(tr("工具调用触发器至少需要选择一个工具", "Select at least one tool for the tool-call trigger"));
      return;
    }
    const body = {
      interval_sec: n,
      on_finding: onFinding,
      on_goal_met: onGoalMet,
      on_task_timeout: onTaskTimeout,
      on_tool_call: onToolCall,
      on_task_create: onTaskCreate,
      interval_message: intervalMsg.trim(),
      finding_message: findingMsg.trim(),
      goal_message: goalMsg.trim(),
      task_timeout_message: taskTimeoutMsg.trim(),
      tool_call_message: toolCallMsg.trim(),
      task_create_message: taskCreateMsg.trim(),
      tool_names: onToolCall ? toolNames : [],
    };
    setSaving(true);
    try {
      if (editingId != null) {
        const cur = triggers.find((x) => x.id === editingId);
        await api.updateTrigger(editingId, { ...body, enabled: cur?.enabled ?? true });
        toast.success(tr("已保存修改", "Changes saved"));
      } else {
        await api.createTrigger(agentKey, { ...body, enabled: true });
        toast.success(tr("已添加触发器", "Trigger added"));
      }
      resetForm();
      reload();
    } catch (e) {
      toast.error((editingId != null ? tr("保存失败：", "Save failed: ") : tr("添加失败：", "Add failed: ")) + (e as Error).message);
    } finally {
      setSaving(false);
    }
  }
  async function toggleEnabled(t: AgentTrigger) {
    try {
      await api.updateTrigger(t.id, {
        enabled: !t.enabled,
        interval_sec: t.interval_sec,
        on_finding: t.on_finding,
        on_goal_met: t.on_goal_met,
        on_task_timeout: t.on_task_timeout,
        on_tool_call: t.on_tool_call,
        on_task_create: t.on_task_create,
        interval_message: t.interval_message,
        finding_message: t.finding_message,
        goal_message: t.goal_message,
        task_timeout_message: t.task_timeout_message,
        tool_call_message: t.tool_call_message,
        task_create_message: t.task_create_message,
        tool_names: t.tool_names,
      });
      reload();
    } catch (e) {
      toast.error(tr("保存失败：", "Save failed: ") + (e as Error).message);
    }
  }
  async function del(id: number) {
    try {
      await api.deleteTrigger(id);
      if (editingId === id) resetForm();
      reload();
    } catch (e) {
      toast.error(tr("删除失败：", "Delete failed: ") + (e as Error).message);
    }
  }

  function condLabel(t: AgentTrigger): string {
    const parts: string[] = [];
    if (t.interval_sec > 0) parts.push(tr(`每 ${t.interval_sec} 秒`, `Every ${t.interval_sec}s`));
    if (t.on_finding) parts.push(tr("发现结果", "Finding"));
    if (t.on_goal_met) parts.push(tr("目标达成", "Goal met"));
    if (t.on_task_timeout) parts.push(tr("任务超时", "Task timeout"));
    if (t.on_tool_call) parts.push(tr(`工具调用（${t.tool_names.length}）`, `Tool calls (${t.tool_names.length})`));
    if (t.on_task_create) parts.push(tr("任务创建", "Task creation"));
    return parts.join(" · ") || tr("（无条件）", "(no conditions)");
  }

  return (
    <div className="grid gap-4">
      <p className="text-muted-foreground text-xs">
        {tr("触发器让自定义智能体自动运行，每次触发开启新对话（在对话页可见）。", "Triggers run this custom agent automatically; each starts a new conversation on the Chat page.")}
        {tr("可选择多种触发条件，系统会在你填写的基础消息中自动附加触发原因及相关任务、发现或目标。", "Select multiple conditions; the system appends the trigger reason and related task, finding or goal to your message.")}
      </p>

      {/* post-trigger handling policy */}
      <div className="grid gap-3 rounded-md border p-3">
        <Label className="text-muted-foreground text-xs">{tr("触发后处理策略（触发器如何排队或合并执行）", "Post-trigger policy (queueing and merging)")}</Label>
        <div className="flex flex-wrap items-center gap-4">
          <div className="grid gap-1">
            <Label className="text-xs">{tr("运行模式", "Run mode")}</Label>
            <Select
              value={runMode}
              onValueChange={(v) => {
                const rm = v as "serial" | "parallel";
                setRunMode(rm);
                saveBehavior({ trigger_run_mode: rm });
              }}
            >
              <SelectTrigger size="sm" className="h-8 w-40">
                <SelectValue />
              </SelectTrigger>
              <SelectContent position="popper">
                <SelectItem value="serial">{tr("串行（排队，逐个执行）", "Serial (queued, one at a time)")}</SelectItem>
                <SelectItem value="parallel">{tr("并行（独立并发会话）", "Parallel (independent concurrent sessions)")}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-1">
            <Label className="text-xs">{tr("合并模式", "Merge mode")}</Label>
            <Select
              value={mergeMode}
              disabled={runMode === "parallel"}
              onValueChange={(v) => {
                const mm = v as "by_task" | "all" | "none";
                setMergeMode(mm);
                saveBehavior({ trigger_merge_mode: mm });
              }}
            >
              <SelectTrigger size="sm" className="h-8 w-44">
                <SelectValue />
              </SelectTrigger>
              <SelectContent position="popper">
                <SelectItem value="by_task">{tr("按任务合并", "Merge by task")}</SelectItem>
                <SelectItem value="all">{tr("全部合并为一个", "Merge all")}</SelectItem>
                <SelectItem value="none">{tr("不合并", "No merge")}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {runMode === "parallel" && (
            <div className="grid gap-1">
              <Label htmlFor="tr-maxpar" className="text-xs">{tr("最大并发数（0 = 不限）", "Maximum concurrency (0 = unlimited)")}</Label>
              <Input
                id="tr-maxpar"
                type="number"
                min={0}
                className="h-8 w-28"
                value={maxParallel}
                onChange={(e) => setMaxParallel(e.target.value)}
                onBlur={() => {
                  const n = Math.max(0, Math.floor(Number(maxParallel) || 0));
                  setMaxParallel(String(n));
                  saveBehavior({ trigger_max_parallel: n });
                }}
              />
            </div>
          )}
        </div>
        <p className="text-muted-foreground text-xs">
          {runMode === "parallel"
            ? tr("并行：每次触发启动独立会话，不合并；超出并发上限的触发器排队等待。", "Parallel: each trigger immediately starts its own concurrent session, no merging; triggers beyond the max concurrency queue for a free slot.")
            : mergeMode === "by_task"
              ? tr("串行 · 按任务合并：每个智能体一次运行一个会话；同一任务的排队事件合并。", "Serial · merge by task: one run at a time per agent; queued events for the same task merge into one conversation.")
              : mergeMode === "all"
                ? tr("串行 · 全部合并：每个智能体一次运行一个会话；出队时将当前排队事件合并为一个对话。", "Serial · merge all: one run at a time per agent; when dequeuing, all currently queued triggers merge into one conversation.")
                : tr("串行 · 不合并：每次触发单独开启对话，逐个执行。", "Serial · no merge: one run at a time per agent; each trigger gets its own conversation.")}
        </p>
      </div>

      {/* create / edit trigger */}
      <div className="grid gap-3 rounded-md border p-3">
        <Label className="text-muted-foreground text-xs">
          {editingId != null
            ? tr(`编辑触发器 #${editingId}（完成后点击保存修改）`, `Edit trigger #${editingId} (click Save changes when done)`)
            : tr("新建触发器（每种条件可设置独立消息）", "New trigger (each condition can have its own user message)")}
        </Label>

        {/* schedule */}
        <div className="grid gap-1.5">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox checked={onInterval} onCheckedChange={(v) => setOnInterval(!!v)} /> {tr("定时触发", "Scheduled trigger")}
          </label>
          {onInterval && (
            <div className="grid gap-1.5">
              <div className="flex items-center gap-2">
                <Label htmlFor="tr-interval" className="text-xs">{tr("每隔", "Every")}</Label>
                <Input id="tr-interval" type="number" min={1} className="h-8 w-24"
                  value={intervalSec} onChange={(e) => setIntervalSec(e.target.value)} />
                <span className="text-muted-foreground text-xs">{tr("秒", "seconds")}</span>
              </div>
              <Textarea className="text-xs" rows={2} value={intervalMsg}
                placeholder={tr("定时触发时发送给智能体的消息，例如检查所有任务", "Scheduled message, e.g. check all tasks")} onChange={(e) => setIntervalMsg(e.target.value)} />
            </div>
          )}
        </div>

        {/* finding */}
        <div className="grid gap-1.5">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox checked={onFinding} onCheckedChange={(v) => setOnFinding(!!v)} /> {tr("发现结果时触发", "Trigger on finding")}
          </label>
          {onFinding && (
            <Textarea className="text-xs" rows={2} value={findingMsg}
              placeholder={tr("发现结果时发送给智能体的消息（系统附加任务和发现详情）", "Message on finding (task and finding details appended)")} onChange={(e) => setFindingMsg(e.target.value)} />
          )}
        </div>

        {/* goal met */}
        <div className="grid gap-1.5">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox checked={onGoalMet} onCheckedChange={(v) => setOnGoalMet(!!v)} /> {tr("目标达成时触发", "Trigger on goal met")}
          </label>
          {onGoalMet && (
            <Textarea className="text-xs" rows={2} value={goalMsg}
              placeholder={tr("目标达成时发送给智能体的消息（系统附加任务和已达成目标）", "Message on goal met (task and achieved goal appended)")} onChange={(e) => setGoalMsg(e.target.value)} />
          )}
        </div>

        {/* task timeout */}
        <div className="grid gap-1.5">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox checked={onTaskTimeout} onCheckedChange={(v) => setOnTaskTimeout(!!v)} /> {tr("任务超时时触发", "Trigger on task timeout")}
          </label>
          {onTaskTimeout && (
            <Textarea className="text-xs" rows={2} value={taskTimeoutMsg}
              placeholder={tr("任务超时时发送给智能体的消息（系统附加任务 ID 和目标）", "Message on task timeout (task ID and goal appended)")} onChange={(e) => setTaskTimeoutMsg(e.target.value)} />
          )}
        </div>

        {/* tool call */}
        <div className="grid gap-1.5">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox checked={onToolCall} onCheckedChange={(v) => setOnToolCall(!!v)} /> {tr("工具调用时触发", "Trigger on tool call")}
          </label>
          {onToolCall && (
            <div className="grid gap-1.5">
              <div className="text-muted-foreground text-xs">
                {tr(`选择要监听的工具（至少一个）；任务执行中每次这些工具调用完成时触发。已选择 ${toolNames.length} 个。`, `Select tools to watch (at least one); each completed call triggers an event. ${toolNames.length} selected.`)}
              </div>
              <div className="max-h-40 overflow-y-auto rounded-md border p-2">
                {tools.length === 0 && <span className="text-muted-foreground text-xs">{tr("（工具列表为空）", "(No tools available)")}</span>}
                <div className="grid gap-1">
                  {tools.map((tool) => (
                    <label key={tool.key} className="flex items-start gap-2 text-xs">
                      <Checkbox className="mt-0.5" checked={toolNames.includes(tool.key)}
                        onCheckedChange={() => toggleTool(tool.key)} />
                      <span className="min-w-0">
                        <span className="font-medium" title={tool.key}>{toolLabel(tool.key, tool.system)}</span>
                        {tool.description && <span className="text-muted-foreground line-clamp-1"> {tool.description}</span>}
                      </span>
                    </label>
                  ))}
                </div>
              </div>
              <Textarea className="text-xs" rows={2} value={toolCallMsg}
                placeholder={tr("工具调用时发送给智能体的消息（系统附加任务信息、工具输入和输出）", "Message on tool call (task details, input and output appended)")} onChange={(e) => setToolCallMsg(e.target.value)} />
            </div>
          )}
        </div>

        {/* task created */}
        <div className="grid gap-1.5">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox checked={onTaskCreate} onCheckedChange={(v) => setOnTaskCreate(!!v)} /> Trigger on task created
          </label>
          {onTaskCreate && (
            <Textarea className="text-xs" rows={2} value={taskCreateMsg}
              placeholder={tr("创建任务时发送给智能体的消息（系统附加任务 ID 和目标）", "Message on task creation (task ID and goal appended)")} onChange={(e) => setTaskCreateMsg(e.target.value)} />
          )}
        </div>

        <div className="flex items-center gap-2">
          <Button size="sm" onClick={submit} disabled={saving}>
            <SaveIcon /> {editingId != null ? tr("保存修改", "Save changes") : tr("添加触发器", "Add trigger")}
          </Button>
          {editingId != null && (
            <Button size="sm" variant="ghost" onClick={resetForm} disabled={saving}>
              <XIcon /> Cancel edit
            </Button>
          )}
        </div>
      </div>

      {/* existing triggers */}
      <div className="grid gap-2">
        <Label className="text-muted-foreground text-xs">{tr("已有触发器", "Existing triggers")}</Label>
        {triggers.length === 0 && <span className="text-muted-foreground text-xs">{tr("（暂无）", "(None)")}</span>}
        {triggers.map((t) => (
          <div
            key={t.id}
            className={cn(
              "flex items-start gap-2 rounded-md border p-2 text-sm",
              editingId === t.id && "border-primary bg-primary/5",
            )}
          >
            <Switch checked={t.enabled} onCheckedChange={() => toggleEnabled(t)} className="mt-0.5" />
            <div className="min-w-0 flex-1">
              <div className="flex flex-wrap items-center gap-1.5">
                <span className="font-medium">{condLabel(t)}</span>
                {!t.enabled && (
                  <Badge variant="outline" className="text-destructive px-1 py-0 text-[9px]">{tr("已停用", "Disabled")}</Badge>
                )}
              </div>
              <div className="text-muted-foreground grid gap-0.5 text-xs">
                {t.interval_sec > 0 && t.interval_message && <div className="line-clamp-1">Schedule: {t.interval_message}</div>}
                {t.on_finding && t.finding_message && <div className="line-clamp-1">Finding: {t.finding_message}</div>}
                {t.on_goal_met && t.goal_message && <div className="line-clamp-1">Goal: {t.goal_message}</div>}
                {t.on_task_timeout && t.task_timeout_message && <div className="line-clamp-1">Timeout: {t.task_timeout_message}</div>}
                {t.on_task_create && t.task_create_message && <div className="line-clamp-1">Task created: {t.task_create_message}</div>}
                {t.on_tool_call && (
                  <>
                    <div className="line-clamp-1">Tools: {t.tool_names.join(", ") || "(none selected)"}</div>
                    {t.tool_call_message && <div className="line-clamp-1">Message: {t.tool_call_message}</div>}
                  </>
                )}
              </div>
            </div>
            <Button variant="ghost" size="icon-sm" className="text-muted-foreground hover:text-foreground"
              onClick={() => startEdit(t)} title={tr("编辑", "Edit")}>
              <PencilIcon className="size-3.5" />
            </Button>
            <Button variant="ghost" size="icon-sm" className="text-muted-foreground hover:text-destructive"
              onClick={() => del(t.id)} title={tr("删除", "Delete")}>
              <Trash2Icon className="size-3.5" />
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}
