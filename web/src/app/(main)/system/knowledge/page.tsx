"use client";

import { useEffect, useState } from "react";
import { BookOpen, RefreshCw, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { http } from "@/lib/api";
import {
  attachmentLabels, sourceClaimLabel,
  type KnowledgeStatus, type KnowledgeList, type KnowledgeEntry, type KnowledgeSource,
} from "@/lib/knowledge-preview";

export default function KnowledgePage() {
  const [status, setStatus] = useState<KnowledgeStatus | null>(null);
  const [list, setList] = useState<KnowledgeList | null>(null);
  const [detail, setDetail] = useState<{entry: KnowledgeEntry; source: KnowledgeSource} | null>(null);
  const [draft, setDraft] = useState("");
  const [query, setQuery] = useState("");
  const [conflicts, setConflicts] = useState(false);
  const [offset, setOffset] = useState(0);
  const [selected, setSelected] = useState("");
  const [refresh, setRefresh] = useState(0);
  const [busy, setBusy] = useState(true);
  const [error, setError] = useState("");
  const [detailError, setDetailError] = useState("");

  useEffect(() => {
    const controller = new AbortController();
    setBusy(true); setError(""); setList(null); setSelected(""); setDetail(null);
    async function load() {
      try {
        const current = await http<KnowledgeStatus>("/knowledge/status", { signal: controller.signal });
        if (controller.signal.aborted) return;
        setStatus(current);
        if (current.enabled) {
          const params = new URLSearchParams({ q: query, conflicts: String(conflicts), offset: String(offset), limit: "20" });
          const result = await http<KnowledgeList>(`/knowledge?${params}`, { signal: controller.signal });
          if (!controller.signal.aborted) setList(result);
        }
      } catch (e) {
        if (!controller.signal.aborted) { setError(e instanceof Error ? e.message : "知识库加载失败"); setStatus(null); }
      } finally {
        if (!controller.signal.aborted) setBusy(false);
      }
    }
    void load();
    return () => controller.abort();
  }, [query, conflicts, offset, refresh]);

  useEffect(() => {
    setDetail(null); setDetailError("");
    if (!selected) return;
    const controller = new AbortController();
    http<{entry: KnowledgeEntry; source: KnowledgeSource}>(`/knowledge/entries/${encodeURIComponent(selected)}`, { signal: controller.signal })
      .then((data) => { if (!controller.signal.aborted) setDetail(data); })
      .catch((e) => { if (!controller.signal.aborted) setDetailError(e instanceof Error ? e.message : "正文加载失败"); });
    return () => controller.abort();
  }, [selected]);

  return (
    <div className="space-y-6 p-4 md:p-6">
      <header className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="flex items-center gap-2 text-2xl font-semibold"><BookOpen aria-hidden className="size-6" />知识库</h1>
          <p className="mt-2 text-sm text-muted-foreground">只读导入预览 · 搜索材料、核对来源、对照数据库与附件正文。</p>
        </div>
        <Button variant="outline" disabled={busy} onClick={() => setRefresh((n) => n + 1)}><RefreshCw className="size-4" />重新加载</Button>
      </header>
      <div className="rounded-lg border border-amber-500/40 bg-amber-500/5 p-4 text-sm">
        外部资料不是当前任务证据。来源的“已验证”标记不代表本项目已复核；本页面不会执行材料、创建任务或写入任务经验。
      </div>
      {error && <div role="alert" className="rounded-lg border border-destructive p-4 text-destructive">加载失败，不能判定为空库：{error}</div>}
      {!busy && status && !status.enabled && (
        <section className="rounded-lg border p-6">
          <h2 className="font-semibold">尚未配置知识库预览</h2>
          <p className="mt-2 text-sm text-muted-foreground">先用离线导出脚本生成目录，再为后端设置 ARTEX_KNOWLEDGE_FILE 并重启。原数据库不会被修改。</p>
        </section>
      )}
      {status?.enabled && status.summary && (
        <>
          <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
            {[["材料总数", status.summary.total], ["有数据库正文", status.summary.with_content], ["正文有差异", status.summary.conflicts], ["附件缺失或不唯一", status.summary.missing]].map(([label, count]) => (
              <div key={label} className="rounded-lg border p-4"><p className="text-sm text-muted-foreground">{label}</p><p className="mt-1 text-2xl font-semibold">{count}</p></div>
            ))}
          </div>
          <p className="break-all text-xs text-muted-foreground">来源：{status.source?.name} · SHA-256：{status.source?.sha256}</p>
          <form className="flex flex-wrap items-center gap-3" onSubmit={(e) => { e.preventDefault(); setQuery(draft.trim()); setOffset(0); }}>
            <Input className="min-w-48 flex-1" aria-label="搜索知识库" placeholder="标题、编号、组件、CVE、标签…" value={draft} maxLength={170} onChange={(e) => setDraft(e.target.value)} />
            <Button type="submit" disabled={busy}><Search className="size-4" />搜索</Button>
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={conflicts} onChange={(e) => { setConflicts(e.target.checked); setOffset(0); }} />只看正文差异</label>
          </form>
        </>
      )}
      {busy && <p role="status" className="text-sm text-muted-foreground">正在加载知识库…</p>}
      {list && (
        <div className="grid items-start gap-4 xl:grid-cols-[minmax(280px,2fr)_minmax(0,3fr)]">
          <section aria-label="知识库搜索结果" className="min-w-0 rounded-lg border">
            <div className="border-b px-4 py-3 text-sm">找到 {list.total} 条材料</div>
            {list.items.length === 0 && <p className="p-6 text-sm text-muted-foreground">没有匹配的材料；这不是任务评估结论。</p>}
            <ul className="divide-y">
              {list.items.map((item) => (
                <li key={item.id}><button type="button" aria-pressed={selected === item.id} onClick={() => setSelected(item.id)} className={`w-full space-y-2 p-4 text-left hover:bg-muted/50 ${selected === item.id ? "bg-muted" : ""}`}>
                  <p className="break-words font-medium">{item.title}</p>
                  <p className="break-all text-xs text-muted-foreground">{item.code}</p>
                  <div className="flex flex-wrap gap-2"><Badge variant="outline">{item.kind || "材料"}</Badge><Badge variant={item.attachment_status === "conflict" ? "destructive" : "secondary"}>{attachmentLabels[item.attachment_status]}</Badge>{item.cve && <Badge variant="outline">{item.cve}</Badge>}</div>
                  <p className="text-xs text-muted-foreground">{sourceClaimLabel(item.source_verified)}</p>
                </button></li>
              ))}
            </ul>
            <div className="flex items-center justify-between gap-2 border-t p-3">
              <Button variant="outline" disabled={busy || offset === 0} onClick={() => setOffset(Math.max(0, offset - 20))}>上一页</Button>
              <span className="text-xs">第 {Math.floor(offset / 20) + 1} / {Math.max(1, Math.ceil(list.total / 20))} 页</span>
              <Button variant="outline" disabled={busy || offset + 20 >= list.total} onClick={() => setOffset(offset + 20)}>下一页</Button>
            </div>
          </section>
          <section aria-label="知识库材料详情" className="min-w-0 space-y-4 rounded-lg border p-4 md:p-6">
            {!selected && <p className="text-sm text-muted-foreground">选择左侧材料，查看正文、来源与附件差异。</p>}
            {selected && !detail && !detailError && <p role="status">正在加载正文…</p>}
            {detailError && <p role="alert" className="text-destructive">{detailError}</p>}
            {detail && <>
              <h2 className="break-words text-lg font-semibold">{detail.entry.title}</h2>
              <p className="text-sm text-muted-foreground">{sourceClaimLabel(detail.entry.source_verified)}</p>
              <dl className="space-y-2 break-all text-sm">
                <div><dt className="font-medium">组件 / 版本</dt><dd>{detail.entry.component || "未提供"} / {detail.entry.versions || "未提供"}</dd></div>
                <div><dt className="font-medium">材料来源</dt><dd>{detail.entry.source || "未提供"}</dd><dd>{detail.entry.source_url}</dd></div>
                <div><dt className="font-medium">来源库 SHA-256</dt><dd className="font-mono text-xs">{detail.source.sha256}</dd></div>
              </dl>
              {detail.entry.description && <Material title="说明" body={detail.entry.description} />}
              {detail.entry.usage && <Material title="原始用法（仅文本）" body={detail.entry.usage} />}
              <Material title="数据库正文" body={detail.entry.content} hash={detail.entry.content_sha256} />
              <div className="space-y-3 border-t pt-4">
                <h3 className="font-semibold">附件对照 · {attachmentLabels[detail.entry.attachment.status]}</h3>
                <p className="break-all text-xs text-muted-foreground">{detail.entry.attachment.path || "未提供附件路径"}</p>
                {detail.entry.attachment.status === "conflict" && <p className="text-sm text-amber-600 dark:text-amber-400">数据库与附件内容不同，两个版本均保留，尚未选择可信版本。</p>}
                {["match", "conflict"].includes(detail.entry.attachment.status) && <Material title="附件正文" body={detail.entry.attachment.content} hash={detail.entry.attachment.sha256} />}
              </div>
            </>}
          </section>
        </div>
      )}
    </div>
  );
}
function Material({ title, body, hash }: { title: string; body: string; hash?: string }) {
  return <section className="space-y-2"><h3 className="font-semibold">{title}</h3>{hash && <p className="break-all font-mono text-xs text-muted-foreground">SHA-256: {hash}</p>}<pre className="max-h-96 overflow-auto whitespace-pre-wrap break-words rounded-md bg-muted p-3 font-mono text-xs leading-relaxed">{body || "（空正文）"}</pre></section>;
}
