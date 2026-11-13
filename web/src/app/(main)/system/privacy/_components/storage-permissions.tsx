"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { http } from "@/lib/api";

type Item = { kind: string; state: string; broad_read: boolean; other_read: boolean; inheritance_enabled: boolean };
type Snapshot = { platform: string; items: Item[] };
const names: Record<string, string> = { config_file: "数据库配置文件", jwt_key: "JWT 签名密钥文件", data_root: "应用数据根目录" };
const namesEn: Record<string, string> = { config_file: "Database configuration file", jwt_key: "JWT signing key file", data_root: "Application data root" };
const statesEn: Record<string, string> = { inspected: "Inspected", partial: "Partial", missing: "Missing", unavailable: "Unavailable", skipped_reparse: "Skipped reparse point", skipped_nonlocal: "Skipped non-local drive", unsupported: "Windows ACL unsupported" };
const states: Record<string, string> = { inspected: "已检查", partial: "部分检查", missing: "不存在", unavailable: "无法读取", skipped_reparse: "跳过重解析点", skipped_nonlocal: "跳过非本地驱动器", unsupported: "当前宿主不支持 Windows ACL" };

export function StoragePermissionsPanel() {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [snapshot, setSnapshot] = React.useState<Snapshot | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  async function inspect() {
    setBusy(true); setError(""); setSnapshot(null);
    try {
      const result = await http<Snapshot>("/privacy/storage-permissions");
      if (!result || !Array.isArray(result.items)) throw new Error("invalid snapshot");
      setSnapshot(result);
    } catch {
      setError(tr("无法读取存储权限诊断，请确认后端版本。", "Unable to read storage permission diagnostics; check backend version."));
    } finally {
      setBusy(false);
    }
  }
  return (
    <Card className="min-w-0">
      <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-3">
        <div className="grid gap-1">
          <CardTitle className="text-base">{tr("Windows 存储权限诊断", "Windows storage permission diagnostics")}</CardTitle>
          <CardDescription>{tr("读取实际配置文件、签名密钥及数据根目录的访问控制规则，不修改 ACL。", "Read access-control rules for the actual configuration file, signing key and data root; do not change ACLs.")}</CardDescription>
        </div>
        <Button type="button" variant="outline" size="sm" disabled={busy} onClick={() => void inspect()}>
          <RefreshCw className={busy ? "size-4 animate-spin" : "size-4"} />
          {busy ? tr("检查中…", "Checking…") : tr("检查存储权限", "Inspect storage permissions")}
        </Button>
      </CardHeader>
      <CardContent className="grid gap-3">
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {snapshot && (
          <div className="divide-y rounded-lg border" aria-live="polite">
            {snapshot.items.map((item) => (
              <div key={item.kind} className="grid gap-1 p-3 text-sm">
                <div className="flex flex-wrap justify-between gap-2">
                  <span className="font-medium">{(locale === "en" ? namesEn : names)[item.kind] ?? tr("存储项目", "Storage item")}</span>
                  <span className="text-muted-foreground">{(locale === "en" ? statesEn : states)[item.state] ?? tr("未知", "Unknown")}</span>
                </div>
                {(item.state === "inspected" || item.state === "partial") && (
                  <>
                    <p className={item.broad_read || item.other_read ? "text-destructive" : "text-muted-foreground"}>
                      {item.broad_read ? tr("发现宽泛读取授权", "Broad read grant found") : tr("未发现已识别的宽泛读取授权", "No recognized broad read grant")} · {item.other_read ? tr("发现其他主体读取授权", "Other principal read grant found") : tr("未发现其他已识别的读取授权", "No other recognized read grant")}
                    </p>
                    <p className="text-muted-foreground">{tr("权限继承：", "Inheritance: ")}{item.inheritance_enabled ? tr("启用", "Enabled") : tr("禁用", "Disabled")}</p>
                  </>
                )}
              </div>
            ))}
          </div>
        )}
        <p className="text-xs leading-relaxed text-muted-foreground">{tr("只检查三处固定位置；数据目录仅检查根目录，不递归验证子文件。报告的是普通允许规则，不计算拒绝规则、组成员及复杂 ACE 的最终有效访问权限。非 Windows 宿主、非本地驱动器或重解析点会明确标记；不返回文件内容、路径或账户标识。", "Inspect three fixed locations; only the data root is checked, not child files. Reports ordinary allow rules, not effective access from deny rules, group membership or complex ACEs. Non-Windows hosts, non-local drives and reparse points are marked. No file contents, paths or account identifiers are returned.")}</p>
      </CardContent>
    </Card>
  );
}
