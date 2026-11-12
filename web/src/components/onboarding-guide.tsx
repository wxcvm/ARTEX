import { onboardingGuides } from "@/lib/onboarding";

const labels = { purpose: "用途", materials: "准备材料", actions: "具体操作", example: "填写示例", success: "成功表现", failure: "失败处理" } as const;

export function OnboardingGuide({ step }: { step: keyof typeof onboardingGuides }) {
  return <details className="rounded-lg border bg-muted/20 p-4">
    <summary className="cursor-pointer text-sm font-medium">详细指南：用途、操作示例和报错处理</summary>
    <dl className="mt-4 space-y-4 text-sm">
      {Object.entries(onboardingGuides[step]).map(([key, text]) => <div key={key}>
        <dt className="font-medium">{labels[key as keyof typeof labels]}</dt>
        <dd className="mt-1 break-words leading-6 text-muted-foreground">{text}</dd>
      </div>)}
    </dl>
  </details>;
}
