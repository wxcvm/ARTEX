"use client";
import Link from "next/link";
import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";
import { api } from "@/lib/api";
import { knowledgeCategory, knowledgeKind, knowledgeCategorySummary } from "@/lib/resource-labels";

export default function KnowledgePage() {
  const locale = useUiLocale();
  const t = (zh: string, en: string) => uiText(locale, zh, en);
  const [query, setQuery] = useState("");
  const [hits, setHits] = useState<NonNullable<Awaited<ReturnType<typeof api.knowledgeSearch>>["hits"]>>([]);
  const [entry, setEntry] = useState<Awaited<ReturnType<typeof api.knowledgeEntry>> | null>(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const [searched, setSearched] = useState(false);
  const sequence = useRef(0);
  async function search() {
    const ticket = ++sequence.current; setBusy(true); setError(""); setEntry(null); setHits([]); setSearched(false);
    try { const result = await api.knowledgeSearch(query.trim()); if (ticket === sequence.current) { setHits(result.hits ?? []); setSearched(true); } }
    catch (e) { if (ticket === sequence.current) setError((e as Error).message); }
    finally { if (ticket === sequence.current) setBusy(false); }
  }
  async function open(id: number) {
    const ticket = ++sequence.current; setError(""); setEntry(null); setBusy(true);
    try { const result = await api.knowledgeEntry(id); if (ticket === sequence.current) setEntry(result); }
    catch (e) { if (ticket === sequence.current) setError((e as Error).message); }
    finally { if (ticket === sequence.current) setBusy(false); }
  }
  return <main className="mx-auto w-full max-w-5xl space-y-5 p-4 sm:p-8"><h1 className="text-2xl font-semibold">{t("浏览知识库", "Knowledge base")}</h1>
    <p className="text-sm text-muted-foreground">{t("只读搜索和查看材料，不需要模型，不会自动执行内容。材料原有的验证标记不代表 ARTEX 已独立复核。", "Read-only search and viewing; no model or automatic execution. Source verification marks do not imply independent verification by ARTEX.")}</p>
    <form className="flex gap-2" onSubmit={e => { e.preventDefault(); void search(); }}><Input aria-label={t("知识库关键词", "Knowledge search keywords")} placeholder={t("输入关键词，例如 Windows、认证、PostgreSQL", "Enter keywords, e.g. Windows, authentication, PostgreSQL")} value={query} onChange={e => setQuery(e.target.value)} /><Button disabled={busy || !query.trim()}>{t("搜索资料", "Search")}</Button></form>
    {error && <p role="alert" className="rounded-lg border p-3 text-sm">{error}</p>}
    {busy && <p role="status">{t("正在读取资料…", "Reading materials…")}</p>}
    {searched && <p className="text-xs text-muted-foreground">{t(`本次显示 ${hits.length} 条（最多 30 条）；需要更多结果时请细化关键词。`, `Showing ${hits.length} results (up to 30). Refine keywords for more results.`)}</p>}
    <div className="grid gap-4 md:grid-cols-[minmax(220px,1fr)_2fr]"><ul className="space-y-2">{hits.map(hit => <li key={hit.id}><button type="button" className="w-full rounded-lg border p-3 text-left text-sm hover:bg-muted" onClick={() => void open(hit.id)}><b>{hit.title}</b>{hit.category && <p className="mt-1 text-xs text-muted-foreground">{locale === "en" ? hit.category : knowledgeCategory(hit.category)}</p>}<p className="mt-1 break-words text-xs text-muted-foreground">{hit.snippet}</p></button></li>)}</ul>
      {entry && <article className="min-w-0 space-y-3 rounded-lg border p-4"><h2 className="text-lg font-semibold">{entry.title}</h2><div className="flex flex-wrap gap-2 text-xs text-muted-foreground">{entry.category && <span>{t("分类：", "Category: ")}{locale === "en" ? entry.category : knowledgeCategory(entry.category)}</span>}{entry.kind && <span>{t("类型：", "Type: ")}{locale === "en" ? entry.kind : knowledgeKind(entry.kind)}</span>}<span>{t("条目编号：", "Entry ID: ")}{entry.id}</span></div>{locale !== "en" && entry.category && knowledgeCategorySummary(entry.category) && <section className="space-y-1 rounded-md border p-3"><h3 className="text-sm font-medium">中文分类导读</h3><p className="text-sm text-muted-foreground">{knowledgeCategorySummary(entry.category)}</p><p className="text-xs text-muted-foreground">这是分类主题说明，不是该条目正文的翻译。</p></section>}<p className="text-xs text-muted-foreground">{t("来源材料原标记：", "Source verification mark: ")}{entry.verified ? t("已标记验证", "Marked verified") : t("未标记验证", "Not marked verified")}{t("；此处仅展示资料。", "; shown for reference only.")}</p>{entry.usage && <section className="space-y-1"><h3 className="text-sm font-medium">{t("原始使用说明", "Original usage notes")}</h3><p className="whitespace-pre-wrap break-words text-sm">{entry.usage}</p></section>}<section className="space-y-1"><h3 className="text-sm font-medium">{t("原始资料正文", "Original material")}</h3><p className="text-xs text-muted-foreground">{t("保留来源中的代码、命令、报文和原文，不自动翻译或修改。", "Source code, commands, messages and original text are preserved without automatic translation or modification.")}</p><pre className="whitespace-pre-wrap break-words text-sm leading-6">{entry.content}</pre></section></article>}</div>
    <Link className="text-sm text-primary underline" href="/welcome">{t("返回首次使用引导", "Back to getting started")}</Link>
  </main>;
}
