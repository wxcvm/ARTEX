#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""母舰域名编号枚举 —— 用「编号前缀 + A 记录枚举」一次性画出攻击者基础设施全图

适用场景：目标站点 CNAME 到形如 nsNN.<mastedomain> 的编号域名，但 apex 域 NXDOMAIN。
本脚本枚举编号前缀，抓出全部节点 IP，并按 /24 聚合输出。

关键坑（我踩过）：只测 example.com / www.example.com / ns1 / ns2 / a / b 会全部 NXDOMAIN，
  → 误判「母舰不存在」。必须枚举「带编号」的前缀（ns1..nsN），记录才会出来。

用法:
    python mastedomain_enum.py <mastedomain> [max_num] [prefix]
    例: python mastedomain_enum.py example.com 400 ns
"""
import socket, sys, concurrent.futures as cf
from collections import defaultdict

def resolve(d):
    try:
        return (d, socket.gethostbyname_ex(d)[2][0])
    except Exception:
        return None

def main(md, maxn=400, prefix='ns', workers=60):
    # 1) 先确认 apex 与常见前缀（预期多为 NXDOMAIN，仅作对照）
    for probe in [md, 'www.'+md, prefix+'1.'+md, prefix+'2.'+md]:
        r = resolve(probe)
        print('  [对照] %-28s -> %s' % (probe, r[1] if r else 'NXDOMAIN'))
    # 2) 枚举编号前缀
    cands = ['%s%d.%s' % (prefix, i, md) for i in range(1, maxn+1)]
    print('\n枚举 %d 个编号前缀 ...' % len(cands))
    with cf.ThreadPoolExecutor(max_workers=workers) as ex:
        hits = [x for x in ex.map(resolve, cands) if x]
    seg = defaultdict(list)
    for d, ip in hits:
        seg['.'.join(ip.split('.')[:3]) + '.0/24'].append((d, ip))
    print('=== %s 母舰节点总数: %d，覆盖 %d 个 /24 ===' % (md, len(hits), len(seg)))
    for s in sorted(seg, key=lambda k: -len(seg[k])):
        print('\n--- %s (%d 台) ---' % (s, len(seg[s])))
        print('   ' + ', '.join('%s=%s' % (d, ip) for d, ip in sorted(seg[s])))
    print('\n建议：① 节点最多的 /24 通常是「母舰核心段」，优先补扫 ② 把每个 nsNN 域名与已知资产对号，'
          '可确认「主站属于哪一号节点」 ③ 用编号把新网段全部登记为新资产。')

if __name__ == '__main__':
    a = sys.argv[1:]
    main(a[0], int(a[1]) if len(a) > 1 else 400, a[2] if len(a) > 2 else 'ns')