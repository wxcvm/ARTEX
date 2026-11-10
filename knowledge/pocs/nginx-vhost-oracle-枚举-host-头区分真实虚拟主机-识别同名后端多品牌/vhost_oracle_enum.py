#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""nginx vhost oracle 枚举：用 Host 头探测共享同一 IP 的真实虚拟主机"""
import subprocess, re, time
SRC="203.0.113.88"; PROXY="http://127.0.0.1:7892"
def probe(host):
    cmd=["curl.exe","-x",PROXY,"-k","-s","-i","--max-time","12",
         "--connect-to","probe.invalid:443:%s:443"%SRC,"-H","Host: "+host,
         "https://probe.invalid/"]
    r=subprocess.run(cmd,capture_output=True,timeout=30).stdout.decode("utf-8","replace")
    t=r.replace("\r\n","\n"); bs=re.split(r"(?m)^HTTP/1\.1 (\d{3})",t)
    code=bs[-2] if len(bs)>=3 else "?"; body=bs[-1].partition("\n\n")[2].strip() if len(bs)>=3 else t
    return code, len(body), body[:120]
# oracle：不存在的 Host -> 200 且 body 为空；真实 vhost -> 非空（404/200 页面/302）
print("baseline(不存在):", probe("zzz-not-exist-12345.example"))
for w in ["admin","agent","chat","www","api","manage","ht","m","oa"]:
    for d in ["example.com","example.com","example.com"]:
        c,n,s=probe("%s.%s"%(w,d))
        if not (c=="200" and n==0):
            print("REAL  %-28s %s len=%d %s"%("%s.%s"%(w,d),c,n,s.replace("\n"," ")))
        time.sleep(1.3)
