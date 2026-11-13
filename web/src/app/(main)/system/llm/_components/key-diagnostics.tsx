"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

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
const storageLabelsEn: Record<string, string> = { empty: "No key stored", plaintext_only: "Plaintext only", encrypted_only: "Ciphertext only", mixed: "Mixed plaintext and ciphertext" };
const integrityLabelsEn: Record<string, string> = { not_applicable: "Not applicable", verified: "Verified", failed: "Authentication failed", unchecked: "Unchecked" };
const integrityLabels: Record<string, string> = {
  not_applicable: "不适用",
  verified: "已验证",
  failed: "认证失败",
  unchecked: "未校验",
};

export function KeyDiagnosticsPanel({ profiles }: { profiles: LLMProfile[] }) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
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
      setError(tr("无法读取模型密钥诊断，请确认后端版本及数据库连接。", "Unable to read key diagnostics. Check the backend version and database connection."));
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
      setMigrationMessage(tr(`已加密迁移 ${result.migrated} 个模型密钥，保留 ${result.skipped_mixed} 个明文与密文并存的配置。`, `Encrypted ${result.migrated} keys; retained ${result.skipped_mixed} mixed profiles.`));
      await inspect();
    } catch (e) {
      setSnapshot(null);
      setError(tr(`迁移未确认成功：${(e as Error).message}。请重新检查密钥保护后再试。`, `Migration could not be confirmed: ${(e as Error).message}. Inspect key protection again before retrying.`));
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
            <h2 className="text-base font-semibold">{tr("模型密钥存储诊断", "Model key storage diagnostics")}</h2>
            <p className="mt-1 text-sm text-muted-foreground">
              {tr("按需检查数据库中的存储覆盖率及当前密钥认证结果，不读取旧明文内容、不修改配置。", "Check database storage coverage and key authentication on demand, without reading legacy plaintext or changing configuration.")}
            </p>
          </div>
          <Button variant="outline" size="sm" disabled={busy} onClick={() => void inspect()}>
            {busy ? <Loader2Icon className="animate-spin" /> : <RefreshCwIcon />}
            {busy ? tr("检查中…", "Checking…") : tr("检查密钥保护", "Inspect key protection")}
          </Button>
        </div>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        {migrationMessage && <p role="status" className="text-sm text-muted-foreground">{migrationMessage}</p>}
        {snapshot && (
          <div className="grid gap-4" aria-live="polite">
            <p className="text-sm">
              {tr("加密密钥：", "Encryption key: ")}{(locale === "en" ? { configured: "Configured", missing: "Missing", invalid: "Invalid" } : { configured: "已配置", missing: "未配置", invalid: "格式无效" })[snapshot.key_state]}
              {" · "}{tr("模型配置：", "Profiles: ")}{snapshot.total}
            </p>
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
              {([
                [tr("仅明文", "Plaintext only"), snapshot.plaintext_only],
                [tr("仅密文", "Ciphertext only"), snapshot.encrypted_only],
                [tr("明文与密文并存", "Mixed"), snapshot.mixed],
                [tr("未保存密钥", "No key stored"), snapshot.empty],
              ] as const).map(([label, count]) => (
                <div key={label} className="rounded-lg border p-3">
                  <p className="text-xs text-muted-foreground">{label}</p>
                  <p className="text-xl font-semibold tabular-nums">{count}</p>
                </div>
              ))}
            </div>
            <p className="text-sm">
              {tr(`密文认证：已验证 ${snapshot.verified} · 失败 ${snapshot.failed} · 未校验 ${snapshot.unchecked}`, `Authentication: ${snapshot.verified} verified · ${snapshot.failed} failed · ${snapshot.unchecked} unchecked`)}
            </p>
            {snapshot.failed > 0 && (
              <p className="text-sm text-destructive">{tr("存在无法认证的密文。请核对原始加密密钥及备份，不要用新密钥覆盖旧数据。", "Some ciphertext cannot be authenticated. Check the original encryption key and backups; do not overwrite old data with a new key.")}</p>
            )}
            {snapshot.unchecked > 0 && (
              <p className="text-sm text-muted-foreground">{tr("当前密钥不可用，已有密文尚未校验，不代表密文已损坏。", "The current key is unavailable; existing ciphertext has not been checked, which does not imply corruption.")}</p>
            )}
            {snapshot.plaintext_only > 0 && (
              <div className="grid gap-2 rounded-lg border p-3">
                <p className="text-sm font-medium">{tr(`迁移旧密钥 · ${snapshot.plaintext_only} 个待处理`, `Migrate legacy keys · ${snapshot.plaintext_only} pending`)}</p>
                <p className="text-xs text-muted-foreground">
                  {tr("仅加密当前预览中的明文密钥，保留模型配置及原始密钥值；已有密文和混合记录不改动。请先单独备份 ARTEX_LLM_SECRET_KEY。", "Encrypt only the plaintext keys in this preview, preserving profiles and original key values. Existing ciphertext and mixed records remain unchanged. Back up ARTEX_LLM_SECRET_KEY separately first.")}
                </p>
                {snapshot.key_state !== "configured" || snapshot.failed > 0 ? (
                  <p className="text-xs text-destructive">{tr("请先配置正确的加密密钥，确保已有密文认证通过。", "Configure the correct encryption key and verify existing ciphertext first.")}</p>
                ) : confirmMigration ? (
                  <div className="flex flex-wrap items-center gap-2">
                    <Button size="sm" disabled={busy} onClick={() => void migrate()}>
                      {tr(`确认加密 ${snapshot.plaintext_only} 个密钥`, `Confirm encryption of ${snapshot.plaintext_only} keys`)}
                    </Button>
                    <Button size="sm" variant="ghost" disabled={busy} onClick={() => setConfirmMigration(false)}>{tr("取消", "Cancel")}</Button>
                  </div>
                ) : (
                  <Button size="sm" variant="outline" className="w-fit" disabled={busy} onClick={() => setConfirmMigration(true)}>
                    {tr("迁移旧密钥", "Migrate legacy keys")}
                  </Button>
                )}
              </div>
            )}
            <div className="divide-y rounded-lg border">
              {snapshot.profiles.map((p) => (
                <div key={p.id} className="flex flex-wrap items-center justify-between gap-2 px-3 py-2 text-sm">
                  <span className="min-w-0 break-all font-medium">{names.get(p.id) ?? tr(`配置 #${p.id}`, `Profile #${p.id}`)}</span>
                  <span className="text-muted-foreground">
                    {(locale === "en" ? storageLabelsEn : storageLabels)[p.storage] ?? tr("未知", "Unknown")} · {(locale === "en" ? integrityLabelsEn : integrityLabels)[p.integrity] ?? tr("未知", "Unknown")}
                  </span>
                </div>
              ))}
              {snapshot.profiles.length === 0 && <p className="p-3 text-sm text-muted-foreground">{tr("没有模型配置。", "No model profiles.")}</p>}
            </div>
            <p className="text-xs leading-relaxed text-muted-foreground">
              {tr("仅密文不代表已验证；明文与密文并存仍有明文副本。空密钥可以是正常的本地模型配置。检查仅涵盖模型配置字段，不覆盖备份、WAL 或运行时内存。结果为本次检查快照，修改配置后可重新检查。", "Ciphertext-only does not mean verified; mixed records still have plaintext copies. Empty keys can be valid for local models. Inspection covers profile fields only, not backups, WAL or runtime memory. Results are a snapshot; inspect again after changes.")}
            </p>
          </div>
        )}
      </CardContent>
    </Card>
    <SensitiveConfigPanel />
    </div>
  );
}
