"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";
import Link from "next/link";
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { http } from "@/lib/api";

type Overview = {
  checked_at: string;
  listener: { scope: string; binding_source: string; public_reachability: string };
  model_keys: {
    state: string; key_state: string; total: number; empty: number;
    plaintext_only: number; encrypted_only: number; mixed: number;
    verified: number; failed: number; unchecked: number;
  };
  config: {
    file_state: string; file_permissions: string; database_source: string;
    profiles_state: string; profiles_checked: number; findings: number; incomplete_findings: number;
  };
  outbound: {
    total: number; retained: number; capacity: number; http_responses: number;
    http_non_2xx: number; transport_failures: number; canceled: number;
    timeouts: number; proxy_configured: number;
  };
};

const labels: Record<string, string> = {
  loopback: "仅本机回环", all_interfaces: "所有网络接口",
  specific_interface: "指定网络接口", unknown: "未知",
  runtime_listener: "实际监听", configured: "已配置", not_inspected: "未检查",
  not_tested: "未验证", complete: "已完成", partial: "部分完成",
  unavailable: "无法读取", missing: "不存在或未配置", invalid: "格式无效",
  readable: "已读取", unreadable: "无法读取",
  environment: "环境变量", file_dsn: "配置文件 DSN", file_fields: "配置文件字段",
  windows_acl_not_inspected: "Windows ACL 未检查",
  group_or_other_read: "存在组或其他用户读取位",
  no_group_or_other_read: "未发现组或其他用户读取位",
};
const labelsEn: Record<string, string> = {"missing":"Missing or unconfigured","readable":"Readable","unreadable":"Unreadable","invalid":"Invalid","not_inspected":"Not inspected","unavailable":"Unavailable","complete":"Complete","partial":"Partial","environment":"Environment","file_dsn":"File DSN","file_fields":"File fields","windows_acl_not_inspected":"Windows ACL not inspected","group_or_other_read":"Group or other read access","no_group_or_other_read":"No group or other read bits found","loopback":"Loopback only","all_interfaces":"All interfaces","specific_interface":"Specific interface","unknown":"Unknown","runtime_listener":"Actual listener","configured":"Configured","not_tested":"Not verified"};

