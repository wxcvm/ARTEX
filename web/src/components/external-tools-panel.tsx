"use client";

import * as React from "react";
import { ArrowUpRightIcon, RefreshCwIcon } from "lucide-react";

import { api } from "@/lib/api";
import type { ExternalToolInventory, ExternalToolStatus } from "@/lib/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

const originName: Record<string, string> = {
  custom: "自定义目录",
  managed: "ARTEX 工具目录",
  bundled: "程序附带",
  path: "系统 PATH",
};

function ToolItem({ tool }: { tool: ExternalToolStatus }) {
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
      {tool.path && (
        <div className="text-muted-foreground break-all font-mono text-xs" title={tool.path}>
          {tool.path}
        </div>
      )}
      {tool.detail && <p className="text-muted-foreground text-xs">{tool.detail}</p>}
      <div className="mt-auto flex items-center justify-between gap-2 pt-2">
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
        {items.length > 0 && (
          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
            {items.map((tool) => <ToolItem key={tool.id} tool={tool} />)}
          </div>
        )}
        <div className="text-muted-foreground space-y-2 border-t pt-4 text-xs leading-relaxed">
          <p>
            搜索顺序：ARTEX_TOOL_DIRS 自定义路径 → ARTEX 本机工具目录 →
            安装目录内的 tools/ 与 data/tools/ → 系统 PATH。Windows 可用分号分隔多个自定义路径。
          </p>
          <p>
            “版本检查通过”仅说明命令能够返回版本，<strong>不代表文件来源或 SHA-256 已被验证</strong>；
            这不是安装、更新或执行扫描的入口。若要让 Agent 调用某个程序，
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
