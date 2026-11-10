# 现象：FOFA 技能报 key 缺失，但 key 其实已配置
# 实测：$env:FOFA_KEY 为空，而 <HOME>\workspace\.env 内有 FOFA_KEY=...

# ---- 每次开工前先执行这段（PowerShell）----
Get-Content "$env:USERPROFILE\workspace\.env" | Where-Object { $_ -match '^FOFA_KEY=' } | ForEach-Object { $env:FOFA_KEY = ($_ -split '=',2)[1] }

# 验证（只打印长度，不泄露 key）
if ($env:FOFA_KEY) { "FOFA_KEY loaded len=$($env:FOFA_KEY.Length)" } else { "MISSING" }

# ---- 之后即可正常调用 fofa_client.py ----
# python runs/fofa_client.py 'domain="example.com"' --pages 2 --size 100 --interval 2.0
