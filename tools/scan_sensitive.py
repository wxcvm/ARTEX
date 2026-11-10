#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""敏感信息自检：扫 knowledge/ 与 docs/ 里的公网 IP、本机路径、客户数据文件。

用法：
    python tools/scan_sensitive.py            # 扫默认目录
    python tools/scan_sensitive.py --strict   # 连已知基础设施 IP 也报（1.1.1.1 / 8.8.8.8 等）
    python tools/scan_sensitive.py path ...   # 扫指定路径

命中即退出码 1，用于 CI / 提交前自检。
"""
import os
import re
import sys

TEXT_EXT = {'.md', '.txt', '.json', '.yaml', '.yml', '.py', '.sh', '.ps1', '.js', '.mjs',
            '.ts', '.tsx', '.http', '.go', '.sql', '.html', '.php', '.java'}
IPV4 = re.compile(r'\b(?:\d{1,3}\.){3}\d{1,3}\b')
HOME = re.compile(r'[A-Za-z]:\\+Users\\+[^\\\s"\'<>|]+')

# 私有 / 文档保留段：不算泄漏
BENIGN_PREFIX = ('10.', '127.', '0.0.0.0', '192.168.', '169.254.', '255.',
                 '203.0.113.', '198.51.100.', '192.0.2.')
# 公网但属于常见基础设施 / 测试常量，默认不报
INFRA = {'1.1.1.1', '8.8.8.8', '8.8.4.4', '1.2.3.4', '1.0.0.1', '9.9.9.9',
         '114.114.114.114', '223.5.5.5'}
FORBIDDEN_FILES = {'assets.db', 'artex.exe', 'artex', 'knowledge.db.bak'}

def is_benign(ip, strict):
    if any(ip.startswith(p) or ip == p for p in BENIGN_PREFIX):
        return True
    oc = ip.split('.')[0]
    # 0.x（非法主机/版本号误匹配）、172.16/12（Docker 与内网常用段）
    if oc == '0':
        return True
    try:
        first, second = int(oc), int(ip.split('.')[1])
        if first == 172 and 16 <= second <= 31:
            return True
    except ValueError:
        pass
    if not strict and ip in INFRA:
        return True
    return False

def scan(root, strict=False):
    hits = []
    for dp, dn, fn in os.walk(root):
        dn[:] = [d for d in dn if d not in {'.git', 'node_modules', 'dist', 'out'}]
        for name in dn + fn:
            if name in FORBIDDEN_FILES:
                hits.append((os.path.join(dp, name), 0, 'FORBIDDEN_FILE', name))
        for f in fn:
            p = os.path.join(dp, f)
            if os.path.splitext(f)[1].lower() not in TEXT_EXT:
                continue
            try:
                text = open(p, encoding='utf-8', errors='replace').read()
            except OSError:
                continue
            for i, line in enumerate(text.splitlines(), 1):
                for m in IPV4.finditer(line):
                    ip = m.group(0)
                    if not is_benign(ip, strict):
                        hits.append((p, i, 'PUBLIC_IP', ip))
                if HOME.search(line):
                    hits.append((p, i, 'LOCAL_PATH', HOME.search(line).group(0)))
    return hits

def main(argv):
    strict = '--strict' in argv
    paths = [a for a in argv if not a.startswith('--')]
    if not paths:
        here = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        paths = [p for p in (os.path.join(here, 'knowledge'), os.path.join(here, 'docs')) if os.path.isdir(p)]
    allhits = []
    for p in paths:
        allhits += scan(p, strict)
    if not allhits:
        print('[ok] 未发现公网 IP / 本机路径 / 客户数据文件')
        return 0
    print('[!] 命中 %d 处：' % len(allhits))
    for p, line, kind, val in allhits[:200]:
        print('  %-14s %s:%s  %s' % (kind, p, line, val))
    if len(allhits) > 200:
        print('  ... 其余 %d 处省略' % (len(allhits) - 200))
    return 1

if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
