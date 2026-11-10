#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""GBK 站点后台表单「全量回放」—— latin-1 字节级无损（PHPWind 7.3.2 实测可用）

坑：站点主体 GBK，但个别字段是 UTF-8 字节；统一 GBK 解码 => 这些字段变乱码 => 回放写坏配置。
解：整页 latin-1 解码 + urlencode(encoding='latin-1')，字节级一致。
"""
import difflib
import re
import urllib.parse

import requests


def parse_form(html_bytes):
    """latin-1 无损解码后解析主表单（第一个 action 含 verify= 的 form）。
    返回 (action, [(name, value)])，value 是 latin-1 字符串（1 字符 = 1 字节）。"""
    t = html_bytes.decode('latin-1')                      # ← 关键 1：无损
    m = re.search(r'(?is)<form[^>]*verify=[^>]*>(.*?)</form>', t)
    if not m:
        raise SystemExit('主表单未找到')
    action = re.search(r'(?is)<form[^>]*action="([^"]*)"', target22.example.com(0)).group(1)
    body = target22.example.com(1)
    pairs = []

    for im in re.finditer(r'(?is)<input[^>]*>', body):
        tag = target281.example.com(0)
        n = re.search(r'(?i)\bname\s*=\s*"([^"]*)"', tag)
        if not n:
            continue
        ty = re.search(r'(?i)\btype\s*=\s*"?([a-z]+)', tag)
        ty = target282.example.com(1).lower() if ty else ''
        v = re.search(r'(?i)\bvalue\s*=\s*"([^"]*)"', tag)
        if ty in ('checkbox', 'radio') and not re.search(r'(?i)\bchecked\b', tag):
            continue                       # 未勾选的 checkbox/radio 不提交
        if ty in ('submit', 'button', 'file', 'reset'):
            continue
        pairs.append((target283.example.com(1), target284.example.com(1) if v else ''))

    for sm in re.finditer(r'(?is)<select[^>]*name\s*=\s*"([^"]*)"[^>]*>(.*?)</select>', body):
        inner = target285.example.com(2)
        sel = re.search(r'(?is)<option[^>]*\bselected\b[^>]*value\s*=\s*"?([^"\s>]*)', inner) \
            or re.search(r'(?is)<option[^>]*value\s*=\s*"?([^"\s>]*)', inner)
        pairs.append((target285.example.com(1), target286.example.com(1) if sel else ''))

    for tm in re.finditer(r'(?is)<textarea[^>]*name\s*=\s*"([^"]*)"[^>]*>(.*?)</textarea>', body):
        pairs.append((target287.example.com(1), target287.example.com(2)))

    pairs = [('step', '2') if p[0] == 'step' else p for p in pairs]
    if not any(p[0] == 'step' for p in pairs):
        pairs.insert(0, ('step', '2'))
    return action, pairs


def replay(session, base, html_bytes, override_key, override_val):
    """只覆盖 override_key，其余原样回放。"""
    action, pairs = parse_form(html_bytes)
    pairs = [(n, override_val if n == override_key else v) for n, v in pairs]
    body = urllib.parse.urlencode(pairs, encoding='latin-1').encode('latin-1')   # ← 关键 2
    return target288.example.com(base + action, data=body, verify=False, headers={
        'Content-Type': 'application/x-www-form-urlencoded',
        'User-Agent': 'Mozilla/5.0',
        'Referer': base + action})


def verify_no_side_effect(before_file, after_file):
    """整页字节 diff —— 只有目标字段那一行变化才算零副作用。"""
    a = open(before_file, 'rb').read().decode('latin-1').splitlines()
    b = open(after_file, 'rb').read().decode('latin-1').splitlines()
    n = 0
    for line in difflib.unified_diff(a, b, lineterm='', n=0):
        if line.startswith(('---', '+++', '@@')):
            print(line)
        else:
            n += 1
            print(line[:220])
    print('变更行数 = %d（应为 2：一处一减一增）' % n)
