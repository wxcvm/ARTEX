"use client";

import * as React from "react";

import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { api } from "@/lib/api";

export function ChangePasswordDialog({
  open,
  onOpenChange,
}: {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
}) {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  const [oldPassword, setOldPassword] = React.useState("");
  const [newPassword, setNewPassword] = React.useState("");
  const [confirmPassword, setConfirmPassword] = React.useState("");
  const [saving, setSaving] = React.useState(false);

  // reset fields whenever the dialog closes
  React.useEffect(() => {
    if (!open) {
      setOldPassword("");
      setNewPassword("");
      setConfirmPassword("");
    }
  }, [open]);

  function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!oldPassword || !newPassword) {
      toast.error(t("请填写当前密码和新密码", "Fill in the current password and the new password"));
      return;
    }
    if (newPassword !== confirmPassword) {
      toast.error(t("两次输入的新密码不一致", "The two new passwords do not match"));
      return;
    }
    setSaving(true);
    api
      .changePassword(oldPassword, newPassword)
      .then(() => {
        toast.success(t("密码已修改", "Password changed"));
        onOpenChange(false);
      })
      .catch((err) => toast.error(`${t("修改失败：", "Change failed: ")}${(err as Error).message}`))
      .finally(() => setSaving(false));
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <form onSubmit={submit}>
          <DialogHeader>
            <DialogTitle>{t("修改密码", "Change password")}</DialogTitle>
            <DialogDescription>
              {t("用户名固定为 ARTEX。请输入当前密码进行验证；修改前签发的登录令牌在过期前仍然有效。", "The username is fixed as ARTEX. Enter the current password to verify; login tokens issued before the change stay valid until they expire.")}
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-3 py-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="cp-old">{t("当前密码", "Current password")}</Label>
              <Input
                id="cp-old"
                type="password"
                autoComplete="current-password"
                value={oldPassword}
                disabled={saving}
                onChange={(e) => setOldPassword(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="cp-new">{t("新密码", "New password")}</Label>
              <Input
                id="cp-new"
                type="password"
                autoComplete="new-password"
                value={newPassword}
                disabled={saving}
                onChange={(e) => setNewPassword(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="cp-confirm">{t("确认新密码", "Confirm new password")}</Label>
              <Input
                id="cp-confirm"
                type="password"
                autoComplete="new-password"
                value={confirmPassword}
                disabled={saving}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={saving}>
              {t("取消", "Cancel")}
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? t("正在修改…", "Changing…") : t("确认修改", "Confirm change")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
