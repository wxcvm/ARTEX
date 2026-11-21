import test from "node:test";
import assert from "node:assert/strict";
import { createSettingsMutationQueue } from "./settings-mutation-queue.ts";

test("serializes independent changes and keeps pending selection over an older response", async () => {
  let saved = { planner: false, worker: false };
  const observed = [];
  let releaseFirst;
  const firstGate = new Promise((resolve) => { releaseFirst = resolve; });
  const queue = createSettingsMutationQueue(async (patch) => {
    observed.push(patch);
    if (observed.length === 1) await firstGate;
    saved = { ...saved, ...patch };
    return { ...saved };
  });

  const planner = queue.save({ planner: true });
  const worker = queue.save({ worker: true });
  await Promise.resolve();
  assert.equal(observed.length, 1, "only one write should be in flight");
  assert.deepEqual(queue.overlay(saved), { planner: true, worker: true });

  releaseFirst();
  const oldResponse = await planner;
  assert.deepEqual(queue.overlay(oldResponse), { planner: true, worker: true });
  const finalResponse = await worker;
  assert.deepEqual(finalResponse, { planner: true, worker: true });
  assert.deepEqual(queue.overlay(finalResponse), finalResponse);
  assert.equal(observed.length, 2);
});

test("a newer toggle of the same field remains visible until its write completes", async () => {
  let releaseFirst;
  const gate = new Promise((resolve) => { releaseFirst = resolve; });
  let writes = 0;
  let stored = { enabled: false };
  const queue = createSettingsMutationQueue(async (patch) => {
    if (++writes === 1) await gate;
    stored = { ...stored, ...patch };
    return { ...stored };
  });

  const first = queue.save({ enabled: true });
  const second = queue.save({ enabled: false });
  assert.equal(queue.overlay(stored).enabled, false);
  releaseFirst();
  assert.equal(queue.overlay(await first).enabled, false);
  assert.deepEqual(await second, { enabled: false });
  assert.equal(writes, 2);
});

test("a failed write does not block later saves or leave a ghost overlay", async () => {
  let count = 0;
  const queue = createSettingsMutationQueue(async (patch) => {
    if (++count === 1) throw new Error("simulated persistence failure");
    return { value: false, ...patch };
  });

  const failed = queue.save({ value: true });
  const next = queue.save({ value: false });
  await assert.rejects(failed, /simulated persistence failure/);
  const result = await next;
  assert.deepEqual(result, { value: false });
  assert.deepEqual(queue.overlay(result), result);
  assert.equal(count, 2);
});
