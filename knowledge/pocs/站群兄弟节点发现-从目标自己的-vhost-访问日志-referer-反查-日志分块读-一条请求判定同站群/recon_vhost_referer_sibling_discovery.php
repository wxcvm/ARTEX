# 站群兄弟节点发现：从目标自己的 vhost 访问日志反查

## 核心思路
一个站群运营方会把几十~几千个"品牌域 + 后台域 + 跳转域"配在同一台（或几台）nginx 上。
nginx 的 `access_log` 是 **per-vhost 一个文件**（宝塔默认 `/www/wwwlogs/<vhost>.log`），
而日志行里有 **Referer** 字段 —— **Referer 里写的全是他们自己别的站的域名**。
⇒ 只要拿到任意一个节点上的文件读权限，就能**从目标自己的数据里挖出目标的全部分身**，
   比扫端口 / 撞子域快一个数量级，而且命中率高（这些都是真实在用的域名，不是猜测）。

## 关键工程问题：日志太大
实测单文件规模：2.28GB / 3.65GB / 1.92GB。**千万不要 `file_get_contents`**（PHP 默认 memory_limit 128M 会直接 Fatal error）。
必须用 `fopen` + `fseek` 只读头部/尾部窗口：

```php
function tailscan($f, $bytes = 1500000) {
    $sz = @filesize($f); if (!$sz) return '';
    $fp = @fopen($f, 'r'); if (!$fp) return '';
    $start = $sz > $bytes ? $sz - $bytes : 0;
    fseek($fp, $start);
    $buf = (string) fread($fp, $bytes);
    fclose($fp);
    if ($start > 0) { $p = strpos($buf, "\n"); $buf = $p === false ? '' : substr($buf, $p + 1); }
    return $buf;
}
$buf = tailscan('/www/wwwlogs/target.com.log', 1500000);
preg_match_all('#"https?://([a-zA-Z0-9._-]+)#', $buf, $m);   // Referer 域名
$refs = array_count_values($m[1]); arsort($refs);
```

## 流程
1. 枚举 `/www/wwwlogs/*.log`（每个是一台 vhost），拿到**该机上所有站**的清单（含已被删 vhost 的死站，日志还在）。
2. 对每个日志 `tailscan` 抓 Referer 域名 + `headscan` 抓最早时间（判历史跨度）。
3. 去重得到"运营方自有域名全家桶"。
4. 用并发 DNS 解析（Python `ThreadPoolExecutor` + `socket.getaddrinfo`，60 线程 200 个域名约 10 秒）。
5. 对解析出的 IP 批量打指纹。
6. 顺手把所有域名用 `Host:` 头打到**本机 127.0.0.1**（`curl CURLOPT_RESOLVE '<d>:443:127.0.0.1'`），
   就能知道"这台机器在替哪些域服务" —— 这是**域名↔IP↔vhost 的一键对齐**，不用猜。

## ★ 一条请求判定"是不是同站群"（指纹判据）
```
GET /example.com?g=Home&m=Article&a=edit_cont&id=5   无凭据
→ HTTP 200 + 正文【恰好 3 字节 "502"】
```
这是 ThinkPHP 站群 `Home/ArticleAction::_initialize()` 的 `else{ echo 502; die; }` 分支。
**9 台不同 IP 全部命中**，比任何 header/server 指纹都硬。

## 实战产出（本演练，176 = 203.0.113.91）
```
① 从 /www/wwwlogs 里读到 11 个 vhost 日志，单文件最大 3.65GB
② tailscan 抽 Referer → 汇总出 1808 个唯一域名，整理出 219 个"运营方自有二级域"
③ 批量解析：68/219 存活（其余 NXDOMAIN 或停在 Sedo 停放 91.195.240.x）
④ 从中定位出全新节点（全部实测 200 + 同一套 TP3）：
     203.0.113.101 / 203.0.113.99 / 203.0.113.100   ← 全新 /24
     203.0.113.45 / 203.0.113.46 / 203.0.113.44
     203.0.113.111
⑤ 域名生成规律（可直接用来预测未公开节点）：
     <品牌前缀+随机盐>.<随机10位>.xyz      例 target407.example.com、target408.example.com
     sdfjdsfd-{1,2,3}.<202xxxxxx>.xyz      序号型，数量最大
     tianzs-{1,2,3}.<2026xxxxx>.xyz
     target409.example.com
     <品牌>.one/target410.example.com、mp.<域>      （target411.example.com / target412.example.com 全解析到同一 IP）
     短跳转页：<任意 .xyz> 返回 90 字节 <script>window.location.href="https://<后台域>/#<品牌代号>#"</script>
              ⇒ **一个 90 字节页面就把"这个域名属于哪个品牌"告诉我们了**
```

## 附带收获（同一台机器上）
- 日志首行经常就是**另一个攻击者**的第一条记录（本例 `203.0.113.22` 从 6/11 起持续打 `/ddbb_admlogin.php` 与 phpMyAdmin 登录页）
  ⇒ 报告里可如实写"该资产同时被多方攻击"
- 日志文件大小本身就是**业务规模证据**（2.28GB / 3.65GB / 1.92GB 三个品牌域）