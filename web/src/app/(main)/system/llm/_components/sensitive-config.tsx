"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";
import { RefreshCwIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { http } from "@/lib/api";

type Finding = { source: string; field: string; code: string; profile_id?: string };
type Status = {
  file_state: string;
  file_permissions: string;
  database_source: string;
  profiles_state: string;
  profiles_checked: number;
  findings: Finding[];
};

const states: Record<string, string> = {
  missing: "不存在或未配置", readable: "已读取", unreadable: "无法读取",
  invalid: "格式无效", not_inspected: "未检查", unavailable: "无法检查",
  complete: "已完成", partial: "部分完成", environment: "环境变量",
  file_dsn: "配置文件 DSN", file_fields: "配置文件字段",
  windows_acl_not_inspected: "Windows ACL 未检查",
  group_or_other_read: "存在组或其他用户读取权限",
  no_group_or_other_read: "未发现组或其他用户读取位",
};
const labels: Record<string, string> = {
  config_file: "配置文件", environment: "环境变量", model_profile: "模型配置",
  "database.password": "数据库密码", "database.dsn": "数据库连接字符串",
  base_url: "模型地址", proxy: "代理地址",
  plaintext_password: "存在明文密码字段",
  url_userinfo: "地址含认证信息",
  sensitive_query: "地址含敏感查询参数",
  http_non_loopback: "非回环模型地址使用 HTTP",
  url_not_inspected: "地址格式未识别，尚未检查",
  query_not_inspected: "查询参数解析不完整",
};
const statesEn: Record<string, string> = {"missing":"Missing or unconfigured","readable":"Readable","unreadable":"Unreadable","invalid":"Invalid","not_inspected":"Not inspected","unavailable":"Unavailable","complete":"Complete","partial":"Partial","environment":"Environment","file_dsn":"File DSN","file_fields":"File fields","windows_acl_not_inspected":"Windows ACL not inspected","group_or_other_read":"Group or other read access","no_group_or_other_read":"No group or other read bits found"};
const labelsEn: Record<string, string> = {"config_file":"Configuration file","environment":"Environment","model_profile":"Model profile","database.password":"Database password","database.dsn":"Database connection string","base_url":"Model URL","proxy":"Proxy URL","plaintext_password":"Plaintext password field","url_userinfo":"URL contains credentials","sensitive_query":"Sensitive query parameter","http_non_loopback":"Non-loopback HTTP model URL","url_not_inspected":"Unrecognized URL format; not inspected","query_not_inspected":"Incomplete query parsing"};

export function SensitiveConfigPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const label = (value: string) => (locale === "en" ? labelsEn : labels)[value] ?? value;
  const state = (value: string) => (locale === "en" ? statesEn : states)[value] ?? tr("未知", "Unknown");
  const [snapshot, setSnapshot] = React.useState<Status | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  async function inspect() {
    setBusy(true);
    setError("");
    setSnapshot(null);
    try {
      const result = await http<Status>("/privacy/sensitive-config");
      if (!result || !Array.isArray(result.findings) || typeof result.profiles_checked !== "number") {
        throw new Error("invalid diagnostic");
      }
      setSnapshot(result);
    } catch {
      setError(tr("无法读取敏感配置诊断，请确认后端版本。", "Unable to read sensitive configuration diagnostics; check backend version."));
    } finally {
      setBusy(false);
    }
  }
  return (
    <Card className="min-w-0">
      <CardContent className="grid gap-4">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h2 className="text-base font-semibold">{tr("敏感配置诊断", "Sensitive configuration diagnostics")}</h2>
            <p className="mt-1 text-sm text-muted-foreground">{tr("检查数据库配置、模型地址及代理中的凭据形式。只返回状态和风险类型，不显示原始地址或密码。", "Inspect credential formats in database, model and proxy configuration; return status and risk categories only, not original addresses or passwords.")}</p>
          </div>
          <Button size="sm" variant="outline" disabled={busy} onClick={() => void inspect()}>
            <RefreshCwIcon className={busy ? "animate-spin" : ""} />
            {busy ? tr("检查中…", "Checking…") : tr("检查敏感配置", "Inspect sensitive configuration")}
          </Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {snapshot && (
          <div className="grid gap-3" aria-live="polite">
            <div className="grid gap-1 text-sm">
              <p>{tr("配置文件：", "Config file: ")}{state(snapshot.file_state)} · {state(snapshot.file_permissions)}</p>
              <p>{tr("数据库配置来源：", "Database source: ")}{state(snapshot.database_source)}</p>
              <p>{tr("已保存模型：", "Saved profiles: ")}{state(snapshot.profiles_state)} · {tr(`已检查 ${snapshot.profiles_checked} 个`, `${snapshot.profiles_checked} checked`)}</p>
              <p>{tr(`发现 ${snapshot.findings.length} 项需关注或未完整检查的配置`, `${snapshot.findings.length} items requiring attention or incomplete inspection`)}</p>
            </div>
            {snapshot.findings.length ? (
              <div className="divide-y rounded-lg border">
                {snapshot.findings.map((f, i) => (
                  <div key={i} className="grid gap-1 px-3 py-2 text-sm">
                    <span className="font-medium">{label(f.code)}</span>
                    <span className="text-muted-foreground">{label(f.source)}{f.profile_id ? tr(" · 配置 #", " · Profile #") + f.profile_id : ""} · {label(f.field)}</span>
                  </div>
                ))}
              </div>
            ) : <p className="text-sm text-muted-foreground">{tr("在已检查范围内未发现已知格式的问题；这不是完整的安全证明。", "No known format issues within the inspected scope; this is not a complete security assurance.")}</p>}
          </div>
        )}
        <p className="text-xs leading-relaxed text-muted-foreground">{tr("认证信息存在并不代表已泄露。环境变量、文件权限、数据库备份及授权访问仍需分别保护。HTTP 提示针对配置地址，不证明实际网络路径；Windows ACL、特殊编码、非 URL 形式连接串及任意未标记秘密不在完整检查范围内。此检查不发送网络探测，不修改配置。", "Credentials being present does not establish leakage. Protect environment variables, file permissions, backups and authorized access separately. HTTP findings concern configured URLs, not actual network routes. Windows ACLs, special encodings, non-URL connection strings and arbitrary unmarked secrets are not fully inspected. No network probes or configuration changes.")}</p>
      </CardContent>
    </Card>
  );
}
