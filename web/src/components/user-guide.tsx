"use client";

import * as React from "react";
import { BookOpenIcon, SearchIcon, XIcon } from "lucide-react";

import { Markdown } from "@/components/markdown";
import { Input } from "@/components/ui/input";
import { guideLinks, parseGuide } from "@/lib/user-guide";

// The source text is embedded when Next.js statically builds /help. No API
// request or second document copy is required at runtime.
export function UserGuide({ markdown }: { markdown: string }) {
  const [search, setSearch] = React.useState("");
  const guide = React.useMemo(() => parseGuide(markdown), [markdown]);
  const sections = React.useMemo(() => guide.sections
    .filter((section) => section.title !== "目录")
    .map((section) => ({ ...section, content: guideLinks(section.content) })), [guide]);
  const query = search.trim().toLocaleLowerCase();
  const visible = query
    ? sections.filter((section) => (section.title + " " + section.content).toLocaleLowerCase().includes(query))
    : sections;

  return (
    <main className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-5 sm:px-6">
      <header className="space-y-2">
        <div className="flex items-center gap-2">
          <BookOpenIcon className="size-6 text-primary" />
          <h1 className="text-2xl font-semibold">ARTEX 中文使用指南</h1>
        </div>
        <p className="text-sm text-muted-foreground">
          安装部署、模型配置、任务操作、Skills、报告与常见问题。
          内容在构建时直接取自仓库文档。
        </p>
        <a className="text-sm text-primary underline underline-offset-2" target="_blank" rel="noreferrer"
          href="https://github.com/wxcvm/ARTEX/blob/main/docs/USER_GUIDE_ZH.md">
          在 GitHub 查看原始手册
        </a>
      </header>

      <label className="relative block max-w-xl">
        <SearchIcon className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input aria-label="搜索使用指南" placeholder="搜索安装、任务、模型、MCP、报错……"
          value={search} onChange={(event) => setSearch(event.target.value)}
          className="pr-9 pl-9" />
        {search && (
          <button type="button" aria-label="清除搜索" onClick={() => setSearch("")}
            className="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground hover:text-foreground">
            <XIcon className="size-4" />
          </button>
        )}
      </label>

      {!query && (
        <>
          <div className="rounded-lg border bg-card p-4">
            <Markdown text={guideLinks(guide.introduction)} internalAnchors />
          </div>
          <nav aria-label="指南章节目录" className="rounded-lg border bg-card p-4">
            <h2 className="mb-3 text-base font-semibold">章节导航</h2>
            <div className="grid gap-2 sm:grid-cols-2">
              {sections.map((section) => (
                <a key={section.id} href={"#" + section.id}
                  className="text-sm text-primary hover:underline">
                  {section.title}
                </a>
              ))}
            </div>
          </nav>
        </>
      )}

      {query && (
        <p aria-live="polite" className="text-sm text-muted-foreground">
          匹配到 {visible.length} 个章节
        </p>
      )}

      {visible.length ? visible.map((section) => (
        <section key={section.id} id={section.id}
          className="scroll-mt-20 space-y-4 rounded-lg border bg-card p-4 sm:p-6">
          <h2 className="text-lg font-semibold">{section.title}</h2>
          <Markdown text={section.content} internalAnchors />
        </section>
      )) : (
        <div className="rounded-lg border p-8 text-center text-sm text-muted-foreground">
          没有找到匹配内容。试试“模型”、“任务”或“安装”。
        </div>
      )}
    </main>
  );
}
