"use client";

import * as React from "react";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";

import { ArchiveIcon, ArrowLeftIcon, BrainIcon, CheckIcon, CircleAlertIcon, CircleHelpIcon, PauseIcon, PlayIcon } from "lucide-react";
import { toast } from "sonner";

import { StatusBadge } from "@/components/status-badge";
import { TaskLLMProfileChain } from "@/components/task-llm-profile-chain";
import { TaskOutcomeCard } from "@/components/task-outcome-card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
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
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverDescription,
  PopoverHeader,
  PopoverTitle,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Spinner } from "@/components/ui/spinner";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { api } from "@/lib/api";
import { guideHref } from "@/lib/user-guide";
import type { LLMProfile, Task } from "@/lib/types";

import { AssetsTab } from "./_tabs/assets-tab";
import { BroadcastTab } from "./_tabs/broadcast-tab";
import { CoverageGraphTab } from "./_tabs/coverage-graph-tab";
import { CredentialsTab } from "./_tabs/credentials-tab";
import { FindingsTab } from "./_tabs/findings-tab";
import { GraphTab } from "./_tabs/graph-tab";
import { InterceptTab } from "./_tabs/intercept-tab";
import { OverviewTab } from "./_tabs/overview-tab";
import { ReportTab } from "./_tabs/report-tab";
import { RetestsTab } from "./_tabs/retests-tab";
import { SessionsTab } from "./_tabs/sessions-tab";

const TABS = [
  { value: "sessions", label: "会话", en: "Sessions" },
  { value: "overview", label: "概览", en: "Overview" },
  { value: "graph", label: "探索图谱", en: "Exploration graph" },
  { value: "broadcast", label: "任务广播", en: "Task broadcasts" },
  { value: "findings", label: "发现结果", en: "Findings" },
  { value: "retests", label: "复测", en: "Retests" },
  { value: "assets", label: "测试资产", en: "Test assets" },
  { value: "credentials", label: "测试凭据", en: "Test credentials" },
  { value: "coverage", label: "资产覆盖", en: "Asset coverage" },
  { value: "intercept", label: "操作审批", en: "Approvals" },
  { value: "report", label: "报告", en: "Report" },
];

function taskProfileIDs(task: Task): string[] {
  if (task.llm_profile_ids && task.llm_profile_ids.length > 0) {
    return task.llm_profile_ids.map(String);
  }
  return task.llm_profile_id ? [String(task.llm_profile_id)] : [];
}

