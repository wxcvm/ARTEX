"use client";

import * as React from "react";
import { Loader2Icon, RefreshCwIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { http } from "@/lib/api";
import type { LLMProfile } from "@/lib/types";
import { SensitiveConfigPanel } from "./sensitive-config";

type KeyDiagnostics = {
  key_state: "configured" | "missing" | "invalid";
  total: number;
  empty: number;
  plaintext_only: number;
  encrypted_only: number;
  mixed: number;
  verified: number;
  failed: number;
  unchecked: number;
  profiles: { id: string; storage: string; integrity: string }[];
};

const storageLabels: Record<string, string> = {
  empty: "未保存密钥",
  plaintext_only: "仅明文",
  encrypted_only: "仅密文",
  mixed: "明文与密文并存",
};
const integrityLabels: Record<string, string> = {
  not_applicable: "不适用",
  verified: "已验证",
  failed: "认证失败",
  unchecked: "未校验",
};

export function KeyDiagnosticsPanel({ profiles }: { profiles: LLMProfile[] }) {
  const [snapshot, setSnapshot] = React.useState<KeyDiagnostics | null>(null);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  const [confirmMigration, setConfirmMigration] = React.useState(false);
  const [migrationMessage, setMigrationMessage] = React.useState("");

  async function inspect() {
    setBusy(true);
    setError("");
    setSnapshot(null);
    setConfirmMigration(false);
    try {
      const result = await http<KeyDiagnostics>("/privacy/model-keys");
      if (!result || !Array.isArray(result.profiles) ||
          !["configured", "missing", "invalid"].includes(result.key_state)) {
        throw new Error("invalid diagnostics");
      }
      setSnapshot(result);
    } catch {
      setError("无法读取模型密钥诊断，请确认后端版本及数据库连接。");
    } finally {
      setBusy(false);
    }
  }

  async function migrate() {
    if (!snapshot || busy || !confirmMigration) return;
    const expectedIDs = snapshot.profiles.filter((p) => p.storage === "plaintext_only").map((p) => p.id);
    if (!expectedIDs.length) return;
    setBusy(true);
    setError("");
    setMigrationMessage("");
    setConfirmMigration(false);
    try {
      const result = await http<{ migrated: number; skipped_mixed: number }>("/privacy/model-keys/migrate", {
        method: "POST",
        body: JSON.stringify({ expected_ids: expectedIDs }),
      });
      setMigrationMessage(`已加密迁移 ${result.migrated} 个模型密钥，保留 ${result.skipped_mixed} 个明文与密文并存的配置。`);
      await inspect();
    } catch (e) {
      setSnapshot(null);
      setError(`迁移未确认成功：${(e as Error).message}。请重新检查密钥保护后再试。`);
    } finally {
      setBusy(false);
    }
  }

  const names = new Map(profiles.map((p) => [String(p.id), p.name]));
  return (
    <div className="grid min-w-0 gap-4">
    <Card className="min-w-0">
      <CardContent className="grid gap-4">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h2 className="text-base font-semibold">模型密钥存储诊断</h2>
            <p className="mt-1 text-sm text-muted-foreground">
              按需检查数据库中的存储覆盖率及当前密钥认证结果，不读取旧明文内容、不修改配置。
            </p>
          </div>
          <Button variant="outline" size="sm" disabled={busy} onClick={() => void inspect()}>
            {busy ? <Loader2Icon className="animate-spin" /> : <RefreshCwIcon />}
            {busy ? "检查中…" : "检查密钥保护"}
          </Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {migrationMessage && <p role="status" className="text-sm text-muted-foreground">{migrationMessage}</p>}
        {snapshot && (
          <div className="grid gap-4" aria-live="polite">
            <p className="text-sm">
              加密密钥：{({ configured: "已配置", missing: "未配置", invalid: "格式无效" })[snapshot.key_state]}
              {" · "}模型配置：{snapshot.total}
            </p>
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
              {([
                ["仅明文", snapshot.plaintext_only],
                ["仅密文", snapshot.encrypted_only],
                ["明文与密文并存", snapshot.mixed],
                ["未保存密钥", snapshot.empty],
              ] as const).map(([label, count]) => (
                <div key={label} className="rounded-lg border p-3">
                  <p className="text-xs text-muted-foreground">{label}</p>
                  <p className="text-xl font-semibold tabular-nums">{count}</p>
                </div>
              ))}
            </div>
            <p className="text-sm">
              密文认证：已验证 {snapshot.verified} · 失败 {snapshot.failed} · 未校验 {snapshot.unchecked}
            </p>
            {snapshot.failed > 0 && (
              <p className="text-sm text-destructive">存在无法认证的密文。请核对原始加密密钥及备份，不要用新密钥覆盖旧数据。</p>
            )}
            {snapshot.unchecked > 0 && (
              <p className="text-sm text-muted-foreground">当前密钥不可用，已有密文尚未校验，不代表密文已损坏。</p>
            )}
            {snapshot.plaintext_only > 0 && (
              <div className="grid gap-2 rounded-lg border p-3">
                <p className="text-sm font-medium">迁移旧密钥 · {snapshot.plaintext_only} 个待处理</p>
                <p className="text-xs text-muted-foreground">
                  仅加密当前预览中的明文密钥，保留模型配置及原始密钥值；已有密文和混合记录不改动。请先单独备份 ARTEX_LLM_SECRET_KEY。
                </p>
                {snapshot.key_state !== "configured" || snapshot.failed > 0 ? (
                  <p className="text-xs text-destructive">请先配置正确的加密密钥，确保已有密文认证通过。</p>
                ) : confirmMigration ? (
                  <div className="flex flex-wrap items-center gap-2">
                    <Button size="sm" disabled={busy} onClick={() => void migrate()}>
                      确认加密 {snapshot.plaintext_only} 个密钥
                    </Button>
                    <Button size="sm" variant="ghost" disabled={busy} onClick={() => setConfirmMigration(false)}>取消</Button>
                  </div>
                ) : (
                  <Button size="sm" variant="outline" className="w-fit" disabled={busy} onClick={() => setConfirmMigration(true)}>
                    迁移旧密钥
                  </Button>
                )}
              </div>
            )}
            <div className="divide-y rounded-lg border">
              {snapshot.profiles.map((p) => (
                <div key={p.id} className="flex flex-wrap items-center justify-between gap-2 px-3 py-2 text-sm">
                  <span className="min-w-0 break-all font-medium">{names.get(p.id) ?? `配置 #${p.id}`}</span>
                  <span className="text-muted-foreground">
                    {storageLabels[p.storage] ?? "未知"} · {integrityLabels[p.integrity] ?? "未知"}
                  </span>
                </div>
              ))}
              {snapshot.profiles.length === 0 && <p className="p-3 text-sm text-muted-foreground">没有模型配置。</p>}
            </div>
            <p className="text-xs leading-relaxed text-muted-foreground">
              仅密文不代表已验证；明文与密文并存仍有明文副本。空密钥可以是正常的本地模型配置。检查仅涵盖模型配置字段，不覆盖备份、WAL 或运行时内存。结果为本次检查快照，修改配置后可重新检查。
            </p>
          </div>
        )}
      </CardContent>
    </Card>
    <SensitiveConfigPanel />
    </div>
  );
}