export function PrivacyOverviewPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const label = (s: string) => (locale === "en" ? labelsEn : labels)[s] ?? tr("未知", "Unknown");
  const [snapshot, setSnapshot] = React.useState<Overview | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");

  async function inspect() {
    setBusy(true);
    setError("");
    setSnapshot(null);
    try {
      const result = await http<Overview>("/privacy/overview");
      if (!result?.listener || !result.model_keys || !result.config || !result.outbound) {
        throw new Error("invalid overview");
      }
      setSnapshot(result);
    } catch {
      setError(tr("无法读取隐私总览，请确认后端版本并重试。", "Unable to read privacy overview; check backend version and retry."));
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card className="min-w-0">
      <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-3">
        <div className="grid gap-1">
          <CardTitle className="text-base">{tr("统一隐私诊断", "Unified privacy diagnostics")}</CardTitle>
          <CardDescription>{tr("按需汇总监听、密钥存储、敏感配置及模型 HTTP 审计。", "Summarize listeners, key storage, sensitive configuration and model HTTP audit on demand.")}</CardDescription>
        </div>
        <Button type="button" variant="outline" size="sm" disabled={busy} onClick={() => void inspect()}>
          <RefreshCw className={busy ? "size-4 animate-spin" : "size-4"} />
          {busy ? tr("检查中…", "Checking…") : tr("检查全部", "Inspect all")}
        </Button>
      </CardHeader>
      <CardContent className="grid gap-4">
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {!snapshot && !error && <p className="text-sm text-muted-foreground">{tr("点击检查全部获取本机诊断快照；不发送外部探测或模型请求。", "Inspect all for a local diagnostic snapshot; no external probes or model requests.")}</p>}
        {snapshot && (
          <div className="grid gap-3" aria-live="polite">
            <p className="text-xs text-muted-foreground">{tr("检查时间：", "Checked at: ")}{new Date(snapshot.checked_at).toLocaleString(locale)}</p>
            <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
              <section className="grid gap-2 rounded-lg border p-4 text-sm">
                <h3 className="font-semibold">{tr("管理接口", "Management interface")}</h3>
                <p>{label(snapshot.listener.scope)} · {label(snapshot.listener.binding_source)}</p>
                <p className="text-muted-foreground">{tr("公网可达性：", "Public reachability: ")}{label(snapshot.listener.public_reachability)}</p>
              </section>
              <section className="grid gap-2 rounded-lg border p-4 text-sm">
                <h3 className="font-semibold">{tr("模型密钥", "Model keys")}</h3>
                {snapshot.model_keys.state !== "complete" ? (
                  <p className="text-muted-foreground">{tr("无法读取数据库，密钥覆盖与完整性未验证。", "Database unavailable; key coverage and integrity not verified.")}</p>
                ) : (
                  <>
                    <p>{tr("密钥：", "Key: ")}{label(snapshot.model_keys.key_state)} · 配置 {snapshot.model_keys.total}</p>
                    <p>仅密文 {snapshot.model_keys.encrypted_only} · 仅明文 {snapshot.model_keys.plaintext_only} · 混合 {snapshot.model_keys.mixed} · 无密钥 {snapshot.model_keys.empty}</p>
                    <p className={snapshot.model_keys.failed > 0 ? "text-destructive" : "text-muted-foreground"}>
                      认证通过 {snapshot.model_keys.verified} · 失败 {snapshot.model_keys.failed} · 未校验 {snapshot.model_keys.unchecked}
                    </p>
                  </>
                )}
                <Link href="/system/llm" className="font-medium underline underline-offset-4">{tr("打开模型配置 / 密钥保护", "Open model configuration / key protection")}</Link>
              </section>
              <section className="grid gap-2 rounded-lg border p-4 text-sm">
                <h3 className="font-semibold">{tr("敏感配置", "Sensitive configuration")}</h3>
                <p>{tr("配置文件：", "Config file: ")}{label(snapshot.config.file_state)} · {label(snapshot.config.file_permissions)}</p>
                <p>{tr("数据库来源：", "Database source: ")}{label(snapshot.config.database_source)}</p>
                <p>{tr("模型检查：", "Profile inspection: ")}{label(snapshot.config.profiles_state)} · {snapshot.config.profiles_checked} 个</p>
                <p>需关注 {snapshot.config.findings} 项 · 其中解析不完整 {snapshot.config.incomplete_findings} 项</p>
                <Link href="/system/llm" className="font-medium underline underline-offset-4">{tr("查看敏感配置诊断", "View sensitive configuration diagnostics")}</Link>
              </section>
              <section className="grid gap-2 rounded-lg border p-4 text-sm">
                <h3 className="font-semibold">{tr("模型网络活动", "Model network activity")}</h3>
                <p>进程累计 {snapshot.outbound.total} 次 · 近期 {snapshot.outbound.retained}/{snapshot.outbound.capacity} 条</p>
                <p>HTTP 响应 {snapshot.outbound.http_responses}（非 2xx：{snapshot.outbound.http_non_2xx}）</p>
                <p>传输失败 {snapshot.outbound.transport_failures} · 超时 {snapshot.outbound.timeouts} · 取消 {snapshot.outbound.canceled}</p>
                <p>近期配置代理 {snapshot.outbound.proxy_configured} 次，不代表实际网络出口。</p>
                <Link href="/system/llm" className="font-medium underline underline-offset-4">{tr("查看模型连接审计", "View model connection audit")}</Link>
              </section>
            </div>
          </div>
        )}
        <p className="text-xs leading-relaxed text-muted-foreground">{tr("各项结果独立判断；未检查或无法读取不等于安全。密钥检查仅涵盖模型配置字段，网络统计仅涵盖当前进程近期 HTTP 请求尝试，非完整流量或任务成功率。不会自动修改配置、监听或数据。", "Interpret each result independently; unavailable or uninspected does not mean safe. Key inspection covers model profile fields only, and network statistics cover recent HTTP attempts in this process, not full traffic or task success. No automatic changes to configuration, listeners or data.")}</p>
      </CardContent>
    </Card>
  );
}
