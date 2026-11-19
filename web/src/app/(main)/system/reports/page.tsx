"use client";

import * as React from "react";
import Link from "next/link";
import { FileTextIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ReportTemplateManager } from "@/components/report-template-manager";
import { api } from "@/lib/api";
import { REPORT_PAGE_SIZE, reportLibraryView, reportOffsetAfterDeletion, type ReportLibraryLoad } from "@/lib/report-library-state";
import type { TaskReportLibraryItem } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export default function ReportCenter(){
  const locale=useUiLocale();
  const tr=(zh:string,en:string)=>uiText(locale,zh,en);
  const [query,setQuery]=React.useState("");
  const [format,setFormat]=React.useState("");
  const [filter,setFilter]=React.useState({q:"",format:""});
  const [load,setLoad]=React.useState<ReportLibraryLoad<TaskReportLibraryItem>|null>(null);
  const [offset,setOffset]=React.useState(0);
  const [acting,setActing]=React.useState(false);
  const actionLock=React.useRef(false);
  const [actionError,setActionError]=React.useState("");
  const [confirm,setConfirm]=React.useState("");
  const [revision,setRevision]=React.useState(0);
  const requestKey=JSON.stringify([filter.q,filter.format,offset]);
  const view=reportLibraryView(load,requestKey,revision);
  const loading=view.status==="loading";
  const page=Math.floor(offset/REPORT_PAGE_SIZE)+1;
  React.useEffect(()=>{
    let live=true;
    setConfirm("");
    setLoad(previous=>({key:requestKey,revision,status:"loading",items:previous?.key===requestKey?previous.items:[]}));
    api.taskReportLibrary(filter.q,filter.format,offset)
      .then(items=>{if(live)setLoad({key:requestKey,revision,status:"ready",items});})
      .catch(e=>{if(live)setLoad(previous=>({key:requestKey,revision,status:"error",items:previous?.key===requestKey?previous.items:[],error:(e as Error).message}));});
    return()=>{live=false;};
  },[filter.q,filter.format,offset,requestKey,revision]);
  function reload(){
    if(actionLock.current)return;
    setActionError("");setConfirm("");setRevision(n=>n+1);
  }
  async function act(item:TaskReportLibraryItem,action:"download"|"remove"){
    if(actionLock.current||!view.canAct||!view.items.some(row=>row.id===item.id&&row.task_id===item.task_id))return;
    actionLock.current=true;
    setActing(true);setActionError("");
    try{
      if(action==="download")await api.downloadTaskReportSnapshot(item.task_id,item.id,item.format);
      else{
        await api.deleteTaskReportSnapshot(item.task_id,item.id);
        setConfirm("");
        setOffset(reportOffsetAfterDeletion(offset,view.items.length));
        setRevision(n=>n+1);
      }
    }
    catch(e){setActionError((e as Error).message);}
    finally{actionLock.current=false;setActing(false);}
  }
  function changePage(next:number){
    if(actionLock.current)return;
    setActionError("");setConfirm("");setOffset(next);
  }
  return <main className="mx-auto w-full max-w-6xl space-y-5 p-4 sm:p-8">
    <header className="space-y-2">
      <h1 className="flex items-center gap-2 text-2xl font-semibold"><FileTextIcon className="size-6"/>{tr("报告中心","Report center")}</h1>
      <p className="text-sm text-muted-foreground">{tr("集中管理任务已保存的报告与可复用模板。已归档任务的报告随归档保存，恢复后会重新出现在这里。","Manage saved task reports and reusable templates. Archived task reports return here after restoration.")}</p>
    </header>
    <section className="space-y-3 rounded-xl border bg-card p-4">
      <h2 className="text-lg font-semibold">{tr("历史报告","Saved reports")}</h2>
      <form className="flex flex-wrap gap-2" onSubmit={e=>{
        e.preventDefault();
        if(actionLock.current)return;
        setActionError("");setConfirm("");setOffset(0);setFilter({q:query.trim(),format});setRevision(n=>n+1);
      }}>
        <Input aria-label={tr("搜索历史报告","Search saved reports")} className="min-w-48 flex-1" value={query} onChange={e=>setQuery(e.target.value)} placeholder={tr("报告名、任务名、模板或任务编号","Report, task, template or task ID")}/>
        <select aria-label={tr("报告格式","Report format")} className="h-9 rounded-md border bg-background px-2 text-sm" value={format} onChange={e=>setFormat(e.target.value)}>
          <option value="">{tr("全部格式","All formats")}</option><option value="docx">Word</option><option value="md">Markdown</option>
        </select>
        <Button disabled={acting}>{tr("搜索","Search")}</Button>
        <Button type="button" variant="outline" disabled={acting||loading} onClick={reload}>{tr("刷新","Refresh")}</Button>
      </form>
      <p className="text-xs text-muted-foreground">{tr(`第 ${page} 页`,`Page ${page}`)} · {filter.format.toUpperCase()||tr("全部格式","All formats")}{filter.q&&<> · {tr("关键词","Keyword")}: {filter.q}</>}</p>
      {loading&&<p role="status" className="text-sm text-muted-foreground">{tr("正在加载报告…","Loading reports…")}</p>}
      {view.status==="error"&&<div className="flex flex-wrap items-center gap-2">
        <p role="alert" className="text-sm text-destructive">{tr("报告加载失败","Could not load reports")}{view.error&&`: ${view.error}`}</p>
        <Button type="button" size="sm" variant="outline" onClick={reload}>{tr("重试","Retry")}</Button>
      </div>}
      {actionError&&<p role="alert" className="text-sm text-destructive">{actionError}</p>}
      {view.stale&&<p className="text-xs text-muted-foreground">{tr("显示当前页上一次的结果；重新加载成功后才能操作。","Showing the last results for this page. Actions require a successful reload.")}</p>}
      {view.empty&&<p className="text-sm text-muted-foreground">{tr("暂无匹配报告。在任务报告页保存报告后即可查看。","No matching reports. Save a report from a task report tab to see it here.")}</p>}
      {view.items.map(item=><article key={item.id} className="flex flex-wrap items-center gap-3 rounded-lg border p-3">
        <div className="min-w-0 flex-1 space-y-1">
          <p className="break-words font-medium">{item.name}</p>
          <Link className="text-xs text-primary hover:underline" href={"/function/tasks/detail?id="+encodeURIComponent(item.task_id)}>{item.task_name} · #{item.task_id}</Link>
          <p className="text-xs text-muted-foreground">{item.format.toUpperCase()} · {item.scope==="confirmed"?tr("仅已确认","Confirmed only"):tr("全部发现","All findings")} · {new Date(item.created_at).toLocaleString(locale)} · {(item.size/1024).toFixed(1)} KiB{item.template_name?" · "+item.template_name:""}</p>
        </div>
        <Button size="sm" variant="outline" disabled={acting||!view.canAct} onClick={()=>void act(item,"download")}>{tr("下载","Download")}</Button>
        {view.canAct&&confirm===item.id?<><Button size="sm" variant="destructive" disabled={acting} onClick={()=>void act(item,"remove")}>{tr("确认删除","Confirm deletion")}</Button><Button size="sm" variant="ghost" disabled={acting} onClick={()=>setConfirm("")}>{tr("取消","Cancel")}</Button></>:<Button size="sm" variant="ghost" disabled={acting||!view.canAct} onClick={()=>setConfirm(item.id)}>{tr("删除","Delete")}</Button>}
      </article>)}
      <div className="flex gap-2"><Button size="sm" variant="outline" disabled={acting||loading||offset===0} onClick={()=>changePage(Math.max(0,offset-REPORT_PAGE_SIZE))}>{tr("上一页","Previous")}</Button><Button size="sm" variant="outline" disabled={acting||!view.canAct||view.items.length<REPORT_PAGE_SIZE} onClick={()=>changePage(offset+REPORT_PAGE_SIZE)}>{tr("下一页","Next")}</Button></div>
    </section>
    <ReportTemplateManager/>
  </main>;
}
