"use client";

import * as React from "react";
import { toast } from "sonner";
import { ArrowDownIcon, ArrowUpIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/lib/api";
import type { FindingScreenshot, FindingTrafficBinding } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export function FindingScreenshotPanel({ findingId, contextTask, readOnly }: {
  findingId: string;
  contextTask?: string;
  readOnly?: boolean;
}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [items, setItems] = React.useState<FindingScreenshot[] | null>(null);
  const [bindings, setBindings] = React.useState<FindingTrafficBinding[]>([]);
  const [error, setError] = React.useState("");
  const [busy, setBusy] = React.useState(false);
  const [refresh, setRefresh] = React.useState(0);
  const [file, setFile] = React.useState<File | null>(null);
  const [caption, setCaption] = React.useState("");
  const [binding, setBinding] = React.useState("");
  const [editing, setEditing] = React.useState<string | null>(null);
  const [editCaption, setEditCaption] = React.useState("");
  const [preview, setPreview] = React.useState<{ blob: Blob; id: string; caption: string } | null>(null);
  const [previewURL, setPreviewURL] = React.useState("");
  const fileInput = React.useRef<HTMLInputElement>(null);

  React.useEffect(() => {
    let active = true;
    setError("");
    api.findingScreenshots(findingId, contextTask)
      .then((data) => { if (active) setItems(data); })
      .catch((e) => { if (active) setError(e instanceof Error ? e.message : "Request failed"); });
    api.findingTraffic(findingId, contextTask)
      .then((data) => { if (active) setBindings(data.bindings ?? []); })
      .catch(() => { if (active) setBindings([]); });
    return () => { active = false; };
  }, [findingId, contextTask, refresh]);

  React.useEffect(() => {
    if (!preview) {
      setPreviewURL("");
      return;
    }
    const url = URL.createObjectURL(preview.blob);
    setPreviewURL(url);
    return () => URL.revokeObjectURL(url);
  }, [preview]);

  async function mutate(action: () => Promise<unknown>) {
    setBusy(true);
    try {
      await action();
      setRefresh((n) => n + 1);
      toast.success(tr("截图证据已更新", "Screenshot evidence updated"));
    } catch (e) {
      toast.error(e instanceof Error ? e.message : tr("操作失败", "Operation failed"));
    } finally {
      setBusy(false);
    }
  }

  async function upload() {
    if (!file) return;
    await mutate(async () => {
      await api.uploadFindingScreenshot(findingId, file, caption, binding, contextTask);
      setFile(null);
      setCaption("");
      setBinding("");
      if (fileInput.current) fileInput.current.value = "";
    });
  }

  async function move(index: number, delta: number) {
    if (!items || index + delta < 0 || index + delta >= items.length || busy) return;
    const ids = items.map((item) => item.id);
    [ids[index], ids[index + delta]] = [ids[index + delta], ids[index]];
    await mutate(async () => {
      setItems(await api.orderFindingScreenshots(findingId, ids, contextTask));
    });
  }

  async function view(item: FindingScreenshot) {
    setBusy(true);
    try {
      const blob = await api.readFindingScreenshot(findingId, item.id, contextTask);
      setPreview({ blob, id: item.id, caption: item.caption });
    } catch (e) {
      toast.error(e instanceof Error ? e.message : tr("加载失败", "Unable to load image"));
    } finally {
      setBusy(false);
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">{tr("截图证据", "Screenshot evidence")} {items ? `(${items.length})` : ""}</CardTitle>
        <CardDescription>{tr("上传真实截图，按需关联已保存的流量证据。留空图注时使用关联证据的说明。", "Upload visual evidence and optionally link saved traffic. An empty caption uses the linked evidence note.")}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-3">
        {error && <div role="alert" className="text-sm text-destructive">{error} <Button size="sm" variant="link" onClick={() => setRefresh((n) => n + 1)}>{tr("重试", "Retry")}</Button></div>}
        {!readOnly && (
          <div className="flex flex-col gap-2 rounded-md border p-3">
            <label className="text-xs text-muted-foreground" htmlFor="finding-screenshot-upload">{tr("选择 PNG / JPEG 截图（8 MiB 以内）", "PNG / JPEG screenshot (up to 8 MiB)")}</label>
            <input id="finding-screenshot-upload" ref={fileInput} type="file" accept="image/png,image/jpeg" disabled={busy} onChange={(e) => setFile(e.target.files?.[0] ?? null)} className="text-sm" />
            <select aria-label={tr("关联流量证据", "Linked traffic evidence")} value={binding} onChange={(e) => setBinding(e.target.value)} disabled={busy} className="min-w-0 rounded-md border bg-background p-2 text-sm">
              <option value="">{tr("不关联流量", "No traffic binding")}</option>
              {bindings.map((b) => <option key={b.id} value={b.id}>#{b.id} · {b.snapshot.method} {b.snapshot.url}</option>)}
            </select>
            <Textarea value={caption} onChange={(e) => setCaption(e.target.value)} placeholder={tr("图注（可选）", "Caption (optional)")} maxLength={4000} rows={2} disabled={busy} />
            <Button size="sm" disabled={!file || busy} onClick={() => void upload()}>{tr("上传截图", "Upload screenshot")}</Button>
          </div>
        )}
        {items === null && !error ? <p className="text-sm text-muted-foreground">{tr("正在加载…", "Loading…")}</p> : null}
        {items?.length === 0 && <p className="text-sm text-muted-foreground">{tr("暂无截图证据", "No screenshots")}</p>}
        {items?.map((item, index) => (
          <div key={item.id} className="flex flex-col gap-2 rounded-md border p-3">
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-sm font-medium">{tr("截图", "Screenshot")} {index + 1}</span>
              {item.binding_id && <span className="text-xs text-muted-foreground">{tr("关联证据", "Evidence")} #{item.binding_id}</span>}
              <span className="flex-1" />
              {!readOnly && <>
                <Button size="sm" variant="ghost" aria-label={tr("上移截图", "Move screenshot up")} disabled={busy || index === 0} onClick={() => void move(index, -1)}><ArrowUpIcon /></Button>
                <Button size="sm" variant="ghost" aria-label={tr("下移截图", "Move screenshot down")} disabled={busy || index === items.length - 1} onClick={() => void move(index, 1)}><ArrowDownIcon /></Button>
              </>}
              <Button size="sm" variant="outline" disabled={busy} onClick={() => void view(item)}>{tr("查看", "View")}</Button>
              {!readOnly && <Button size="sm" variant="outline" disabled={busy} onClick={() => { setEditing(item.id); setEditCaption(item.caption); }}>{tr("编辑图注", "Edit caption")}</Button>}
              {!readOnly && <Button size="sm" variant="ghost" disabled={busy} onClick={() => void mutate(async () => {
                await api.deleteFindingScreenshot(findingId, item.id, contextTask);
                if (preview?.id === item.id) setPreview(null);
              })}>{tr("删除", "Delete")}</Button>}
            </div>
            {editing === item.id ? (
              <>
                <Textarea value={editCaption} onChange={(e) => setEditCaption(e.target.value)} rows={2} maxLength={4000} disabled={busy} />
                <div className="flex gap-2">
                  <Button size="sm" disabled={busy} onClick={() => void mutate(async () => {
                    await api.editFindingScreenshot(findingId, item.id, editCaption, contextTask);
                    setEditing(null);
                    if (preview?.id === item.id) setPreview({ ...preview, caption: editCaption });
                  })}>{tr("保存", "Save")}</Button>
                  <Button size="sm" variant="outline" disabled={busy} onClick={() => setEditing(null)}>{tr("取消", "Cancel")}</Button>
                </div>
              </>
            ) : <p className="whitespace-pre-wrap break-words text-sm">{item.caption || tr("（无图注）", "(No caption)")}</p>}
          </div>
        ))}
        {preview && previewURL && <figure className="flex flex-col gap-2 rounded-md border p-3">
          <img src={previewURL} alt={preview.caption || tr("截图证据", "Screenshot evidence")} className="max-h-[60vh] w-auto max-w-full object-contain" />
          <figcaption className="break-words text-sm">{preview.caption}</figcaption>
          <Button size="sm" variant="outline" onClick={() => setPreview(null)}>{tr("关闭预览", "Close preview")}</Button>
        </figure>}
      </CardContent>
    </Card>
  );
}
