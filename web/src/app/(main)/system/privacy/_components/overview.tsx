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
  execution_network: {
    mode: string; configuration: string; verification: string; scope: string;
    final_exit: string; actual_egress: string; observed_at: string | null;
    fail_closed: string; verify_before_task: boolean; covered_processes: string[] | null;
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
  off: "关闭", B: "模式 B", other: "其他模式",
  configured_only: "已配置，未验证", exit_ip_match: "出口 IP 匹配",
  exit_ip_mismatch: "出口 IP 不匹配", probe_failed: "出口探测失败",
  matched_expected: "匹配预期出口", mismatch: "出口不匹配",
  matches_expected: "匹配预期出口", not_verified: "未验证",
};
const labelsEn: Record<string, string> = {"missing":"Missing or unconfigured","readable":"Readable","unreadable":"Unreadable","invalid":"Invalid","not_inspected":"Not inspected","unavailable":"Unavailable","complete":"Complete","partial":"Partial","environment":"Environment","file_dsn":"File DSN","file_fields":"File fields","windows_acl_not_inspected":"Windows ACL not inspected","group_or_other_read":"Group or other read access","no_group_or_other_read":"No group or other read bits found","loopback":"Loopback only","all_interfaces":"All interfaces","specific_interface":"Specific interface","unknown":"Unknown","runtime_listener":"Actual listener","configured":"Configured","not_tested":"Not verified","off":"Off","B":"Mode B","other":"Other mode","configured_only":"Configured, not verified","exit_ip_match":"Exit IP matched","exit_ip_mismatch":"Exit IP mismatch","probe_failed":"Egress probe failed","matched_expected":"Matches expected exit","mismatch":"Exit mismatch","matches_expected":"Matches expected exit","not_verified":"Not verified"};

export function PrivacyOverviewPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const label = (s: string) => (locale === "en" ? labelsEn : labels)[s] ?? tr("未知", "Unknown");
  const [snapshot, setSnapshot] = React.useState<Overview | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  const [expectedExitIP, setExpectedExitIP] = React.useState("");
  const [probeBusy, setProbeBusy] = React.useState(false);
  const [probeResult, setProbeResult] = React.useState<{ ok: boolean; message: string } | null>(null);

  async function verifyModeB() {
    setProbeBusy(true);
    setProbeResult(null);
    try {
      const result = await http<{ ok: boolean; status: string }>("/privacy/execution-network/verify", {
        method: "POST",
        body: JSON.stringify({ expected_exit_ip: expectedExitIP.trim() }),
      });
      setProbeResult({
        ok: result.ok,
        message: result.ok
          ? tr("隔离执行器观察到的出口 IP 与预期一致。逐跳路径仍未独立验收。", "The isolated executor observed the expected exit IP. Individual hops are not independently verified.")
          : tr("观察到的出口 IP 与预期不一致；请检查链路，不要启动网络任务。", "The observed exit IP did not match. Check the chain before starting network tasks."),
      });
      await inspect();
    } catch (err) {
      setProbeResult({
        ok: false,
        message: err instanceof Error ? err.message : tr("出口验证失败。", "Egress verification failed."),
      });
      await inspect();
    } finally {
      setProbeBusy(false);
    }
  }

  async function inspect() {
    setBusy(true);
    setError("");
    setSnapshot(null);
    try {
      const result = await http<Overview>("/privacy/overview");
      if (!result?.listener || !result.model_keys || !result.config || !result.outbound || !result.execution_network) {
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
              <section className="grid gap-2 rounded-lg border p-4 text-sm">
                <h3 className="font-semibold">{tr("ARTEX 任务出口 · 模式 B", "ARTEX task egress · Mode B")}</h3>
                <p>{tr("模式：", "Mode: ")}{label(snapshot.execution_network.mode)} · {tr("配置：", "Configuration: ")}{label(snapshot.execution_network.configuration)}</p>
                <p>{tr("出口验证：", "Egress verification: ")}{label(snapshot.execution_network.verification)} · {tr("最终出口：", "Final exit: ")}{label(snapshot.execution_network.final_exit)}</p>
                <p>{tr("失败关闭：", "Fail-closed: ")}{label(snapshot.execution_network.fail_closed)} · {tr("任务开始前验证：", "Verify before task: ")}{snapshot.execution_network.verify_before_task ? tr("已启用", "Enabled") : tr("未启用", "Not enabled")}</p>
                <p className="text-muted-foreground">{tr("覆盖范围：", "Covered scope: ")}{(snapshot.execution_network.covered_processes ?? []).join(", ") || tr("无", "None")}</p>
                {snapshot.execution_network.observed_at && <p className="text-muted-foreground">{tr("最近观测：", "Last observed: ")}{new Date(snapshot.execution_network.observed_at).toLocaleString(locale)}</p>}
                {snapshot.execution_network.mode === "B" && (
                  <div className="mt-2 grid gap-2 border-t pt-3">
                    <label className="grid gap-1">
                      <span>{tr("预期最终出口 IPv4（VPS1 或 VPS2）", "Expected final exit IPv4 (VPS1 or VPS2)")}</span>
                      <input
                        className="w-full rounded-md border bg-background px-3 py-2 font-mono text-sm"
                        type="text"
                        inputMode="numeric"
                        autoComplete="off"
                        placeholder={tr("填写预期出口公网 IPv4", "Expected public exit IPv4")}
                        value={expectedExitIP}
                        onChange={(event) => setExpectedExitIP(event.target.value)}
                      />
                    </label>
                    <Button type="button" size="sm" variant="outline"
                      disabled={probeBusy || busy || !expectedExitIP.trim()}
                      onClick={() => void verifyModeB()}>
                      {probeBusy ? tr("正在验证隔离执行器出口…", "Verifying isolated executor egress…") : tr("验证 ARTEX 任务出口", "Verify ARTEX task egress")}
                    </Button>
                    {probeResult && <p role="status" className={probeResult.ok ? "text-emerald-600" : "text-destructive"}>{probeResult.message}</p>}
                    <p className="text-xs text-muted-foreground">{tr("仅在点击按钮后，从 Linux Docker tool-executor 经固定内部网关发起一次 HTTPS 出口比对；不自动探测、不修改配置。它不证明每一跳、主机 DNS/IPv6、控制面或外部 MCP 已隔离。", "Only after clicking, the Linux Docker tool-executor makes one HTTPS exit comparison through the fixed private gateway. No automatic probe or config change. This does not prove each hop, host DNS/IPv6, control-plane or external MCP isolation.")}</p>
                  </div>
                )}
              </section>
            </div>
          </div>
        )}
        <p className="text-xs leading-relaxed text-muted-foreground">{tr("各项结果独立判断；未检查或无法读取不等于安全。密钥检查仅涵盖模型配置字段，网络统计仅涵盖当前进程近期 HTTP 请求尝试，非完整流量或任务成功率。不会自动修改配置、监听或数据。", "Interpret each result independently; unavailable or uninspected does not mean safe. Key inspection covers model profile fields only, and network statistics cover recent HTTP attempts in this process, not full traffic or task success. No automatic changes to configuration, listeners or data.")}</p>
      </CardContent>
    </Card>
  );
}
