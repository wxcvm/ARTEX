// Bring the socket up under capture, dump the complete STOMP conversation,
// then probe cross-role destinations through the official stompClient.
const fs = require('fs');
const CDP_PORT = 9222, OUT = 'C:\\rt\\dgha\\runs\\ocs_ws_conv.jsonl';
const sleep = ms => new Promise(r => setTimeout(r, ms));

(async () => {
  const list = await (await fetch(`http://127.0.0.1:${CDP_PORT}/json/list`)).json();
  const page = list.find(t => t.type === 'page');
  const ws = new WebSocket(page.webSocketDebuggerUrl);
  let id = 0; const pending = new Map();
  const send = (m, p = {}) => { const i = ++id; ws.send(JSON.stringify({ id: i, method: m, params: p })); return new Promise(r => pending.set(i, r)); };
  await new Promise(r => ws.onopen = r);
  const out = fs.createWriteStream(OUT, { flags: 'w' });
  const seen = [];
  ws.onmessage = ev => {
    const m = JSON.parse(ev.data);
    if (m.id && pending.has(m.id)) { pending.get(m.id)(m); pending.delete(m.id); return; }
    if (m.method && m.method.startsWith('Network.webSocket')) {
      const p = m.params || {}, r = p.response || {};
      const rec = { m: m.method.replace('Network.webSocket', ''), url: p.url || '', op: r.opcode, payload: r.payloadData !== undefined ? String(r.payloadData) : undefined, err: p.errorMessage };
      if (rec.payload === undefined && !rec.url && !rec.err) return;
      seen.push(rec); out.write(JSON.stringify(rec) + '\n');
      console.log(`[${rec.m}] ${JSON.stringify(rec.payload !== undefined ? rec.payload : (rec.url || rec.err))}`);
    }
  };
  const evalx = async (expr) => {
    const r = await send('Runtime.evaluate', { expression: expr, returnByValue: true, awaitPromise: true, allowUnsafeEvalBlockedByCSP: true });
    return r.result?.exceptionDetails ? 'EXC ' + String(r.result.exceptionDetails.exception?.description).slice(0, 300) : r.result?.result?.value;
  };
  const q = async (l, e) => { const v = await evalx(e); console.log(`\n### ${l}\n` + String(typeof v === 'string' ? v : JSON.stringify(v)).slice(0, 1300)); return v; };

  await send('Network.enable', { maxResourceBufferSize: 1 << 26, maxTotalBufferSize: 1 << 28 });
  await send('Runtime.enable');
  const V = 'document.querySelector("#app").__vue__';

  console.log('=== bringing socket up under capture ===');
  await q('call transceiver.init()', `(()=>{const t=${V}.$root.transceiver;try{t.init();return "init() ok"}catch(e){return "ERR "+e}})()`);
  await sleep(9000);
  await q('state', `(()=>{const t=${V}.$root.transceiver;return JSON.stringify({connected:t.stompClient&&t.stompClient.connected,url:t.socket&&t.socket.url,subs:t.stompClient&&Object.keys(t.stompClient.subscriptions||{})})})()`);

  // ---- probe: send to visitor destination (baseline) ----
  console.log('\n=== baseline send to /app/visitor/ack ===');
  await q('send /app/visitor/ack', `(()=>{const t=${V}.$root.transceiver;try{t.stompClient.send("/app/visitor/ack",{},JSON.stringify({msgType:101,test:"baseline"}));return "sent"}catch(e){return "ERR "+e}})()`);
  await sleep(4000);

  // ---- probe: cross-role destinations ----
  const DESTS = [
    '/app/staff/ack', '/app/staff/list', '/app/staff/sessions', '/app/staff/auth',
    '/app/visitor/list', '/app/visitor/sessions', '/app/visitor/history',
    '/app/sessions', '/app/session/list', '/app/staff',
  ];
  for (const d of DESTS) {
    console.log(`\n=== probe dest ${d} ===`);
    await q('send ' + d, `(()=>{const t=${V}.$root.transceiver;try{t.stompClient.send("${d}",{},JSON.stringify({msgType:102,cid:101096}));return "sent"}catch(e){return "ERR "+e}})()`);
    await sleep(1800);
  }

  // ---- probe: subscribe to staff/staff-wide topics ----
  const TOPICS = ['/topic/staff/receive', '/user/topic/staff/receive', '/topic/staff', '/topic/visitor/receive', '/topic/all', '/user/topic/visitor/receive'];
  for (const tp of TOPICS) {
    console.log(`\n=== subscribe ${tp} ===`);
    await q('subscribe ' + tp, `(()=>{const t=${V}.$root.transceiver;try{t.stompClient.subscribe("${tp}",function(m){window.__leak=(window.__leak||[]);window.__leak.push(m.body)});return "subscribed"}catch(e){return "ERR "+e}})()`);
    await sleep(1500);
  }
  await sleep(9000);
  await q('any leaked messages?', `JSON.stringify((window.__leak||[]).slice(0,5))`);
  await q('subs now', `(()=>{const t=${V}.$root.transceiver;return JSON.stringify(Object.keys(t.stompClient.subscriptions||{}))})()`);

  console.log('\n=== total frames: ' + seen.length + ' ===');
  process.exit(0);
})();