function TaskLLMControl({ task, profiles, onUpdated }: { task: Task; profiles: LLMProfile[]; onUpdated: () => void }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [open, setOpen] = React.useState(false);
  const [profileIDs, setProfileIDs] = React.useState<string[]>(() => taskProfileIDs(task));
  const [activeProfileID, setActiveProfileID] = React.useState(
    task.active_llm_profile_id ? String(task.active_llm_profile_id) : (taskProfileIDs(task)[0] ?? ""),
  );
  const [saving, setSaving] = React.useState(false);
  const popoverContentRef = React.useRef<HTMLDivElement>(null);

  const chain = taskProfileIDs(task);
  const exhausted = task.llm_failover_state === "chain_exhausted";
  // The chain can be edited in any state, terminal included: after the task ends the main agent conversation still uses this chain,
  // and a broken model on it must be replaceable, or a completed task could no longer be interacted with.
  const terminal = ["done", "failed", "timeout"].includes(task.status);
  // A null active profile on an exhausted, non-empty chain is a persisted end
  // cursor. Keep the status display honest; choosing the first profile is only
  // the editor's reset draft and does not mean it is currently active.
  let activeID = chain[0] ?? "";
  if (exhausted) activeID = "";
  if (task.active_llm_profile_id) activeID = String(task.active_llm_profile_id);
  const activeProfile = profiles.find((profile) => profile.id === activeID);
  const currentLabel = exhausted
    ? tr("配置链已耗尽", "Model chain exhausted")
    : (activeProfile?.name ?? (activeID ? tr(`配置 #${activeID}`, `Profile #${activeID}`) : tr("跟随默认配置", "Follow default")));
  const activeIndex = chain.indexOf(activeID);
  const backupCount = activeIndex >= 0 ? Math.max(0, chain.length - activeIndex - 1) : 0;
  const currentTitle = [currentLabel, activeProfile?.model, backupCount > 0 ? `${backupCount} backups` : ""]
    .filter(Boolean)
    .join(" · ");
  let editorDescription = tr("顺序或当前配置的修改将在下一次模型调用时生效。", "Order or active profile changes apply on the next model call.");
  if (terminal) editorDescription = tr("任务已结束，修改仅影响后续主智能体对话。", "The task has ended; changes affect only future main-agent conversations.");
  let saveLabel = tr("保存", "Save");
  if (exhausted) saveLabel = tr("保存并重置", "Save and reset");
  if (saving) saveLabel = tr("正在保存", "Saving");

  const syncDraft = React.useCallback(() => {
    const next = taskProfileIDs(task);
    setProfileIDs(next);
    setActiveProfileID(task.active_llm_profile_id ? String(task.active_llm_profile_id) : (next[0] ?? ""));
  }, [task]);

  const handleOpenChange = (next: boolean) => {
    setOpen(next);
    if (next) syncDraft();
  };

  const handleProfileIDsChange = (next: string[]) => {
    setProfileIDs(next);
    setActiveProfileID((current) => (next.includes(current) ? current : (next[0] ?? "")));
  };

  const save = async () => {
    setSaving(true);
    try {
      const result = await api.updateTaskLLMProfiles(
        task.id,
        profileIDs.map(Number),
        activeProfileID ? Number(activeProfileID) : undefined,
      );
      if (result.switch_event) {
        toast.info(result.switch_event.summary, { id: `task-${task.id}-llm-${result.switch_event.seq}` });
      } else {
        toast.success(
          result.reopened_intents > 0
            ? `LLM configs updated and ${result.reopened_intents} quota-blocked intents reopened`
            : "LLM configs updated",
        );
      }
      setOpen(false);
      onUpdated();
    } catch (error) {
      toast.error(tr("更新失败：", "Update failed: ") + (error as Error).message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>
        <Button
          size="sm"
          variant={exhausted ? "destructive" : "outline"}
          aria-label={tr("查看或切换任务模型配置", "View or switch task model")}
          title={currentTitle}
        >
          <BrainIcon data-icon="inline-start" />
          <span className="hidden max-w-36 truncate lg:inline">{currentLabel}</span>
          {backupCount > 0 && <span className="hidden text-muted-foreground xl:inline">+{backupCount}</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent ref={popoverContentRef} align="start" className="w-[min(28rem,calc(100vw-2rem))] gap-4 p-4">
        <PopoverHeader>
          <PopoverTitle>{tr("任务模型配置链", "Task model chain")}</PopoverTitle>
          <PopoverDescription>{editorDescription}</PopoverDescription>
        </PopoverHeader>

        {exhausted && (
          <Alert variant="destructive">
            <CircleAlertIcon />
            <AlertTitle>{tr("配置链额度耗尽", "Model chain quota exhausted")}</AlertTitle>
            <AlertDescription>
              {task.llm_failover_reason ?? tr("选中的配置均被判定为额度耗尽。保存配置链可重置失败状态。", "All selected profiles are out of quota. Save the chain to reset failures.")}
            </AlertDescription>
          </Alert>
        )}

        <TaskLLMProfileChain
          profiles={profiles}
          value={profileIDs}
          onValueChange={handleProfileIDsChange}
          activeProfileId={activeProfileID}
          onActiveProfileChange={setActiveProfileID}
          inputId="task-llm-profiles"
          disabled={saving}
          portalContainer={popoverContentRef}
        />

        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" size="sm" onClick={() => setOpen(false)}>
            {tr("关闭", "Close")}
          </Button>
          <Button type="button" size="sm" onClick={save} disabled={saving}>
            {saving && <Spinner data-icon="inline-start" />}
            {saveLabel}
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}

function TaskDetailInner() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const router = useRouter();
  const searchParams = useSearchParams();
  const id = searchParams.get("id") ?? "";
  const [task, setTask] = React.useState<Task | null>(null);
  const [paused, setPaused] = React.useState(false);
  const [loaded, setLoaded] = React.useState(false);
  const [tab, setTab] = React.useState("sessions");
  const [interceptPendingCount, setInterceptPendingCount] = React.useState(0);
  const [profiles, setProfiles] = React.useState<LLMProfile[]>([]);
  const [archiving, setArchiving] = React.useState(false);

  React.useEffect(() => {
    api
      .llmProfiles()
      .then(setProfiles)
      .catch(() => setProfiles([]));
  }, []);

  React.useEffect(() => {
    let alive = true;
    const load = () =>
      api
        .interceptTask(id)
        .then((rows) => {
          if (alive) setInterceptPendingCount(rows.filter((r) => r.status === "pending").length);
        })
        .catch(() => {
          // Polling is best-effort; the next interval retries automatically.
        });
    void load();
    const t = setInterval(load, 5000);
    return () => {
      alive = false;
      clearInterval(t);
    };
  }, [id]);

  const taskLoadInFlight = React.useRef<string | null>(null);
  const load = React.useCallback(() => {
    if (taskLoadInFlight.current === id) return;
    taskLoadInFlight.current = id;
    Promise.all([api.task(id), api.stats(id).catch(() => null)])
      .then(([task, s]) => {
        if (taskLoadInFlight.current !== id) return;
        const base = { ...task };
        const at = s?.active_task;
        if (at) {
          base.in_flight = at.in_flight;
          base.goals_total = at.goals_total;
          base.goals_met = at.goals_met;
          base.engine_mode = s?.engine_mode ?? at.engine_mode;
          base.paused = at.paused;
        }
        setTask(base);
        setPaused(at?.paused ?? base.paused ?? false);
      })
      .catch(() => {
        // Keep the last rendered task state during a transient poll failure.
      })
      .finally(() => {
        if (taskLoadInFlight.current === id) {
          taskLoadInFlight.current = null;
          setLoaded(true);
        }
      });
  }, [id]);
  React.useEffect(() => {
    load();
    const timer = setInterval(load, 5000);
    return () => clearInterval(timer);
  }, [load]);

  async function togglePause() {
    if (task && ["done", "failed", "timeout"].includes(task.status)) return;
    const next = !paused;
    try {
      await api.controlTask(id, next ? "pause" : "resume");
      setPaused(next);
      toast.success(next ? tr("探索已暂停", "Exploration paused") : tr("探索已恢复", "Exploration resumed"));
    } catch (e) {
      toast.error(tr("操作失败：", "Operation failed: ") + (e as Error).message);
    }
  }

  async function archiveTask() {
    if (!task || archiving) return;
    setArchiving(true);
    try {
      await api.archiveTask(task.id);
      toast.success(tr("任务已加入归档队列", "Task queued for archiving"));
      router.push("/function/tasks");
    } catch (error) {
      toast.error(`归档失败：${(error as Error).message}`);
      setArchiving(false);
    }
  }

  if (!task) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center gap-3 p-10 text-center">
        <p className="text-muted-foreground">{loaded ? tr(`任务 ${id} 已删除、已归档或不存在`, `Task ${id} was deleted, archived or does not exist`) : tr("正在加载…", "Loading…")}</p>
        {loaded && (
          <Button asChild variant="outline">
            <Link href="/function/tasks">
              <ArrowLeftIcon /> {tr("返回任务列表", "Back to task list")}
            </Link>
          </Button>
        )}
      </div>
    );
  }

  const completed = task.status === "done";
  const terminal = ["done", "failed", "timeout"].includes(task.status);
  const archiveLifecycleEligible = terminal || paused || task.status === "paused";
  const canArchive = archiveLifecycleEligible && !task.archive_blocked_by_task_id;
  let archiveDisabledReason = task.queued ? tr("排队任务需先暂停", "Pause queued tasks first") : tr("运行中任务需先暂停", "Pause running tasks first");
  if (archiveLifecycleEligible && task.archive_blocked_by_task_id) {
    archiveDisabledReason = tr(`此任务被未归档任务 #${task.archive_blocked_by_task_id} 直接继承；请先归档依赖任务`, `Task #${task.archive_blocked_by_task_id} directly inherits this task; archive the dependent task first`);
  }
  const engineMode = paused ? "paused" : (task.engine_mode ?? "idle");
  let controlVariant: "default" | "secondary" | "outline" = "outline";
  let controlIcon = <PauseIcon data-icon="inline-start" />;
  let controlLabel = tr("暂停", "Pause");
  if (terminal) {
    controlVariant = "secondary";
    controlIcon = <CheckIcon data-icon="inline-start" />;
    controlLabel = completed ? tr("已完成", "Completed") : tr("已结束", "Ended");
  } else if (paused) {
    controlVariant = "default";
    controlIcon = <PlayIcon data-icon="inline-start" />;
    controlLabel = tr("恢复", "Resume");
  }
  const archiveTrigger = (
    <Button
      size="icon-sm"
      variant="ghost"
      disabled={!canArchive || archiving}
      aria-label={canArchive ? tr("归档任务", "Archive task") : archiveDisabledReason}
    >
      {archiving ? <Spinner /> : <ArchiveIcon />}
    </Button>
  );

  return (
    <Tabs value={tab} onValueChange={setTab} className="flex min-w-0 flex-1 flex-col gap-0">
      {/* Top fixed area */}
      <header className="sticky top-0 z-10 flex min-w-0 flex-col gap-3 border-b border-border/70 bg-background px-4 py-4 shadow-sm lg:px-7">
        <div className="flex min-w-0 flex-wrap items-center gap-2">
          <SidebarTrigger className="-ml-1" />
          <Button asChild variant="ghost" size="icon" className="size-7">
            <Link href="/function/tasks">
              <ArrowLeftIcon />
            </Link>
          </Button>
          <h1 className="min-w-[10rem] flex-1 truncate text-base font-semibold tracking-tight sm:text-lg" title={task.name || task.description}>
            {task.name || task.description}
          </h1>
          <code className="hidden rounded bg-muted px-1.5 py-0.5 font-mono text-xs text-muted-foreground sm:inline">
            {task.id}
          </code>
          <Button asChild variant="ghost" size="icon" className="size-7 shrink-0"
            title={tr("任务详情操作指南", "Task detail guide")}>
            <Link href={guideHref("/function/tasks/detail")} aria-label={tr("任务详情操作指南", "Task detail guide")}>
              <CircleHelpIcon />
            </Link>
          </Button>
          <Separator orientation="vertical" className="mx-1 hidden h-4 sm:block" />
          <TaskLLMControl task={task} profiles={profiles} onUpdated={load} />
          <StatusBadge domain={terminal ? "task" : "engine"} value={terminal ? task.status : engineMode} dot />
          {canArchive ? (
            <AlertDialog>
              <AlertDialogTrigger asChild>{archiveTrigger}</AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>{tr(`归档任务 #${task.id}？`, `Archive task #${task.id}?`)}</AlertDialogTitle>
                  <AlertDialogDescription>
                    {tr("任务图谱、关联记录、独占资产和流量、工作文件及模型历史将压缩到冷存储。归档完成后可在任务列表的「已归档」中恢复。", "The task graph, related records, exclusive assets and traffic, working files and LLM history are compressed into cold storage. Restore from the Archived tab.")}
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>{tr("取消", "Cancel")}</AlertDialogCancel>
                  <AlertDialogAction onClick={() => void archiveTask()}>{tr("确认归档", "Confirm archive")}</AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          ) : (
            <Tooltip>
              <TooltipTrigger asChild>
                <span className="inline-flex">{archiveTrigger}</span>
              </TooltipTrigger>
              <TooltipContent>{archiveDisabledReason}</TooltipContent>
            </Tooltip>
          )}
          <Button size="sm" variant={controlVariant} onClick={togglePause} disabled={terminal}>
            {controlIcon}
            {controlLabel}
          </Button>
        </div>
        <p className="line-clamp-2 max-w-4xl text-xs leading-relaxed text-muted-foreground">{task.goal}</p>
        {/* Tabs */}
        <div className="no-scrollbar min-w-0 overflow-x-auto pb-1">
          <TabsList variant="default" className="h-auto min-w-max gap-1 rounded-lg bg-muted/55 p-1">
            {TABS.map((t) => (
              <TabsTrigger key={t.value} value={t.value} className="rounded-md px-3 py-2 text-sm">
                {locale === "en" ? t.en : t.label}
                {t.value === "intercept" && interceptPendingCount > 0 && (
                  <span className="ml-1.5 inline-flex h-4 min-w-[16px] items-center justify-center rounded-full bg-amber-500 px-1 text-[10px] font-semibold leading-none text-white">
                    {interceptPendingCount > 99 ? "99+" : interceptPendingCount}
                  </span>
                )}
              </TabsTrigger>
            ))}
          </TabsList>
        </div>
      </header>

      {/* Tab content */}
      <div className="min-w-0 flex-1 bg-background p-4 sm:p-5 lg:p-7">
        {task.id === id && <TaskOutcomeCard key={id} task={task} />}
        <TabsContent value="sessions" className="mt-0">
          <SessionsTab taskId={id} taskPaused={task?.status === "paused"} />
        </TabsContent>
        <TabsContent value="overview" className="mt-0">
          <OverviewTab taskId={id} />
        </TabsContent>
        <TabsContent value="graph" className="mt-0">
          <GraphTab taskId={id} />
        </TabsContent>
        <TabsContent value="broadcast" className="mt-0">
          <BroadcastTab taskId={id} />
        </TabsContent>
        <TabsContent value="findings" className="mt-0">
          <FindingsTab taskId={id} />
        </TabsContent>
        <TabsContent value="retests" className="mt-0">
          <RetestsTab key={id} taskId={id} />
        </TabsContent>
        <TabsContent value="assets" className="mt-0">
          <AssetsTab taskId={id} />
        </TabsContent>
        <TabsContent value="credentials" className="mt-0">
          <CredentialsTab taskId={id} />
        </TabsContent>
        <TabsContent value="coverage" className="mt-0">
          <CoverageGraphTab taskId={id} coverageEnabled={task?.coverage_enabled !== false} />
        </TabsContent>
        <TabsContent value="intercept" className="mt-0">
          <InterceptTab taskId={id} />
        </TabsContent>
        <TabsContent value="report" className="mt-0">
          <ReportTab key={id} taskId={id} />
        </TabsContent>
      </div>
    </Tabs>
  );
}

// useSearchParams must sit under a Suspense boundary for static export.
export default function TaskDetailPage() {
  return (
    <React.Suspense fallback={null}>
      <TaskDetailInner />
    </React.Suspense>
  );
}
