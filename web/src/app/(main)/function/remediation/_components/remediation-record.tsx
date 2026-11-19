"use client";

import Link from "next/link";
import { useEffect, useRef, useState, type FormEvent } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { api, ApiResponseError } from "@/lib/api";
import { mergeRemediationEvents, remediationActions, remediationDraft, remediationLimits, remediationMutation } from "@/lib/remediation";
import { statusLabel } from "@/lib/status";
import type { RemediationAction, RemediationDetail, RemediationEvent } from "@/lib/types";
import { uiText, type UiLocale } from "@/lib/ui-locale";
import { useUiLocale } from "@/lib/ui-locale.client";
import { actionLabel, remediationTime, stageGuidance, stageLabel } from "./labels";

type Read = {
  phase: "loading" | "ready" | "missing" | "error";
  data?: RemediationDetail;
  message?: string;
};
const errorMessage = (error: unknown) => error instanceof Error ? error.message : undefined;

function EventRow({ event, locale }: { event: RemediationEvent; locale: UiLocale }) {
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  return (
    <li className="space-y-1 rounded-lg border p-3 text-sm">
      <p className="font-medium">{actionLabel(locale, event.action)}</p>
      <p className="text-xs text-muted-foreground">
        {remediationTime(event.time, locale)} · {event.actor} · {tr("记录版本", "Record version")} {event.version}
      </p>
      <p className="text-xs text-muted-foreground">{stageLabel(locale, event.from_stage)} → {stageLabel(locale, event.to_stage)}</p>
      {event.note && <p className="whitespace-pre-wrap break-words">{event.note}</p>}
      {event.reference && <p className="break-words text-xs">{tr("材料引用：", "Material reference: ")}{event.reference}</p>}
      {event.plan && <div className="space-y-1">
        <p className="text-xs text-muted-foreground">
          {tr("计划快照", "Plan snapshot")}
          {event.owner && <> · {tr("责任人：", "Owner: ")}{event.owner}</>}
          {event.due_date && <> · {tr("期限：", "Due: ")}{event.due_date}</>}
        </p>
        <p className="whitespace-pre-wrap break-words">{event.plan}</p>
      </div>}
    </li>
  );
}

