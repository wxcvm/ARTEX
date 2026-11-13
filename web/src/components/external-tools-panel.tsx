"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

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
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
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
      tr(`从官方 GitHub 下载 ${tool.name} ${release.version}（${sizeMiB} MiB）？\n下载后校验 SHA-256，再安装到 ARTEX 用户工具目录；不会覆盖自定义目录、修改 PATH 或绑定 Agent。`, `Download ${tool.name} ${release.version} (${sizeMiB} MiB) from official GitHub?\nVerify SHA-256, then install in the ARTEX user tools directory; custom directories, PATH and agent bindings are unchanged.`)
    )) return;
    setInstalling(true);
    setError("");
    try {
      await api.installExternalTool(tool.id, release.version, release.sha256);
      toast.success(tr(`${tool.name} ${release.version} 已安装并通过下载摘要校验`, `${tool.name} ${release.version} installed and download digest verified`));
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
      tr(`将 ${tool.name} 从 ${tool.managed_version} 恢复到 ${tool.previous_version}？当前版本会保留为上一版；正在运行此程序时请先结束运行。`, `Restore ${tool.name} from ${tool.managed_version} to ${tool.previous_version}? The current version becomes the previous version; stop the running program first.`)
    )) return;
    setManaging(true);
    setError("");
    try {
      const result = await api.rollbackExternalTool(tool.id);
      toast.success(tr(`${tool.name} 已恢复到 ${result.version}`, `${tool.name} restored to ${result.version}`));
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
      tr(`卸载 ARTEX 管理的 ${tool.name} ${tool.managed_version} 和保留的上一版？\n自定义目录及 PATH 中的工具不会被删除。`, `Uninstall ARTEX-managed ${tool.name} ${tool.managed_version} and its retained previous version?\nTools in custom directories and PATH are not deleted.`)
    )) return;
    setManaging(true);
    setError("");
    try {
      await api.uninstallExternalTool(tool.id);
      toast.success(tr(`已卸载 ARTEX 管理的 ${tool.name}`, `Uninstalled ARTEX-managed ${tool.name}`));
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
          <Badge variant="secondary">{tr("版本检查通过", "Version check passed")}</Badge>
        ) : tool.status === "unverified" ? (
          <Badge variant="outline">{tr("已发现 · 检查失败", "Found · check failed")}</Badge>
        ) : (
          <Badge variant="outline">{tr("未发现", "Not found")}</Badge>
        )}
        {tool.origin && (
          <span className="text-muted-foreground ml-auto text-xs">
            {locale === "en" ? ({ custom: "Custom directory", managed: "ARTEX tool directory", bundled: "Bundled", path: "System PATH" } as Record<string, string>)[tool.origin] ?? tool.origin : originName[tool.origin] ?? tool.origin}
          </span>
        )}
      </div>
      <p className="text-muted-foreground text-xs">{tool.description}</p>
      {tool.version && <div className="break-all font-mono text-xs">{tool.version}</div>}
      {tool.managed_version && (
        <div className="space-y-2">
          <div className="text-xs">
            {tr("ARTEX 管理版本：", "ARTEX managed version: ")}<span className="font-mono">{tool.managed_version}</span>
            {tool.previous_version && (
              <span className="text-muted-foreground ml-2">
                {tr("上一版", "Previous version")} <span className="font-mono">{tool.previous_version}</span>
              </span>
            )}
          </div>
          {supported && (
            <div className="flex flex-wrap gap-2">
              {tool.previous_version && (
                <Button variant="outline" size="sm" disabled={managing || installing}
                  onClick={() => void rollback()}>{tr("恢复上一版", "Restore previous version")}</Button>
              )}
              <Button variant="outline" size="sm" disabled={managing || installing}
                onClick={() => void uninstall()}>
                {managing ? tr("处理中…", "Working…") : tr("卸载管理版本", "Uninstall managed version")}
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
        <p className="text-muted-foreground text-xs">{tr("自定义目录优先于 ARTEX 管理目录，界面显示的是实际优先匹配的版本。", "Custom directories take precedence over ARTEX-managed directories; the display shows the version actually selected.")}</p>
      )}
      <div className="mt-auto space-y-2 pt-2">
        {supported ? (
          <>
            <Button size="sm" variant="outline" disabled={checking || installing || managing} onClick={() => void checkRelease()}>
              <RefreshCwIcon className={checking ? "animate-spin" : ""} />
              {checking ? tr("查询官方发布…", "Checking official releases…") : tr("查看官方最新版", "Check latest official release")}
            </Button>
            {release && (
              <div className="space-y-2 rounded-md bg-muted/50 p-2 text-xs">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span>{tr("官方", "Official")} {release.version} · {(release.size_bytes / 1048576).toFixed(1)} MiB</span>
                  <Button size="sm" disabled={installing || managing || alreadyManaged} onClick={() => void install()}>
                    <DownloadIcon className={installing ? "animate-pulse" : ""} />
                    {alreadyManaged ? tr("已安装本版本", "Version already installed") : installing ? tr("下载并校验中…", "Downloading and verifying…") : tool.managed_version ? tr("更新到此版本", "Update to this version") : tr("下载并安装", "Download and install")}
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
            {tr("当前仅 Windows x64 后端提供自动安装；本机检测与手动配置仍可使用。", "Automatic installation is currently available only on Windows x64 backends; local detection and manual setup remain available.")}
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
          {tr("官方项目", "Official project")} <ArrowUpRightIcon className="size-3" />
        </a>
        {tool.docker_reference_version && (
          <span className="text-muted-foreground text-[11px]">
            {tr("Docker 参考版本", "Docker reference version")} {tool.docker_reference_version}
          </span>
        )}
      </div>
    </div>
  );
}

export function ExternalToolsPanel({ onOpenCustom }: { onOpenCustom: () => void }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
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
            {tr("检测 ARTEX 后端所在主机的已安装程序，不是当前浏览器或手机的工具。", "Detect programs on the ARTEX backend host, not tools on the current browser or phone.")}
            {inventory && (" " + inventory.host + tr(" · 版本检查成功 ", " · Version checks passed ") + recognized + "/" + items.length)}
          </CardDescription>
        </div>
        <Button size="sm" variant="outline" onClick={() => void refresh()} disabled={loading}>
          <RefreshCwIcon className={loading ? "animate-spin" : ""} />
          {tr("重新检测", "Detect again")}
        </Button>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && <p role="alert" className="text-destructive text-sm">{tr("检测失败：", "Detection failed: ")}{error}</p>}
        {loading && !inventory && <p className="text-muted-foreground text-sm">{tr("正在检测本机外部工具…", "Detecting local external tools…")}</p>}
        <OfflineToolsPanel host={inventory?.host ?? ""} managedCount={items.filter((tool) => tool.managed_version).length} onRefresh={refresh} />
        {items.length > 0 && (
          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
            {items.map((tool) => <ToolItem key={tool.id} tool={tool} host={inventory?.host ?? ""} onRefresh={refresh} />)}
          </div>
        )}
        <div className="text-muted-foreground space-y-2 border-t pt-4 text-xs leading-relaxed">
          <p>
            {tr("搜索顺序：ARTEX_TOOL_DIRS 自定义路径 → ARTEX 本机工具目录 → 安装目录内的 tools/ 与 data/tools/ → 系统 PATH。Windows 可用分号分隔多个自定义路径。", "Search order: ARTEX_TOOL_DIRS custom paths → local ARTEX tools → tools/ and data/tools/ in the installation → system PATH. Separate multiple Windows custom paths with semicolons.")}
          </p>
          <p>
            {tr("“版本检查通过”仅说明命令能够返回版本；手动安装的文件不代表来源已验证。官方按需安装使用 GitHub Release 的 SHA-256 校验下载内容，但不等于发布者代码签名。更新仅替换 ARTEX 管理版本并保留上一版，可恢复或卸载；不会删除自定义目录文件或自动执行扫描。可在自定义工具中明确配置 Agent 调用，避免重复创建。", "Version check passed only means the command returned a version; manually installed files are not source-verified. Official on-demand installs verify the GitHub Release SHA-256, not publisher signatures. Updates replace only ARTEX-managed versions and retain the previous version for restore or uninstall. Custom files are not deleted and scans are not run automatically. Configure agent invocation explicitly in custom tools to avoid duplicate entries.")}
          </p>
          <Button size="sm" variant="outline" onClick={onOpenCustom}>
            {tr("前往自定义工具", "Open custom tools")}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
