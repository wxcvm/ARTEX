// ARTEX Traffic Capture — background service worker.
//
// Captures the network activity of one tab through the chrome.debugger API
// (Network domain) and pushes text request/response pairs to the ARTEX
// ingest endpoint in batches. State lives in chrome.storage.session so the
// popup and the (short-lived) service worker share one view; the queue is
// deliberately in-memory — a browser restart drops unsent entries, which is
// the acceptable failure mode for a manual capture session.

const FLUSH_AT = 10; // entries per batch
const MAX_BODY = 256 * 1024; // matches the server-side ingest cap
const SKIP_MIME = ["image/", "font/", "audio/", "video/", "application/octet-stream"];

let activeFlush = null;
const pendingResponses = new Set();
let stopping = false;

let session = {
  capturing: false,
  tabId: null,
  captureId: "",
  queue: [], // pending ExternalExchange objects
  requests: new Map(), // debugger requestId -> partial request info
  ingested: 0,
  failed: 0,
  lastError: "",
};

async function saveState() {
  await chrome.storage.session.set({
    capturing: session.capturing,
    tabId: session.tabId,
    captureId: session.captureId,
    ingested: session.ingested,
    failed: session.failed,
    lastError: session.lastError,
    queued: session.queue.length,
  });
}

async function getConfig() {
  const { serverUrl = "", token = "" } = await chrome.storage.local.get(["serverUrl", "token"]);
  return { serverUrl: serverUrl.replace(/\/+$/, ""), token };
}

async function getActiveTabId() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  return tab ? tab.id : null;
}

async function start() {
  if (session.capturing) {
    const tabId = await getActiveTabId();
    return tabId === session.tabId
      ? { ok: true }
      : { ok: false, error: `Already recording tab ${session.tabId}. Stop before switching tabs.` };
  }
  if (stopping) return { ok: false, error: "Previous capture is stopping." };
  const { serverUrl, token } = await getConfig();
  if (!serverUrl || !token) {
    return { ok: false, error: "Set the ARTEX URL and token first." };
  }
  const tabId = await getActiveTabId();
  if (tabId == null) return { ok: false, error: "No active tab." };
  try {
    await chrome.debugger.attach({ tabId }, "1.3");
  } catch (e) {
    return { ok: false, error: `attach failed: ${e.message}` };
  }
  try {
    await chrome.debugger.sendCommand({ tabId }, "Network.enable", {});
  } catch (e) {
    // Attaching is a separate operation; a failed enable must release it.
    try { await chrome.debugger.detach({ tabId }); } catch (_) {}
    return { ok: false, error: `Network.enable failed: ${e.message}` };
  }
  // Unsent entries are intentionally kept in memory until flushed; starting
  // another tab must not silently discard evidence from the previous capture.
  session = { ...session, capturing: true, tabId, captureId: crypto.randomUUID(), requests: new Map(), ingested: 0, failed: 0, lastError: "" };
  chrome.action.setBadgeText({ text: "REC" });
  chrome.action.setBadgeBackgroundColor({ color: "#c2410c" });
  await saveState();
  return { ok: true };
}

async function stop() {
  if (stopping) return { ok: false, error: "Previous capture is stopping." };
  stopping = true;
  const tabId = session.capturing ? session.tabId : null;
  // Freeze ingestion before waiting for outstanding body reads and upload.
  session.capturing = false;
  session.tabId = null;
  const unfinished = retainUnfinishedRequests();
  try {
    await Promise.allSettled([...pendingResponses, ...unfinished]);
    if (tabId != null) {
      try {
        await chrome.debugger.sendCommand({ tabId }, "Network.disable", {});
        await chrome.debugger.detach({ tabId });
      } catch (_) {
        // The tab may already be gone. Local capture still stops.
        try { await chrome.debugger.detach({ tabId }); } catch (_) {}
      }
    }
    await flush();
    chrome.action.setBadgeText({ text: "" });
    await saveState();
    return { ok: true, queued: session.queue.length };
  } finally {
    stopping = false;
  }
}

function skipBody(mimeType) {
  const mt = (mimeType || "").toLowerCase();
  return SKIP_MIME.some((p) => mt.startsWith(p));
}

function truncate(s) {
  if (s.length <= MAX_BODY) return s;
  return s.slice(0, MAX_BODY) + "\n... [truncated by extension]";
}

async function queueExchange(info, response, skipResponseBody = false) {
  // Chrome CDP Network.Headers is a JSON object, not a HAR header array.
  const respHeaders = Object.fromEntries(
    Object.entries(response.headers || {}).map(([name, value]) => [name.toLowerCase(), String(value)])
  );
  const mime = respHeaders["content-type"] || "";

  let reqBody = truncate(info.postData || "");
  let respBody = "";
  if (!skipResponseBody && !skipBody(mime)) {
    try {
      const body = await chrome.debugger.sendCommand(
        { tabId: session.tabId },
        "Network.getResponseBody",
        { requestId: info.requestId }
      );
      if (body && !body.base64Encoded) {
        respBody = truncate(body.body || "");
      }
    } catch (_) {
      // body already discarded by the browser; capture the pair without bodies
    }
  }

  session.queue.push({
    captureId: info.captureId,
    exchangeId: String(info.requestId),
    method: info.method,
    url: info.url,
    status: response.status || 0,
    requestHeaders: info.headers || {},
    requestBody: reqBody,
    responseHeaders: respHeaders,
    responseBody: respBody,
  });
  if (session.queue.length >= FLUSH_AT) await flush();
}