export function RemediationRecord({ id, onChanged }: { id: string; onChanged: () => void }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [read, setRead] = useState<Read>({ phase: "loading" });
  const [revision, setRevision] = useState(0);
  const [draft, setDraft] = useState<ReturnType<typeof remediationDraft> | null>(null);
  const [selected, setSelected] = useState<RemediationAction | null>(null);
  const [conflict, setConflict] = useState(false);
  const [notice, setNotice] = useState<"saved" | "conflict" | null>(null);
  const [writeError, setWriteError] = useState("");
  const [historyError, setHistoryError] = useState("");
  const [working, setWorking] = useState<"write" | "history" | null>(null);
  const activeRequest = useRef(false);
  const item = read.data?.item;
  const actions = item ? remediationActions(item) : [];
  const action = selected ?? actions[0];
  const fresh = read.phase === "ready";
  const reload = () => {
    setRead(previous => ({ phase: "loading", data: previous.data }));
    setHistoryError("");
    setRevision(n => n + 1);
  };

  useEffect(() => {
    let live = true;
    api.remediationDetail(id).then(data => {
      if (!live) return;
      setRead({ phase: "ready", data });
      setDraft(previous => previous ?? remediationDraft(data.item));
    }).catch((error: unknown) => {
      if (!live) return;
      if (error instanceof ApiResponseError && error.status === 404) {
        setRead({ phase: "missing" });
        return;
      }
      setRead(previous => ({ phase: "error", data: previous.data, message: errorMessage(error) }));
    });
    return () => { live = false; };
  }, [id, revision]);

  async function save(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!item || !draft || !action || activeRequest.current) return;
    const input = remediationMutation(item, action, draft, fresh && !conflict);
    if (!input.ok) {
      setWriteError({
        unavailable: tr("当前阶段不允许执行这个操作。", "This action is unavailable at the current stage."),
        note: tr("请填写人工处理说明。", "Enter a human review note."),
        reference: tr("请填写验收材料或内部工单的引用。", "Enter a reference to acceptance material or an internal ticket."),
        plan: tr("请填写整改计划。", "Enter the remediation plan."),
        scope: tr("请先读取并核对最新记录。", "Read and review the latest record first."),
        length: tr("输入超过长度限制。", "An input exceeds its length limit."),
      }[input.error]);
      return;
    }
    activeRequest.current = true;
    setWorking("write"); setWriteError(""); setNotice(null);
    try {
      const data = await api.addRemediationEvent(id, input.value);
      setRead({ phase: "ready", data });
      setDraft(remediationDraft(data.item));
      setSelected(null);
      setNotice("saved");
      onChanged();
    } catch (error: unknown) {
      if (error instanceof ApiResponseError && error.status === 409) {
        setConflict(true);
        setNotice("conflict");
        reload();
        onChanged();
      } else if (error instanceof ApiResponseError && error.status === 404) {
        setRead({ phase: "missing" });
        onChanged();
      } else {
        setWriteError(errorMessage(error) ?? tr("未能保存整改记录。", "The remediation record could not be saved."));
      }
    } finally {
      activeRequest.current = false;
      setWorking(null);
    }
  }

  async function older() {
    const current = read.data;
    if (!current?.next_before || !fresh || activeRequest.current) return;
    activeRequest.current = true;
    setWorking("history"); setHistoryError("");
    try {
      const next = await api.remediationDetail(id, current.next_before);
      if (next.item.version !== current.item.version || next.item.evidence_version !== current.item.evidence_version || next.item.read_only !== current.item.read_only) {
        setConflict(true);
        setNotice("conflict");
        reload();
        onChanged();
        return;
      }
      setRead({
        phase: "ready",
        data: { ...current, item: next.item, events: mergeRemediationEvents(current.events, next.events), next_before: next.next_before },
      });
    } catch (error: unknown) {
      if (error instanceof ApiResponseError && error.status === 404) {
        setRead({ phase: "missing" });
        onChanged();
      } else {
        setHistoryError(errorMessage(error) ?? tr("未能读取更早的记录。", "Older records could not be loaded."));
      }
    } finally {
      activeRequest.current = false;
      setWorking(null);
    }
  }

  return (
    <section className="min-w-0 space-y-4 rounded-xl border bg-card p-4" aria-label={tr("整改详情", "Remediation detail")}>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h2 className="text-lg font-semibold">{tr("整改详情", "Remediation detail")} <code className="text-sm text-muted-foreground">#{id}</code></h2>
        <Button size="sm" variant="outline" onClick={reload} disabled={read.phase === "loading" || working !== null}>
          {tr("重新读取记录", "Reload record")}
        </Button>
      </div>
      {read.phase === "loading" && <p role="status" className="text-sm text-muted-foreground">{tr("正在读取详情…", "Loading record…")}</p>}
      {read.phase === "missing" && <p role="status" className="text-sm text-muted-foreground">
        {tr("找不到该记录，或当前上下文无权查看。", "This record does not exist or is unavailable in this context.")}
      </p>}
      {read.phase === "error" && <div role="alert" className="space-y-1 text-sm">
        <p className="text-destructive">{tr("本次未能读取详情。", "The record could not be loaded.")}</p>
        {read.message && <p className="break-words">{read.message}</p>}
        {item && <p className="text-muted-foreground">{tr("仅显示上次成功读取的内容，暂不可提交变更。", "Showing the last successful response. Changes are disabled.")}</p>}
      </div>}
      {item && <>
        <div className="space-y-2 rounded-lg border p-3">
          <p className="break-words font-semibold">{item.title || `#${item.finding_id}`}</p>
          <p className="text-sm">{stageLabel(locale, item.stage)} · {statusLabel("severity", item.severity, locale)} · {statusLabel("finding", item.finding_status, locale)}</p>
          <p className="text-xs text-muted-foreground">
            {tr("记录版本", "Record version")} {item.version} · {tr("当前证据版本", "Current evidence version")} {item.evidence_version} · {tr("提交证据版本", "Submitted evidence version")} {item.submitted_evidence_version ?? "—"} · {tr("已验收证据版本", "Accepted evidence version")} {item.accepted_evidence_version ?? "—"}
          </p>
          {read.data?.events[0] && <p className="text-xs text-muted-foreground">{tr("最近登记：", "Latest event: ")}{remediationTime(read.data.events[0].time, locale)}</p>}
          <p className="text-sm">{stageGuidance(locale, item.stage)}</p>{item.reference && <p className="break-words text-sm">{tr("当前验收材料：", "Current acceptance material: ")}{item.reference}</p>}
          {item.needs_review && <p className="text-sm font-medium text-amber-700 dark:text-amber-300">
            {tr("证据版本发生变化。先人工复核，旧提交不能作为当前版本的验收依据。", "Evidence has changed. Review it manually; a previous submission is not valid for accepting the current version.")}
          </p>}
          {item.stage === "pending_acceptance" && item.submitted_evidence_version !== item.evidence_version && (
            <p className="text-xs text-muted-foreground">{tr("只能退回，并由整改人员按当前版本重新提交。", "Return the record for a fresh submission against the current version.")}</p>
          )}
          {item.read_only && <p className="text-sm font-medium">{tr("该记录为只读，不能追加整改事件。", "This record is read-only. No new remediation events can be added.")}</p>}
          <Link className="block text-sm text-primary underline underline-offset-2" href={`/function/findings/detail?id=${encodeURIComponent(item.finding_id)}`}>
            {tr("打开已有发现与证据", "Open existing finding and evidence")}
          </Link>
        </div>
        {item.plan && <div className="rounded-lg border p-3 text-sm">
          <h3 className="font-semibold">{tr("当前整改计划", "Current remediation plan")}</h3>
          <p className="mt-1 text-xs text-muted-foreground">{tr("责任人：", "Owner: ")}{item.owner || "—"} · {tr("期限：", "Due: ")}{item.due_date || "—"}</p>
          <p className="mt-2 whitespace-pre-wrap break-words">{item.plan}</p>
        </div>}
        {notice === "saved" && <p role="status" className="text-sm">{tr("已登记，原始证据未更改。", "Recorded. The original evidence was not changed.")}</p>}
        {notice === "conflict" && <div role="alert" className="space-y-2 rounded-lg border border-amber-500 p-3 text-sm">
          <p>{tr("记录或证据版本已变化，未提交原操作。草稿仍保留；请检查最新详情后再手动提交。", "The record or evidence version changed. Nothing was submitted. Your draft is retained; review the latest record before submitting manually.")}</p>
          {conflict && <Button type="button" size="sm" variant="outline" disabled={!fresh || working !== null || item.read_only}
            onClick={() => { setConflict(false); setNotice(null); }}>
            {tr("我已核对最新版本，继续编辑", "I reviewed the latest version; continue editing")}
          </Button>}
        </div>}
        {!item.read_only && draft && action && <form onSubmit={save} className="space-y-3 rounded-lg border p-3">
          <fieldset className="space-y-3" disabled={working !== null}>
            <legend className="text-sm font-semibold">{tr("人工登记", "Record a human action")}</legend>
            <div className="space-y-1">
              <Label htmlFor="remediation-action">{tr("选择操作", "Action")}</Label>
              <select id="remediation-action" className="block max-w-full rounded-lg border bg-background p-2 text-sm"
                value={action} onChange={event => { setSelected(event.target.value as RemediationAction); setWriteError(""); }}>
                {!actions.includes(action) && <option value={action} disabled>{actionLabel(locale, action)} · {tr("当前不可用", "Unavailable")}</option>}
                {actions.map(value => <option key={value} value={value}>{actionLabel(locale, value)}</option>)}
              </select>
            </div>
            {action === "update_plan" && <>
              <div className="space-y-1"><Label htmlFor="remediation-plan">{tr("整改计划", "Remediation plan")}</Label>
                <Textarea id="remediation-plan" maxLength={remediationLimits.plan} required value={draft.plan} onChange={event => setDraft({ ...draft, plan: event.target.value })} />
              </div>
              <div className="flex flex-wrap gap-3">
                <div className="min-w-32 flex-1 space-y-1"><Label htmlFor="remediation-owner">{tr("责任人（选填）", "Owner (optional)")}</Label>
                  <Input id="remediation-owner" maxLength={remediationLimits.owner} value={draft.owner} onChange={event => setDraft({ ...draft, owner: event.target.value })} />
                </div>
                <div className="min-w-32 flex-1 space-y-1"><Label htmlFor="remediation-due">{tr("期限（选填）", "Due date (optional)")}</Label>
                  <Input id="remediation-due" type="date" value={draft.dueDate} onChange={event => setDraft({ ...draft, dueDate: event.target.value })} />
                </div>
              </div>
            </>}
            {action === "submit" && <div className="space-y-1">
              <Label htmlFor="remediation-reference">{tr("验收材料引用", "Acceptance material reference")}</Label>
              <Input id="remediation-reference" aria-describedby="remediation-reference-help" maxLength={remediationLimits.reference} required value={draft.reference} onChange={event => setDraft({ ...draft, reference: event.target.value })} />
              <p id="remediation-reference-help" className="text-xs text-muted-foreground">{tr("填写内部工单或文档编号，不要粘贴令牌或密钥。", "Use an internal ticket or document ID; do not paste tokens or keys.")}</p>
            </div>}
            <div className="space-y-1">
              <Label htmlFor="remediation-note">{action === "update_plan" ? tr("处理说明（选填）", "Action note (optional)") : tr("处理说明", "Action note")}</Label>
              <Textarea id="remediation-note" maxLength={remediationLimits.note} required={action !== "update_plan"} value={draft.note} onChange={event => setDraft({ ...draft, note: event.target.value })} />
            </div>
            <Button type="submit" disabled={!fresh || conflict || !actions.includes(action) || working !== null}>
              {working === "write" ? tr("正在保存…", "Saving…") : actionLabel(locale, action)}
            </Button>
          </fieldset>
          {writeError && <p role="alert" className="break-words text-sm text-destructive">{writeError}</p>}
        </form>}
        <section className="space-y-2" aria-label={tr("整改事件历史", "Remediation event history")}>
          <h3 className="text-sm font-semibold">{tr("整改事件历史", "Remediation event history")}</h3>
          {read.data && read.data.events.length > 0 ? <ol className="space-y-2">
            {read.data.events.map(event => <EventRow key={event.id} event={event} locale={locale} />)}
          </ol> : fresh && <p className="text-sm text-muted-foreground">{tr("尚无人工登记事件。", "No human actions recorded yet.")}</p>}
          {historyError && <p role="alert" className="break-words text-sm text-destructive">{historyError}</p>}
          {!!read.data?.next_before && <Button type="button" variant="outline" size="sm" disabled={!fresh || working !== null} onClick={() => void older()}>
            {working === "history" ? tr("正在读取历史…", "Loading history…") : tr("读取更早事件", "Load earlier events")}
          </Button>}
        </section>
      </>}
    </section>
  );
}
