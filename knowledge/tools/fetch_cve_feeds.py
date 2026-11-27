#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""抓取公开 CVE 情报源，产出可直接喂给 import_cve_arsenal.py 的合并索引。

数据源（均为公开、无需鉴权的官方接口）：
  * CISA KEV —— 已被在野利用的漏洞目录（JSON 馈源）
  * NVD API 2.0 —— CVSSv3 CRITICAL（>=9.0）条目

产出：
  knowledge/absorb/cve-arsenal/kev.json           原始 KEV 馈源
  knowledge/absorb/cve-arsenal/nvd-critical.json  原始 NVD 响应
  knowledge/absorb/cve-arsenal/INDEX.json         合并去重后的统一条目
  knowledge/absorb/cve-arsenal/FETCH.md           本次抓取的元数据

限速：NVD 无 API key 时 5 请求 / 30 秒，脚本内 sleep 7s；KEV 单请求。
用法：
    python knowledge/tools/fetch_cve_feeds.py --limit 2000
    python knowledge/tools/fetch_cve_feeds.py --limit 2000 --offline   # 只用已落盘的原始文件合并
"""
import argparse
import json
import os
import sys
import time
import urllib.request

HERE = os.path.dirname(os.path.abspath(__file__))
KB = os.path.dirname(HERE)
OUT = os.path.join(KB, 'absorb', 'cve-arsenal')

KEV_URL = ('https://www.cisa.gov/sites/default/files/feeds/'
           'known_exploited_vulnerabilities.json')
NVD_URL = ('https://services.nvd.nist.gov/rest/json/cves/2.0'
           '?cvssV3Severity=CRITICAL&resultsPerPage=2000&startIndex=%d')

UA = 'ARTEX-kb-fetch/1.0 (+https://github.com/wxcvm/ARTEX)'


def fetch(url, timeout=90):
    req = urllib.request.Request(url, headers={'User-Agent': UA})
    with urllib.request.urlopen(req, timeout=timeout) as r:
        return r.read()


def load_kev(offline):
    path = os.path.join(OUT, 'kev.json')
    if offline or not os.path.exists(path):
        if offline:
            raise SystemExit('offline mode but %s missing' % path)
        raw = fetch(KEV_URL)
        os.makedirs(OUT, exist_ok=True)
        with open(path, 'wb') as f:
            f.write(raw)
        print('kev.json    : %d bytes (fetched)' % len(raw))
    with open(path, encoding='utf-8') as f:
        return json.load(f)


def load_nvd(offline, pages):
    path = os.path.join(OUT, 'nvd-critical.json')
    if offline:
        if not os.path.exists(path):
            raise SystemExit('offline mode but %s missing' % path)
        with open(path, encoding='utf-8') as f:
            return json.load(f)
    vulns = []
    total = None
    for p in range(pages):
        url = NVD_URL % (p * 2000)
        raw = fetch(url)
        doc = json.loads(raw)
        total = doc.get('totalResults', total)
        batch = doc.get('vulnerabilities') or []
        vulns.extend(batch)
        print('nvd page %d  : %d items (of %s)' % (p, len(batch), total))
        if len(batch) < 2000:
            break
        time.sleep(7)
    doc = {'vulnerabilities': vulns, 'totalResults': total}
    os.makedirs(OUT, exist_ok=True)
    with open(path, 'w', encoding='utf-8', newline='\n') as f:
        json.dump(doc, f, ensure_ascii=False)
    return doc


def cvss_of(cve):
    metrics = cve.get('metrics') or {}
    for key in ('cvssMetricV31', 'cvssMetricV30', 'cvssMetricV40', 'cvssMetricV2'):
        arr = metrics.get(key)
        if arr:
            d = arr[0].get('cvssData') or {}
            return (d.get('baseScore'), d.get('baseSeverity'), d.get('vectorString'))
    return (None, None, None)


def desc_of(cve):
    for d in cve.get('descriptions') or []:
        if d.get('lang') == 'en':
            return d.get('value') or ''
    return ''


def refs_of(cve):
    out = []
    for r in cve.get('references') or []:
        u = r.get('url')
        if u:
            out.append(u)
    return out


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--limit', type=int, default=2000, help='合并后保留条目上限')
    ap.add_argument('--pages', type=int, default=3, help='NVD 分页数（每页 2000）')
    ap.add_argument('--offline', action='store_true', help='不联网，只用已落盘原始文件')
    args = ap.parse_args()

    kev = load_kev(args.offline)
    nvd = load_nvd(args.offline, args.pages)

    merged = {}
    stamp = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())

    for v in kev.get('vulnerabilities') or []:
        cid = (v.get('cveID') or '').strip().upper()
        if not cid:
            continue
        merged[cid] = {
            'cve': cid,
            'title': v.get('vulnerabilityName') or cid,
            'vendor': v.get('vendorProject') or '',
            'product': v.get('product') or '',
            'severity': 'critical',
            'cvss': None,
            'vector': '',
            'ransomware': (v.get('knownRansomwareCampaignUse') or 'Unknown'),
            'date_added': v.get('dateAdded') or '',
            'description': v.get('shortDescription') or '',
            'required_action': v.get('requiredAction') or '',
            'source': 'cisa-kev',
            'source_url': 'https://www.cisa.gov/known-exploited-vulnerabilities-catalog',
            'references': [],
            'exploited': True,
        }

    for item in nvd.get('vulnerabilities') or []:
        cve = item.get('cve') or {}
        cid = (cve.get('id') or '').strip().upper()
        if not cid:
            continue
        score, sev, vec = cvss_of(cve)
        if cid in merged:
            if score is not None:
                merged[cid]['cvss'] = score
                merged[cid]['vector'] = vec or ''
            refs = refs_of(cve)
            if refs:
                merged[cid]['references'] = refs[:6]
            continue
        desc = desc_of(cve)
        merged[cid] = {
            'cve': cid,
            'title': desc[:160].rstrip('.') or cid,
            'vendor': '',
            'product': '',
            'severity': (sev or 'critical').lower(),
            'cvss': score,
            'vector': vec or '',
            'ransomware': 'Unknown',
            'date_added': (cve.get('published') or '')[:10],
            'description': desc,
            'required_action': '',
            'source': 'nvd',
            'source_url': 'https://nvd.nist.gov/vuln/detail/' + cid,
            'references': refs_of(cve)[:6],
            'exploited': False,
        }

    rows = list(merged.values())
    # 排序：KEV 在野利用优先，其次 CVSS 降序，再按 CVE 编号
    rows.sort(key=lambda r: (0 if r['exploited'] else 1,
                             -(r['cvss'] or 0.0), r['cve']))
    rows = rows[:args.limit]

    os.makedirs(OUT, exist_ok=True)
    idx = os.path.join(OUT, 'INDEX.json')
    with open(idx, 'w', encoding='utf-8', newline='\n') as f:
        json.dump({'generated_at': stamp, 'count': len(rows), 'entries': rows},
                  f, ensure_ascii=False, indent=1)

    kev_n = sum(1 for r in rows if r['exploited'])
    crit_n = sum(1 for r in rows if (r['cvss'] or 0) >= 9.0)
    rans_n = sum(1 for r in rows if r['ransomware'] == 'Known')

    md = os.path.join(OUT, 'FETCH.md')
    with open(md, 'w', encoding='utf-8', newline='\n') as f:
        f.write('# CVE 情报抓取记录\n\n')
        f.write('- 抓取时间：%s\n' % stamp)
        f.write('- KEV 原始条目：%d\n' % len(kev.get('vulnerabilities') or []))
        f.write('- NVD 原始条目：%d（totalResults=%s）\n' % (
            len(nvd.get('vulnerabilities') or []), nvd.get('totalResults')))
        f.write('- 合并去重后保留：%d\n' % len(rows))
        f.write('- 其中 KEV 在野利用：%d\n' % kev_n)
        f.write('- 其中 CVSS >= 9.0：%d\n' % crit_n)
        f.write('- 其中标注勒索组织利用：%d\n' % rans_n)
        f.write('\n数据源：CISA KEV（公开 JSON 馈源）、NVD API 2.0（公开 REST）。\n')
        f.write('所有条目均未在本仓库环境实测，属索引情报，命中只代表候选。\n')

    print('\nmerged      : %d entries' % len(rows))
    print('kev         : %d' % kev_n)
    print('cvss>=9.0   : %d' % crit_n)
    print('ransomware  : %d' % rans_n)
    print('index       : %s' % idx)
    return 0


if __name__ == '__main__':
    sys.exit(main())
