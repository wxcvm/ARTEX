import assert from "node:assert/strict";
import test from "node:test";

import { parseIDs, parseKeywords } from "../app/(main)/system/notify/_components/channel-fields.ts";

test("notification ID filters accept English and Chinese commas", () => {
  assert.deepEqual(parseIDs("1,2，3 4\n5"), [1, 2, 3, 4, 5]);
});

test("notification keyword filters keep multiword names and accept Chinese commas", () => {
  assert.deepEqual(parseKeywords("SQL injection，XSS,\nSSRF"), ["SQL injection", "XSS", "SSRF"]);
});
