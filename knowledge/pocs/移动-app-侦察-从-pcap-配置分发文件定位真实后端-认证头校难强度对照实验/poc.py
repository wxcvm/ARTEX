# 手法一：从抓包解出 App「配置分发文件」定位真实后端
移动 App 常把服务器地址放在一个静态小文件里（App 启动时拉取），该文件响应多为 gzip 的 text/plain，
而抓包里 HTTP 是明文（okhttp 默认对这类配置走 http）。步骤：
1) 先判断 pcap 链路层类型：Android tcpdump 常见 linktype=276（DLT_LINUX_SLL2，**20 字节头**）或 113（16 字节）。
   注意：linktype=276 时若按 16 字节解析，会得到空结果——这是最容易踩的坑。
2) 按 TCP 流重组（四元组双向配对），只保留服务端 payload 非零的流。
   判据：若某流"客户端发了完整 HTTP 请求、服务端只回 RST、payload=0"，说明该端点是死的，
   不要把它当"成功调用"，更不要据此推导签名算法。
3) 对服务端响应体做 gzip/zlib 解压（先处理 Transfer-Encoding: chunked），从中 grep 域名/IP。
4) 配置文件通常形如 JSON：{"baseUrlMap":{"base_urls":[...],"upLoad_urls":[...],"im_urls":[...]},
   "map":{方法名->后端分组}}。这份 map 直接给出全部接口方法名，可省去反编译。
本靶标实测：GET http://203.0.113.9/a.txt → gzip → 完整后端清单，
定位到抓包时正在使用的生产后端（并暴露出一个**上传 Servlet** 与**管理包 example.com.Manager**）。

# 手法二：认证头校验强度对照实验（判断"需登录"是真鉴权还是缺参数）
对同一路径，固定其它所有变量，只改认证头，跑 5 组：
  ① 不带任何头 ② 只带 token ③ token + sign ④ 只带 sign ⑤ 新 timestamp + 旧 sign
判读：
 - 若 5 组响应**字节完全一致** → 网关根本没校验 token/sign，此时不要再逆向签名算法（省下大量时间）。
 - 若 8005「请先登录」的正文里带 `Required Long parameter 'xxx' is not present`（Spring 风格），
   说明真正的门槛是**缺参数**而不是鉴权：补上该参数往往直接返回数据，并顺带形成未授权/越权漏洞。
 - 若错误文案是版本类（如「请下载最新版App」）且与参数无关，说明登录通道被**强制升级拦截**，
   需要逆向最新版客户端协议才能登录，不要在此浪费尝试次数（还可能触发风控）。
本靶标实测：5 组头响应完全一致（api-token/sign 形同虚设），8005 实为缺 userId；
登录接口则被版本拦截，明文口令也登不进。
