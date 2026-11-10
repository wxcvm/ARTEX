# 目标形态
App 用统一入口 `POST /api`（或 `POST /<业务路径>`），业务参数整体加密后作为**单个参数**提交：
```
POST /api  Content-Type: application/x-www-form-urlencoded
param=<加密后的业务参数>
```
真正的业务路由（如 `user/loginV591`、`user/recharge/getList`）写在**密文内部**，因此"必须用最新版 App"这类版本拦截可以完全绕开。

# 还原步骤
## 1. 先确认 APK 是否有业务明文 dex
```
python -c "import zipfile;z=zipfile.ZipFile('x.apk');print([n for n in z.namelist() if n.endswith('.dex')])"
python key_grep.py        # 逐个 dex 统计业务串（包名 / 协议关键字 / login 接口名）出现次数
```
加固包的特征：classes*.dex 里**只有第三方 SDK**，完全没有业务包名字符串 → 业务 dex 被抽取加密，需换用已 dump 的设备侧 dex（如 `patched59_signed.apk`、自设备导出的 `installed_app.apk`）。

## 2. 用自写 DEX 反汇编器做字符串 xref
```
python dex_xref.py <业务dex> "AES/CBC"          # 谁用了这个变换模式
python dex_xref.py <业务dex> "param"            # 谁把密文塞进 param
python dex_xref.py <业务dex> "<登录接口名>"      # 登录参数从哪来
```
关键实现要点（否则解析必错）：
- `class_data_item` 的 4 个 size **连续读在最前面**，再依次读字段/方法列表；
- `method_idx_diff`/`field_idx_diff` **累加**（direct 与 virtual 各自从 0 起）；
- `access_flags` 是 uleb，`ACC_CONSTRUCTOR=0x10000` 要 3 字节；
- 字符串索引 >0xFFFF 走 `const-string/jumbo(0x1b)` + 4 字节索引。

## 3. 读工具类的字段初值拿密钥
`<clinit>`/`<init>` 里通常直接 `const-string` 写密钥与 IV，再 `iput-object` 存字段。找到 `Cipher.getInstance(...)` 调用点即可确认：
```
const-string  v3, "AES/CBC/PKCS5Padding"
invoke-static {v3}, javax.crypto.Cipher->getInstance
...
new-instance  v5, javax.crypto.spec.SecretKeySpec
invoke-direct {v5, v4, v6}, SecretKeySpec-><init>(byte[], "AES")
new-instance  v6, javax.crypto.spec.IvParameterSpec
invoke-direct {v6, v7}, IvParameterSpec-><init>(byte[])
invoke-virtual {v3, v1, v5, v6}, javax.crypto.Cipher->init   # v1 = ENCRYPT_MODE
```
实例：`target23.example.com.boxbox.l.h4` 的字段 `b`/`c` 即 Key/IV。

## 4. 用真机抓包做**零请求**验证
拿到 key/iv 后先解密已有 pcap 的全部密文（自写 TCP 流重组），明文能读出可读 JSON 才算解对；注意抓包主机上的失败流（RST、服务端 payload=0）**其客户端密文依然有效**，可用于验证。
```
python pcap_aes_verify.py     # 7/7 成功 → {"apiUserId":"618189","uri":"user/getMyInfoV310",...}
```

## 5. 明文格式与 Base64 变体
- 明文常是 fastjson `toJSONString(obj, PrettyFormat)`：`{\n\t"k":v,\n\t...\n}`（数字不带引号）；
- 密文常做 URL-safe 变体替换：`+→-`、`/→_`、`=→@`（部分实现只做前两个），写解密时三种都试。

## 6. 登录参数的构成规律（常见三层校验）
| 层 | 字段 | 绕过/取值方法 |
|---|---|---|
| 口令 | `password` = MD5(明文)，个别为 SHA1 | 无盐哈希可离线反查 |
| 设备 | `deviceId` / `machineCode` | 常与抓包请求头 `uuid` 同值；**可从服务端未授权接口直接读回档案中的 machineCode** |
| 验证码 | `imgCode`（阿里云/易盾 validate） | 最难，需 SDK 环境；服务端若只对"风控触发"的请求要它，则设备指纹对齐后它成为唯一门槛 |

**差分实验定位卡点**（省尝试次数）：用一个必然失败的输入做对照——
- `phone=目标账号 + 正确口令` → 若返回"账号登录异常/验证码错误"而非"账户密码错误"，即说明**口令已通过**；
- 空 `imgCode` 通常先报"缺少参数"，此类失败**不消耗**密码尝试次数，可以放心探参数名。

# 复用价值
任何用"整体加密参数 + 业务 uri 在密文里"的 App 后端，都可照此流程在 1 小时内拿到可编程客户端，从而摆脱 App 版本限制、自由构造任意业务请求（越权/未授权测试的前提）。
