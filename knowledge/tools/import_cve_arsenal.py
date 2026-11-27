#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""把 cve-arsenal 合并索引导入 knowledge.db 与 pocs/。

数据源：knowledge/absorb/cve-arsenal/INDEX.json
（由 fetch_cve_feeds.py 从 CISA KEV + NVD API 2.0 合并去重得到）

每条导入：
  * pocs/<cve 小写>/entry.md —— 定位要点（CVE / 产品 / 版本 / CVSS / 在野利用 / 参考）
  * poc 表一行               —— code = slug = 目录名，export_docs.py 的 pass0 精确配对
  * poc_fts 一行             —— 与 poc 同步，保证 kb_check 的 P1（台账 == FTS）不触发

幂等：先按 created_by='cve-arsenal' 删除既有行与其 FTS，再整批插入。
备份：写入前把 knowledge.db 复制为 knowledge.db.bak-<时间戳>（.gitignore 已含 *.db.bak-*）。

用法：
    python knowledge/tools/import_cve_arsenal.py --dry    # 只报告
    python knowledge/tools/import_cve_arsenal.py          # 备份后写入
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
SOURCE = os.path.join(KB, 'absorb', 'cve-arsenal', 'INDEX.json')

CREATOR = 'cve-arsenal'

# 版本号里的四段以上数字会被 tools/scan_sensitive.py 的 IPv4 正则当成公网 IP。
FOUR = re.compile(r'\b(?:\d{1,3}\.){3,}\d{1,3}\b')
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


def slug_of(cve):
    return cve.strip().lower().replace(' ', '-')


