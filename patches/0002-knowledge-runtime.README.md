# 0002 · 知识库运行时（kb_search / kb_get）

**补丁**：`patches/0002-knowledge-runtime.patch`（6 文件，+350/−4）
**状态**：未应用到源码。本机无 Go 工具链，`go build` 需要你本地验证。

## 它解决什么

`knowledge/` 里的 `knowledge.db`（448 条、FTS5 索引）在**当前 fork 里没有任何代码读它** ——
`knowledge_dir` / `kb_search` / `kb_get` 只存在于 0.4.2 发行版二进制。这个补丁把运行时补上：

```
config.json  knowledge_dir  ──┐
ARTEX_KB_DIR（环境变量，优先） ─┴─→ kb.Open() ─→ agent.SetKnowledge() ─→ kb_search / kb_get
```

没有 `knowledge.db` 时**一切照旧**：`kb.Open` 返回 `(nil,nil)`，两个工具不注册，日志打一行。

## 改动清单

| 文件 | 改动 |
|---|---|
| `config/config.go` | `Config.KnowledgeDir` + `KnowledgeDir()`（`ARTEX_KB_DIR` > `knowledge_dir` > 空） |
| `kb/kb.go`（新，234 行） | 只读打开 `knowledge.db`；FTS5 AND → OR → LIKE 三级降级；`Search` / `Get` / `Tier` |
| `agent/knowledge.go`（新，84 行） | 进程级句柄 + `kb_search` / `kb_get` 两个只读工具 + `withKB()` 包装 |
| `agent/tools.go` | `PlannerTools()` 走 `withKB(...)` |
| `agent/tools_insert.go` | `WorkerTools()` 走 `withKB(...)` |
| `cmd/artex/main.go` | 启动时 `kb.Open(config.KnowledgeDir())` + `agent.SetKnowledge` + 一行日志 |

## 设计要点

1. **只读**：`file:...?mode=ro&_pragma=busy_timeout(2000)`，不会有写锁、不会改库。
2. **注入安全**：查询按空白切词，每个词单独加引号，用 `AND` 连接 —— 模型无法注入 `NEAR` / `*` / `-` 之类的 FTS 运算符（那会直接报错或全表扫）。
3. **三级降级**：`AND` 无结果 → `OR` → `LIKE`（返回 `degraded:true`）。别人给的旧库没有 `poc_fts` 也能用。
4. **分档**：`negative`（标题/分类含 判死·负面·证否·全否·不成立）/ `ammo`（verified=1）/ `playbook`（methodology 或 markdown）/ `reference`。
5. **只挂给执行面**：planner / worker 走 `withKB`；mainagent 不挂（人机接口不需要翻库）。

## 验证（本机已完成 / 待你完成）

### 已完成（Python 在真库上跑等价 SQL）

`python D:\渗透\ARTEX-audit\verify_kb_sql.py` + `test_fts.py`：

| 查询 | 结果 |
|---|---|
| `thinkphp` | 5 条（#1 精确版本识别、#375 模板写入 RCE、#2 SQL 回显、#369 Cookie 伪造、#379 filter_exp 失效） |
| `1panel 验证码` | AND 命中 3 条（#131 / #129 / #130） |
| `phpwind 写入` | AND 命中 5 条（#75 头像上传写文件、#249 action 枚举、#200 bank_sort、#250 CSS 编辑、#206 passport_client） |
| `suo5` | 5 条 |
| 不存在的词 | 0 条 → 走 LIKE，仍 0 |
| `Get(1)` | 命中，正文 1686 字符 |
| 分档分布 | ammo 356 / reference 60 / negative 31 / playbook 1 |

### 待你完成

```bash
git apply patches/0002-knowledge-runtime.patch
go build -tags embedui -o artex ./cmd/artex

# 造一个最小可跑配置（或复用现有 config.json，加一行）
#   "knowledge_dir": "D:/渗透/ARTEX-src/knowledge"
ARTEX_KB_DIR=D:/渗透/ARTEX-src/knowledge ./artex -addr 127.0.0.1:8787
```

启动日志应出现其中之一：

```
[kb] knowledge base: D:/渗透/ARTEX-src/knowledge/knowledge.db (fts=true)
[kb] no knowledge.db under ""; kb_* tools disabled          ← 没配
[kb] knowledge base disabled: <error>                        ← 库坏了/被占
```

然后建个任务，让 worker 打一句 `kb_search("1panel 验证码")`，应返回 3 条、`tier=ammo`，
`kb_get(131)` 能取到正文。

## 已知边界

- `verified` 列是 INTEGER，Go 侧按 int 读再转 bool（`driver.Bool` 也能转，但少一层依赖）。
- FTS 分词器是 `unicode61`：中文按连续汉字成一个 token，所以**整句短语**匹配不到，
  必须切词（补丁已做 AND/OR）。
- 库文件 6 MB，`sql.Open` 惰性连接；`Ping` 时才会真正打开，失败只记日志不阻断启动。
- 不监听文件变化：知识库在运行中被替换需要重启进程。
