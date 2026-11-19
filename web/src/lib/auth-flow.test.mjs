import assert from "node:assert/strict";
import test from "node:test";
import { resolveSignInDestination, shouldClearSession } from "./auth-flow.ts";

test("wrong sign-in or current password does not reset the session or reload the form", () => {
  assert.equal(shouldClearSession("/auth/login", 401, "Incorrect username or password", null, null), false);
  assert.equal(shouldClearSession("/auth/change-password", 401, "Current password is wrong", "same", "same"), false);
});

test("a genuine protected 401 clears only the session that issued the request", () => {
  assert.equal(shouldClearSession("/auth/change-password", 401, "Unauthorized", "same", "same"), true);
  assert.equal(shouldClearSession("/llm", 401, "Unauthorized", "old", "new"), false);
  assert.equal(shouldClearSession("/llm", 503, "Unavailable", "same", "same"), false);
});

test("a model status error falls back to the guide only for the current session", async () => {
  assert.equal(await resolveSignInDestination("current", () => "current", async () => ({ configured: true })), "/function/tasks");
  assert.equal(await resolveSignInDestination("current", () => "current", async () => ({ configured: false })), "/welcome");
  assert.equal(await resolveSignInDestination("current", () => "current", async () => { throw new Error("Offline"); }), "/welcome");
  let token = "expired";
  assert.equal(await resolveSignInDestination("expired", () => token, async () => {
    token = null;
    throw new Error("Unauthorized");
  }), null);
});

test("a late success for a replaced session does not navigate", async () => {
  let token = "old";
  const destination = resolveSignInDestination("old", () => token, async () => {
    await Promise.resolve();
    return { configured: true };
  });
  token = "new";
  assert.equal(await destination, null);
});
