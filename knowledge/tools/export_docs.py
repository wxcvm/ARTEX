#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""从 knowledge.db + pocs/ 重新生成 INDEX.md 与 MAP.md。

    python knowledge/tools/export_docs.py
"""
import difflib
import os
import re
import sqlite3

HERE = os.path.dirname(os.path.abspath(__file__))
KB = os.path.dirname(HERE)
POCS = os.path.join(KB, 'pocs')
DB = os.path.join(KB, 'knowledge.db')


def base(s):
    s = (s or '').replace('<target>', '').replace('tgt', '')
    return re.sub(r'[^0-9a-zA-Z\u4e00-\u9fff]', '', s).lower()


def grams(s):
    return {s[i:i + 2] for i in range(len(s) - 1)} if len(s) > 1 else {s}


def overlap(a, b):
    A, B = grams(a), grams(b)
    return len(A & B) / min(len(A), len(B)) if A and B else 0.0


def main():
    con = sqlite3.connect(DB)
    rows = con.execute('select id,title,category,kind,verified,tags,component from poc order by id').fetchall()
    dirs = sorted(d for d in os.listdir(POCS) if os.path.isdir(os.path.join(POCS, d)))

    # --- INDEX.md ---
    cats = {}
    kinds = {}
    for r in rows:
        cats[r[2] or 'uncategorized'] = cats.get(r[2] or 'uncategorized', 0) + 1
        kinds[r[3] or '?'] = kinds.get(r[3] or '?', 0) + 1
    verified = sum(1 for r in rows if r[4] == 1)
    md = ['# 知识库索引', '',
          '共 **%d** 条，其中已验证 **%d** 条。数据源 `knowledge.db`（SQLite + FTS5，`knowledge_dir` 直接加载）。'
          % (len(rows), verified), '',
          '生成：`python knowledge/tools/export_docs.py`；体检：`python knowledge/tools/kb_check.py`', '',
          '## 分类分布', '', '| category | 条数 |', '|---|---|']
    for k, v in sorted(cats.items(), key=lambda x: -x[1]):
        md.append('| `%s` | %d |' % (k, v))
    md += ['', '## 形态分布', '', '| kind | 条数 |', '|---|---|']
    for k, v in sorted(kinds.items(), key=lambda x: -x[1]):
        md.append('| `%s` | %d |' % (k, v))
    md += ['', '## 全量条目', '', '| id | 标题 | category | kind | 已验证 | tags |', '|---|---|---|---|---|---|']
    for i, t, c, k, v, tags, comp in rows:
        md.append('| %d | %s | %s | %s | %s | %s |' % (
            i, (t or '').replace('|', '/'), c or '', k or '', 'Y' if v == 1 else '', (tags or '')[:70].replace('|', '/')))
    open(os.path.join(KB, 'INDEX.md'), 'w', encoding='utf-8', newline='\n').write('\n'.join(md) + '\n')
    print('INDEX.md:', len(rows), 'rows')

    # --- MAP.md（四轮匹配：code → 归一化精确 → difflib → 二元组）---
    assign, used = {}, set()
    # pass0: 回填行的 code 就是目录名
    codes = dict(con.execute('select id, code from poc').fetchall())
    dirset = set(dirs)
    for i, c in codes.items():
        if c in dirset and c not in used and i not in assign:
            assign[i] = (c, 'code', 1.0)
            used.add(c)
    for d in dirs:
        if d in used:
            continue
        for i, t, *_ in rows:
            if i in assign:
                continue
            bt = base(t)
            if bt and bt == base(d):
                assign[i] = (d, 'exact', 1.0)
                used.add(d)
                break
    cands = []
    for i, t, *_ in rows:
        if i in assign:
            continue
        bt = base(t)
        for d in dirs:
            if d in used:
                continue
            r = difflib.SequenceMatcher(None, bt, base(d)).ratio()
            if r >= 0.74:
                cands.append((r, i, d))
    cands.sort(reverse=True)
    for r, i, d in cands:
        if i in assign or d in used:
            continue
        assign[i] = (d, 'fuzzy', round(r, 3))
        used.add(d)
    cands = []
    for i, t, *_ in rows:
        if i in assign:
            continue
        for d in dirs:
            if d in used:
                continue
            r = overlap(base(t), base(d))
            if r >= 0.62:
                cands.append((r, i, d))
    cands.sort(reverse=True)
    for r, i, d in cands:
        if i in assign or d in used:
            continue
        assign[i] = (d, 'overlap', round(r, 3))
        used.add(d)

    byid = {r[0]: r for r in rows}
    orphan_db = [byid[r[0]] for r in rows if r[0] not in assign]
    orphan_dir = [d for d in dirs if d not in used]

    md = ['# 知识库对账（`knowledge.db` ↔ `pocs/`）', '',
          '| 项 | 数量 |', '|---|---|',
          '| `knowledge.db` 条目 | **%d** |' % len(rows),
          '| `pocs/` 原文目录 | **%d** |' % len(dirs),
          '| 已配对 | **%d** |' % len(assign),
          '| 仅有台账（正文在 db 的 `content`） | **%d** |' % len(orphan_db),
          '| 仅有目录（db 无条目） | **%d** |' % len(orphan_dir),
          '', '> 匹配：归一化精确 → difflib ≥0.74 → 字符二元组重叠 ≥0.62，全局贪心一对一。',
          '', '## 一、仅有台账', '']
    md += ([ '| %d | %s | %s |' % (r[0], (r[1] or '').replace('|', '/'), r[2] or '') for r in orphan_db]
           if orphan_db else ['（无）'])
    md += ['', '## 二、仅有目录', '']
    md += (['- `%s`' % d for d in orphan_dir] if orphan_dir else ['（无）'])
    md += ['', '## 三、配对明细', '', '| id | 标题 | category | 目录 | 方式 | 分数 |', '|---|---|---|---|---|---|']
    for i in sorted(assign):
        d, how, sc = assign[i]
        r = byid[i]
        md.append('| %d | %s | %s | `%s` | %s | %s |' % (i, (r[1] or '').replace('|', '/'), r[2] or '', d, how, sc))
    open(os.path.join(KB, 'MAP.md'), 'w', encoding='utf-8', newline='\n').write('\n'.join(md) + '\n')
    print('MAP.md: matched %d / orphan_db %d / orphan_dir %d' % (len(assign), len(orphan_db), len(orphan_dir)))


if __name__ == '__main__':
    main()
