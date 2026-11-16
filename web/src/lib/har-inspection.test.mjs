import assert from "node:assert/strict";
import test from "node:test";
import { inspectHar } from "./har-inspection.ts";

test("HAR inventory distinguishes missing status, invalid status and omitted text", () => {
  const request = { method: "GET", url: "https://example.invalid/private?token=SECRET" };
  const result = inspectHar({ log: { entries: [
    { request, startedDateTime: "2026-11-16T10:00:00Z", response: { status: 200, content: { text: "" } } },
    { request, response: { status: 0, content: { size: 100 } } },
    { request, startedDateTime: "2026-11-16T10:00:00+08:00", response: { status: 204 } },
    { request, startedDateTime: "2026-11-16T10:00:00", response: { status: "200", content: { text: "PRIVATE", encoding: "base64" } } },
    { request },
    { request: { method: "GET", url: "file:///tmp/fixture" } },
    { request: { method: "GET", url: "not a URL" } },
    null,
  ] } });
  assert.deepEqual(result, {
    entries: 8, httpEntries: 5, invalidEntries: 3,
    observedStatus: 2, unobservedStatus: 2, invalidStatus: 1,
    responseTextPresent: 2, responseTextNotEmbedded: 3, encodedResponseText: 1,
    requestPostDataNotEmbedded: 0, entriesWithDuplicateHeaders: 0,
    captureTimePresent: 2, captureTimeUnknown: 3,
    declaredPages: 0, invalidPages: 0, entriesWithKnownPage: 0, entriesWithUnknownPage: 0,
    entriesWithAmbiguousPage: 0, entriesWithoutPage: 5,
  });
  assert.doesNotMatch(JSON.stringify(result), /SECRET|PRIVATE|example\.invalid/);
});

test("lossy header maps and form metadata are identified without exposing values", () => {
  const request = {
    method: "POST", url: "https://example.invalid/",
    headers: [{ name: "Cookie", value: "SECRET" }, { name: "cookie", value: "PRIVATE" }],
    postData: { params: [{ name: "password", value: "PRIVATE" }] },
  };
  const result = inspectHar({ log: { entries: [
    { request, response: { status: 400, headers: [{ name: "Set-Cookie", value: "x" }, { name: "set-cookie", value: "y" }] } },
    { request: { ...request, headers: [], postData: { text: "" } }, response: { status: 404, content: { text: "" } } },
  ] } });
  assert.equal(result.entriesWithDuplicateHeaders, 1);
  assert.equal(result.requestPostDataNotEmbedded, 1);
  assert.equal(result.responseTextPresent, 1);
  assert.doesNotMatch(JSON.stringify(result), /password|SECRET|PRIVATE/);
});

test("invalid HAR containers error and empty archives remain explicit", () => {
  for (const value of [null, {}, { log: {} }, { log: { entries: {} } }]) {
    assert.throws(() => inspectHar(value), /log.entries/);
  }
  const result = inspectHar({ log: { entries: [] } });
  assert.equal(result.entries, 0);
  assert.equal(result.httpEntries, 0);
  assert.equal(result.responseTextNotEmbedded, 0);
});


test("browser page attribution distinguishes missing, unknown and ambiguous page references without leaking metadata", () => {
  const request = { method: "GET", url: "https://example.invalid/private" };
  const result = inspectHar({log:{pages:[{id:"PAGE_PRIVATE",title:"SECRET"},{id:"DUP"},{id:"DUP"},{id:""},{id:123}],entries:[
    {request,pageref:"PAGE_PRIVATE"}, {request,pageref:"DUP"}, {request,pageref:"UNKNOWN"},
    {request}, {request,pageref:42}, {request,pageref:""}, {request:{method:"GET",url:"file:///tmp/test"},pageref:"PAGE_PRIVATE"},
  ]}});
  assert.equal(result.declaredPages,5);
  assert.equal(result.invalidPages,2);
  assert.equal(result.entriesWithKnownPage,1);
  assert.equal(result.entriesWithAmbiguousPage,1);
  assert.equal(result.entriesWithUnknownPage,2);
  assert.equal(result.entriesWithoutPage,2);
  assert.equal(result.entriesWithKnownPage+result.entriesWithAmbiguousPage+result.entriesWithUnknownPage+result.entriesWithoutPage,result.httpEntries);
  assert.doesNotMatch(JSON.stringify(result), /PAGE_PRIVATE|SECRET|DUP|UNKNOWN|example\\.invalid/);
});
