"use client";

export default function GlobalError({ reset }: { reset: () => void }) {
  return (
    <html lang="zh-CN">
      <head>
        <title>ARTEX · 界面恢复</title>
      </head>
      <body style={{ margin: 0, background: "#171c26", color: "#f7f8fb", fontFamily: "system-ui, sans-serif" }}>
        <main role="alert" style={{ minHeight: "100vh", display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", gap: 16, padding: 24, textAlign: "center" }}>
          <h1 style={{ margin: 0, fontSize: 22 }}>ARTEX 界面加载异常</h1>
          <p style={{ maxWidth: 480, lineHeight: 1.6, color: "#bfc9da" }}>
            页面布局遇到了意外问题，可以尝试重新渲染或重新加载。
          </p>
          <div style={{ display: "flex", gap: 12, flexWrap: "wrap", justifyContent: "center" }}>
            <button type="button" onClick={reset} style={{ padding: "10px 20px", cursor: "pointer" }}>
              重试
            </button>
            <button type="button" onClick={() => window.location.reload()} style={{ padding: "10px 20px", cursor: "pointer" }}>
              重新加载界面
            </button>
          </div>
        </main>
      </body>
    </html>
  );
}
