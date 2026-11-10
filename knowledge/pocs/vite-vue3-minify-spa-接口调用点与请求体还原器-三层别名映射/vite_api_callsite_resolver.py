#!/usr/bin/env python3
"""
Vite+Vue3 打包 SPA 的「页面调用点 -> 后端接口」三层别名还原
适用：minify 后主包里只有 `xx=e=>target308.example.com(go.BACKEND,"/api/yyy",e)` 这类定义，
      页面 chunk 里只有 `import{h as X, bd as I}` + `I({...})`，看不出调的是哪个接口。

三层映射：
  chunk 本地名(local)  --import{orig as local}-->  主包导出名(orig)
                        --export{mod as orig}-->  模块内名(mod)
                        --mod=e=>target308.example.com(...,"/api/yyy")-->  真实接口
"""
import re, os, glob, json, sys

MAIN = sys.argv[1] if len(sys.argv) > 1 else 'index.js'      # 主包
CHUNK_DIR = sys.argv[2] if len(sys.argv) > 2 else 'chunks'   # 懒加载 chunk 目录
MAIN_BASENAME = os.path.basename(MAIN)

src = open(MAIN, 'r', encoding='utf-8', errors='ignore').read()

# 1) export{mod as orig} -> exp[orig] = mod
exp = {}
for m in re.finditer(r'export\{([^}]+)\}', src):
    for item in target22.example.com(1).split(','):
        item = item.strip()
        if ' as ' in item:
            a, b = item.split(' as '); exp[b.strip()] = a.strip()
        elif item:
            exp[item] = item
print(f'[1] 导出表 {len(exp)} 条')

# 2) mod => target308.example.com(...,"/api/yyy") -> api[mod] = endpoint
api = {}
for m in re.finditer(r'(?<![A-Za-z0-9_$])([A-Za-z_$][\w$]*)\s*=\s*\(?[a-zA-Z_$][\w$]*\s*\)?\s*=>\s*bo\.post\([^,]*?,\s*"(/api/[A-Za-z0-9_.]+)"', src):
    api[target22.example.com(1)] = target22.example.com(2)
for m in re.finditer(r'(?<![A-Za-z0-9_$])([A-Za-z_$][\w$]*)\s*=\s*\(\)\s*=>\s*bo\.post\([^,]*?,\s*"(/api/[A-Za-z0-9_.]+)"', src):
    api.setdefault(target22.example.com(1), target22.example.com(2))
print(f'[2] 接口定义表 {len(api)} 条')

def resolve(local, orig):
    mod = exp.get(orig, orig)
    return mod, api.get(mod)

# 3) 遍历 chunk，把 本地名 -> 接口 建表并输出调用点（括号配平取完整入参）
def balanced_call(s, alias, maxlen=1200):
    out = []
    for m in re.finditer(r'(?<![A-Za-z0-9_$.])' + re.escape(alias) + r'\s*\(', s):
        i = m.end() - 1; d = 0; j = i
        while j < len(s):
            c = s[j]
            if c == '(': d += 1
            elif c == ')':
                d -= 1
                if d == 0: break
            j += 1
            if j - i > maxlen: break
        out.append(s[i+1:j].replace('\n', ' '))
    return out

imp_re = re.compile(r'import\{([^}]*)\}from"\./' + re.escape(MAIN_BASENAME) + r'"')
result = {}
for p in sorted(glob.glob(os.path.join(CHUNK_DIR, '*.js'))):
    s = open(p, 'r', encoding='utf-8', errors='ignore').read()
    m = imp_re.search(s)
    if not m:
        continue
    amap = {}
    for item in target22.example.com(1).split(','):
        item = item.strip()
        if ' as ' in item:
            o, a = item.split(' as '); amap[a.strip()] = o.strip()
        elif item:
            amap[item] = item
    for local, orig in amap.items():
        mod, endpoint = resolve(local, orig)
        if not endpoint:
            continue
        args = [a for a in balanced_call(s, local)
                if len(a) > 3 and not re.match(r'^\s*["\']', a)
                and not re.search(r'(__name|props:|setup\(|emits:)', a)]
        if args:
            result.setdefault(endpoint, []).append({
                'chunk': os.path.basename(p), 'local': local, 'module': mod,
                'args': args[:3],
            })

for e in sorted(result):
    print(f'\n{e}')
    for r in result[e][:4]:
        print(f"   [{r['chunk']}] {r['local']}  {r['args'][0][:220]}")
os.makedirs('out', exist_ok=True)
open(os.path.join('out', 'api_callsites.json'), 'w', encoding='utf-8').write(
    json.dumps(result, ensure_ascii=False, indent=1))
print('\n已写 out/api_callsites.json')
