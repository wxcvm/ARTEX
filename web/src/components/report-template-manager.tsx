"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";
import type { ReportTemplate, ReportTemplateInput } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

const empty=():ReportTemplateInput=>({name:"",brand:"",cover_title:"",subtitle:"",prepared_for:"",author:"",confidentiality:"",footer:""});

export function ReportTemplateSelect({value,onChange}:{value:string;onChange:(v:string)=>void}){
  const locale=useUiLocale();
  const t=(zh:string,en:string)=>uiText(locale,zh,en);
  const [items,setItems]=React.useState<ReportTemplate[]>([]);
  const [error,setError]=React.useState("");
  React.useEffect(()=>{
    let active=true;
    api.reportTemplates().then(r=>{if(active)setItems(r);}).catch(e=>{if(active)setError((e as Error).message);});
    return()=>{active=false;};
  },[]);
  return <div className="flex flex-wrap items-center gap-2">
    <select value={value} onChange={e=>onChange(e.target.value)} aria-label={t("报告模板","Report template")} className="h-8 max-w-64 rounded-md border bg-background px-2 text-sm">
      <option value="">{t("默认专业模板","Default professional")}</option>
      {value&&!items.some(x=>x.id===value)&&<option value={value}>{t("已选模板（待确认）","Selected template (pending)")}</option>}
      {items.map(x=><option value={x.id} key={x.id}>{x.name}</option>)}
    </select>
    {error&&<span role="alert" className="text-xs text-destructive">{error}</span>}
  </div>;
}

export function ReportTemplateManager(){
  const locale=useUiLocale();
  const t=(zh:string,en:string)=>uiText(locale,zh,en);
  const [items,setItems]=React.useState<ReportTemplate[]>([]);
  const [draft,setDraft]=React.useState<ReportTemplateInput|null>(null);
  const [id,setId]=React.useState("");
  const [open,setOpen]=React.useState(false);
  const [busy,setBusy]=React.useState(false);
  const [error,setError]=React.useState("");
  const [confirm,setConfirm]=React.useState(false);
  async function load(){
    setBusy(true);setError("");
    try{setItems(await api.reportTemplates());}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function save(){
    if(!draft||busy)return;
    setBusy(true);setError("");
    try{
      const item=await api.saveReportTemplate(draft,id||undefined);
      setItems(old=>[...old.filter(x=>x.id!==item.id),item].sort((a,b)=>a.name.localeCompare(b.name)));
      setId(item.id);setDraft(item);setConfirm(false);
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function remove(){
    if(!id||busy)return;
    setBusy(true);setError("");
    try{await api.deleteReportTemplate(id);setItems(old=>old.filter(x=>x.id!==id));setId("");setDraft(null);setConfirm(false);}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  async function changeLogo(file?: File){
    if(!id||busy)return;
    setBusy(true);setError("");
    try{
      const item=await api.setReportTemplateLogo(id,file);
      setItems(old=>old.map(x=>x.id===item.id?item:x));
    }catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  const hasLogo=items.find(x=>x.id===id)?.has_logo;
  const fields:{key:keyof ReportTemplateInput;zh:string;en:string}[]=[
    {key:"name",zh:"模板名称",en:"Template name"},
    {key:"brand",zh:"品牌 / 机构",en:"Brand / organization"},
    {key:"cover_title",zh:"封面标题",en:"Cover title"},
    {key:"subtitle",zh:"副标题",en:"Subtitle"},
    {key:"prepared_for",zh:"委托方",en:"Prepared for"},
    {key:"author",zh:"编制人",en:"Prepared by"},
    {key:"confidentiality",zh:"保密说明",en:"Confidentiality"},
    {key:"footer",zh:"页脚文字",en:"Footer text"},
  ];
  return <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
    <CardHeader>
      <CardTitle>{t("报告模板管理","Report templates")}</CardTitle>
      <CardDescription>{t("保存可复用的专业 Word 封面与页脚，导出时选用。原始发现与证据不变。","Save reusable Word cover and footer settings. Findings and evidence remain unchanged.")}</CardDescription>
    </CardHeader>
    <CardContent className="space-y-3">
      <Button size="sm" variant="outline" disabled={busy} onClick={()=>{if(open){setOpen(false);return;}setOpen(true);void load();}}>{open?t("收起","Collapse"):t("管理模板","Manage templates")}</Button>
      {open&&<>
        <div className="flex flex-wrap gap-2">
          <select aria-label={t("选择已有模板","Select saved template")} value={id} disabled={busy} className="h-9 min-w-0 flex-1 rounded-md border bg-background px-2 text-sm" onChange={e=>{const item=items.find(x=>x.id===e.target.value);setId(e.target.value);setDraft(item??null);setConfirm(false);}}>
            <option value="">{t("选择模板","Select template")}</option>
            {items.map(x=><option key={x.id} value={x.id}>{x.name}</option>)}
          </select>
          <Button size="sm" variant="outline" disabled={busy} onClick={()=>{setId("");setDraft(empty());setConfirm(false);}}>{t("新建","New")}</Button>
        </div>
        {draft&&<div className="space-y-2">
          {fields.map(field=><label key={field.key} className="block space-y-1 text-xs">{t(field.zh,field.en)}<Input disabled={busy} value={draft[field.key]} onChange={e=>setDraft({...draft,[field.key]:e.target.value})}/></label>)}
          <p className="text-xs text-muted-foreground">{t("除模板名称外均可留空；封面标题留空时使用报告原始标题。","All fields except the name are optional. An empty cover title uses the original report title.")}</p>
          <div className="flex flex-wrap gap-2">
            <Button size="sm" disabled={busy||!draft.name.trim()} onClick={()=>void save()}>{t("保存","Save")}</Button>
            {id&&(confirm?<><Button size="sm" variant="destructive" disabled={busy} onClick={()=>void remove()}>{t("确认删除","Confirm deletion")}</Button><Button size="sm" variant="ghost" onClick={()=>setConfirm(false)}>{t("取消","Cancel")}</Button></>:<Button size="sm" variant="ghost" disabled={busy} onClick={()=>setConfirm(true)}>{t("删除","Delete")}</Button>)}
          </div>
        </div>}
        {id&&draft&&<div className="space-y-2 rounded-md border p-3">
          <p className="text-sm font-medium">{t("报告品牌 Logo","Report brand logo")} · {hasLogo?t("已设置","Configured"):t("未设置","Not configured")}</p>
          <p className="text-xs text-muted-foreground">{t("PNG / JPEG，最大 2 MiB。图片嵌入 Word 封面，历史报告不受修改影响。","PNG / JPEG, up to 2 MiB. Embedded in the Word cover; saved reports are unchanged.")}</p>
          <Input type="file" accept="image/png,image/jpeg" disabled={busy} onChange={e=>{const file=e.target.files?.[0];e.target.value="";if(file)void changeLogo(file);}}/>
          {hasLogo&&<Button size="sm" variant="outline" disabled={busy} onClick={()=>void changeLogo()}>{t("移除 Logo","Remove logo")}</Button>}
        </div>}
        {error&&<p role="alert" className="text-xs text-destructive">{error}</p>}
      </>}
    </CardContent>
  </Card>;
}
