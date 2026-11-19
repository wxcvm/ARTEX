import assert from "node:assert/strict";
import test from "node:test";
import { reportLibraryView, reportOffsetAfterDeletion } from "./report-library-state.ts";

const ready = { key: "all:0", revision: 0, status: "ready", items: ["report-1"] };

test("a previous filter or page cannot supply rows or actions for the current view", () => {
  for (const key of ["md:0", "all:50"]) {
    const view = reportLibraryView(ready, key, 0);
    assert.deepEqual(view.items, []);
    assert.equal(view.status, "loading");
    assert.equal(view.canAct, false);
    assert.equal(view.empty, false);
  }
});

test("refresh and same-page errors may retain rows as stale, but never enable actions", () => {
  const refresh = reportLibraryView(ready, ready.key, 1);
  const failed = reportLibraryView({ ...ready, revision: 1, status: "error", error: "503" }, ready.key, 1);
  for (const view of [refresh, failed]) {
    assert.equal(view.stale, true);
    assert.equal(view.canAct, false);
    assert.equal(view.empty, false);
  }
  assert.equal(failed.error, "503");
});

test("only a successful empty response displays the no-results state", () => {
  const empty = { ...ready, items: [] };
  assert.equal(reportLibraryView(empty, ready.key, 0).empty, true);
  assert.equal(reportLibraryView({ ...empty, status: "error" }, ready.key, 0).empty, false);
  assert.equal(reportLibraryView(null, ready.key, 0).empty, false);
});

test("deleting the last row steps back one page without producing a negative offset", () => {
  assert.equal(reportOffsetAfterDeletion(100, 1), 50);
  assert.equal(reportOffsetAfterDeletion(50, 2), 50);
  assert.equal(reportOffsetAfterDeletion(0, 1), 0);
});
