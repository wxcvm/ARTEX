# 思路
白标（white-label）商城通常只改品牌与主题，**不改代码**。前端 bundle 里常残留上游项目的**兜底默认值/内部标识**，用它就能把"未知自研系统"溯源为开源项目，从而把"在目标上盲试"换成**离线源码审计**（零目标负载、可精确定位）。

# 定位手法（按命中率排序）
1. **品牌名兜底值**：`String(config?.brand?.site_name || "").trim() || "上游项目名"`
   → 本例在 `assets/useRegister-*.js` / `useLogin-*.js` 里抓到 `..."":""?"tgt2403"`，直接暴露上游 = tgt2403（独角数卡 Next）。
2. **app_version**：公开配置接口（`/api/v1/public/config`、`/api/config`、`/config.json` 等）常回版本号 → 本例 `app_version="example.com"`，用于**锁定 tag**。
3. **独有路由/字段名**：`guest_password`、`order_password`、`api-credential`、`reseller/site-config` 等罕用命名，直接丢进 GitHub 代码搜索。
4. **静态资源里的文件头注释/许可证**：如 `Copyright (C) ... The Squid Software Foundation` 之类。
5. **i18n 语言包**：常保留上游 key 命名空间（`auth.register.errors.*`）。

# 落地步骤
```bash
# 1) 前端抓 chunk 后，在数据层文件里搜兜底默认值
grep -ohE '\|\|[[:space:]]*"[A-Za-z][A-Za-z0-9_-]{3,40}"' chunks/*.js | sort -u
# 2) 确认上游项目 + 版本
git ls-remote --tags https://github.com/<org>/<repo> | grep -i "1.4.7"
# 3) 离线取源码（本机出网不稳时用干净出口跳板）
curl -sSL "https://codeload.github.com/<org>/<repo>/tar.gz/refs/tags/example.com" -o src.tgz && tar xzf src.tgz
# 4) 审计重点（按"能否在目标上验证"排序）
#    a. 默认密钥/占位配置：config.yml.example 里的 secret 是否为占位值
grep -rn -i "change-in-production\|changeme\|your-secret" --include=*.example --include=*.yml .
#    b. 路由挂载：找漏挂鉴权中间件的路由（对照 public 与 protected 组）
grep -rn "\.Use(\|Group(" internal/app/httpserver/routes_*.go
#    c. 自认证头/凭据拼接逻辑是否真校验（本例 Guest 头最终落到 SQL 哈希匹配）
#    d. 验证码生成是否用 crypto/rand、是否有限次
```

# 本例战果（tgt6403-com / example.com）
- 溯源成功：**tgt2403 example.com**（Go，tag commit 47cc6207c6796ccfda0526a3708be39e1b59eba3），1152 个 .go 文件离线审计。
- 得到 3 条可验证假设：JWT 默认密钥（目标实测**证伪**）、Guest 头越权（源码+实测**双重证伪**：口令经 hashGuestCredential 参与 SQL 匹配）、send-verify-code 无限流（**candidate**）。
- 关键收获是**止损**：把原本要在目标上盲试的 3 条线，用 0~4 个请求就定性。