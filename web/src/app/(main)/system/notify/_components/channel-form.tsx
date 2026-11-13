"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";


import { CheckIcon } from "lucide-react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import type { NotificationFilter } from "@/lib/types";

// asText / inputType 是本文件内的取值辅助（与控件渲染强相关），不放在 channel-fields。
import { type FieldDef, type FieldKind, SEVERITY_OPTIONS } from "./channel-fields";

// asText 把任意配置值渲染成输入框可用的字符串。
// config 来自 JSON，值可能是 string / number / boolean / array / null，
// 这里只关心「能不能塞进文本框」，具体序列化由 buildConfig 负责。
function asText(v: unknown): string {
  if (typeof v === "string") return v;
  if (v === null || v === undefined) return "";
  return String(v);
}

// inputType 把字段类型映射到 input 的 type 属性。
function inputType(kind: FieldKind): "text" | "password" | "number" {
  if (kind === "password") return "password";
  if (kind === "number") return "number";
  return "text";
}

// ConfigField 按字段定义渲染对应的控件。
//
// 掩码字段的处理是这里唯一的讲究：输入框**不显示**掩码值本身，只显示一行
// 「已保存」提示。这样界面上就只有一个规则——框里有字就是用户填的，
// 空框就是空值。若把 "__masked__:…abc123" 塞进输入框，用户会以为那是要自己
// 删掉的占位文本，反而更容易误清凭据。
const notificationEnglish: Record<string, string> = {
  "Webhook 地址": "Webhook URL",
  "加签密钥": "Signing secret",
  "签名校验密钥": "Signature verification secret",
  "目标 URL": "Destination URL",
  "请求方法": "Request method",
  "POST（带请求体）": "POST (with body)",
  "PUT（带请求体）": "PUT (with body)",
  "PATCH（带请求体）": "PATCH (with body)",
  "GET（不带请求体）": "GET (no body)",
  "自定义请求头": "Custom headers",
  "请求体模板": "Request body template",
  "API 地址": "API URL",
  "SMTP 服务器": "SMTP server",
  "端口": "Port",
  "账号": "Username",
  "密码 / 授权码": "Password / app password",
  "发件人": "Sender",
  "收件人": "Recipients",
  "隐式 TLS": "Implicit TLS",
  "机器人安全设置选「加签」时填写；选「自定义关键词」或未开启安全设置则留空": "Set when signing is enabled in bot security settings; otherwise leave empty.",
  "机器人开启「签名校验」时填写，否则留空": "Set when bot signature verification is enabled; otherwise leave empty.",
  "每行 KEY=VALUE，例如 Authorization=Bearer xxx": "One KEY=VALUE per line, e.g. Authorization=Bearer xxx",
  "留空用官方地址；自建 Bot API 反代时填写": "Leave empty for official URL; enter a custom Bot API reverse proxy URL if used.",
  "587 走 STARTTLS；465 请把「隐式 TLS」打开": "587 uses STARTTLS; enable implicit TLS for 465.",
  "多个地址用逗号分隔": "Separate multiple addresses with commas.",
  "465 端口打开；587 保持关闭（会自动 STARTTLS）": "Enable for 465; disable for 587 (automatic STARTTLS).",
  "不限": "Any severity",
  "低危及以上": "Low and above",
  "中危及以上": "Medium and above",
  "高危及以上": "High and above",
  "仅严重": "Critical only",
  "钉钉": "DingTalk",
  "飞书": "Lark",
  "企业微信": "WeCom",
  "通用 Webhook": "Generic webhook",
  "邮件": "Email",
  "留空用内置默认模板。变量：{{.Title}} {{.Batch}} {{.Count}} {{.HomeURL}} {{.SentAt}}，以及 range .Items 下的 .Name/.VulnClass/.Severity/.Summary/.Assets/.DetailURL/.StatusLabel。插入字符串请用 {{json .Xxx}} 而不是 {{.Xxx}}，否则标题里的引号会破坏 JSON。": "Empty uses the built-in template. Variables: {{.Title}} {{.Batch}} {{.Count}} {{.HomeURL}} {{.SentAt}}, and .Name/.VulnClass/.Severity/.Summary/.Assets/.DetailURL/.StatusLabel inside range .Items. Insert strings using {{json .Xxx}} rather than {{.Xxx}} to keep quotes from breaking JSON."
};
export function notificationText(value: string, locale: "zh-CN" | "en") { return locale === "en" ? notificationEnglish[value] ?? value : value; }

