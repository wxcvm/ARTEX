"use client";

import * as React from "react";
import { ArrowUpRightIcon, DownloadIcon, RefreshCwIcon } from "lucide-react";
import { toast } from "sonner";

import { api } from "@/lib/api";
import { OfflineToolsPanel } from "@/components/offline-tools-panel";
import type { ExternalToolInventory, ExternalToolRelease, ExternalToolStatus } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

const originName: Record<string, string> = {
  custom: "自定义目录",
  managed: "ARTEX 工具目录",
  bundled: "程序附带",
  path: "系统 PATH",
};

function ToolItem({
  tool, host, onRefresh,
}: {
  tool: ExternalToolStatus;
  host: string;
  onRefresh: () => Promise<void>;
}) {
  const [release, setRelease] = React.useState<ExternalToolRelease | null>(null);
  const [checking, setChecking] = React.useState(false);
  const [installing, setInstalling] = React.useState(false);
  const [managing, setManaging] = React.useState(false);
  const [error, setError] = React.useState("");
  const supported = host === "windows/amd64";
  const alreadyManaged = !!release && tool.managed_version === release.version;

  async function checkRelease() {
    setChecking(true);
    setError("");
    try {
      setRelease(await api.externalToolRelease(tool.id));
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setChecking(false);
    }
  }

  async function install() {
    if (!release) return;
    const sizeMiB = (release.size_bytes / 1048576).toFixed(1);
    if (!window.confirm(
      "从官方 GitHub 下载 " + tool.name + " " + release.version + "（" + sizeMiB + " MiB）？\n" +
      "下载后校验 SHA-256，再安装到 ARTEX 用户工具目录；不会覆盖自定义目录、修改 PATH 或绑定 Agent。"
    )) return;
    setInstalling(true);
    setError("");
    try {
      await api.installExternalTool(tool.id, release.version, release.sha256);
      toast.success(tool.name + " " + release.version + " 已安装并通过下载摘要校验");
      await onRefresh();
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setInstalling(false);
    }
  }

  async function rollback() {
    if (!tool.previous_version || !tool.managed_version) return;
    if (!window.confirm(
      "将 " + tool.name + " 从 " + tool.managed_version + " 恢复到 " + tool.previous_version +
      "？当前版本会保留为上一版；正在运行此程序时请先结束运行。"
    )) return;
    setManaging(true);
    setError("");
    try {
      const result = await api.rollbackExternalTool(tool.id);
      toast.success(tool.name + " 已恢复到 " + result.version);
      await onRefresh();
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setManaging(false);
    }
  }

  async function uninstall() {
    if (!tool.managed_version) return;
    if (!window.confirm(
      "卸载 ARTEX 管理的 " + tool.name + " " + tool.managed_version +
      " 和保留的上一版？\n自定义目录及 PATH 中的工具不会被删除。"
    )) return;
    setManaging(true);
    setError("");
    try {
      await api.uninstallExternalTool(tool.id);
      toast.success("已卸载 ARTEX 管理的 " + tool.name);
      setRelease(null);
      await onRefresh();
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setManaging(false);
    }
  }

  return (
    <div className="flex flex-col gap-2 rounded-lg border p-4">
      <div className="flex flex-wrap items-center gap-2">
        <span className="font-mono text-sm font-semibold">{tool.name}</span>
        {tool.status === "installed" ? (
          <Badge variant="secondary">版本检查通过</Badge>
        ) : tool.status === "unverified" ? (
          <Badge variant="outline">已发现 · 检查失败</Badge>
        ) : (
          <Badge variant="outline">未发现</Badge>
        )}
        {tool.origin && (
          <span className="text-muted-foreground ml-auto text-xs">
            {originName[tool.origin] ?? tool.origin}
          </span>
        )}
      </div>
      <p className="text-muted-foreground text-xs">{tool.description}</p>
      {tool.version && <div className="break-all font-mono text-xs">{tool.version}</div>}
      {tool.managed_version && (
        <div className="space-y-2">
          <div className="text-xs">
            ARTEX 管理版本：<span className="font-mono">{tool.managed_version}</span>
            {tool.previous_version && (
              <span className="text-muted-foreground ml-2">
                上一版 <span className="font-mono">{tool.previous_version}</span>
              </span>
            )}
          </div>
          {supported && (
            <div className="flex flex-wrap gap-2">
              {tool.previous_version && (
                <Button variant="outline" size="sm" disabled={managing || installing}
                  onClick={() => void rollback()}>恢复上一版</Button>
              )}
              <Button variant="outline" size="sm" disabled={managing || installing}
                onClick={() => void uninstall()}>
                {managing ? "处理中…" : "卸载管理版本"}
              </Button>
            </div>
          )}
        </div>
      )}
      {tool.path && (
        <div className="text-muted-foreground break-all font-mono text-xs" title={tool.path}>
          {tool.path}
        </div>
      )}
      {tool.detail && <p className="text-muted-foreground text-xs">{tool.detail}</p>}
      {tool.managed_version && tool.origin === "custom" && (
        <p className="text-muted-foreground text-xs">自定义目录优先于 ARTEX 管理目录，界面显示的是实际优先匹配的版本。</p>
      )}
      <div className="mt-auto space-y-2 pt-2">
        {supported ? (
          <>
            <Button size="sm" variant="outline" disabled={checking || installing || managing} onClick={() => void checkRelease()}>
              <RefreshCwIcon className={checking ? "animate-spin" : ""} />
              {checking ? "查询官方发布…" : "查看官方最新版"}
            </Button>
            {release && (
              <div className="space-y-2 rounded-md bg-muted/50 p-2 text-xs">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span>官方 {release.version} · {(release.size_bytes / 1048576).toFixed(1)} MiB</span>
                  <Button size="sm" disabled={installing || managing || alreadyManaged} onClick={() => void install()}>
                    <DownloadIcon className={installing ? "animate-pulse" : ""} />
                    {alreadyManaged ? "已安装本版本" : installing ? "下载并校验中…" : tool.managed_version ? "更新到此版本" : "下载并安装"}
                  </Button>
                </div>
                <p className="text-muted-foreground break-all font-mono" title={release.sha256}>
                  GitHub SHA-256：{release.sha256}
                </p>
              </div>
            )}
          </>
        ) : (
          <p className="text-muted-foreground text-xs">
            当前仅 Windows x64 后端提供自动安装；本机检测与手动配置仍可使用。
          </p>
        )}
        {error && <p role="alert" className="text-destructive text-xs">{error}</p>}
      </div>
      <div className="flex items-center justify-between gap-2 pt-1">
        <a
          href={tool.source}
          target="_blank"
          rel="noopener noreferrer"
          className="text-primary inline-flex items-center gap-1 text-xs hover:underline"
        >
          官方项目 <ArrowUpRightIcon className="size-3" />
        </a>
        {tool.docker_reference_version && (
          <span className="text-muted-foreground text-[11px]">
            Docker 参考版本 {tool.docker_reference_version}
          </span>
        )}
      </div>
    </div>
  );
}

