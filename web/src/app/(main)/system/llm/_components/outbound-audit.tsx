"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";
import { RefreshCwIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { http } from "@/lib/api";

type OutboundEvent = {
  sequence: number;
  at: string;
  scheme: string;
  host: string;
  port: string;
  proxy_configured: boolean;
  status: number;
  result: string;
  elapsed_ms: number;
};

type OutboundSnapshot = {
  capacity: number;
  total: number;
  events: OutboundEvent[];
};

const resultLabelsEn: Record<string, string> = { http_response: "HTTP response", transport_error: "Transport failure", canceled: "Canceled", timeout: "Timeout" };
const resultLabels: Record<string, string> = {
  http_response: "收到 HTTP 响应",
  transport_error: "传输失败",
  canceled: "请求取消",
  timeout: "请求超时",
};

export function OutboundAuditPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [snapshot, setSnapshot] = React.useState<OutboundSnapshot | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");

  async function refresh() {
    setBusy(true);
    setError("");
    try {
      const result = await http<OutboundSnapshot>("/privacy/model-outbound");
      if (!result || !Array.isArray(result.events) || !Number.isSafeInteger(result.total)) {
        throw new Error("invalid audit snapshot");
      }
      setSnapshot(result);
    } catch {
      setError(tr("无法读取出站审计，请确认 ARTEX 后端版本。", "Unable to read outbound audit. Check the ARTEX backend version."));
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card className="min-w-0">
      <CardContent className="grid gap-4">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h2 className="text-base font-semibold">{tr("模型出站连接审计", "Model outbound connection audit")}</h2>
            <p className="mt-1 text-sm text-muted-foreground">
              {tr("查看模型 HTTP 请求的目标主机、代理配置、响应状态及耗时。只读取现有记录，不发送模型请求。", "Inspect destination hosts, proxy configuration, response status and duration of model HTTP requests. Reads existing records only; sends no model requests.")}
            </p>
          </div>
          <Button size="sm" variant="outline" disabled={busy} onClick={() => void refresh()}>
            <RefreshCwIcon className={busy ? "animate-spin" : ""} />
            {busy ? tr("读取中…", "Reading…") : tr("刷新审计", "Refresh audit")}
          </Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {snapshot && (
          <div className="grid gap-3" aria-live="polite">
            <p className="text-sm text-muted-foreground">
              {tr(`本进程已完成 ${snapshot.total} 次 HTTP 请求尝试 · 保留最近 ${snapshot.events.length} / ${snapshot.capacity} 条`, `${snapshot.total} HTTP attempts in this process · retaining ${snapshot.events.length} / ${snapshot.capacity} recent records`)}
            </p>
            {snapshot.events.length === 0 ? (
              <p className="rounded-lg border border-dashed p-4 text-sm text-muted-foreground">{tr("尚未记录模型 HTTP 请求。", "No model HTTP requests recorded.")}</p>
            ) : (
              <div className="divide-y rounded-lg border">
                {snapshot.events.map((e) => (
                  <div key={e.sequence} className="grid gap-1 px-3 py-3 text-sm">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <span className="min-w-0 break-all font-medium">
                        {e.scheme ? `${e.scheme}://` : ""}{e.host || tr("未知目标", "Unknown destination")}{e.port ? tr(` · 端口 ${e.port}`, ` · port ${e.port}`) : ""}
                      </span>
                      <span className="text-muted-foreground">#{e.sequence} · {new Date(e.at).toLocaleString(locale)}</span>
                    </div>
                    <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
                      <span>{e.proxy_configured ? tr("已配置代理", "Proxy configured") : tr("未配置模型代理", "No model proxy")}</span>
                      <span>{(locale === "en" ? resultLabelsEn : resultLabels)[e.result] ?? tr("未知结果", "Unknown result")}</span>
                      {e.status > 0 && <span>HTTP {e.status}</span>}
                      <span>{e.elapsed_ms} ms</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
        <p className="text-xs leading-relaxed text-muted-foreground">
          {tr("这是应用层 HTTP 尝试记录，不是 TCP 连接清单。代理状态只表示模型配置，不能证明实际网络出口；耗时截至响应头或传输错误，不包含完整流式响应。记录仅保留在当前后端进程内，不包含认证头、请求路径、查询参数或正文。外部工具及其他进程的流量不在此范围内。", "These are application-level HTTP attempts, not TCP connections. Proxy status indicates model configuration, not the actual egress route. Duration ends at response headers or transport error, not the full stream. Records exist only in the current backend process and exclude authentication headers, paths, queries and bodies. External tools and other processes are out of scope.")}
        </p>
      </CardContent>
    </Card>
  );
}
