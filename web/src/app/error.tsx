"use client";

export default function PageError({ reset }: { reset: () => void }) {
  return (
    <main className="flex min-h-dvh flex-col items-center justify-center gap-4 bg-background p-6 text-center text-foreground" role="alert">
      <h1 className="text-xl font-semibold">ARTEX 页面暂时无法显示</h1>
      <p className="max-w-lg text-sm text-muted-foreground">
        页面渲染发生异常。可以重试当前页面；如果仍然无法恢复，请重新加载界面。
      </p>
      <div className="flex flex-wrap justify-center gap-3">
        <button type="button" onClick={reset} className="rounded-md bg-primary px-4 py-2 text-primary-foreground">
          重试
        </button>
        <button type="button" onClick={() => window.location.reload()} className="rounded-md border px-4 py-2">
          重新加载界面
        </button>
      </div>
    </main>
  );
}
