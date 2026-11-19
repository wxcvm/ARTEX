import { uiText, type UiLocale } from "@/lib/ui-locale";
import type { RemediationAction, RemediationStage } from "@/lib/types";

const stages = {
  pending_review: ["待复核", "Pending review"],
  pending_fix: ["待修复", "Pending fix"],
  pending_acceptance: ["待验收", "Pending acceptance"],
  accepted: ["已验收", "Accepted"],
} as const;
const actions = {
  review: ["登记复核", "Record review"],
  update_plan: ["保存整改计划", "Save remediation plan"],
  submit: ["提交验收材料", "Submit for acceptance"],
  accept: ["确认验收", "Accept remediation"],
  return: ["退回整改", "Return for fixes"],
  reopen: ["重新打开整改", "Reopen remediation"],
} as const;
const guidance = {
  pending_review: ["先由人工核对已有证据、适用范围与记录，再登记结论。", "Manually review the existing evidence, scope and record before recording your conclusion."],
  pending_fix: ["记录责任人、期限与修复计划；完成修复后附上内部工单或验收材料编号。", "Record the owner, due date and fix plan. Once the fix is complete, provide an internal ticket or acceptance evidence reference."],
  pending_acceptance: ["人工核对提交材料与当前证据版本，再确认验收或写明退回原因。", "Manually compare the submitted material with the current evidence version, then accept it or explain why it is returned."],
  accepted: ["这表示人工完成了该版本的验收，不代表后续证据已复核。情况变化时可重新打开。", "This is a human acceptance of that evidence version, not a review of later changes. Reopen if circumstances change."],
} as const;

export const stageLabel = (locale: UiLocale, stage: RemediationStage) => uiText(locale, stages[stage][0], stages[stage][1]);
export const actionLabel = (locale: UiLocale, action: RemediationAction) => uiText(locale, actions[action][0], actions[action][1]);
export const stageGuidance = (locale: UiLocale, stage: RemediationStage) => uiText(locale, guidance[stage][0], guidance[stage][1]);

export function remediationTime(value: string | undefined, locale: UiLocale): string {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.valueOf()) ? "—" : date.toLocaleString(locale);
}
