"use client";
import * as React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";
import type { TaskReportSnapshot } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export function TaskReportHistory({taskId,scope,locale,layout,template}:{taskId:string;scope:"all"|"confirmed";locale:string;layout:"simple"|"professional";template:string}){
  const uiLocale=useUiLocale();
  const tr=(zh:string,en:string)=>uiText(uiLocale,zh,en);
  const [items,setItems]=React.useState<TaskReportSnapshot[]>([]);
  const [offset,setOffset]=React.useState(0);
  const [open,setOpen]=React.useState(false);
  const [name,setName]=React.useState("");
  const [format,setFormat]=React.useState<"md"|"docx">("docx");
  const [busy,setBusy]=React.useState(false);
  const [error,setError]=React.useState("");
  const [message,setMessage]=React.useState("");
  const [confirm,setConfirm]=React.useState("");
  async function load(start=0){
    setBusy(true);setError("");
    try{setItems(await api.taskReportSnapshots(taskId,start));setOffset(start);}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function save(){
    setBusy(true);setError("");setMessage("");
    try{
      const item=await api.saveTaskReportSnapshot(taskId,format,scope,locale,layout,template,name);
      setMessage(tr("已保存报告：","Saved report: ")+item.name);
      setItems(await api.taskReportSnapshots(taskId,0));setOffset(0);setName("");
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function remove(id:string){
    setBusy(true);setError("");
    try{await api.deleteTaskReportSnapshot(taskId,id);setConfirm("");setItems(await api.taskReportSnapshots(taskId,offset));}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function download(item:TaskReportSnapshot){
    setBusy(true);setError("");
    try{await api.downloadTaskReportSnapshot(taskId,item.id,item.format);}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  return <section className="space-y-3 rounded-lg border p-3">
    <div className="flex flex-wrap items-center gap-2">
      <h3 className="flex-1 text-sm font-medium">{tr("已保存报告","Saved reports")}</h3>
      <Button size="sm" variant="outline" disabled={busy} onClick={()=>{if(open){setOpen(false);return;}setOpen(true);void load();}}>{open?tr("收起","Collapse"):tr("管理报告","Manage reports")}</Button>
    </div>
    <p className="text-xs text-muted-foreground">{tr("保存当前任务的报告快照。后续修改发现或模板不会改变已保存文件；快照随任务归档。","Save an immutable report for this task. Later edits to findings or templates do not change it; snapshots follow task archives.")}</p>
    {open&&<>
      <div className="flex flex-wrap items-center gap-2">
        <Input className="min-w-36 flex-1" value={name} disabled={busy} onChange={e=>setName(e.target.value)} placeholder={tr("报告名称（可选）","Report name (optional)")}/>
        <select aria-label={tr("保存格式","Save format")} className="h-9 rounded-md border bg-background px-2 text-sm" value={format} disabled={busy} onChange={e=>setFormat(e.target.value==="md"?"md":"docx")}><option value="docx">Word</option><option value="md">Markdown</option></select>
        <Button size="sm" disabled={busy} onClick={()=>void save()}>{tr("保存当前报告","Save report")}</Button>
      </div>
      <p className="text-xs text-muted-foreground">{tr("按当前任务的最新数据生成，使用上方已选报告范围与 Word 模板。","Generated from current task data using the selected scope and Word template above.")}</p>
      {message&&<p role="status" className="text-xs">{message}</p>}
      {error&&<p role="alert" className="text-xs text-destructive">{error}</p>}
      {items.length===0&&<p className="text-xs text-muted-foreground">{tr("暂无已保存报告","No saved reports")}</p>}
      {items.map(item=><article key={item.id} className="flex flex-wrap items-center gap-2 rounded-md border p-2">
        <div className="min-w-0 flex-1">
          <p className="break-words text-sm font-medium">{item.name}</p>
          <p className="text-xs text-muted-foreground">{item.format.toUpperCase()} · {item.scope==="confirmed"?tr("仅已确认","Confirmed"):tr("全部","All")} · {new Date(item.created_at).toLocaleString()} · {(item.size/1024).toFixed(1)} KiB{item.template_name?" · "+item.template_name:""}</p>
        </div>
        <Button size="sm" variant="outline" disabled={busy} onClick={()=>void download(item)}>{tr("下载","Download")}</Button>
        {confirm===item.id?<><Button size="sm" variant="destructive" disabled={busy} onClick={()=>void remove(item.id)}>{tr("确认删除","Confirm delete")}</Button><Button size="sm" variant="ghost" onClick={()=>setConfirm("")}>{tr("取消","Cancel")}</Button></>:<Button size="sm" variant="ghost" disabled={busy} onClick={()=>setConfirm(item.id)}>{tr("删除","Delete")}</Button>}
      </article>)}
      <div className="flex gap-2"><Button size="sm" variant="ghost" disabled={busy} onClick={()=>void load(offset)}>{tr("刷新","Refresh")}</Button><Button size="sm" variant="ghost" disabled={busy||offset===0} onClick={()=>void load(Math.max(0,offset-50))}>{tr("上一页","Previous")}</Button><Button size="sm" variant="ghost" disabled={busy||items.length<50} onClick={()=>void load(offset+50)}>{tr("下一页","Next")}</Button></div>
    </>}
  </section>;
}
