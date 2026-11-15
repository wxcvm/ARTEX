import assert from "node:assert/strict";
import test from "node:test";
import { beginFileRead, updateFileEditor } from "./skill-file-editor.ts";

test("late reads and errors cannot replace another file or a new visit to the same file", () => {
  const a = { skill: "a", path: "SKILL.md" };
  const b = { skill: "b", path: "SKILL.md" };
  const again = { ...a };
  const current = beginFileRead(b);
  assert.equal(updateFileEditor(current, a, { status: "ready", content: "old" }), current);
  assert.equal(updateFileEditor(current, a, { status: "error" }), current);
  const reopened = beginFileRead(again);
  assert.equal(updateFileEditor(reopened, a, { status: "ready", content: "old" }), reopened);
  assert.equal(updateFileEditor(null, a, { status: "ready" }), null);
});

test("late save completion cannot clear another file's edits or saving state", () => {
  const a = { skill: "a", path: "SKILL.md" };
  const b = { skill: "b", path: "SKILL.md" };
  const current = updateFileEditor(beginFileRead(b), b, { status: "ready", content: "new", dirty: true, saving: true });
  assert.equal(updateFileEditor(current, a, { dirty: false, saving: false }), current);
  assert.deepEqual(updateFileEditor(current, b, { dirty: false, saving: false }), { ...current, dirty: false, saving: false });
});
