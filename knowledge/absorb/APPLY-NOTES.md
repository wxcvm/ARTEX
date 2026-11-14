# 补丁应用记录（2026-11-14）

本文件记录 `patches/0001`、`patches/0002` 真正落到源码时的实际情况——补丁文件保留在
`patches/` 里作为历史，源码已经包含它们的改动。

## 工具链

| 项 | 值 |
|---|---|
| go | go1.27.4 windows/amd64（winget `GoLang.Go`） |
| GOPROXY | `https://goproxy.cn,https://proxy.golang.org,direct` |
| 基线编译 | `go build ./...` 通过（在打补丁之前先验证过） |

## 应用前的两处前置修正

1. **`knowledge/pocs/easypay-…/poc.go` → `poc.go.txt`**
   `go build ./...` 会把带 `.go` 文件的目录当包，而路径里有中文（`invalid char '支'`）。
   该文件是 PoC 片段不是可编译包，改扩展名后 `./...` 恢复正常。语料里其余非代码后缀同理不受影响。

2. **`gofmt -l` 全仓报警**：仓库文件是 CRLF，gofmt 会把每一行都判成格式差异；
   本次新增的 `kb/kb.go`、`agent/knowledge.go`、`kb/kb_test.go` 是 LF，不受影响。

## 编译与测试

```
git apply patches/0001-guard-hooks.patch patches/0002-knowledge-runtime.patch   # exit 0
go build ./...                                                                  # exit 0
go build ./cmd/artex                                                            # 44.6 MB
go vet ./kb/... ./config/... ./cmd/...                                          # exit 0
go test ./kb/ -count=1                                                          # ok
go test ./server/ -run 'ShippedSkills|SkillFrontmatter|YamlScalar' -count=1      # ok
go test ./config/ -count=1                                                      # ok
```

`kb/kb_test.go` 覆盖：缺失库返回 `(nil,nil)`、FTS 命中、AND→OR 降级、无 FTS 时 LIKE 降级、
`Get` 命中/缺失、空查询报错、引号注入不炸。

## 台账回填

`knowledge/tools/backfill_db.py` 把 60 个只有原文目录、没有台账行的 POC 补进 `poc` 表：

| 阶段 | 台账行数 | FTS 行数 |
|---|---|---|
| 初始（发行包） | 448 | 448 |
| 第一次（宽松匹配） | 486 | 486 |
| 第二次（严格匹配 + 前缀规则） | **508** | **508** |

- `poc_fts` 是**独立 FTS5 表**（不是 external-content），`INSERT INTO poc_fts(poc_fts) VALUES('rebuild')` 对它无效；
  必须 `DELETE` + 按列 `INSERT ... SELECT` 重灌——脚本已按此实现。
- `poc.code` 有 UNIQUE 约束，插入用 `INSERT OR IGNORE`；匹配先用 `code` 再用标题归一化，避免把已有行重复插入。
- 每次写库前自动 `knowledge.db.bak-<ts>`，本次三份备份在验证后已删除。

## 对账现状

`knowledge/MAP.md`：508 行台账 / 429 个目录 / 配对 428 / 仅有台账 80（纯文本条目）/
仅有目录 **1**（`microsoft-google-oauth-refresh_token-…-tgt7737-org-13030…`，台账里的同名行
在脱敏前叫 `aigemini-org-130`，是同一份记录的两个名字，保留不重建）。
