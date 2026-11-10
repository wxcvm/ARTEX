#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Decode the STATICAES key material hidden in the tails of the resource responses.

Algorithm (webpack module 1483 parseHiddenData):
    n = len(bytes); a = min(n, 2048); r = n - a
    out[i] = 90 ^ bytes[r+i]          # XOR 0x5A over the LAST min(n,2048) bytes
    text = utf8(out)
    markers: [P1] [P2] [S] [AES]  -> value = text after LAST occurrence of the tag, trimmed
Reads runs/tgt2905_res_*.json (fetched with curl) when present.
"""
import base64
import glob
import json
import os
import sys

try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

TAGS = ["[P1]", "[P2]", "[S]", "[AES]"]
found = {}

for f in sorted(glob.glob(r"C:\rt\dgha\runs\tgt2905_res_*.json")):
    try:
        j = json.load(open(f, encoding="utf-8"))
    except Exception as e:
        print("%-34s not json (%s)" % (os.path.basename(f), e))
        continue
    b64 = (j.get("data") or {}).get("body")
    if not b64:
        print("%-34s no .data.body" % os.path.basename(f))
        continue
    raw = base64.b64decode(b64 + "=" * (-len(b64) % 4))
    n = len(raw)
    a = min(n, 2048)
    r = n - a
    out = bytes(90 ^ b for b in raw[r:r + a])
    txt = out.decode("utf-8", "replace")
    hits = [t for t in TAGS if t in txt]
    print("%-34s bytes=%-7d window=%-5d off=%-6d tags=%s" % (os.path.basename(f), n, a, r, hits or "-"))
    for t in hits:
        i = txt.rfind(t)
        v = txt[i + len(t):].strip()
        found[t] = v
        print("      %-6s -> %r (len %d)" % (t, v[:300], len(v)))
    if hits:
        open(f.replace(".json", "_tail.txt"), "w", encoding="utf-8").write(txt)

print("\n=== recovered ===")
for t in TAGS:
    v = found.get(t)
    print("%-6s len=%-4d %r" % (t, len(v) if v else 0, (v[:300] if v else "")))

if found:
    json.dump(found, open(r"C:\rt\dgha\runs\tgt2905_static_key.json", "w", encoding="utf-8"),
              ensure_ascii=False, indent=1)
    print("\nwritten runs/tgt2905_static_key.json")
    if found.get("[AES]"):
        print("AES KEY  =", found["[AES]"])
    if found.get("[S]"):
        print("SALT     =", found["[S]"])
    print("P1+P2    =", (found.get("[P1]", "") + found.get("[P2]", ""))[:200])
