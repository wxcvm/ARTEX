#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""knowledge/ 语料体检。

语料是"打过的记录"，不是可构建的工程：一个 .py 里可能贴的是 JS、用法说明或 Python 2。
所以这里只把**台账/FTS 不一致**当 P1（阻断），其余按类别记入 P2。

检查项：
  1. 目录结构   —— 每个 POC 目录至少 1 个文件；0 字节文件
  2. Python     —— py_compile；失败按 NOT_PYTHON / PY2 / SYNTAX 归因
  3. Shell      —— .sh 首行 #!
  4. 重复内容   —— 全文哈希分组
  5. 台账       —— knowledge.db 行数 == poc_fts 行数

用法：
    python knowledge/tools/kb_check.py            # 打印 + 退出码
    python knowledge/tools/kb_check.py --write    # 同时写 knowledge/QA.md
退出码：0 = 无 P1；1 = 有 P1（库缺失或 FTS 与台账不一致）
"""
import hashlib
import os
import py_compile
import re
import sqlite3
import sys
import tempfile
import warnings

HERE = os.path.dirname(os.path.abspath(__file__))
KB = os.path.dirname(HERE)
POCS = os.path.join(KB, 'pocs')
DB = os.path.join(KB, 'knowledge.db')

TEXT_EXT = {'.py', '.sh', '.js', '.php', '.http', '.md', '.txt', '.yaml', '.yml', '.ps1', '.json', '.go', '.java'}
PY2_PRINT = re.compile(r'^\s*print\s+[^\(\s]', re.M)
PY_HINT = re.compile(r'^\s*(#!|import |from |def |class |#)', re.M)


def walk_files(root):
    for dp, dn, fn in os.walk(root):
        dn[:] = [d for d in dn if not d.startswith('.')]
        for f in fn:
            yield os.path.join(dp, f)


def classify_py_failure(path, err):
    try:
        text = open(path, encoding='utf-8', errors='replace').read()
    except OSError:
        return 'SYNTAX', ''
    if not PY_HINT.search(text):
        return 'NOT_PYTHON', 'file extension says .py but the body is not Python (notes / pasted JS / usage text)'
    if PY2_PRINT.search(text):
        return 'PY2', 'Python 2 syntax'
    line = str(err).splitlines()
    msg = line[-1] if line else ''
    return 'SYNTAX', msg[:140]


def main(argv):
    write = '--write' in argv
    p1, p2 = [], []
    stats = {}

    if not os.path.isdir(POCS):
        print('[x] missing', POCS)
        return 1

    dirs = sorted(d for d in os.listdir(POCS) if os.path.isdir(os.path.join(POCS, d)))
    files = list(walk_files(POCS))
    stats['poc_dirs'] = len(dirs)
    stats['poc_files'] = len(files)
    stats['poc_bytes'] = sum(os.path.getsize(f) for f in files)

    for d in dirs:
        p = os.path.join(POCS, d)
        if not any(os.path.isfile(os.path.join(p, x)) for x in os.listdir(p)):
            p2.append(('EMPTY_DIR', d, ''))
    for f in files:
        if os.path.getsize(f) == 0:
            p2.append(('EMPTY_FILE', os.path.relpath(f, POCS), ''))

    py_files = [f for f in files if f.lower().endswith('.py')]
    py_kinds = {}
    with warnings.catch_warnings():
        warnings.simplefilter('ignore')
        for f in py_files:
            try:
                py_compile.compile(f, cfile=os.path.join(tempfile.gettempdir(), 'kbchk.pyc'), doraise=True)
            except py_compile.PyCompileError as e:
                kind, why = classify_py_failure(f, e)
                py_kinds[kind] = py_kinds.get(kind, 0) + 1
                p2.append((kind, os.path.relpath(f, POCS), why))
    stats['py_files'] = len(py_files)
    stats['py_not_runnable'] = sum(py_kinds.values())
    stats['py_kind_detail'] = ', '.join('%s=%d' % (k, v) for k, v in sorted(py_kinds.items())) or '-'

    sh_files = [f for f in files if f.lower().endswith('.sh')]
    for f in sh_files:
        try:
            first = open(f, 'rb').readline()
        except OSError:
            continue
        if not first.startswith(b'#!'):
            p2.append(('NO_SHEBANG', os.path.relpath(f, POCS), first.decode('utf-8', 'replace').strip()[:60]))
    stats['sh_files'] = len(sh_files)

    seen, dup = {}, 0
    for f in files:
        if os.path.splitext(f)[1].lower() not in TEXT_EXT:
            continue
        try:
            h = hashlib.sha256(open(f, 'rb').read()).hexdigest()
        except OSError:
            continue
        rel = os.path.relpath(f, POCS)
        if h in seen:
            dup += 1
            p2.append(('DUPLICATE', rel, 'same as ' + seen[h]))
        else:
            seen[h] = rel
    stats['duplicate_files'] = dup

    if os.path.exists(DB):
        con = sqlite3.connect('file:%s?mode=ro' % DB.replace('\\', '/'), uri=True)
        cur = con.cursor()
        rows = cur.execute('SELECT count(*) FROM poc').fetchone()[0]
        try:
            fts = cur.execute('SELECT count(*) FROM poc_fts').fetchone()[0]
        except sqlite3.Error:
            fts = -1
        con.close()
        stats['db_rows'], stats['fts_rows'] = rows, fts
        if fts != rows:
            p1.append(('FTS_MISMATCH', 'knowledge.db', 'poc=%d fts=%d' % (rows, fts)))
    else:
        p1.append(('NO_DB', 'knowledge.db', 'missing'))

    lines = ['# knowledge/ 语料体检', '',
             '> 语料是记录而非工程：`.py` 里可能贴的是 JS / 用法说明 / Python 2。',
             '> 只有台账与 FTS 不一致算阻断（P1）。', '',
             '| 项 | 值 |', '|---|---|']
    for k, v in stats.items():
        lines.append('| %s | %s |' % (k, v))
    lines += ['', '## P1（阻断）', '']
    lines += ['- %s `%s` %s' % x for x in p1] or ['（无）']
    lines += ['', '## P2（记录）', '']
    if p2:
        lines += ['| 类型 | 文件 | 备注 |', '|---|---|---|']
        lines += ['| %s | `%s` | %s |' % (a, b, c) for a, b, c in p2[:120]]
        if len(p2) > 120:
            lines += ['', '（其余 %d 条省略）' % (len(p2) - 120)]
    else:
        lines.append('（无）')
    lines += ['', '## 说明', '',
              '- `PY2` / `NOT_PYTHON`：上游记录的原始形态，保留原样（改名会破坏目录与台账的对应）。',
              '- `NO_SHEBANG`：不影响阅读，直接 `bash <file>` 或补首行即可。',
              '- `DUPLICATE`：同一段脚本在多个目录各存一份，属历史记录。', '']

    out = '\n'.join(lines)
    print(out)
    if write:
        open(os.path.join(KB, 'QA.md'), 'w', encoding='utf-8').write(out)
        print('[written] knowledge/QA.md')
    return 1 if p1 else 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