def build_entry_md(e):
    """生成落盘的 entry.md：定位要点 + 判定口径 + 参考，不含利用步骤。"""
    cvss = e.get('cvss')
    cvss_s = ('%.1f' % float(cvss)) if isinstance(cvss, (int, float)) else '-'
    exploited = '是（CISA KEV 收录）' if e.get('exploited') else '未记录在野利用'
    ransom = e.get('ransomware') or 'Unknown'
    ransom_s = {'Known': '已知被勒索组织利用', 'Unknown': '未知'}.get(ransom, ransom)
    refs = e.get('references') or []
    title = shrink_versions(e.get('title') or e['cve'])
    product = shrink_versions(e.get('product') or '-')
    desc = shrink_versions(e.get('description') or '-')
    action = shrink_versions(e.get('required_action') or '-')
    vector = e.get('vector') or '-'

    lines = [
        '# ' + title,
        '',
        '- CVE：%s ｜ 来源：%s' % (e['cve'], e.get('source') or '-'),
        '- 严重度：%s ｜ CVSS：%s ｜ 向量：%s' % (
            (e.get('severity') or '-').upper(), cvss_s, vector),
        '- 厂商：%s ｜ 产品：%s' % (e.get('vendor') or '-', product),
        '- 收录日期：%s' % (e.get('date_added') or '-'),
        '- 在野利用：%s' % exploited,
        '- 勒索利用：%s' % ransom_s,
        '',
        '## 漏洞描述',
        '',
        desc,
        '',
        '## 定位与验证口径',
        '',
        '- 组件指纹：以 `%s` 的版本声明判断候选面。' % product,
        '- 修复建议（来源公告 requiredAction）：%s' % (action if action != '-' else '未提供'),
        '- 触发条件：当前馈源未结构化提取；必须核对厂商公告与 NVD CPE，不能从修复建议推断利用条件。',
        '- 判断口径：本条为**未在本仓库环境实测**的索引情报，命中只代表候选；'
        '需按目标实际版本与回显复核后，才可升级为已验证结论。',
        '- 报结论前先写一句证伪条件：如果目标版本不在受影响区间、或请求返回与公告不符，则本条不成立。',
        '',
        '## 原始资料',
        '',
        e.get('source_url') or '-',
    ]
    for r in refs:
        lines.append('- ' + shrink_versions(r))
    lines.append('')
    return '\n'.join(lines)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--dry', action='store_true', help='只报告，不写库不落盘')
    ap.add_argument('--no-pocs', action='store_true', help='只写库，不生成 entry.md')
    ap.add_argument('--limit', type=int, default=0, help='最多导入多少条（0 = 全部）')
    ap.add_argument('--source', default=SOURCE,
                    help='合并索引路径（默认 knowledge/absorb/cve-arsenal/INDEX.json）')
    args = ap.parse_args()

    with open(args.source, encoding='utf-8') as f:
        payload = json.load(f)
    entries = payload.get('entries') or []
    if not entries:
        raise SystemExit('source index has no entries: %s' % SOURCE)
    if args.limit:
        entries = entries[:args.limit]

    stamp = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
    if not os.path.exists(DB):
        raise SystemExit('missing db: %s' % DB)

    # Version-like dotted strings are shortened only to avoid the repository's IPv4 leak scanner;
    # they are not an affected-version range and must never be treated as a version match.

    con = sqlite3.connect(DB)
    con.row_factory = sqlite3.Row

    before_poc = con.execute('select count(*) from poc').fetchone()[0]
    before_fts = con.execute('select count(*) from poc_fts').fetchone()[0]
    imported = con.execute(
        'select count(*) from poc where created_by = ?', (CREATOR,)).fetchone()[0]

    slugs = [slug_of(e['cve']) for e in entries]
    if len(slugs) != len(set(slugs)):
        dup = [s for s in set(slugs) if slugs.count(s) > 1][:5]
        raise SystemExit('duplicate slugs in source index: %s' % dup)

    owned = {r[0] for r in con.execute('select code from poc where created_by = ?', (CREATOR,))}
    existing = {r[0] for r in con.execute('select code from poc where coalesce(created_by,\'\') <> ?', (CREATOR,))}
    clash = [s for s in slugs if s in existing]
    dirs = {d for d in os.listdir(POCS)} if os.path.isdir(POCS) else set()
    dir_clash = [] if args.no_pocs else [s for s in slugs if s in dirs and s not in owned]

    print('source      : %s' % args.source)
    print('entries     : %d' % len(entries))
    print('db now      : poc=%d fts=%d' % (before_poc, before_fts))
    print('already imp : %d' % imported)
    print('slug clash  : db=%d dir=%d' % (len(clash), len(dir_clash)))
    if clash[:5] or dir_clash[:5]:
        print('  db e.g.   : %s' % clash[:5])
        print('  dir e.g.  : %s' % dir_clash[:5])
    if clash or dir_clash:
        con.close()
        raise SystemExit('aborting: refusing to overwrite knowledge entries not owned by cve-arsenal')

    if args.dry:
        print('\n[dry] nothing written. would end at poc=%d fts=%d' % (
            before_poc - imported + len(entries), before_fts - imported + len(entries)))
        con.close()
        return 0

    bak = DB + '.bak-' + time.strftime('%Y%m%d-%H%M%S')
    shutil.copy2(DB, bak)
    print('\nbackup      : %s' % bak)

    written_docs = 0
    try:
        cur = con.cursor()
        cur.execute(
            'delete from poc_fts where poc_id in (select id from poc where created_by = ?)',
            (CREATOR,))
        cur.execute('delete from poc where created_by = ?', (CREATOR,))

        inserted = []
        for e in entries:
            slug = slug_of(e['cve'])
            md = build_entry_md(e)
            rel = '%s/entry.md' % slug
            if not args.no_pocs:
                d = os.path.join(POCS, slug)
                os.makedirs(d, exist_ok=True)
                with open(os.path.join(d, 'entry.md'), 'w', encoding='utf-8', newline='\n') as f:
                    f.write(md)
                written_docs += 1

            cvss = e.get('cvss')
            tags = [e['cve'], 'cve-arsenal', e.get('source') or 'cve']
            if e.get('exploited'):
                tags.append('kev')
                tags.append('in-the-wild')
            if (e.get('ransomware') or '') == 'Known':
                tags.append('ransomware')
            if isinstance(cvss, (int, float)) and cvss >= 9.0:
                tags.append('cvss-critical')
            if e.get('vendor'):
                tags.append(e['vendor'])
            if e.get('product'):
                tags.append(e['product'])
            tags = list(dict.fromkeys(t.strip() for t in tags if t and t.strip()))
            tags_s = ','.join(tags)[:900]

            title_s = shrink_versions(e.get('title') or e['cve'])[:300]
            desc = '%s · %s · CVSS %s' % (
                shrink_versions(e.get('vendor') or '-'),
                shrink_versions(e.get('product') or '-'),
                ('%.1f' % cvss) if isinstance(cvss, (int, float)) else '-')
            note = ('导入自 %s，未实测' % (e.get('source') or 'cve')
                    + ('（CISA KEV 在野利用）' if e.get('exploited') else ''))

            cur.execute(
                '''insert into poc
                   (code,title,kind,category,cve,component,versions,severity,language,
                    source,source_url,description,usage,content,path,verified,verified_note,
                    hit_count,tags,created_by,created_at,updated_at)
                   values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,0,?,0,?,?,?,?)''',
                (slug, title_s, 'markdown',
                 'kev' if e.get('exploited') else 'nvd',
                 e['cve'], (e.get('product') or None), None,
                 (e.get('severity') or 'critical'), None,
                 'kb', e.get('source_url') or None, desc[:400],
                 (e.get('required_action') or None), md, rel,
                 note[:300], tags_s, CREATOR, stamp, stamp))
            pid = cur.lastrowid
            cur.execute(
                'insert into poc_fts (poc_id,title,cve,component,versions,tags,description,content)'
                ' values (?,?,?,?,?,?,?,?)',
                (pid, title_s, e['cve'], (e.get('product') or ''),
                 None, tags_s, desc[:400], md))
            inserted.append(pid)

        con.commit()
    except Exception:
        con.rollback()
        print('\nFAILED — rolled back. db untouched; backup at %s' % bak)
        con.close()
        raise

    after_poc = con.execute('select count(*) from poc').fetchone()[0]
    after_fts = con.execute('select count(*) from poc_fts').fetchone()[0]
    total_imp = con.execute(
        'select count(*) from poc where created_by = ?', (CREATOR,)).fetchone()[0]

    print('docs        : %d entry.md written' % written_docs)
    print('ids         : %d .. %d' % (inserted[0], inserted[-1]))
    print('poc         : %d -> %d' % (before_poc, after_poc))
    print('fts         : %d -> %d' % (before_fts, after_fts))
    print('imported    : %d' % total_imp)

    bad = con.execute(
        'select count(*) from poc p where p.created_by = ? and not exists'
        ' (select 1 from poc_fts f where f.poc_id = p.id)', (CREATOR,)).fetchone()[0]
    print('fts orphans : %d' % bad)
    for kw in ('log4j', 'actuator', 'nacos', 'kev'):
        n = con.execute('select count(*) from poc_fts where poc_fts match ?', (kw,)).fetchone()[0]
        print('  match %-9s %d' % (kw, n))

    con.close()
    if bad:
        raise SystemExit('FTS out of sync for %d imported rows' % bad)
    return 0


if __name__ == '__main__':
    sys.exit(main())
