"use client";

import * as React from "react";
import Link from "next/link";
import { FileTextIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ReportTemplateManager } from "@/components/report-template-manager";
import { api } from "@/lib/api";
import type { TaskReportLibraryItem } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export default function ReportCenter(){
  const locale=useUiLocale();
  const tr=(zh:string,en:string)=>uiText(locale,zh,en);
  const [query,setQuery]=React.useState("");
  const [format,setFormat]=React.useState("");
  const [filter,setFilter]=React.useState({q:"",format:""});
  const [items,setItems]=React.useState<TaskReportLibraryItem[]>([]);
  const [offset,setOffset]=React.useState(0);
  const [busy,setBusy]=React.useState(false);
  const [error,setError]=React.useState("");
  const [confirm,setConfirm]=React.useState("");
  const [revision,setRevision]=React.useState(0);
  React.useEffect(()=>{
    let live=true;
    setBusy(true);setError("");setConfirm("");
    api.taskReportLibrary(filter.q,filter.format,offset)
      .then(r=>{if(live)setItems(r);})
      .catch(e=>{if(live)setError((e as Error).message);})
      .finally(()=>{if(live)setBusy(false);});
    return()=>{live=false;};
  },[filter,offset,revision]);
  async function download(item:TaskReportLibraryItem){
    setBusy(true);setError("");
    try{await api.downloadTaskReportSnapshot(item.task_id,item.id,item.format);}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function remove(item:TaskReportLibraryItem){
    setBusy(true);setError("");
    try{await api.deleteTaskReportSnapshot(item.task_id,item.id);setConfirm("");setRevision(n=>n+1);}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  return <main className="mx-auto w-full max-w-6xl space-y-5 p-4 sm:p-8">
    <header className="space-y-2">
      <h1 className="flex items-center gap-2 text-2xl font-semibold"><FileTextIcon className="size-6"/>{tr("报告中心","Report center")}</h1>
      <p className="text-sm text-muted-foreground">{tr("集中管理任务已保存的报告与可复用模板。已归档任务的报告随归档保存，恢复后会重新出现在这里。","Manage saved task reports and reusable templates. Archived task reports return here after restoration.")}</p>
    </header>
    <section className="space-y-3 rounded-xl border bg-card p-4">
      <h2 className="text-lg font-semibold">{tr("历史报告","Saved reports")}</h2>
      <form className="flex flex-wrap gap-2" onSubmit={e=>{e.preventDefault();setOffset(0);setFilter({q:query.trim(),format});}}>
        <Input className="min-w-48 flex-1" value={query} onChange={e=>setQuery(e.target.value)} placeholder={tr("报告名、任务名、模板或任务编号","Report, task, template or task ID")}/>
        <select aria-label={tr("报告格式","Report format")} className="h-9 rounded-md border bg-background px-2 text-sm" value={format} onChange={e=>setFormat(e.target.value)}>
          <option value="">{tr("全部格式","All formats")}</option><option value="docx">Word</option><option value="md">Markdown</option>
        </select>
        <Button disabled={busy}>{tr("搜索","Search")}</Button>
        <Button type="button" variant="outline" disabled={busy} onClick={()=>setRevision(n=>n+1)}>{tr("刷新","Refresh")}</Button>
      </form>
      {error&&<p role="alert" className="text-sm text-destructive">{error}</p>}
      {!items.length&&!busy&&<p className="text-sm text-muted-foreground">{tr("暂无匹配报告。在任务报告页保存报告后即可查看。","No matching reports. Save a report from a task report tab to see it here.")}</p>}
      {items.map(item=><article key={item.id} className="flex flex-wrap items-center gap-3 rounded-lg border p-3">
        <div className="min-w-0 flex-1 space-y-1">
          <p className="break-words font-medium">{item.name}</p>
          <Link className="text-xs text-primary hover:underline" href={"/function/tasks/detail?id="+encodeURIComponent(item.task_id)}>{item.task_name} · #{item.task_id}</Link>
          <p className="text-xs text-muted-foreground">{item.format.toUpperCase()} · {item.scope==="confirmed"?tr("仅已确认","Confirmed only"):tr("全部发现","All findings")} · {new Date(item.created_at).toLocaleString()} · {(item.size/1024).toFixed(1)} KiB{item.template_name?" · "+item.template_name:""}</p>
        </div>
        <Button size="sm" variant="outline" disabled={busy} onClick={()=>void download(item)}>{tr("下载","Download")}</Button>
        {confirm===item.id?<><Button size="sm" variant="destructive" disabled={busy} onClick={()=>void remove(item)}>{tr("确认删除","Confirm deletion")}</Button><Button size="sm" variant="ghost" onClick={()=>setConfirm("")}>{tr("取消","Cancel")}</Button></>:<Button size="sm" variant="ghost" disabled={busy} onClick={()=>setConfirm(item.id)}>{tr("删除","Delete")}</Button>}
      </article>)}
      <div className="flex gap-2"><Button size="sm" variant="outline" disabled={busy||offset===0} onClick={()=>setOffset(n=>Math.max(0,n-50))}>{tr("上一页","Previous")}</Button><Button size="sm" variant="outline" disabled={busy||items.length<50} onClick={()=>setOffset(n=>n+50)}>{tr("下一页","Next")}</Button></div>
    </section>
    <ReportTemplateManager/>
  </main>;
}
