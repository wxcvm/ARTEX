import { readFileSync } from "node:fs";
import { test } from "node:test";
import assert from "node:assert/strict";
import vm from "node:vm";

const workerCode = readFileSync(new URL("../../../extension/background.js", import.meta.url), "utf8");

function fixture() {
  const config = { serverUrl: "http://127.0.0.1:8787", token: "fixture-token" };
  const stored = {};
  const events = {};
  const calls = { attached: [], detached: [], commands: [], uploads: [], badge: "" };
  let tabId = 7;
  let rejectEnable = false;
  let rejectUpload = false;
  let loseResponse = false;
  let uuidSequence = 0;
  let listener;
  const chrome = {
    storage: {
      local: { get: async () => ({ ...config }) },
      session: { set: async (value) => Object.assign(stored, value) },
    },
    tabs: { query: async () => [{ id: tabId }] },
    debugger: {
      attach: async (target) => { calls.attached.push(target.tabId); },
      detach: async (target) => { calls.detached.push(target.tabId); },
      sendCommand: async (_target, command, params) => {
        calls.commands.push(command);
        if (command === "Network.enable" && rejectEnable) throw new Error("injected enable failure");
        if (command === "Network.getResponseBody") return { body: `response-${params.requestId}`, base64Encoded: false };
        return {};
      },
      onEvent: { addListener: (fn) => { events.event = fn; } },
      onDetach: { addListener: (fn) => { events.detach = fn; } },
    },
    action: {
      setBadgeText: ({ text }) => { calls.badge = text; },
      setBadgeBackgroundColor: () => {},
    },
    runtime: { onMessage: { addListener: (fn) => { listener = fn; } } },
  };
  const fetch = async (_url, options) => {
    if (rejectUpload) throw new Error("injected offline");
    const data = JSON.parse(options.body);
    calls.uploads.push(data.exchanges);
    if (loseResponse) throw new Error("injected lost response");
    return { ok: true, json: async () => ({ ingested: data.exchanges.length, duplicates: 0, failed: 0 }) };
  };
  vm.runInNewContext(workerCode, { chrome, fetch, crypto: { randomUUID: () => `00000000-0000-4000-8000-${String(++uuidSequence).padStart(12, "0")}` } }, { filename: "extension/background.js" });
  const message = (type) => new Promise((resolve, reject) => {
    try {
      if (listener({ type }, {}, resolve) !== true) reject(new Error("message listener not async"));
    } catch (error) { reject(error); }
  });
  const entry = (id) => {
    events.event({ tabId: 7 }, "Network.requestWillBeSent", {
      requestId: String(id),
      request: { method: "GET", url: `https://example.test/${id}`, headers: {} },
    });
    events.event({ tabId: 7 }, "Network.responseReceived", {
      requestId: String(id), response: { status: 200, headers: [] },
    });
    events.event({ tabId: 7 }, "Network.loadingFinished", { requestId: String(id) });
  };
  return {
    calls, stored, events, message, entry,
    setTab: (id) => { tabId = id; },
    setEnableFailure: (value) => { rejectEnable = value; },
    setUploadFailure: (value) => { rejectUpload = value; },
    setLostResponse: (value) => { loseResponse = value; },
  };
}

test("failed Network.enable releases the debugger and never marks capture active", async () => {
  const x = fixture();
  x.setEnableFailure(true);
  const start = await x.message("start");
  assert.equal(start.ok, false);
  assert.match(start.error, /Network.enable failed/);
  assert.deepEqual(x.calls.attached, [7]);
  assert.deepEqual(x.calls.detached, [7]);
  assert.equal((await x.message("getState")).capturing, false);
  x.setEnableFailure(false);
  assert.equal((await x.message("start")).ok, true);
  assert.equal((await x.message("stop")).ok, true);
});

test("capture remains on its original tab until explicitly stopped", async () => {
  const x = fixture();
  assert.equal((await x.message("start")).ok, true);
  x.setTab(8);
  const other = await x.message("start");
  assert.equal(other.ok, false);
  assert.match(other.error, /Already recording tab 7/);
  assert.equal((await x.message("getState")).tabId, 7);
  await x.message("stop");
  assert.equal(x.calls.badge, "");
});

test("stop awaits outstanding responses and drains all batches before detach", async () => {
  const x = fixture();
  await x.message("start");
  for (let id = 0; id < 23; id++) x.entry(id);
  const stopped = await x.message("stop");
  assert.equal(stopped.ok, true);
  assert.equal(stopped.queued, 0);
  assert.equal(x.calls.uploads.flat().length, 23);
  assert.equal((await x.message("getState")).ingested, 23);
  assert.equal((await x.message("getState")).capturing, false);
  assert.ok(x.calls.commands.includes("Network.disable"));
  assert.deepEqual(x.calls.detached, [7]);
});

