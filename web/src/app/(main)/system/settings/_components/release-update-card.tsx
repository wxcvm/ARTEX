"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import type { ReleaseUpdate } from "@/lib/types";
import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

export function ReleaseUpdateCard(){
  const locale=useUiLocale();
  const t=(zh:string,en:string)=>uiText(locale,zh,en);
  const [current,setCurrent]=React.useState("");
  const [result,setResult]=React.useState<ReleaseUpdate|null>(null);
  const [busy,setBusy]=React.useState(false);
  const [error,setError]=React.useState("");
  React.useEffect(()=>{
    let active=true;
    api.health().then(r=>{if(active)setCurrent(r.version);}).catch(()=>{});
    return ()=>{active=false;};
  },[]);
  async function check(){
    setBusy(true);setError("");
    try{const r=await api.checkReleaseUpdate();setResult(r);setCurrent(r.current);}
    catch(e){setError((e as Error).message);}
    finally{setBusy(false);}
  }
  return <Card className="mb-4 min-w-0 break-inside-avoid rounded-xl border-border/70 bg-card/90 shadow-sm xl:mb-5">
    <CardHeader>
      <CardTitle>{t("版本与更新","Version and updates")}</CardTitle>
      <CardDescription>{t("从 ARTEX 官方仓库检查稳定版本，不会自动下载或重启。","Check stable releases from the official ARTEX repository. No automatic download or restart.")}</CardDescription>
    </CardHeader>
    <CardContent className="space-y-3">
      <p className="text-sm">{t("当前版本：","Current: ")}{current||"—"}</p>
      <Button size="sm" variant="outline" disabled={busy} onClick={()=>void check()}>{busy?t("正在检查…","Checking…"):t("检查更新","Check for updates")}</Button>
      {error&&<p role="alert" className="break-words text-xs text-destructive">{error}</p>}
      {result&&<section className="space-y-2 rounded-md border p-3">
        <p className="text-sm font-medium">{t("最新稳定版本：","Latest stable: ")}{result.latest}{result.source==="mirror"&&<span className="ml-2 rounded bg-muted px-1.5 py-0.5 text-[11px] text-muted-foreground">{t("社区镜像","mirror")}</span>}</p>
        {result.note&&<p className="break-words text-xs text-muted-foreground">{result.note}</p>}
        <p className="text-xs text-muted-foreground">{!result.comparable?t("当前为开发或非标准版本，请手动确认。","Development or nonstandard build; compare manually."):result.has_update?t("有新版本可用。","An update is available."):t("当前已是最新稳定版或更新版本。","Current build is up to date or newer.")}</p>
        <a className="block text-sm text-primary underline" href={result.release_url} target="_blank" rel="noreferrer">{t("查看官方发布说明","Official release notes")}</a>
        {result.assets.map(asset=><a key={asset.name} href={asset.url} target="_blank" rel="noreferrer" className="block break-all text-xs text-primary underline">{asset.name} · {Math.max(0.01,asset.size/1048576).toFixed(2)} MiB</a>)}
        <p className="text-xs text-muted-foreground">{t("Windows 桌面用户请选择 Setup 安装器；源码或 Docker 部署继续使用现有更新流程。","Windows desktop users should choose Setup. Source and Docker deployments keep the existing update workflow.")}</p>
      </section>}
    </CardContent>
  </Card>;
}
