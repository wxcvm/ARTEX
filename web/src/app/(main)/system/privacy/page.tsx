"use client";

import * as React from "react";
import Link from "next/link";
import { AlertTriangle, CheckCircle2, Globe2, LockKeyhole, RefreshCw, ShieldCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { PrivacyOverviewPanel } from "./_components/overview";
import { StoragePermissionsPanel } from "./_components/storage-permissions";

type Exposure = Awaited<ReturnType<typeof api.privacyExposure>>;

const SCOPE_LABELS: Record<Exposure["scope"], string> = {
  loopback: "仅本机监听",
  all_interfaces: "监听所有网络接口",
  specific_interface: "监听指定网络接口",
  unknown: "无法确定监听范围",
};

export default function PrivacyPage() {
  const [exposure, setExposure] = React.useState<Exposure | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");

  const refresh = React.useCallback(() => {
    setLoading(true);
    setError("");
    api.privacyExposure()
      .then(setExposure)
      .catch((err: unknown) => {
        setExposure(null);
        setError(err instanceof Error ? err.message : "无法获取监听信息");
      })
      .finally(() => setLoading(false));
  }, []);

  React.useEffect(() => {
    refresh();
  }, [refresh]);

  const loopback = exposure?.scope === "loopback";
  const ScopeIcon = loopback ? CheckCircle2 : AlertTriangle;

  return (
    <div className="flex flex-1 flex-col gap-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">隐私与暴露面</h1>
          <p className="text-sm text-muted-foreground">本机诊断 · 不联系外部探测服务</p>
        </div>
        <Button type="button" variant="outline" size="sm" onClick={refresh} disabled={loading}>
          <RefreshCw className="size-4" /> {loading ? "读取中…" : "刷新"}
        </Button>
      </div>

      <PrivacyOverviewPanel />
      <StoragePermissionsPanel />

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <ShieldCheck className="size-4" /> ARTEX 管理接口
          </CardTitle>
          <CardDescription>后端优先报告实际监听地址；未观测时使用配置地址，不扫描公网或修改网络。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {error ? <p role="alert" className="text-sm text-destructive">读取失败：{error}</p> : null}
          {exposure ? (
            <>
              <div className="flex flex-wrap items-center gap-2">
                <ScopeIcon className={loopback ? "size-5 text-emerald-600" : "size-5 text-amber-600"} />
                <span className="font-medium">{SCOPE_LABELS[exposure.scope]}</span>
              </div>
              <div className="rounded-md border p-3 text-sm">
                <p className="text-muted-foreground">HTTP 监听地址</p>
                <code className="break-all font-mono">{exposure.listen_address || "未提供"}</code>
              </div>
              <p className="text-sm text-muted-foreground">
                {loopback
                  ? "应用自身监听在本机回环地址；仍可能被反向代理、隧道或其他转发服务间接发布。"
                  : "非本机回环监听或监听范围未知，可能允许其他设备连接；是否实际可达还取决于主机网络和防火墙。"}
              </p>
            </>
          ) : !error ? <p className="text-sm text-muted-foreground">{loading ? "正在读取…" : "暂无结果"}</p> : null}
        </CardContent>
      </Card>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Globe2 className="size-4" /> 公网可达性
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            <p className="font-medium">未验证</p>
            <p className="text-muted-foreground">
              监听地址不能说明公网是否可以直接访问。Docker 端口发布、Windows 防火墙、路由器端口映射、
              VPS 安全组、反向代理和 IPv6 入站规则都需要分别检查。
            </p>
            <p className="text-muted-foreground">这个页面不会向外部服务发送探测请求，也不会声称“绝对安全”。</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <LockKeyhole className="size-4" /> 数据与外部连接
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            <p className="text-muted-foreground">
              HTTP 管理接口受 ARTEX 登录认证保护。第三方 LLM、Web 搜索、代理及流量录制取决于实际配置。
              本页不会推断这些服务是否保存了日志。
            </p>
            <Link href="/system/settings" className="font-medium underline underline-offset-4">
              查看流量录制、搜索和代理设置
            </Link>
            <p className="text-muted-foreground">Purge/DSH 的实际监听与远程认证由对应宿主管理，暂未在此页面读取。</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
