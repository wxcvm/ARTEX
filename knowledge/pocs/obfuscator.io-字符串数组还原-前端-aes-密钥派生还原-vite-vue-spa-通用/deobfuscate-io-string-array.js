#!/usr/bin/env node
/*
 * obfuscator.io 字符串数组还原器
 * 原理：bundle 把字符串放进数组 $d()，用 RC4 变体函数 Sa(e,t) 按索引解密，
 *       并由一个 IIFE 重排数组顺序（顺序错则解出乱码）。
 * 正确做法：把 [重排IIFE] + [function Sa] + [function $d] 三段原样提取出来，
 *          在 Node 里 require 执行，让混淆代码自己解自己。
 */
const fs = require("fs");

// —— 第 1 步：定位三个片段的边界 ——
const src = fs.readFileSync(process.argv[2], "utf8");
const i_sa       = src.indexOf("function Sa(e,t)");              // 解码器
const i_shuffle  = src.lastIndexOf("(function(e,t){", i_sa);     // 数组重排 IIFE
const i_d        = src.indexOf("function $d()", i_sa);           // 字符串数组
const i_dend     = src.indexOf("$d()}", i_d) + 4;

const out = [
  src.slice(i_shuffle, i_sa),      // 重排 IIFE（必须包含，否则数组顺序错误）
  src.slice(i_sa, i_d),            // Sa
  src.slice(i_d, i_dend),          // $d
  "module.exports={Sa,$d};",
].join("\n");
fs.writeFileSync("./_decoder.js", out);

// —— 第 2 步：调用约定 ——
// 源码中的包装函数形式各异，先反推出真实调用：
//   Ht(e,t)   => Sa(e+877, t)
//   o(a,r)    => Sa(r-79,  a)     （注意参数序相反！）
//   n(o,a)    => Ht(o-163,a)
// 关键：Sa 内部第一件事是 e = e - 768，所以「想取数组下标 k」= Sa(k+768, key)
const { Sa } = require("./_decoder.js");
function Ht(e, t) { return Sa(e - -877, t); }

// —— 第 3 步：批量还原 ——
// 从目标类/函数的源码片段里正则抓所有 Ht(数字,"key") 调用并解码
const seg = src.slice(src.indexOf("class xN"), src.indexOf("class xN") + 1600);
for (const c of seg.matchAll(/Ht\((-?\d+)\s*,\s*"([^"]*)"\)/g)) {
  console.log(c[0], "=>", JSON.stringify(Ht(parseInt(c[1]), c[2])));
}

/* ================= 实战产出（某博彩平台 lobby 前端） =================
 * 还原出的加密类 xN 全部密钥：
 *   getSK()  mm = "tgJt2ol!^E" + ")E8(%P"                  => tgJt2ol!^E)E8(%P
 *   IV       Ql = "5421698523" + "412578"                   => 5421698523412578
 *   oss 密钥 Xl = atob("dGhhbmtzLHBpZzJjbG91ZA==")          => thanks,pig4cloud
 *   RK       gm = "fc361cdb77"+"0aebc2126c"+"c0dac989c8"+"96"
 * 加密算法：
 *   encrypt(data,key) = AES-CBC + ZeroPadding, iv=Ql, 输出 base64
 *   ossDecrypt(ct)    = AES-ECB, key=Xl       ← 用它解开了 /cocos/config_data.json
 * 密钥派生 gk(token, a, username, keyType)：
 *   k1: EN(token)        = md5(token + md5(token)).slice(0,16)     // md5hex 小写
 *   k2: UN(token,user)   = md5(token + username).slice(0,16)
 * 请求体：{"encryptString":"<base64>"}
 */

/* ================= 附带产出：AES-ECB 解密配置下发文件 =================
 * import base64; from Crypto.Cipher import AES
 * ct = base64.b64decode(open("config_data.json").read().strip())
 * pt = AES.new(b"thanks,pig4cloud", AES.MODE_ECB).decrypt(ct)
 * pt = pt[:-pt[-1]] if 0 < pt[-1] <= 16 else pt
 * // 解出明文含 api_domain / oss_domain / h5_domain / siteCode 等完整基础设施清单
 */
