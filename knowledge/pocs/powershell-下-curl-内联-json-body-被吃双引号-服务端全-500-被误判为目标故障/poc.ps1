## 现象与判据
在 Windows PowerShell 里用 curl.exe 内联传 JSON body：

    curl.exe -s -X POST "http://TARGET/api/auth/login" -H "Content-Type: application/json" -d '{"username":"admin","password":"x"}'

PowerShell 的单引号**不会**保护内部双引号不被 curl.exe 处理，实测 curl 实际发出的
Content-Length=29 而正确 body 应为 37 字节 —— 双引号被吞掉，请求体畸形。

**判据（三步确认）**：
1. 用 `-v` 看实际发送的 `> Content-Length:`，与手算的 body 长度比对，不一致即中招：
   `curl.exe -s -v -X POST "$T/api/auth/login" -H "Content-Type: application/json" -d '{...}' 2>&1 | Select-String '^> Content-Length'`
2. 畸形体在 Express(express.json) 这类框架下常表现为**统一 500 Internal Server Error**
   （而非标准的 400 Bad Request），长度固定（如 21 字节 "Internal Server Error"）。
3. 换文件方式传同一 body，若状态码/错误文案变化，则 100% 是本地传参问题，不是目标故障。

## 危害（为什么值得单独记一条）
本坑会把"所有 POST 接口都 500"误读为**目标后端依赖故障**（本例中曾被误判为「Redis 挂了 → 应用不可用」），
从而放弃整条测试线、误报目标不可用、或写出错误的排除结论。

## 正确做法（Windows 下统一使用）
把 body 写文件，再用 `--data-binary @file`：

    Set-Content -Path body.json -Value '{"username":"admin","password":"x"}' -NoNewline -Encoding ascii
    curl.exe -s -i -X POST "$T/api/auth/login" -H "Content-Type: application/json" --data-binary "@body.json"

Python 更稳（json.dumps 不受 shell 影响）：
    python -c "import json,urllib.request as u; b=json.dumps({'username':'admin','password':'x'}).encode(); r=u.urlopen(u.Request('http://T/api/auth/login',data=b,headers={'Content-Type':'application/json'})); print(r.status,r.read())"

## 同类风险
- bash/zsh 下单引号是安全的，本坑**特指 PowerShell**（cmd.exe 双引号转义规则也不同）。
- 任何"需要 JSON / 长 body / 含特殊字符"的请求都应走文件或 Python，不要内联。
- PowerShell 里 URL 含 `&` 也必须加引号，否则被当作命令分隔符。