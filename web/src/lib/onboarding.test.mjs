import test from "node:test";
import assert from "node:assert/strict";
import { onboardingGuides, onboardingSteps, modelErrorGuide, testAndEnableModel } from "./onboarding.ts";

test("four steps each have a complete six-part guide", () => {
  assert.equal(onboardingSteps.length, 4);
  for (const guide of Object.values(onboardingGuides)) {
    assert.deepEqual(Object.keys(guide), ["purpose", "materials", "actions", "example", "success", "failure"]);
    for (const text of Object.values(guide)) assert.ok(text.length > 20);
  }
});
test("model test, save and activation run in order", async () => {
  const calls = [];
  await testAndEnableModel({ test: async () => { calls.push("test"); return { ok: true, reply: "hello" }; }, save: async () => { calls.push("save"); return "7"; }, activate: async id => { calls.push(id); } });
  assert.deepEqual(calls, ["test", "save", "7"]);
});
for (const result of [{ ok: false, error: "401" }, { ok: true, reply: " " }]) test("failed or empty model reply never persists configuration: " + JSON.stringify(result), async () => {
  let saved = false;
  await assert.rejects(testAndEnableModel({ test: async () => result, save: async () => { saved = true; }, activate: async () => assert.fail("must not activate") }));
  assert.equal(saved, false);
});
test("save failure never activates; activation failure propagates", async () => {
  await assert.rejects(testAndEnableModel({ test: async () => ({ ok: true, reply: "hi" }), save: async () => { throw Error("save failed"); }, activate: async () => assert.fail() }), /save failed/);
  await assert.rejects(testAndEnableModel({ test: async () => ({ ok: true, reply: "hi" }), save: async () => "1", activate: async () => { throw Error("activation failed"); } }), /activation failed/);
});
test("actionable errors distinguish credentials, URL and availability", () => {
  assert.match(modelErrorGuide("401"), /API Key/);
  assert.match(modelErrorGuide("404"), /地址/);
  assert.match(modelErrorGuide("timeout"), /服务/);
});
