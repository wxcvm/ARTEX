import assert from "node:assert/strict";
import test from "node:test";
import { downloadFilename } from "./download-filename.ts";

const cases = [
  ["UTF-8 wins over legacy", `attachment; filename="fallback.md"; filename*=UTF-8''%E6%BC%8F%E6%B4%9E%E6%8A%A5%E5%91%8A.md`, "漏洞报告.md"],
  ["extended first and mixed case", `attachment; FILENAME*=utf-8'zh'%E6%8A%A5%E5%91%8A%20final.md; filename="fallback.md"`, "报告 final.md"],
  ["spaces", `attachment; filename="report final.md"`, "report final.md"],
  ["quoted semicolon", `attachment; filename="report; final.md"`, "report; final.md"],
  ["escaped quotes", String.raw`attachment; filename="report \"final\".md"`, 'report "final".md'],
  ["quoted parameter is not another parameter", `attachment; other="x; filename*=UTF-8''wrong.md"; filename="right.md"`, "right.md"],
  ["bad escape falls back", `attachment; filename="legacy.md"; filename*=UTF-8''bad%ZZ.md`, "legacy.md"],
  ["bad UTF-8 falls back", `attachment; filename="legacy.md"; filename*=UTF-8''bad%C3%28.md`, "legacy.md"],
  ["unsupported charset falls back", `attachment; filename="legacy.md"; filename*=ISO-8859-1''%E9.md`, "legacy.md"],
  ["blank extended falls back", `attachment; filename="legacy.md"; filename*=UTF-8''%20`, "legacy.md"],
  ["empty header", "", "download"],
  ["empty filename", `attachment; filename=""`, "download"],
  ["bare token and trailing parameters", `attachment; filename=report.csv; size=123`, "report.csv"],
  ["literal plus and percent", `attachment; filename*=UTF-8''50%25+%2B%2520.txt`, "50%++%20.txt"],
  ["Unicode legacy compatibility", `attachment; filename="报告.md"`, "报告.md"],
  ["control characters fall back", `attachment; filename="legacy.md"; filename*=UTF-8''report%00.md`, "legacy.md"],
  ["path components are stripped", `attachment; filename*=UTF-8''..%2F..%5Creport.md`, "report.md"],
  ["dot name is unusable", `attachment; filename=".."`, "download"],
];

for (const [description, header, expected] of cases) {
  test(description, () => assert.equal(downloadFilename(header, "download"), expected));
}
