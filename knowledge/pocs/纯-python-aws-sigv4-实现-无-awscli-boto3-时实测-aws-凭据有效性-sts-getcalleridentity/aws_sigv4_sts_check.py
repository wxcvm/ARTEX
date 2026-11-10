import hashlib, hmac, datetime, urllib.request, urllib.error

def _sign(key, msg):
    return hmac.new(key, msg.encode("utf-8"), hashlib.sha256).digest()

def _sig_key(secret, datestamp, region, service):
    k = _sign(("AWS4" + secret).encode("utf-8"), datestamp)
    k = _sign(k, region); k = _sign(k, service); k = _sign(k, "aws4_request")
    return k

def sts_get_caller_identity(ak, sk, session_token=None, region="us-east-1", timeout=25):
    host, service = "sts.amazonaws.com", "sts"
    now = datetime.datetime.now(datetime.timezone.utc)
    amzdate = now.strftime("%Y%m%dT%H%M%SZ"); datestamp = now.strftime("%Y%m%d")
    body = "Action=GetCallerIdentity&Version=2011-06-15"
    payload_hash = hashlib.sha256(body.encode()).hexdigest()
    headers = {"content-type": "application/x-www-form-urlencoded",
               "host": host, "x-amz-date": amzdate}
    if session_token:
        headers["x-amz-security-token"] = session_token
    names = sorted(headers)
    canonical_headers = "".join("%s:%s\n" % (n, headers[n].strip()) for n in names)
    signed_headers = ";".join(names)
    canonical_request = "\n".join(["POST", "/", "", canonical_headers, signed_headers, payload_hash])
    scope = "%s/%s/%s/aws4_request" % (datestamp, region, service)
    sts = "\n".join(["AWS4-HMAC-SHA256", amzdate, scope,
                     hashlib.sha256(canonical_request.encode()).hexdigest()])
    sig = hmac.new(_sig_key(sk, datestamp, region, service),
                   sts.encode("utf-8"), hashlib.sha256).hexdigest()
    send = dict(headers)
    send["Authorization"] = ("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s"
                             % (ak, scope, signed_headers, sig))
    req = urllib.request.Request("https://" + host + "/", data=body.encode(), headers=send, method="POST")
    op = urllib.request.build_opener(urllib.request.ProxyHandler({}))  # 必须显式禁用本机代理
    try:
        with op.open(req, timeout=timeout) as r:
            return r.status, r.read(20000).decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read(20000).decode("utf-8", "replace")
    except Exception as e:
        return -1, str(e)

if __name__ == "__main__":
    import sys
    print(sts_get_caller_identity(*sys.argv[1:4]))
