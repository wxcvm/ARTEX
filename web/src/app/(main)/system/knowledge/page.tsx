"use client";
import Link from "next/link";
import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";
import { knowledgeCategory, knowledgeKind } from "@/lib/resource-labels";

export default function KnowledgePage() {
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
  return <main className="mx-auto w-full max-w-5xl space-y-5 p-4 sm:p-8"><h1 className="text-2xl font-semibold">浏览知识库</h1>
    <p className="text-sm text-muted-foreground">只读搜索和查看材料，不需要模型，不会自动执行内容。材料原有的验证标记不代表 ARTEX 已独立复核。</p>
    <form className="flex gap-2" onSubmit={e => { e.preventDefault(); void search(); }}><Input aria-label="知识库关键词" placeholder="输入关键词，例如 Windows、认证、PostgreSQL" value={query} onChange={e => setQuery(e.target.value)} /><Button disabled={busy || !query.trim()}>搜索资料</Button></form>
    {error && <p role="alert" className="rounded-lg border p-3 text-sm">{error}</p>}
    {busy && <p role="status">正在读取资料…</p>}
    {searched && <p className="text-xs text-muted-foreground">本次显示 {hits.length} 条（最多 30 条）；需要更多结果时请细化关键词。</p>}
    <div className="grid gap-4 md:grid-cols-[minmax(220px,1fr)_2fr]"><ul className="space-y-2">{hits.map(hit => <li key={hit.id}><button type="button" className="w-full rounded-lg border p-3 text-left text-sm hover:bg-muted" onClick={() => void open(hit.id)}><b>{hit.title}</b>{hit.category && <p className="mt-1 text-xs text-muted-foreground">{knowledgeCategory(hit.category)}</p>}<p className="mt-1 break-words text-xs text-muted-foreground">{hit.snippet}</p></button></li>)}</ul>
      {entry && <article className="min-w-0 space-y-3 rounded-lg border p-4"><h2 className="text-lg font-semibold">{entry.title}</h2><div className="flex flex-wrap gap-2 text-xs text-muted-foreground">{entry.category && <span>分类：{knowledgeCategory(entry.category)}</span>}{entry.kind && <span>类型：{knowledgeKind(entry.kind)}</span>}<span>条目编号：{entry.id}</span></div><p className="text-xs text-muted-foreground">来源材料原标记：{entry.verified ? "已标记验证" : "未标记验证"}；此处仅展示资料。</p>{entry.usage && <section className="space-y-1"><h3 className="text-sm font-medium">原始使用说明</h3><p className="whitespace-pre-wrap break-words text-sm">{entry.usage}</p></section>}<section className="space-y-1"><h3 className="text-sm font-medium">原始资料正文</h3><p className="text-xs text-muted-foreground">保留来源中的代码、命令、报文和原文，不自动翻译或修改。</p><pre className="whitespace-pre-wrap break-words text-sm leading-6">{entry.content}</pre></section></article>}</div>
    <Link className="text-sm text-primary underline" href="/welcome">返回首次使用引导</Link>
  </main>;
}
