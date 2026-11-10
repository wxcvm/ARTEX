# 目标画像
Vue3 + Vite 构建的 SPA（`/assets/index-<hash>.js`），业务接口全部在懒加载 chunk 里。

# 步骤（全部只读 GET；VPS 出口执行，本机 Windows 出网被 TUN 劫持不可信）
1) 抓首页拿入口 JS
   curl -sS -A "$UA" "https://TARGET/" -o page.html
   grep -oE 'src="[^"]+\.js"' page.html

2) 抓入口 JS，提取 chunk 清单（Vite 把 chunk 名内联在入口 bundle 里）
   curl -sS -A "$UA" "https://TARGET/assets/index-XXXX.js" -o index.js
   grep -oE '[A-Za-z0-9_.-]+-[A-Za-z0-9_-]{8}\.js' index.js | sort -u > chunks.txt

3) **关键**：chunk 名就是路由/模块名，按业务关键字挑选，别全下（本例 199 个 chunk，只下 41 个）
   grep -E '^(use|use[A-Z]|.*Console|.*Dashboard|.*Finance|.*Orders|.*Payment|.*Checkout|.*Login|.*Register|.*Profile|Guest.*|Reseller.*)-' chunks.txt \
     | awk '{p=$0; sub(/-[A-Za-z0-9_-]{8}\.js$/,"",p); if(!seen[p]++) print}' > sel.txt   # 同名两套 hash 去重
   while read -r f; do curl -sS -A "$UA" "https://TARGET/assets/$f" -o "chunks/$f"; sleep 1; done < sel.txt

4) 提取真实端点（**比在入口 bundle 里找 /api 有效得多**——axios 调用都在数据层 use*.js）
   grep -ohE '\.(get|post|put|patch|delete)\(("|`)/[^"`]{1,90}' chunks/*.js | sort -u
   # 相对路径即可；baseURL 在入口 JS 里，本例 `Mc="/api/v1"`

5) 区分「路由存在但需鉴权」与「路由不存在」（判 401 vs 404 的观测量必须一致）
   curl -sS -o out -w '%{http_code}' "https://TARGET/api/v1/me"           # 200 + {"status_code":401,...} = 存在需鉴权
   curl -sS -o out -w '%{http_code}' "https://TARGET/api/v1/me/orders"    # 404 page not found      = 不存在

# 本例战果
入口 JS 只含 60 条路由字面量，几乎无 /api；换到 chunk 层后一次提出 **50+ 个真实端点**
（含 /reseller/* 17 条财务与提现、/me/* 用户态、/guest/orders 游客订单、/api-credential）。