export function ConfigField({
  def,
  value,
  isSecret,
  onChange,
}: {
  def: FieldDef;
  value: unknown;
  isSecret: boolean;
  onChange: (v: unknown) => void;
}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const id = `n-cfg-${def.key}`;
  const raw = asText(value);
  // 后端回显的掩码值：形如 "__masked__:…abc123"，尾部是原值的可辨识片段。
  const masked = isSecret && raw.startsWith("__masked__");
  const maskedTail = masked ? (raw.split("…")[1] ?? "") : "";

  if (def.kind === "switch") {
    return (
      <div className="flex items-center gap-2 text-sm">
        <Switch checked={value === true} onCheckedChange={onChange} aria-label={notificationText(def.label, locale)} />
        {notificationText(def.label, locale)}
        {def.help && <span className="text-muted-foreground">（{notificationText(def.help, locale)}）</span>}
      </div>
    );
  }

  if (def.kind === "select") {
    return (
      <div className="grid gap-2">
        <Label>{notificationText(def.label, locale)}</Label>
        <Select value={raw || def.options?.[0]?.value} onValueChange={onChange}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {(def.options ?? []).map((o) => (
              <SelectItem key={o.value} value={o.value}>
                {notificationText(o.label, locale)}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    );
  }

  // 控件按字段类型分派。用 if 链而不是嵌套三元，是因为这里要区分四种控件，
  // 三层三元读起来已经要停下来数括号了。
  function control() {
    if (def.kind === "textarea" || def.kind === "kv") {
      return (
        <Textarea
          id={id}
          className="font-mono"
          placeholder={def.placeholder}
          value={masked ? "" : raw}
          onChange={(e) => onChange(e.target.value)}
        />
      );
    }
    if (def.kind === "list") {
      return (
        <Input
          id={id}
          value={Array.isArray(value) ? (value as string[]).join(", ") : raw}
          onChange={(e) => onChange(e.target.value)}
          placeholder={def.placeholder}
        />
      );
    }
    return (
      <Input
        id={id}
        className={def.kind === "text" ? "font-mono" : ""}
        type={inputType(def.kind)}
        placeholder={def.placeholder}
        value={masked ? "" : raw}
        onChange={(e) => onChange(e.target.value)}
      />
    );
  }

  const hint = masked ? (
    <p className="text-muted-foreground flex items-center gap-1 text-xs">
      <CheckIcon className="size-3" />
      {tr(`已保存${maskedTail ? `（尾号 ${maskedTail}）` : ""} · 填入新值即覆盖，清空则删除该项`, `Saved${maskedTail ? ` (ending ${maskedTail})` : ""} · Enter a new value to replace; clearing removes it`)}
    </p>
  ) : (
    def.help && <p className="text-muted-foreground text-xs">{notificationText(def.help, locale)}</p>
  );

  return (
    <div className="grid gap-2">
      <Label htmlFor={id}>{notificationText(def.label, locale)}</Label>
      {control()}
      {hint}
    </div>
  );
}

// FilterSummary 把过滤条件摘要成一行，让卡片不用展开就能看出这个渠道推什么。
export function FilterSummary({ filter }: { filter: NotificationFilter }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const parts: string[] = [];
  if (filter.min_severity) {
    parts.push(notificationText(SEVERITY_OPTIONS.find((o) => o.value === filter.min_severity)?.label ?? filter.min_severity, locale));
  }
  if (filter.vulnclass_include?.length) parts.push(tr(`类型含 ${filter.vulnclass_include.length} 词`, `${filter.vulnclass_include.length} included keywords`));
  if (filter.vulnclass_exclude?.length) parts.push(tr(`排除 ${filter.vulnclass_exclude.length} 词`, `${filter.vulnclass_exclude.length} excluded keywords`));
  if (filter.task_ids?.length) parts.push(tr(`${filter.task_ids.length} 个任务`, `${filter.task_ids.length} tasks`));
  if (filter.asset_ids?.length) parts.push(tr(`${filter.asset_ids.length} 个资产`, `${filter.asset_ids.length} assets`));
  if (filter.on_status_change) parts.push(tr("含状态变更", "Includes status changes"));
  if (parts.length === 0) {
    return <p className="text-muted-foreground text-sm">{tr("全部漏洞", "All findings")}</p>;
  }
  return <p className="text-muted-foreground text-sm">{parts.join(" · ")}</p>;
}
