#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""LLM API key 零消费只读验活：判定 key 是否有效 / 余额 / 付费层 or 免费层。

铁律：只允许只读端点，绝不调用 /chat/completions、/messages、/generateContent 等任何计费接口。
已知只读端点：
  DeepSeek    GET https://example.com/user/balance   -> is_available + balance_infos(CNY/USD 余额，唯一能读出真金白银的官方只读端点)
  DeepSeek    GET https://example.com/models
  OpenRouter  GET https://example.com/api/v1/key        -> is_free_tier / limit / usage（判免费层）
  OpenAI 兼容 GET {base_url}/models                       -> 仅判 key 有效性
用法：改 PLAN 与 KEY，或从 channels 快照里抽 key 灌入；每 key 最多 1 请求、间隔 2s。
"""
import json, sys, time, urllib.request, urllib.error

PROXY = 'http://127.0.0.1:7892'          # 出境代理；不需要则设 None
ENDPOINTS = [
    ('deepseek',   'https://example.com/user/balance'),
    ('models',     'https://example.com/models'),
    ('openrouter', 'https://example.com/api/v1/key'),
]
KEYS = ['sk-REPLACE_ME']
TAKE_PER_KEY = 1          # 每个 key 请求数上限
SLEEP = 2.0

opener = urllib.request.build_opener(
    urllib.request.ProxyHandler({'http': PROXY, 'https': PROXY}) if PROXY else urllib.request.BaseHandler())

for key in KEYS:
    for name, url in ENDPOINTS[:TAKE_PER_KEY]:
        req = urllib.request.Request(url, headers={
            'Authorization': 'Bearer ' + key,
            'User-Agent': 'Mozilla/5.0', 'Accept': 'application/json'})
        try:
            with opener.open(req, timeout=25) as r:
                print(name, r.status, r.read().decode('utf-8', 'replace')[:800])
        except urllib.error.HTTPError as e:
            print(name, e.code, e.read().decode('utf-8', 'replace')[:400])
        except Exception as e:
            print(name, 'ERR', repr(e))
        time.sleep(SLEEP)