export function ExternalToolsPanel({ onOpenCustom }: { onOpenCustom: () => void }) {
  const [inventory, setInventory] = React.useState<ExternalToolInventory | null>(null);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState("");

  const refresh = React.useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setInventory(await api.externalTools());
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    void refresh();
  }, [refresh]);

  const items = inventory?.tools ?? [];
  const recognized = items.filter((tool) => tool.status === "installed").length;

  return (
    <Card>
      <CardHeader className="gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div className="space-y-1.5">
          <CardTitle>外部工具检测</CardTitle>
          <CardDescription>
            检测 ARTEX 后端所在主机的已安装程序，不是当前浏览器或手机的工具。
            {inventory && (" " + inventory.host + " · 版本检查成功 " + recognized + "/" + items.length)}
          </CardDescription>
        </div>
        <Button size="sm" variant="outline" onClick={() => void refresh()} disabled={loading}>
          <RefreshCwIcon className={loading ? "animate-spin" : ""} />
          重新检测
        </Button>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && <p role="alert" className="text-destructive text-sm">检测失败：{error}</p>}
        {loading && !inventory && <p className="text-muted-foreground text-sm">正在检测本机外部工具…</p>}
        <OfflineToolsPanel host={inventory?.host ?? ""} managedCount={items.filter((tool) => tool.managed_version).length} onRefresh={refresh} />
        {items.length > 0 && (
          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
            {items.map((tool) => <ToolItem key={tool.id} tool={tool} host={inventory?.host ?? ""} onRefresh={refresh} />)}
          </div>
        )}
        <div className="text-muted-foreground space-y-2 border-t pt-4 text-xs leading-relaxed">
          <p>
            搜索顺序：ARTEX_TOOL_DIRS 自定义路径 → ARTEX 本机工具目录 →
            安装目录内的 tools/ 与 data/tools/ → 系统 PATH。Windows 可用分号分隔多个自定义路径。
          </p>
          <p>
            “版本检查通过”仅说明命令能够返回版本；手动安装的文件不代表来源已验证。
            官方按需安装会使用 GitHub Release 提供的 SHA-256 校验下载内容，但<strong>不等于发布者代码签名</strong>。
            更新只替换 ARTEX 管理目录中的版本并保留上一版；可恢复旧版或卸载受管理的程序，不会删除你自定义目录中的文件，也不会自动执行扫描。若要让 Agent 调用某个程序，
            可在现有自定义工具中明确配置，避免自动创建重复的 Agent 工具。
          </p>
          <Button size="sm" variant="outline" onClick={onOpenCustom}>
            前往自定义工具
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
