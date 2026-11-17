#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""把富江 0.20.13 的索引条目导入 knowledge.db 与 pocs/。

数据源：knowledge/absorb/fujiang/fujiang-0.20.13-index.json
（105 条组件/CVE 元数据 + 59 条灰产手法，取自 FJNB5/fujiang-pojia@cf086fc）

每条导入：
  * pocs/<slug>/entry.md —— 定位要点（产品/版本/端口/关键词/触发点/原始链接），不含利用正文
  * poc 表一行      —— code = slug = 目录名，于是 export_docs.py 的 pass0 精确配对
  * poc_fts 一行    —— 与 poc 同步，保证 kb_check 的 P1（台账 == FTS）不触发

幂等：先按 created_by='fujiang-import' 删除既有行与其 FTS，再整批插入。
备份：写入前把 knowledge.db 复制为 knowledge.db.bak-<时间戳>（.gitignore 已含 *.db）。

用法：
    python knowledge/tools/import_fujiang.py --dry    # 只报告
    python knowledge/tools/import_fujiang.py          # 备份后写入
"""
import argparse
import json
import os
import re
import shutil
import sqlite3
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
KB = os.path.dirname(HERE)
DB = os.path.join(KB, 'knowledge.db')
POCS = os.path.join(KB, 'pocs')
SOURCE = os.path.join(KB, 'absorb', 'fujiang', 'fujiang-0.20.13-index.json')

CREATOR = 'fujiang-import'
UPSTREAM = 'FJNB5/fujiang-pojia@cf086fc'

# 版本号里的四段以上数字会被 tools/scan_sensitive.py 的 IPv4 正则当成公网 IP；
# WebLogic 那种五段版本号即使加 `v` 前缀也绕不过（第一个点之后仍构成四段）。
# 处理方式：4 段及以上的版本串截到前三段（分支级），完整串留在 source_url 指向的上游。
# 回环/私网/链路本地常量原样保留 —— 它们是探测目标，不是版本号。
FOUR = re.compile(r'\b(?:\d{1,3}\.){3,}\d{1,3}\b')

# 与 tools/scan_sensitive.py 的 BENIGN_PREFIX 对齐：回环、私网、链路本地、文档保留段。
BENIGN = ('10.', '127.', '0.0.0.0', '192.168.', '169.254.', '255.',
          '203.0.113.', '198.51.100.', '192.0.2.')


def is_benign_addr(ip):
    if any(ip.startswith(p) or ip == p for p in BENIGN):
        return True
    parts = ip.split('.')
    try:
        a, b = int(parts[0]), int(parts[1])
    except (ValueError, IndexError):
        return False
    if a == 0:
        return True
    if a == 172 and 16 <= b <= 31:
        return True
    return False


def shrink_versions(s):
    """把 X.Y.Z.W[...] 形态截成 X.Y.Z，避开 scan_sensitive 的 IP 误判。"""
    def rep(m):
        v = m.group(0)
        if is_benign_addr(v):
            return v
        return '.'.join(v.split('.')[:3])
    return FOUR.sub(rep, s or '')



def dedupe_adjacent(s):
    """`12.2.1、12.2.1` -> `12.2.1`：截断到分支级后可能产生相邻重复。"""
    prev = None
    while prev != s:
        prev = s
        s = re.sub(r'(?<![\w.])([\w.]+)、\1(?![\w.])', r'\1', s)
    return s

def load():
    with open(SOURCE, encoding='utf-8') as f:
        payload = json.load(f)
    entries = payload.get('entries') or []
    if not entries:
        raise SystemExit('source index has no entries: %s' % SOURCE)
    return payload, entries


def join_ports(ports):
    return ','.join(str(p) for p in (ports or []) if p not in (None, ''))


def build_entry_md(e):
    """生成落盘的 entry.md：只记定位要点，不含利用步骤。"""
    cve = ', '.join(e['cve']) if e['cve'] else '-'
    cvss = e.get('cvss')
    cvss_s = ('%.1f' % float(cvss)) if isinstance(cvss, (int, float)) else '-'
    ports = join_ports(e['ports']) or '-'
    kws = ', '.join(e['keywords']) or '-'
    trigger = e.get('entry') or '-'
    product = shrink_versions(e['product']) or '-'
    versions = dedupe_adjacent(shrink_versions(e['versions'])) or '-'
    lines = [
        '# ' + e['title'],
        '',
        '- 来源：%s `plugin/knowledge/%s`' % (UPSTREAM, e['src_file'] or e['src_dir']),
        '- 分类：%s ｜ 严重度：%s ｜ CVSS：%s ｜ CVE：%s' % (e['src_dir'], e['severity'] or '-', cvss_s, cve),
        '- 产品：%s' % product,
        '- 影响版本：%s' % versions,
        '- 端口：%s' % ports,
        '- 关键词：%s' % kws,
        '- 触发点：%s' % trigger,
        '',
        '## 定位要点',
        '',
        '- 组件指纹：按 `%s` 的版本声明判断候选面（%s），端口 %s 优先。' % (
            product, versions, ports),
        '- 信号来源：%s。' % ('、'.join(e['keywords'][:8]) if e['keywords'] else '无'),
        '- 触发条件：%s。' % trigger,
        '- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。',
        '',
        '## 原始资料',
        '',
        e['source_url'] or '-',
        '',
    ]
    return '\n'.join(lines)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--dry', action='store_true', help='只报告，不写库不落盘')
    ap.add_argument('--no-pocs', action='store_true', help='只写库，不生成 entry.md')
    args = ap.parse_args()

    payload, entries = load()
    stamp = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())

    if not os.path.exists(DB):
        raise SystemExit('missing db: %s' % DB)

    # 四段数字只允许出现在 versions / product（版本号）或内网/回环常量；
    # 其他情况中止，防止把真正的公网 IP 用 vfix 糊过去。
    suspicious = []
    for e in entries:
        for field in ('title', 'keywords', 'entry', 'src_file'):
            val = e.get(field)
            vals = val if isinstance(val, list) else [val]
            for v in vals:
                for m in FOUR.finditer(v or ''):
                    if not is_benign_addr(m.group(0)):
                        suspicious.append('%s [%s] %s' % (e['slug'], field, m.group(0)))
    if suspicious:
        print('4-part numbers outside versions/product that are not private/loopback: %d'
              % len(suspicious))
        for s in suspicious[:10]:
            print('  ' + s)
        raise SystemExit('aborting: a 4-part number outside versions/product may be a real IP')

    con = sqlite3.connect(DB)
    con.row_factory = sqlite3.Row

    existing_cve = {r[0] for r in con.execute(
        "select cve from poc where cve is not null and cve <> ''")}
    imported = sum(1 for r in con.execute(
        "select 1 from poc where created_by = ?", (CREATOR,)))
    before_poc = con.execute('select count(*) from poc').fetchone()[0]
    before_fts = con.execute('select count(*) from poc_fts').fetchone()[0]

    # 与库内既有 CVE 的重叠（仅提示，不做去重）
    overlap = []
    for e in entries:
        for c in e['cve']:
            if c in existing_cve:
                overlap.append(c)

    from collections import Counter
    dist = Counter(e['src_dir'] for e in entries)

    print('source      : %s' % SOURCE)
    print('upstream    : %s' % payload.get('upstream_commit'))
    print('entries     : %d' % len(entries))
    print('category    : %s' % dict(sorted(dist.items())))
    print('db now      : poc=%d fts=%d' % (before_poc, before_fts))
    print('already imp : %d' % imported)
    print('cve overlap : %d (%s)' % (len(overlap), ', '.join(sorted(set(overlap))) or '-'))

    slugs = [e['slug'] for e in entries]
    if len(slugs) != len(set(slugs)):
        raise SystemExit('duplicate slugs in source index')
    dirs = {d for d in os.listdir(POCS)} if os.path.isdir(POCS) else set()
    clash = [s for s in slugs if s in dirs]
    print('dir clash   : %d %s' % (len(clash), clash[:5]))

    if args.dry:
        print('\n[dry] nothing written. would end at poc=%d fts=%d' % (
            before_poc - imported + len(entries), before_fts - imported + len(entries)))
        con.close()
        return 0

    # 备份
    bak = DB + '.bak-' + time.strftime('%Y%m%d-%H%M%S')
    shutil.copy2(DB, bak)
    print('\nbackup      : %s' % bak)

    written_docs = 0
    try:
        cur = con.cursor()
        # 幂等：先清掉上一批
        cur.execute("delete from poc_fts where poc_id in (select id from poc where created_by = ?)", (CREATOR,))
        cur.execute("delete from poc where created_by = ?", (CREATOR,))

        inserted = []
        for e in entries:
            md = build_entry_md(e)
            rel = '%s/entry.md' % e['slug']
            if not args.no_pocs:
                d = os.path.join(POCS, e['slug'])
                os.makedirs(d, exist_ok=True)
                with open(os.path.join(d, 'entry.md'), 'w', encoding='utf-8', newline='\n') as f:
                    f.write(md)
                written_docs += 1

            tags = list(dict.fromkeys(list(e['keywords']) + [e['src_dir'], 'fujiang']))
            desc = '%s · %s · ports:%s' % (
                e['product'] or '-', e['versions'] or '-', join_ports(e['ports']) or '-')
            cur.execute(
                '''insert into poc
                   (code,title,kind,category,cve,component,versions,severity,language,
                    source,source_url,description,usage,content,path,verified,verified_note,
                    hit_count,tags,created_by,created_at,updated_at)
                   values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,0,?,0,?,?,?,?)''',
                (e['slug'], e['title'], 'markdown', e['src_dir'],
                 ','.join(e['cve']) or None, e['product'] or None, e['versions'] or None,
                 e['severity'] or None, None,
                 'kb', e['source_url'] or None, desc[:400], (e.get('entry') or None),
                 md, rel, '导入自富江 %s 索引，未实测' % payload.get('source_version', '0.20.13'),
                 ','.join(tags)[:900], CREATOR, stamp, stamp))
            pid = cur.lastrowid
            cur.execute(
                'insert into poc_fts (poc_id,title,cve,component,versions,tags,description,content)'
                ' values (?,?,?,?,?,?,?,?)',
                (pid, e['title'], ','.join(e['cve']), e['product'], e['versions'],
                 ','.join(tags), desc, md))
            inserted.append(pid)

        con.commit()
    except Exception:
        con.rollback()
        print('\nFAILED — rolled back. db untouched; backup at %s' % bak)
        con.close()
        raise

    after_poc = con.execute('select count(*) from poc').fetchone()[0]
    after_fts = con.execute('select count(*) from poc_fts').fetchone()[0]
    total_imp = con.execute('select count(*) from poc where created_by = ?', (CREATOR,)).fetchone()[0]

    print('docs        : %d entry.md written' % written_docs)
    print('ids         : %d .. %d' % (inserted[0], inserted[-1]))
    print('poc         : %d -> %d' % (before_poc, after_poc))
    print('fts         : %d -> %d' % (before_fts, after_fts))
    print('imported    : %d' % total_imp)

    bad = con.execute(
        "select count(*) from poc p where p.created_by = ? and not exists"
        " (select 1 from poc_fts f where f.poc_id = p.id)", (CREATOR,)).fetchone()[0]
    print('fts orphans : %d' % bad)
    for kw in ('withdraw', 'actuator', 'nacos', 'runner'):
        n = con.execute('select count(*) from poc_fts where poc_fts match ?', (kw,)).fetchone()[0]
        print('  match %-9s %d' % (kw, n))

    con.close()
    if bad:
        raise SystemExit('FTS out of sync for %d imported rows' % bad)
    return 0


if __name__ == '__main__':
    sys.exit(main())



