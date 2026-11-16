"use client";

import * as React from "react";
import { toast } from "sonner";
import { PlusIcon, UploadIcon, SearchIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/lib/api";
import type { KnowledgeDocument, KnowledgeRetrieval } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

type Draft = { id?: string; title: string; category: string; content: string };
const blank = (): Draft => ({title: "", category: "", content: ""});

export function KnowledgeDocumentLibrary({taskId}: {taskId?: string}) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [items,setItems] = React.useState<KnowledgeDocument[]>([]);
  const [query,setQuery] = React.useState("");
  const [category,setCategory] = React.useState("");
  const [filters,setFilters] = React.useState({q:"",category:""});
  const [includeGlobal,setIncludeGlobal] = React.useState(false);
  const [offset,setOffset] = React.useState(0);
  const [more,setMore] = React.useState(false);
  const [loading,setLoading] = React.useState(true);
  const [busy,setBusy] = React.useState(false);
  const [error,setError] = React.useState("");
  const [revision,setRevision] = React.useState(0);
  const [selected,setSelected] = React.useState<KnowledgeDocument | null>(null);
  const [draft,setDraft] = React.useState<Draft | null>(null);
  const [file,setFile] = React.useState<File | null>(null);
  const [uploadCategory,setUploadCategory] = React.useState("");
  const [deleting,setDeleting] = React.useState(false);
  const [history,setHistory] = React.useState<KnowledgeRetrieval[] | null>(null);
  const [historyBusy,setHistoryBusy] = React.useState(false);
  const [historyError,setHistoryError] = React.useState("");
  const input = React.useRef<HTMLInputElement>(null);
  const request = React.useRef(0);

  React.useEffect(() => {
    let active = true;
    setLoading(true);
    setError("");
    api.knowledgeDocuments(taskId, includeGlobal, filters.q, filters.category, offset)
      .then(page => {if(active){setItems(page.items);setMore(page.has_more);}})
      .catch(e => {if(active){setItems([]);setMore(false);setError((e as Error).message);}})
      .finally(() => {if(active)setLoading(false);});
    return () => {active=false;};
  }, [taskId,includeGlobal,filters,offset,revision]);

  async function loadHistory(){
    setHistoryBusy(true);setHistoryError("");
    try{setHistory(await api.knowledgeRetrievals(taskId));}
    catch(e){setHistoryError((e as Error).message);}
    finally{setHistoryBusy(false);}
  }
  function refresh(){setOffset(0);setRevision(n=>n+1);}
  function clearSelection(){request.current++;setSelected(null);setDraft(null);setDeleting(false);}
  async function open(doc: KnowledgeDocument){
    const ticket=++request.current;
    setBusy(true);setError("");setDraft(null);setSelected(null);setDeleting(false);
    try{
      const data=await api.knowledgeDocument(doc.id,taskId,includeGlobal);
      if(ticket===request.current)setSelected(data);
    }catch(e){if(ticket===request.current)setError((e as Error).message);}
    finally{setBusy(false);}
  }
  const canEdit=selected && (selected.task_id || undefined)===taskId;
  async function save(){
    if(!draft||busy)return;
    setBusy(true);setError("");
    try{
      const saved=await api.saveKnowledgeDocument({title:draft.title,category:draft.category,content:draft.content},taskId,draft.id);
      setSelected(saved);setDraft(null);refresh();
      toast.success(tr("文档已保存","Document saved"));
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function upload(){
    if(!file||busy)return;
    setBusy(true);setError("");
    try{
      const doc=await api.uploadKnowledgeDocument(file,uploadCategory,taskId);
      clearSelection();setSelected(doc);setFile(null);
      if(input.current)input.current.value="";
      refresh();toast.success(tr("文档已导入","Document imported"));
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function remove(){
    if(!selected||!canEdit||busy)return;
    setBusy(true);setError("");
    try{
      await api.deleteKnowledgeDocument(selected.id,taskId);
      clearSelection();refresh();toast.success(tr("文档已删除","Document deleted"));
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }

  return <Card>
    <CardHeader>
      <CardTitle>{taskId?tr("任务知识库","Task knowledge"):tr("全局文档库","Global documents")}</CardTitle>
      <CardDescription>{taskId?tr("仅属于此任务的资料随任务归档。可同时检索全局文档；全局资料在此只读。","Task documents follow the task archive. Optionally search global documents; global entries are read-only here."):tr("管理自有资料，与内置只读知识库分开保存。本地检索，不需要模型或外部服务。","Manage your documents separately from the bundled library. Local search requires no model or external service.")}</CardDescription>
    </CardHeader>
    <CardContent className="space-y-4">
      <form className="flex flex-wrap gap-2" onSubmit={e=>{e.preventDefault();clearSelection();setOffset(0);setFilters({q:query.trim(),category:category.trim()});}}>
        <Input className="min-w-[160px] flex-1" aria-label={tr("搜索文档","Search documents")} placeholder={tr("搜索标题、分类或正文","Title, category or text")} value={query} onChange={e=>setQuery(e.target.value)} disabled={busy}/>
        <Input className="w-36" aria-label={tr("分类筛选","Category filter")} placeholder={tr("分类（可选）","Category (optional)")} value={category} onChange={e=>setCategory(e.target.value)} disabled={busy}/>
        <Button disabled={busy||loading}><SearchIcon/>{tr("检索","Search")}</Button>
        {taskId&&<label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={includeGlobal} disabled={busy} onChange={e=>{clearSelection();setOffset(0);setIncludeGlobal(e.target.checked);}}/>{tr("包含全局资料","Include global")}</label>}
        <Button type="button" variant="outline" disabled={busy} onClick={()=>{clearSelection();setDraft(blank());}}><PlusIcon/>{tr("新建","New document")}</Button>
      </form>
      <div className="flex flex-wrap items-center gap-2 rounded-md border p-3">
        <input ref={input} type="file" accept=".txt,.md,.csv,.json,.docx" disabled={busy} aria-label={tr("上传知识文档","Upload knowledge document")} className="min-w-0 flex-1 text-sm" onChange={e=>setFile(e.target.files?.[0]??null)}/>
        <Input className="w-36" aria-label={tr("上传分类","Upload category")} placeholder={tr("分类（可选）","Category (optional)")} value={uploadCategory} onChange={e=>setUploadCategory(e.target.value)} disabled={busy}/>
        <Button variant="outline" disabled={!file||busy} onClick={()=>void upload()}><UploadIcon/>{tr("导入","Import")}</Button>
        <p className="w-full text-xs text-muted-foreground">TXT / MD / CSV / JSON / DOCX · {tr("文件最大 8 MiB，抽取文本最大 4 MiB","8 MiB file, 4 MiB extracted text")}</p>
      </div>
      {error&&<p role="alert" className="text-sm text-destructive">{error}</p>}
      <div className="grid gap-4 lg:grid-cols-[minmax(220px,1fr)_2fr]">
        <div className="space-y-2">
          {loading?<p className="text-sm text-muted-foreground">{tr("加载中…","Loading…")}</p>:items.length===0?<p className="text-sm text-muted-foreground">{tr("暂无匹配文档","No matching documents")}</p>:items.map(doc=><button type="button" key={doc.id} disabled={busy} className="block w-full rounded-md border p-3 text-left hover:bg-muted" onClick={()=>void open(doc)}>
            <div className="break-words text-sm font-medium">{doc.title}</div>
            <div className="mt-1 text-xs text-muted-foreground">{doc.task_id?tr("任务","Task"):tr("全局","Global")}{doc.category?" · "+doc.category:""}</div>
            <p className="mt-1 line-clamp-2 break-words text-xs text-muted-foreground">{doc.content}</p>
          </button>)}
          <div className="flex gap-2">
            <Button variant="outline" size="sm" disabled={loading||busy||offset===0} onClick={()=>setOffset(n=>Math.max(0,n-50))}>{tr("上一页","Previous")}</Button>
            <Button variant="outline" size="sm" disabled={loading||busy||!more} onClick={()=>setOffset(n=>n+50)}>{tr("下一页","Next")}</Button>
          </div>
        </div>
        {draft?<div className="min-w-0 space-y-3 rounded-md border p-3">
          <Input aria-label={tr("文档标题","Document title")} placeholder={tr("标题","Title")} value={draft.title} onChange={e=>setDraft({...draft,title:e.target.value})} disabled={busy}/>
          <Input aria-label={tr("文档分类","Document category")} placeholder={tr("分类","Category")} value={draft.category} onChange={e=>setDraft({...draft,category:e.target.value})} disabled={busy}/>
          <Textarea aria-label={tr("文档正文","Document text")} placeholder={tr("文档正文","Document text")} rows={12} value={draft.content} onChange={e=>setDraft({...draft,content:e.target.value})} disabled={busy}/>
          <div className="flex gap-2"><Button disabled={busy||!draft.title.trim()||!draft.content.trim()} onClick={()=>void save()}>{tr("保存","Save")}</Button><Button variant="outline" disabled={busy} onClick={()=>setDraft(null)}>{tr("取消","Cancel")}</Button></div>
        </div>:selected?<article className="min-w-0 space-y-3 rounded-md border p-3">
          <h3 className="break-words text-lg font-semibold">{selected.title}</h3>
          <p className="text-xs text-muted-foreground">{selected.task_id?tr("任务文档","Task document"):tr("全局文档","Global document")} · {selected.category||tr("未分类","Uncategorized")}{selected.source_filename?" · "+selected.source_filename:""}</p>
          {canEdit&&<div className="flex flex-wrap gap-2"><Button variant="outline" size="sm" disabled={busy} onClick={()=>{setDeleting(false);setDraft({id:selected.id,title:selected.title,category:selected.category,content:selected.content??""});}}>{tr("编辑","Edit")}</Button>
            {!deleting?<Button variant="ghost" size="sm" disabled={busy} onClick={()=>setDeleting(true)}>{tr("删除","Delete")}</Button>:<><Button variant="destructive" size="sm" disabled={busy} onClick={()=>void remove()}>{tr("确认删除","Confirm deletion")}</Button><Button variant="outline" size="sm" disabled={busy} onClick={()=>setDeleting(false)}>{tr("取消","Cancel")}</Button></>}
          </div>}
          <pre className="max-h-[65vh] overflow-auto whitespace-pre-wrap break-words text-sm leading-6">{selected.content}</pre>
        </article>:<p className="text-sm text-muted-foreground">{tr("选择文档查看正文，或创建新资料。","Select a document to view its text, or create one.")}</p>}
      </div>
      <section className="space-y-3 rounded-md border p-3">
        <div className="flex flex-wrap items-center gap-2">
          <h3 className="flex-1 text-sm font-medium">{tr("智能体检索历史","Agent retrieval history")}</h3>
          <Button size="sm" variant="outline" disabled={historyBusy} onClick={()=>void loadHistory()}>{historyBusy?tr("读取中…","Loading…"):history?tr("刷新","Refresh"):tr("查看历史","View history")}</Button>
          {history&&<Button size="sm" variant="ghost" disabled={historyBusy} onClick={()=>setHistory(null)}>{tr("收起","Collapse")}</Button>}
        </div>
        <p className="text-xs text-muted-foreground">{tr("最近 50 次本作用域的检索，保存关键词和命中来源，不保存文档正文。","Last 50 searches in this scope, with queries and source references, not document bodies.")}</p>
        {historyError&&<p role="alert" className="text-sm text-destructive">{historyError}</p>}
        {history?.length===0&&<p className="text-sm text-muted-foreground">{tr("暂无检索记录","No retrievals yet")}</p>}
        {history&&<div className="max-h-[48vh] space-y-2 overflow-auto">{history.map(entry=><article key={entry.id} className="space-y-1 rounded-md border p-2">
          <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            <span>{entry.agent_key||tr("未知智能体","Unknown agent")}</span>
            <span>{new Date(entry.created_at).toLocaleString()}</span>
            {entry.degraded&&<span>{tr("降级检索","Degraded retrieval")}</span>}
          </div>
          <p className="break-words text-sm font-medium">{entry.query}</p>
          <p className="text-xs text-muted-foreground">{tr("命中","Hits")}：{entry.hits.length}</p>
          {entry.hits.map(hit=><div key={hit.scope+":"+hit.id} className="break-words text-xs">
            {hit.scope==="bundled"?tr("内置","Bundled"):hit.scope==="global"?tr("全局","Global"):tr("任务","Task")} · {hit.title} <span className="text-muted-foreground">#{hit.id}</span>
          </div>)}
        </article>)}</div>}
      </section>
    </CardContent>
  </Card>;
}
