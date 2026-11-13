"use client";

import * as React from "react";
import { toast } from "sonner";
import { Bot, PlusIcon, Trash2Icon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { AgentEditor } from "@/components/agent-editor";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { api } from "@/lib/api";
import type { Agent } from "@/lib/types";

// AgentGridCard is one clickable tile opening the agent's editor drawer. Custom
// (non-builtin) agents get a delete button.
function AgentGridCard({
  agent,
  onOpen,
  onDeleted,
}: {
  agent: Agent;
  onOpen: () => void;
  onDeleted: () => void;
}) {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  async function del() {
    try {
      await api.deleteAgent(agent.key);
      toast.success(t(`已删除智能体「${agent.name}」`, `Deleted agent "${agent.name}"`));
      onDeleted();
    } catch (e) {
      toast.error(t("删除失败：", "Delete failed: ") + (e as Error).message);
    }
  }
  return (
    <div className="hover:border-primary/50 group relative flex flex-col gap-2 rounded-lg border p-4 transition-colors">
      <button type="button" onClick={onOpen} className="flex flex-col gap-2 text-left">
        <div className="flex flex-wrap items-center gap-2">
          <Bot className="text-muted-foreground size-4" />
          <span className="text-sm font-medium">{agent.name}</span>
          <span className="text-muted-foreground font-mono text-xs">{agent.key}</span>
          {agent.builtin ? (
            <Badge variant="secondary" className="px-1.5 py-0 text-[10px]">
              {t("内置", "Built-in")}
            </Badge>
          ) : (
            <Badge variant="outline" className="px-1.5 py-0 text-[10px]">
              {t("自定义", "Custom")}
            </Badge>
          )}
          {!agent.enabled && (
            <Badge variant="outline" className="text-destructive px-1.5 py-0 text-[10px]">
              {t("已停用", "Disabled")}
            </Badge>
          )}
        </div>
        <p className="text-muted-foreground line-clamp-2 min-h-8 text-xs">
          {agent.description || t("（暂无说明）", "(No description)")}
        </p>
        <div className="text-muted-foreground flex flex-wrap gap-1.5 text-[10px]">
          <span className="rounded border px-1.5 py-0.5">MCP {agent.mcp_count ?? 0}</span>
          <span className="rounded border px-1.5 py-0.5">{t("技能", "Skills")} {agent.skill_count ?? 0}</span>
          <span className="rounded border px-1.5 py-0.5">{t("工具", "Tools")} {agent.tool_count ?? 0}</span>
        </div>
      </button>
      {!agent.builtin && (
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button
              variant="ghost"
              size="icon-sm"
              className="text-muted-foreground hover:text-destructive absolute top-2 right-2 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <Trash2Icon className="size-3.5" />
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t(`删除智能体「${agent.name}」？`, `Delete agent "${agent.name}"?`)}</AlertDialogTitle>
              <AlertDialogDescription>
                {t("此操作也会删除该智能体的提示词、变量、可见资源及工具绑定，且无法撤销。", "This permanently deletes the agent's prompts, variables, visible resources and tool bindings.")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{t("取消", "Cancel")}</AlertDialogCancel>
              <AlertDialogAction onClick={del}>{t("删除", "Delete")}</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </div>
  );
}

function CreateAgentDialog({ onCreated }: { onCreated: (key: string) => void }) {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  const [open, setOpen] = React.useState(false);
  const [key, setKey] = React.useState("");
  const [name, setName] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [busy, setBusy] = React.useState(false);

  async function create() {
    setBusy(true);
    try {
      const a = await api.createAgent(key.trim(), name.trim(), description.trim());
      toast.success(t(`已创建智能体「${a.name}」`, `Created agent "${a.name}"`));
      setOpen(false);
      setKey("");
      setName("");
      setDescription("");
      onCreated(a.key);
    } catch (e) {
      toast.error(t("创建失败：", "Create failed: ") + (e as Error).message);
    } finally {
      setBusy(false);
    }
  }

  const keyOk = /^[a-z][a-z0-9_]*$/.test(key.trim());
  const canCreate = keyOk && name.trim().length > 0 && !busy;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <PlusIcon /> {t("创建智能体", "Create agent")}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("创建自定义智能体", "Create custom agent")}</DialogTitle>
          <DialogDescription>
            {t("创建对话助手。标识符用于内部识别，创建后不可修改；名称和说明用于展示。", "Create a conversational assistant. The identifier is immutable; the name and description are for display.")}
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-2">
          <div className="grid gap-1.5">
            <Label htmlFor="agent-key">{t("标识符", "Identifier")}</Label>
            <Input
              id="agent-key"
              placeholder={t("例如 research_helper", "e.g. research_helper")}
              value={key}
              onChange={(e) => setKey(e.target.value)}
              className="font-mono"
            />
            {key.length > 0 && !keyOk && (
              <span className="text-destructive text-xs">{t("必须以小写字母开头，仅支持小写字母、数字和下划线", "Start with a lowercase letter; use lowercase letters, digits and underscores only")}</span>
            )}
          </div>
          <div className="grid gap-1.5">
            <Label htmlFor="agent-name">{t("名称", "Name")}</Label>
            <Input
              id="agent-name"
              placeholder={t("例如：研究助手", "e.g. Research assistant")}
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
          <div className="grid gap-1.5">
            <Label htmlFor="agent-desc">{t("说明", "Description")}</Label>
            <Textarea
              id="agent-desc"
              placeholder={t("用一句话介绍智能体的用途", "Describe the agent's purpose in one sentence")}
              rows={2}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>
        </div>
        <DialogFooter>
          <Button onClick={create} disabled={!canCreate}>
            {t("创建", "Create")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export default function AgentsPage() {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  const [agents, setAgents] = React.useState<Agent[]>([]);
  const [editKey, setEditKey] = React.useState<string | null>(null);

  const reload = React.useCallback(() => {
    api.agents().then(setAgents).catch(() => setAgents([]));
  }, []);
  React.useEffect(() => {
    reload();
  }, [reload]);

  const editing = agents.find((a) => a.key === editKey) ?? null;

  return (
    <div className="flex flex-1 flex-col gap-4 md:gap-6">
      <div className="flex items-start justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">{t("智能体", "Agents")}</h1>
          <p className="text-muted-foreground text-sm">
            {t("管理内置智能体的提示词和配置，以及创建和管理自定义对话助手", "Manage built-in agent prompts and settings, or create and manage custom assistants.")}
          </p>
        </div>
        <CreateAgentDialog
          onCreated={(key) => {
            reload();
            setEditKey(key);
          }}
        />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t("智能体列表", "Agent list")}</CardTitle>
          <CardDescription>{t(`共 ${agents.length} 个`, `${agents.length} agents`)}</CardDescription>
        </CardHeader>
        <CardContent>
          {agents.length === 0 ? (
            <p className="text-muted-foreground py-6 text-center text-sm">{t("（暂无智能体）", "(No agents)")}</p>
          ) : (
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {agents.map((a) => (
                <AgentGridCard key={a.key} agent={a} onOpen={() => setEditKey(a.key)} onDeleted={reload} />
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Sheet open={!!editing} onOpenChange={(o) => !o && setEditKey(null)}>
        <SheetContent
          side="right"
          className="flex flex-col gap-0 p-0 data-[side=right]:w-[45vw] data-[side=right]:sm:max-w-[45vw]"
        >
          {editing && (
            <>
              <SheetHeader className="px-4">
                <SheetTitle className="flex items-center gap-2">
                  {editing.name}
                  <span className="text-muted-foreground font-mono text-xs">{editing.key}</span>
                  {!editing.builtin && (
                    <Badge variant="outline" className="px-1.5 py-0 text-[10px]">
                      {t("自定义", "Custom")}
                    </Badge>
                  )}
                </SheetTitle>
                <SheetDescription>{editing.description || t("提示词、配置、可见资源和工具绑定", "Prompts, settings, visible resources and tool bindings")}</SheetDescription>
              </SheetHeader>
              <AgentEditor agentKey={editing.key} onSaved={reload} />
            </>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