async function flush() {
  // A manual flush or Stop must drain every pending batch, not just the first
  // ten. Concurrent producers share one flush and do not duplicate uploads.
  if (activeFlush) return activeFlush;
  activeFlush = flushBatches().finally(() => { activeFlush = null; });
  return activeFlush;
}

async function flushBatches() {
  if (session.queue.length === 0) return;
  const { serverUrl, token } = await getConfig();
  if (!serverUrl || !token) {
    session.lastError = "missing server URL or token";
    await saveState();
    return;
  }
  while (session.queue.length > 0) {
    const batch = session.queue.splice(0, FLUSH_AT);
    try {
      const res = await fetch(`${serverUrl}/api/traffic/ingest`, {
        method: "POST",
        headers: { "Authorization": `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ exchanges: batch }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${(await res.text()).slice(0, 200)}`);
      const out = await res.json();
      if ((out.ingested || 0) + (out.duplicates || 0) + (out.failed || 0) !== batch.length) {
        throw new Error("server receipt does not account for every exchange");
      }
      session.ingested += (out.ingested || 0) + (out.duplicates || 0);
      session.failed += out.failed || 0;
      session.lastError = out.failed ? `${out.failed} entries rejected by the server` : "";
    } catch (e) {
      // Preserve unsent entries for manual retry without persisting bodies to disk.
      session.queue = batch.concat(session.queue);
      session.lastError = e.message;
      break;
    }
  }
  await saveState();
}

function trackPending(promise) {
  pendingResponses.add(promise);
  void promise.finally(() => pendingResponses.delete(promise));
}

// Stop and debugger detach have no terminal Network event for requests still
// in flight. Keep the attempted exchange, but never invent a response status
// or body. A status observed earlier is preserved as such.
function retainUnfinishedRequests() {
  const unfinished = [...session.requests.values()];
  session.requests.clear();
  return unfinished.map((info) =>
    queueExchange(info, info.response || { status: 0, headers: {} }, true)
      .catch((e) => { session.lastError = e.message; })
  );
}

chrome.debugger.onEvent.addListener((source, method, params) => {
  if (!session.capturing || !source.tabId || source.tabId !== session.tabId) return;
  if (method === "Network.requestWillBeSent") {
    session.requests.set(params.requestId, {
      requestId: params.requestId,
      captureId: session.captureId,
      method: params.request.method,
      url: params.request.url,
      headers: params.request.headers,
      postData: params.request.postData || "",
    });
  } else if (method === "Network.responseReceived") {
    const info = session.requests.get(params.requestId);
    if (info) info.response = params.response;
  } else if (method === "Network.loadingFinished" || method === "Network.loadingFailed") {
    const info = session.requests.get(params.requestId);
    session.requests.delete(params.requestId);
    if (!info) return;
    // Preserve failed/cancelled requests without inventing a 200 response.
    // If a status was actually observed before the failure, retain it.
    const failed = method === "Network.loadingFailed";
    const response = info.response || { status: 0, headers: {} };
    const pending = queueExchange(info, response, failed).catch((e) => {
      session.lastError = e.message;
      return saveState();
    });
    trackPending(pending);
  }
});

chrome.debugger.onDetach.addListener((source) => {
  if (session.capturing && source.tabId === session.tabId) {
    // DevTools opened, or the tab navigated somewhere undebaggable: stop
    // cleanly instead of silently dropping events.
    session.capturing = false;
    session.tabId = null;
    const unfinished = retainUnfinishedRequests();
    session.lastError = "debugger detached (DevTools opened?)";
    chrome.action.setBadgeText({ text: "" });
    void Promise.allSettled([...pendingResponses, ...unfinished]).then(async () => {
      await flush();
      // A successful upload must not erase why capture unexpectedly ended.
      if (!session.lastError) session.lastError = "debugger detached (DevTools opened?)";
      await saveState();
    });
  }
});

chrome.runtime.onMessage.addListener((msg, _sender, sendResponse) => {
  (async () => {
    switch (msg?.type) {
      case "getState": {
        sendResponse({
          capturing: session.capturing,
          tabId: session.tabId,
          captureId: session.captureId,
          ingested: session.ingested,
          failed: session.failed,
          queued: session.queue.length,
          lastError: session.lastError,
          configured: !!(await getConfig()).serverUrl,
        });
        return;
      }
      case "start":
        sendResponse(await start());
        return;
      case "stop":
        sendResponse(await stop());
        return;
      case "flushNow":
        await flush();
        sendResponse(await stateSnapshot());
        return;
      default:
        sendResponse({ ok: false, error: "unknown message" });
    }
  })();
  return true; // async sendResponse
});

async function stateSnapshot() {
  return {
    capturing: session.capturing,
    tabId: session.tabId,
    captureId: session.captureId,
    ingested: session.ingested,
    failed: session.failed,
    queued: session.queue.length,
    lastError: session.lastError,
  };
}
