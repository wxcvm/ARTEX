#!/usr/bin/env python3
# 盲探 UEditor 零鉴权面（纯 GET，不写任何数据）
import http.client, sys

def probe(ip, host, path, timeout=12):
    c = http.client.HTTPConnection(ip, 80, timeout=timeout)
    c.request("GET", path, headers={"Host": host, "User-Agent": "Mozilla/5.0", "Connection": "close"})
    r = c.getresponse(); b = r.read(1500); c.close()
    return r.status, len(b), b.decode('utf-8', 'replace')

IP, HOST = "TARGET_IP", "TARGET_HOST"
BASE = "/Public/admin/ueditor/php/controller.php"
for act in ("config", "listimage", "listfile"):
    st, ln, body = probe(IP, HOST, "%s?action=%s" % (BASE, act))
    print("[%s] %s len=%d %s" % (act, st, ln, body[:160]))
for d in ("/ueditor/", "/ueditor/php/", "/ueditor/php/upload/image/", "/ueditor/php/upload/file/"):
    st, ln, _ = probe(IP, HOST, d)
    print("[dir] %-34s %s  %s" % (d, st, "存在(403=禁列目录)" if st == 403 else ("不存在" if st == 404 else "?")))
print("其他路径：/Public/admin/ueditor/php/config.json（直读配置）")
