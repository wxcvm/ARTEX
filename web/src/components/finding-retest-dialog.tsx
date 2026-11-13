"use client";

import * as React from "react";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText, type UiLocale } from "@/lib/ui-locale";

import { RotateCcwIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldDescription, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/lib/api";
import type { FindingRetest } from "@/lib/types";

interface FindingRetestDialogProps {
  findingId: string;
  findingName?: string;
  onClose: () => void;
  onStarted?: (retest: FindingRetest) => void;
}

// Mounted only while open; notes cleared on close. The list and detail share the submit lock and error handling; stays on the current page after starting.
export function FindingRetestDialog({ findingId, findingName, onClose, onStarted }: FindingRetestDialogProps) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const notesId = React.useId();
  const [notes, setNotes] = React.useState("");
  const [submitting, setSubmitting] = React.useState(false);
  const submitLock = React.useRef(false);

  async function start() {
    if (submitLock.current) return;
    submitLock.current = true;
    setSubmitting(true);
    try {
      const result = await api.startFindingRetest(findingId, notes.trim());
      onStarted?.(result.retest);
      onClose();
      toast.success(result.created ? tr("复测已启动，点击「复测中」查看会话", "Retest started; click Retesting to view the session") : tr("此漏洞正在复测，请查看已有会话", "This finding is already being retested; view its session"));
    } catch (e) {
      toast.error(`${tr("启动复测失败：", "Failed to start retest: ")}${(e as Error).message}`);
    } finally {
      submitLock.current = false;
      setSubmitting(false);
    }
  }

  return (
    <Dialog open onOpenChange={(open) => !open && !submitLock.current && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{tr(`复测漏洞 #${findingId}`, `Retest finding #${findingId}`)}</DialogTitle>
          <DialogDescription className="break-words">
            {findingName ? <span className="mb-2 block">{findingName}</span> : null}
            {tr("复测智能体", "Retest agent")}
            {tr("读取原始证据和测试约束，在独立会话中进行针对性验证。复测成功且确认修复后，漏洞状态自动变为「已修复」；其他结论保留原状态。", "Read the original evidence and constraints for targeted verification in a separate session. A successful retest confirming the fix marks the finding Fixed; other verdicts retain its state.")}
          </DialogDescription>
        </DialogHeader>
        <FieldGroup>
          <Field data-disabled={submitting}>
            <FieldLabel htmlFor={notesId}>{tr("备注（可选）", "Notes (optional)")}</FieldLabel>
            <Textarea
              id={notesId}
              value={notes}
              maxLength={4000}
              rows={4}
              disabled={submitting}
              onChange={(e) => setNotes(e.target.value)}
              placeholder={tr("例如使用原测试账号验证原接口，修复版本为 v2。", "e.g. verify the original endpoint with the test account; fix version v2.")}
            />
            <FieldDescription>{tr("填写修复版本、测试条件或本次复测的限制。", "Enter the fix version, test conditions or limitations.")}</FieldDescription>
          </Field>
        </FieldGroup>
        <DialogFooter>
          <Button variant="outline" disabled={submitting} onClick={onClose}>
            {tr("取消", "Cancel")}
          </Button>
          <Button disabled={submitting} onClick={() => void start()}>
            {submitting ? <Spinner data-icon="inline-start" /> : <RotateCcwIcon data-icon="inline-start" />}
            {submitting ? tr("正在创建…", "Creating…") : tr("启动复测", "Start retest")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
