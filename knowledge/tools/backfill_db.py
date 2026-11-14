#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""把 pocs/ 里没有台账行的目录回填进 knowledge.db。

背景：knowledge.db 来自发行包（448 行），pocs/ 是合并后的 429 个目录；约 60 个目录
只有原文没有台账行 → agent 的 kb_search 永远搜不到它们。本脚本按目录名归一化匹配，
缺的补一行（title/kind/language/content/tags），然后重建 FTS。

用法：
    python knowledge/tools/backfill_db.py --dry     # 只报告要补多少
    python knowledge/tools/backfill_db.py           # 备份后写入
"""
import argparse
import datetime
import difflib
import os
import re
import shutil
import sqlite3
import time

HERE = os.path.dirname(os.path.abspath(__file__))
KB = os.path.dirname(HERE)
POCS = os.path.join(KB, 'pocs')
DB = os.path.join(KB, 'knowledge.db')

TEXT_EXT = {'.md', '.txt', '.py', '.sh', '.js', '.php', '.http', '.yaml', '.yml', '.ps1', '.json', '.go', '.java', '.html'}
KIND_BY_EXT = {'.py': 'script', '.sh': 'script', '.ps1': 'script', '.js': 'script', '.php': 'poc',
               '.http': 'http', '.md': 'markdown', '.txt': 'markdown', '.yaml': 'template', '.yml': 'template',
               '.json': 'payload', '.go': 'exp', '.java': 'exp', '.html': 'markdown'}
LANG_BY_EXT = {'.py': 'python', '.sh': 'bash', '.ps1': 'powershell', '.js': 'javascript', '.php': 'php',
               '.http': 'http', '.go': 'go', '.java': 'java'}
STOP = {'the', 'and', 'for', 'with', 'from', '未', '的', '与', '和', '及'}


def norm(s):
    s = (s or '').replace('<target>', '').replace('tgt', '')
    return re.sub(r'[^0-9a-zA-Z\u4e00-\u9fff]', '', s).lower()


def dir_files(d):
    out = []
    root = os.path.join(POCS, d)
    for dp, dn, fn in os.walk(root):
        for f in fn:
            out.append(os.path.join(dp, f))
    return sorted(out)


def content_of(files):
    parts, total = [], 0
    for f in files:
        if os.path.splitext(f)[1].lower() not in TEXT_EXT:
            continue
        try:
            t = open(f, encoding='utf-8', errors='replace').read()
        except OSError:
            continue
        head = '### %s\n' % os.path.basename(f)
        body = t[:32768]
        parts.append(head + body)
        total += len(body)
        if total >= 65536:
            break
    return '\n\n'.join(parts)[:65536]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--dry', action='store_true')
    args = ap.parse_args()

    con = sqlite3.connect(DB)
    cur = con.cursor()
    rows = cur.execute('select id, title, code from poc').fetchall()
    known = [(r[0], norm(r[1])) for r in rows]
    codes = {(r[2] or '').strip() for r in rows}

    dirs = sorted(d for d in os.listdir(POCS) if os.path.isdir(os.path.join(POCS, d)))
    missing = []
    for d in dirs:
        if d in codes:          # 目录名就是台账里的 code：已有行，只是标题不同
            continue
        nd = norm(d)
        hit = False
        for _, nk in known:
            if not nd or not nk:
                continue
            # 严格：相等 / 前缀（标题是目录名去掉主机段后的形态）/ 极高相似度
            if (nd == nk
                    or (nd.startswith(nk) and len(nd) - len(nk) <= 40)
                    or (nk.startswith(nd) and len(nk) - len(nd) <= 40)
                    or difflib.SequenceMatcher(None, nd, nk).ratio() >= 0.9):
                hit = True
                break
        if not hit:
            missing.append(d)

    print('台账 %d 行 / 目录 %d 个 / 待补 %d' % (len(rows), len(dirs), len(missing)))
    if args.dry:
        for d in missing[:20]:
            print('  +', d[:70])
        return 0

    backup = DB + '.bak-' + time.strftime('%Y%m%d-%H%M%S')
    if missing:
        shutil.copy2(DB, backup)
        print('备份:', backup)

    now = datetime.datetime.now().isoformat(timespec='seconds')
    added = 0
    for d in missing:
        files = dir_files(d)
        if not files:
            continue
        exts = [os.path.splitext(f)[1].lower() for f in files]
        main_ext = next((e for e in exts if e in KIND_BY_EXT), exts[0])
        tags = ','.join([t for t in re.split(r'[-_/ ]+', d) if len(t) >= 2 and t.lower() not in STOP][:8])
        cur.execute(
            '''INSERT OR IGNORE INTO poc(code,title,kind,category,language,source,description,content,verified,tags,created_by,created_at,updated_at)
               VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)''',
            (d, d, KIND_BY_EXT.get(main_ext, 'poc'), 'absorbed', LANG_BY_EXT.get(main_ext, ''),
             'absorbed', d.replace('-', ' '), content_of(files), 0, tags, 'absorb', now, now))
        added += cur.rowcount

    # poc_fts 是独立 FTS5 表（不是 external-content），'rebuild' 对它无效，
    # 必须按列显式重灌一次。
    cur.execute('DELETE FROM poc_fts')
    cur.execute('''INSERT INTO poc_fts(poc_id, title, cve, component, versions, tags, description, content)
                   SELECT id, IFNULL(title,''), IFNULL(cve,''), IFNULL(component,''), IFNULL(versions,''),
                          IFNULL(tags,''), IFNULL(description,''), IFNULL(content,'')
                     FROM poc''')
    con.commit()
    total = cur.execute('select count(*) from poc').fetchone()[0]
    fts = cur.execute('select count(*) from poc_fts').fetchone()[0]
    con.close()
    print('新增 %d 行；台账 %d 行 / FTS %d 行' % (added, total, fts))
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