test("offline stop keeps unsent entries in memory for the next manual retry", async () => {
  const x = fixture();
  await x.message("start");
  x.setUploadFailure(true);
  x.entry(0);
  const stopped = await x.message("stop");
  assert.equal(stopped.ok, true);
  assert.equal(stopped.queued, 1);
  assert.match(x.stored.lastError, /injected offline/);
  x.setUploadFailure(false);
  await x.message("start");
  await x.message("stop");
  assert.equal(x.calls.uploads.flat().length, 1);
  assert.equal((await x.message("getState")).queued, 0);
});

test("unexpected debugger detach stops capture and flushes pending evidence", async () => {
  const x = fixture();
  await x.message("start");
  x.entry(1);
  x.events.detach({ tabId: 7 });
  await new Promise((resolve) => setImmediate(resolve));
  const state = await x.message("getState");
  assert.equal(state.capturing, false);
  assert.match(state.lastError, /debugger detached/);
  assert.equal(x.calls.uploads.flat().length, 1);
});

test("each capture gets a new group and previously queued records keep their original attribution", async () => {
  const x = fixture();
  await x.message("start");
  const first = (await x.message("getState")).captureId;
  x.setUploadFailure(true);
  x.entry(1);
  await x.message("stop");
  assert.equal((await x.message("getState")).queued, 1);
  x.setUploadFailure(false);
  await x.message("start");
  const second = (await x.message("getState")).captureId;
  assert.notEqual(first, second);
  x.entry(2);
  await x.message("stop");
  const entries = x.calls.uploads.flat();
  assert.equal(entries.length, 2);
  assert.equal(entries[0].captureId, first);
  assert.equal(entries[1].captureId, second);
});


test("stop retains requests without terminal events, without inventing successful responses", async () => {
  const x = fixture();
  await x.message("start");
  for (const requestId of ["no-headers", "headers-seen"]) {
    x.events.event({ tabId: 7 }, "Network.requestWillBeSent", {
      requestId, request: { method: "POST", url: `https://example.test/${requestId}`, headers: {}, postData: "started" },
    });
  }
  x.events.event({ tabId: 7 }, "Network.responseReceived", {
    requestId: "headers-seen", response: { status: 202, headers: [{ name: "content-type", value: "text/plain" }] },
  });
  const stopped = await x.message("stop");
  assert.equal(stopped.queued, 0);
  const entries = x.calls.uploads.flat();
  assert.equal(entries.length, 2);
  assert.equal(entries.find(e => e.exchangeId === "no-headers").status, 0);
  assert.equal(entries.find(e => e.exchangeId === "headers-seen").status, 202);
  assert.ok(entries.every(e => e.responseBody === ""));
  assert.ok(entries.every(e => e.requestBody === "started"));
  assert.equal((await x.message("getState")).ingested, 2);
});

test("unexpected debugger detach also retains requests still in flight", async () => {
  const x = fixture();
  await x.message("start");
  x.events.event({ tabId: 7 }, "Network.requestWillBeSent", {
    requestId: "detached", request: { method: "GET", url: "https://example.test/pending", headers: {} },
  });
  x.events.detach({ tabId: 7 });
  await new Promise(resolve => setImmediate(resolve));
  assert.equal(x.calls.uploads.flat().length, 1);
  assert.equal(x.calls.uploads.flat()[0].status, 0);
  assert.equal(x.calls.uploads.flat()[0].exchangeId, "detached");
  assert.match((await x.message("getState")).lastError, /debugger detached/);
});

test("failed network requests are recorded with unobserved status rather than silently discarded", async () => {
  const x = fixture();
  await x.message("start");
  const id = (await x.message("getState")).captureId;
  x.events.event({ tabId: 7 }, "Network.requestWillBeSent", {
    requestId: "failed-1", request: { method: "GET", url: "https://example.test/offline", headers: {} },
  });
  x.events.event({ tabId: 7 }, "Network.loadingFailed", {
    requestId: "failed-1", errorText: "net::ERR_CONNECTION_REFUSED",
  });
  const stopped = await x.message("stop");
  assert.equal(stopped.ok, true);
  const recorded = x.calls.uploads.flat();
  assert.equal(recorded.length, 1);
  assert.equal(recorded[0].captureId, id);
  assert.equal(recorded[0].exchangeId, "failed-1");
  assert.equal(recorded[0].status, 0);
  assert.equal(recorded[0].responseBody, "");
});

test("failed transfer with lost HTTP reply retains stable request identity for manual retry", async () => {
  const x = fixture();
  await x.message("start");
  x.setLostResponse(true);
  x.entry(9);
  const stopped = await x.message("stop");
  assert.equal(stopped.queued, 1);
  assert.equal(x.calls.uploads.length, 1);
  x.setLostResponse(false);
  await x.message("start");
  await x.message("stop");
  assert.equal(x.calls.uploads.length, 2);
  assert.equal(x.calls.uploads[0][0].exchangeId, x.calls.uploads[1][0].exchangeId);
  assert.equal(x.calls.uploads[0][0].captureId, x.calls.uploads[1][0].captureId);
  assert.equal((await x.message("getState")).queued, 0);
});
