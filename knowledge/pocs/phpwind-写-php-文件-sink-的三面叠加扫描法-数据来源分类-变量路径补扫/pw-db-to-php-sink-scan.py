#!/usr/bin/env python3
# 换分类维度：不看「内容是否含请求变量」，看「内容的数据来源」
# 来源A=请求变量(已判死) B=数据库字段/远程响应 ★本轮目标 C=配置缓存回读 D=常量
import os, re, json
ROOT = r"runs/pw-admin-sinks/utf8"
WRITE_FNS = ["writeover","file_put_contents","fwrite","fputs","fopen","rename","copy"]

def extract_calls(text):
    """括号平衡抓取，返回 [(lineno, func, argstring)]"""
    out = []
    for m in re.finditer(r'\b(' + '|'.join(WRITE_FNS) + r')\s*\(', text, re.I):
        i = m.end(); depth = 1; start = i; n = len(text); in_s = None
        while i < n and depth > 0:
            c = text[i]
            if in_s:
                if c == '\\': i += 2; continue
                if c == in_s: in_s = None
            else:
                if c in ('"',"'"): in_s = c
                elif c == '(': depth += 1
                elif c == ')':
                    depth -= 1
                    if depth == 0: break
            i += 1
        out.append((text.count('\n',0,m.start())+1, target22.example.com(1), text[start:i]))
    return out

# 三个扫描面必须叠加，缺一漏类：
# ① writeover/file_put_contents 的路径参数【字面量】含 .php
# ② fopen+fwrite/fputs 的路径是【变量或方法调用】—— 例 $this->filepath、$this->getCacheFileName()
#    这类逃过所有「路径含 .php」的正则，是最大的盲区
# ③ rename/copy 顺带核对

# 变量赋值链回溯（递归 2-3 层），自动判定是否有转义函数
ESCAPERS = ['pw_var_export','db_cv','addslashes','addcslashes','key_cv',
            'pwEscape','intval','htmlspecialchars','preg_quote','strip_tags']
def audit(rel, expr, depth=0, seen=None, before_line=None, root=""):
    if seen is None: seen = set()
    if depth > 3: return set()
    found = set()
    lines = open(os.path.join(root, rel), encoding='utf-8', errors='replace').read().split('\n')
    for v in set(re.findall(r'\$([A-Za-z_]\w*)', expr)):
        if (rel, v) in seen: continue
        seen.add((rel, v))
        pat = re.compile(r'\$' + re.escape(v) + r'\s*\.?=\s*(?!=)')
        for i, ln in enumerate(lines):
            if before_line and i+1 > before_line: continue
            if pat.search(ln):
                found |= {e for e in ESCAPERS if e in ln}
                found |= audit(rel, ln.split('=',1)[-1], depth+1, seen, i+1, root)
    return found
