"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api, http } from "@/lib/api";
import { MOCK } from "@/lib/mock/enabled";
import type { KnowledgeEmbeddingConfig, KnowledgeIndexStatus } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export function KnowledgeEmbeddingSettings({onSaved}: {onSaved:()=>void}) {
  const locale=useUiLocale();
  const tr=(zh:string,en:string)=>uiText(locale,zh,en);
  const [config,setConfig]=React.useState<KnowledgeEmbeddingConfig|null>(null);
  const [open,setOpen]=React.useState(false);
  const [busy,setBusy]=React.useState(false);
  const [message,setMessage]=React.useState("");
  const [error,setError]=React.useState("");
  async function load(){
    setOpen(true);setBusy(true);setError("");
    try{setConfig(await api.knowledgeEmbedding());}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function test(){
    if(!config)return;
    setBusy(true);setError("");setMessage("");
    try{
      const result=await api.testKnowledgeEmbedding(config);
      if(!result.ok)throw new Error(result.error||tr("连接失败","Connection failed"));
      setMessage(tr("连接成功 · 向量维度 ","Connected · dimensions ")+result.dimension);
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function save(){
    if(!config)return;
    setBusy(true);setError("");setMessage("");
    try{
      setConfig(await api.saveKnowledgeEmbedding(config));
      setMessage(tr("已保存设置","Settings saved"));onSaved();
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  return <section className="space-y-3 rounded-md border p-3">
    <div className="flex flex-wrap items-center gap-2">
      <h3 className="flex-1 text-sm font-medium">{tr("本地语义检索","Local semantic search")}</h3>
      <Button size="sm" variant="outline" disabled={busy} onClick={()=>open?setOpen(false):void load()}>{open?tr("收起","Collapse"):tr("配置","Configure")}</Button>
    </div>
    <p className="text-xs text-muted-foreground">{tr("可选。使用同一台机器上的 LM Studio 或 Ollama 生成向量；未开启或模型不可用时仍可关键词检索。主动建立单文档或批量索引后，智能体才会使用已索引的资料。","Optional. Use LM Studio or Ollama on this machine. Keyword search remains available when disabled or unavailable. Agents use semantic results only after you explicitly index documents.")}</p>
    {open&&config&&<div className="space-y-3">
      <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={config.enabled} disabled={busy} onChange={e=>setConfig({...config,enabled:e.target.checked})}/>{tr("启用本地 Embedding","Enable local embedding")}</label>
      <div className="flex flex-wrap gap-2">
        <Button size="sm" variant="outline" disabled={busy} onClick={()=>setConfig({...config,base_url:"http://127.0.0.1:1234/v1"})}>LM Studio</Button>
        <Button size="sm" variant="outline" disabled={busy} onClick={()=>setConfig({...config,base_url:"http://127.0.0.1:11434/v1"})}>Ollama</Button>
      </div>
      <label className="block space-y-1 text-sm">{tr("本地 API 地址","Local API base URL")}<Input value={config.base_url} placeholder="http://127.0.0.1:1234/v1" disabled={busy} onChange={e=>setConfig({...config,base_url:e.target.value})}/></label>
      <label className="block space-y-1 text-sm">{tr("Embedding 模型名称","Embedding model name")}<Input value={config.model} placeholder="text-embedding-nomic-embed-text-v1.5" disabled={busy} onChange={e=>setConfig({...config,model:e.target.value})}/></label>
      <div className="flex flex-wrap gap-2">
        <Button size="sm" variant="outline" disabled={busy||!config.model.trim()||!config.base_url.trim()} onClick={()=>void test()}>{tr("测试连接","Test connection")}</Button>
        <Button size="sm" disabled={busy||(config.enabled&&(!config.model.trim()||!config.base_url.trim()))} onClick={()=>void save()}>{tr("保存设置","Save settings")}</Button>
      </div>
      <p className="text-xs text-muted-foreground">{tr("服务地址仅允许本机回环。测试仅发送固定测试文本；建立索引会将所选文档正文发往该本地服务。更换模型或编辑文档后需要重新索引。","Loopback endpoints only. Testing sends fixed text; indexing sends the selected document to that local service. Reindex after changing models or editing a document.")}</p>
    </div>}
    {open&&error&&<p role="alert" className="text-sm text-destructive">{error}</p>}
    {open&&message&&<p role="status" className="text-sm">{message}</p>}
  </section>;
}

export function KnowledgeIndexControl({documentId,taskId}:{documentId:string;taskId?:string}) {
  const locale=useUiLocale();
  const tr=(zh:string,en:string)=>uiText(locale,zh,en);
  const [status,setStatus]=React.useState<KnowledgeIndexStatus|null>(null);
  const [busy,setBusy]=React.useState(false);
  const [error,setError]=React.useState("");
  React.useEffect(()=>{
    let active=true;
    api.knowledgeIndex(documentId,taskId)
      .then(s=>{if(active)setStatus(s);})
      .catch(e=>{if(active)setError((e as Error).message);});
    return ()=>{active=false;};
  },[documentId,taskId]);
  async function build(){
    setBusy(true);setError("");
    try{setStatus(await api.buildKnowledgeIndex(documentId,taskId));}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  return <div className="space-y-2 rounded-md border p-2">
    <div className="flex flex-wrap items-center gap-2">
      <span className="flex-1 text-xs">{!status?tr("语义索引状态","Semantic index status"):!status.enabled?tr("本地语义检索未启用","Local semantic search disabled"):status.chunks>0?tr("已索引片段：","Indexed chunks: ")+status.chunks:tr("尚未索引或需要重建","Not indexed or needs rebuild")}</span>
      {status?.enabled&&<Button size="sm" variant="outline" disabled={busy} onClick={()=>void build()}>{busy?tr("索引中…","Indexing…"):status.chunks?tr("重建索引","Rebuild index"):tr("建立索引","Build index")}</Button>}
    </div>
    {status?.enabled&&<p className="text-xs text-muted-foreground">{status.model} · {tr("仅当前文档；正文发送至本机模型服务。","This document only; text is sent to the local model service.")}</p>}
    {error&&<p role="alert" className="text-xs text-destructive">{error}</p>}
  </div>;
}


type KnowledgeBatchStatus = {
  state: "idle" | "running" | "stopping" | "finished" | "cancelled" | "error";
  model?: string; current?: string; succeeded: number; failed: number;
  reused: number; embedded: number; errors?: string[];
};
function batchActive(s: KnowledgeBatchStatus | null) { return s?.state === "running" || s?.state === "stopping"; }

export function KnowledgeBatchIndex({taskId,configRevision,onComplete}:{taskId?:string;configRevision:number;onComplete:()=>void}){
  const locale=useUiLocale();
  const tr=(zh:string,en:string)=>uiText(locale,zh,en);
  const [enabled,setEnabled]=React.useState(false);
  const [status,setStatus]=React.useState<KnowledgeBatchStatus|null>(null);
  const [busy,setBusy]=React.useState(false);
  const [error,setError]=React.useState("");
  const [refresh,setRefresh]=React.useState(0);
  const wasActive=React.useRef(false);
  const completion=React.useRef(onComplete);
  completion.current=onComplete;
  const path="/knowledge/documents/index-batch?"+new URLSearchParams(taskId?{task:taskId}:{});
  const active=batchActive(status);
  React.useEffect(()=>{
    let live=true;
    api.knowledgeEmbedding().then(c=>{if(live)setEnabled(c.enabled);}).catch(()=>{if(live)setEnabled(false);});
    return ()=>{live=false;};
  },[configRevision,taskId]);
  React.useEffect(()=>{
    if(MOCK){setStatus({state:"idle",succeeded:0,failed:0,reused:0,embedded:0});return;}
    let live=true;
    let timer: ReturnType<typeof setTimeout> | undefined;
    async function poll(){
      let delay=5000;
      try{
        const next=await http<KnowledgeBatchStatus>(path);
        if(!live)return;
        setStatus(next);setError("");
        const running=batchActive(next);
        if(wasActive.current&&!running)completion.current();
        wasActive.current=running;
        if(running)delay=1000;
      }catch(e){if(live)setError((e as Error).message);}
      if(live)timer=setTimeout(poll,delay);
    }
    void poll();
    return ()=>{live=false;if(timer)clearTimeout(timer);};
  },[path,refresh]);
  async function change(method:"POST"|"DELETE"){
    setBusy(true);setError("");
    try{
      const next=await http<KnowledgeBatchStatus>(path,{method});
      setStatus(next);
      if(method==="POST"){wasActive.current=batchActive(next);if(!wasActive.current)completion.current();}
      setRefresh(n=>n+1);
    }
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  const state=status?.state;
  const stateText=state==="running"?tr("正在索引","Indexing"):state==="stopping"?tr("正在取消","Cancelling"):state==="finished"?tr("本轮结束","Finished"):state==="cancelled"?tr("已取消","Cancelled"):state==="error"?tr("执行失败","Failed"):tr("尚未启动","Not started");
  return <section className="space-y-2 rounded-md border p-3">
    <div className="flex flex-wrap items-center gap-2">
      <h3 className="flex-1 text-sm font-medium">{tr("后台增量批量索引","Background incremental indexing")}</h3>
      <Button size="sm" variant="outline" disabled={!enabled||busy||active||!status} onClick={()=>void change("POST")}>{tr("索引本库待处理文档","Index pending documents")}</Button>
      {active&&<Button size="sm" variant="ghost" disabled={busy||state==="stopping"} onClick={()=>void change("DELETE")}>{tr("取消索引","Cancel indexing")}</Button>}
    </div>
    <p className="text-xs text-muted-foreground">{tr("后台串行处理当前知识库，离开页面不影响运行。已完成文档保留；取消或服务重启后可重新启动，跳过已索引文档。","Processes this library on the server, even after leaving the page. Completed indexes are retained; restart the job after cancellation or a service restart to skip already indexed documents.")}</p>
    {status&&state!=="idle"&&<p role="status" className="text-xs">{stateText} · {status.model} · {tr("成功 ","Success ")}{status.succeeded}{tr("，失败 ",", failed ")}{status.failed}{tr("；复用片段 ","; reused ")}{status.reused}{tr("，新生成 ",", embedded ")}{status.embedded}{status.current?" · "+status.current:""}</p>}
    {error&&<p role="alert" className="text-xs text-destructive">{error}</p>}
    {!!status?.errors?.length&&<div role="alert" className="max-h-40 space-y-1 overflow-auto text-xs text-destructive">{status.errors.map((msg,i)=><p key={i}>{msg}</p>)}</div>}
  </section>;
}
