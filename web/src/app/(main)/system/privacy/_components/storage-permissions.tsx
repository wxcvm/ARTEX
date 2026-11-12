"use client";

import * as React from "react";
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { http } from "@/lib/api";

type Item = { kind: string; state: string; broad_read: boolean; other_read: boolean; inheritance_enabled: boolean };
type Snapshot = { platform: string; items: Item[] };
const names: Record<string, string> = { config_file: "数据库配置文件", jwt_key: "JWT 签名密钥文件", data_root: "应用数据根目录" };
const states: Record<string, string> = { inspected: "已检查", partial: "部分检查", missing: "不存在", unavailable: "无法读取", skipped_reparse: "跳过重解析点", skipped_nonlocal: "跳过非本地驱动器", unsupported: "当前宿主不支持 Windows ACL" };

export function StoragePermissionsPanel() {
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
      setError("无法读取存储权限诊断，请确认后端版本。");
    } finally {
      setBusy(false);
    }
  }
  return (
    <Card className="min-w-0">
      <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-3">
        <div className="grid gap-1">
          <CardTitle className="text-base">Windows 存储权限诊断</CardTitle>
          <CardDescription>读取实际配置文件、签名密钥及数据根目录的访问控制规则，不修改 ACL。</CardDescription>
        </div>
        <Button type="button" variant="outline" size="sm" disabled={busy} onClick={() => void inspect()}>
          <RefreshCw className={busy ? "size-4 animate-spin" : "size-4"} />
          {busy ? "检查中…" : "检查存储权限"}
        </Button>
      </CardHeader>
      <CardContent className="grid gap-3">
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {snapshot && (
          <div className="divide-y rounded-lg border" aria-live="polite">
            {snapshot.items.map((item) => (
              <div key={item.kind} className="grid gap-1 p-3 text-sm">
                <div className="flex flex-wrap justify-between gap-2">
                  <span className="font-medium">{names[item.kind] ?? "存储项目"}</span>
                  <span className="text-muted-foreground">{states[item.state] ?? "未知"}</span>
                </div>
                {(item.state === "inspected" || item.state === "partial") && (
                  <>
                    <p className={item.broad_read || item.other_read ? "text-destructive" : "text-muted-foreground"}>
                      {item.broad_read ? "发现宽泛读取授权" : "未发现已识别的宽泛读取授权"} · {item.other_read ? "发现其他主体读取授权" : "未发现其他已识别的读取授权"}
                    </p>
                    <p className="text-muted-foreground">权限继承：{item.inheritance_enabled ? "启用" : "禁用"}</p>
                  </>
                )}
              </div>
            ))}
          </div>
        )}
        <p className="text-xs leading-relaxed text-muted-foreground">只检查三处固定位置；数据目录仅检查根目录，不递归验证子文件。报告的是普通允许规则，不计算拒绝规则、组成员及复杂 ACE 的最终有效访问权限。非 Windows 宿主、非本地驱动器或重解析点会明确标记；不返回文件内容、路径或账户标识。</p>
      </CardContent>
    </Card>
  );
}
