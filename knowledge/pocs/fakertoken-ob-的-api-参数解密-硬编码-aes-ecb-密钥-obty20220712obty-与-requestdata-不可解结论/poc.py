# 结论（实测）
fakerToken 返回的 `loginUrl` 里 `api=<base64>` 参数 = **apiDomain 的 AES-128-ECB/Pkcs7 密文**，
密钥硬编码在 OB 客户端主包：`CRYPTO_DECRYPT_KEY_URL_API = "OBTY20220712OBTY"`（16 字节 ASCII）。

# 一行解密（python，需 pycryptodome）
```python
import base64
from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad
key = b"OBTY20220712OBTY"
s   = "DSTPJZSZ56drn/CnMOQuyiMigwz58KbNupOQCX96zpE="   # loginUrl 里的 api= 值
pt  = unpad(AES.new(key, AES.MODE_ECB).decrypt(base64.b64decode(s)), 16)
print(pt.decode())        # -> https://target307.example.com
```

# 一行加密（反向：构造该参数）
`base64.b64encode(AES.new(key, AES.MODE_ECB).encrypt(pad(b"https://target317.example.com", 16))).decode()`

# 同族密钥（同一份客户端里）
- `panda1234_1234ob` → OSS 数据解密：`get_oss_decrypt_str(H, K || DECRYPT_KEY)`（AES-ECB/Pkcs7）
- `OBTY20220712OBTY` → API 参数/请求体加密：`api_encrypt(H)`（AES-ECB/Pkcs7）
两者都是从 `CryptoJS.enc.Utf8.parse("...")` 直接取的裸 16 字节，**无 KDF、无盐、无 IV**。

# 定位方法（换站群复用）
1. 抓客户端主包：`grep -o 'Utf8\.parse("[^"]*")' bundle.js | sort -u` → 拿到全部 AES 密钥。
2. `grep -o 'CryptoJS\.AES\.\(encrypt\|decrypt\)[^;]*' bundle.js` → 看每个密钥对应哪个字段（api_encrypt / get_oss_decrypt_str / ...）。
3. 对响应里形如 `api=`、`xxx=<base64>` 的短参数（32 字节密文 → 44 字符 base64）直接套上面一行解密。

# 未解开项（避免重复投入）
`requestData`（704 字符 → base64 528 字节、熵 7.55、33 个 AES 块）**不在**这套密钥下：
已穷尽 AES{ECB,CBC,CFB,OFB,CTR}×{两密钥}×{iv=key,iv=0}×{Pkcs7,NoPadding} + raw/zlib/gzip（含偏移扫魔数）+ 双 base64 + 40 候选密钥字典 → 0 命中。
RB88 与 OB 两个前端产物里**只有上述两个 AES 密钥**；RB88 侧对 OB 响应只做 `atob→pako.inflate→decodeURIComponent→JSON`（无 AES）。
→ 该常量的密钥在签发侧，离线不可解，需签发侧泄露或更完整的客户端产物。

# 边界
纯离线，未对厂商或目标发起任何请求；未做在线爆破/密钥恢复攻击。
