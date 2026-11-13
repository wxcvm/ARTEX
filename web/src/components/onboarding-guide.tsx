"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { onboardingGuides, onboardingGuidesEn } from "@/lib/onboarding";

const labels = { purpose: "用途", materials: "准备材料", actions: "具体操作", example: "填写示例", success: "成功表现", failure: "失败处理" } as const;

export function OnboardingGuide({ step }: { step: keyof typeof onboardingGuides }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  return <details className="rounded-lg border bg-muted/20 p-4">
    <summary className="cursor-pointer text-sm font-medium">{tr("详细指南：用途、操作示例和报错处理", "Detailed guide: purpose, examples and troubleshooting")}</summary>
    <dl className="mt-4 space-y-4 text-sm">
      {Object.entries(locale === "en" ? onboardingGuidesEn[step] : onboardingGuides[step]).map(([key, text]) => <div key={key}>
        <dt className="font-medium">{locale === "en" ? ({ purpose: "Purpose", materials: "Requirements", actions: "Steps", example: "Example", success: "Success", failure: "Troubleshooting" } as const)[key as keyof typeof labels] : labels[key as keyof typeof labels]}</dt>
        <dd className="mt-1 break-words leading-6 text-muted-foreground">{text}</dd>
      </div>)}
    </dl>
  </details>;
}
