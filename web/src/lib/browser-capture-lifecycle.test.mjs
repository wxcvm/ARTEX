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
    return { ok: true, json: async () => ({ ingested: data.exchanges.length, failed: 0 }) };
  };
  vm.runInNewContext(workerCode, { chrome, fetch }, { filename: "extension/background.js" });
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
