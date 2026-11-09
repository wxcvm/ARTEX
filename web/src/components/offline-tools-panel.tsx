"use client";

import * as React from "react";
import { DownloadIcon, PackageIcon, UploadIcon } from "lucide-react";
import { toast } from "sonner";

import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";

export function OfflineToolsPanel({
  host, managedCount, onRefresh,
}: {
  host: string;
  managedCount: number;
  onRefresh: () => Promise<void>;
}) {
  const fileInput = React.useRef<HTMLInputElement>(null);
  const [exporting, setExporting] = React.useState(false);
  const [importing, setImporting] = React.useState(false);
  const [error, setError] = React.useState("");
  const supported = host === "windows/amd64";

  async function exportBundle() {
    setExporting(true);
    setError("");
    try {
      const blob = await api.exportOfflineTools();
      const url = URL.createObjectURL(blob);
      const element = document.createElement("a");
      element.href = url;
      element.download = "ARTEX-Tools-Offline-windows-amd64.zip";
      document.body.appendChild(element);
      element.click();
      element.remove();
      window.setTimeout(() => URL.revokeObjectURL(url), 60_000);
      toast.success("已生成 ARTEX 离线工具包，可拷贝到另一台电脑");
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setExporting(false);
    }
  }

  async function importBundle(file: File | null) {
    if (!file) return;
    if (!window.confirm(
      "导入 " + file.name + "（" + (file.size / 1048576).toFixed(1) + " MiB）？\n" +
      "ARTEX 会校验包内每个程序的 SHA-256，然后更新受管理的工具；" +
      "原有版本保留用于回滚，自定义工具不会被覆盖。\n" +
      "离线包仅有内容完整性校验，没有独立的发布者签名，请使用可信来源的工具包。"
    )) return;

    setImporting(true);
    setError("");
    try {
      const result = await api.importOfflineTools(file);
      toast.success("离线导入完成：" + result.installed.length + " 个工具");
    } catch (e) {
      setError((e as Error).message + "。如有部分工具已更新，请刷新查看。");
    } finally {
      setImporting(false);
      await onRefresh();
    }
  }

  return (
    <div className="space-y-3 rounded-lg border p-4">
      <div className="flex items-center gap-2 font-medium">
        <PackageIcon className="size-4" />
        离线工具包
      </div>
      <p className="text-muted-foreground text-xs leading-relaxed">
        在联网的 Windows x64 主机上导出 ARTEX 已管理的工具，拷贝 ZIP 到另一台
        Windows x64 主机后导入。无需离线机器访问 GitHub，也不会更改系统 PATH 或 Agent 绑定。
      </p>
      <input
        ref={fileInput}
        className="hidden"
        type="file"
        accept=".zip,application/zip"
        aria-label="选择 ARTEX 离线工具 ZIP"
        onChange={(event) => {
          const file = event.currentTarget.files?.[0] ?? null;
          event.currentTarget.value = "";
          void importBundle(file);
        }}
      />
      <div className="flex flex-wrap gap-2">
        <Button size="sm" variant="outline" disabled={!supported || managedCount === 0 || exporting || importing}
          onClick={() => void exportBundle()}>
          <DownloadIcon className={exporting ? "animate-pulse" : ""} />
          {exporting ? "正在生成…" : "导出已管理工具 ZIP"}
        </Button>
        <Button size="sm" variant="outline" disabled={!supported || importing || exporting}
          onClick={() => fileInput.current?.click()}>
          <UploadIcon className={importing ? "animate-pulse" : ""} />
          {importing ? "正在验证并导入…" : "导入离线工具 ZIP"}
        </Button>
      </div>
      {supported && <p className="text-muted-foreground text-xs">
        当前可导出 {managedCount} 个管理版本；仅导出 ARTEX 管理的当前版，不包含 PATH、自定义工具或历史备份。
      </p>}
      {!supported && <p className="text-muted-foreground text-xs">
        离线导入和导出目前支持 Windows x64 后端；其他系统仍可使用现有工具检测。
      </p>}
      {error && <p role="alert" className="text-destructive text-sm">{error}</p>}
    </div>
  );
}
