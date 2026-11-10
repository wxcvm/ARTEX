// 在 Node 里执行页面混淆脚本，还原被混淆的 window 全局配置
//
// 【原理】
// 这类前端用「字符串表 + 索引解密函数」形态混淆（_0x599bf4(idx) / _0x232979(idx,'key')），
// 底层是标准 rc4/xxtea 变体。但它们**必须在页面里执行**才能拿到明文 ——
// 所以最省事的做法不是静态去混淆，而是**伪装运行环境、把脚本原样跑一遍**，
// 让它自己把 window 配置对象填好，然后直接读对象。
//
// 【实测】example.com 首页 54503 字节，一次跑通还原出 window._env_ 全部 90+ 字段
//        含 24 个后端服务地址（含彩票 adminapi）、148 个 App 包名、3 个 S3 桶
//
const fs = require('fs');

const html = fs.readFileSync('example.com', 'utf8');

// 1. 抓出所有 <script> 内联块，挑最长的（含解密表与解密函数的那块）
const blocks = [...html.matchAll(/<script[^>]*>([\s\S]*?)<\/script>/g)].map(m => m[1]);
let code = blocks.filter(b => b.includes('_0x23ce7e') || b.includes('_0x8d4bd5'))
                .sort((a, b) => b.length - a.length)[0]
           || blocks.sort((a, b) => b.length - a.length)[0];

// 2. 只保留到 window[...]={...} 配置对象结束（去掉后面的 DOM 操作，避免报错）
const idx = code.indexOf('window[_0x');
if (idx > -1) {
  const end = code.indexOf('};', idx);
  code = code.slice(0, end > -1 ? end + 2 : code.length);
}

// 3. 造假运行环境，把脚本原样跑一遍，读回 window
const sandbox = {
  window: {}, document: { head: {}, body: {} },
  navigator: { userAgent: 'node' }, console,
};
const fn = new Function('window', 'document', 'navigator', 'console',
                        code + '\n;return window;');
const w = fn(sandbox.window, sandbox.document, sandbox.navigator, console);

// 4. 打印还原结果
for (const [k, v] of Object.entries(w)) {
  if (v && typeof v === 'object') {
    console.log('\n[' + k + ']');
    for (const [k2, v2] of Object.entries(v)) console.log('   ' + k2 + ' = ' + v2);
  } else {
    console.log(k + ' = ' + v);
  }
}
